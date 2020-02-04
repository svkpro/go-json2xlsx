package main

// The profiles is just over HTTP, so we just have a single transport.go.

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"
)

type xlsxResponse struct {
	Data string `json:"data"`
}

type XlsxRequest struct {
	Data XlsxPayloadData `json:"data"`
}

type XlsxPayloadData struct{
	Sheet string `json:"Sheet"`
	Value string `json:"Value"`
}

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods("GET").Path("/xlsx/{file}").Handler(httptransport.NewServer(
		endpoints.GetXlsxEndpoint,
		decodeGetXlsxRequest,
		getEncodeResponse,
		options...,
	))

	router.Methods("POST").Path("/make-xlsx").Handler(httptransport.NewServer(
		endpoints.PostMakeXlsxEndpoint,
		decodePostMakeXlsxRequest,
		postEncodeResponse,
		options...,
	))

	return router
}

func decodeGetXlsxRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := mux.Vars(r)

	return q, nil
}

func getEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	path, ok := response.(string)
	if !ok {
		panic("Bad response generation!")
	}
	f, err := os.Open(path)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		bufferedReader := bufio.NewReader(f)
		w.Header().Set("Content-Type", "application/vnd.ms-excel")
		w.Header().Add("Content-Transfer-Encoding", "binary")
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}

	return nil
}

func postEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(response)
}

func decodePostMakeXlsxRequest(_ context.Context, r *http.Request) (interface{}, error) {
	validate = validator.New()
	var request XlsxRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if err := ValidateStruct(request); err != nil {
		return nil, err
	}

	return request, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
