package main

// The profiles is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods("POST").Path("/count").Handler(httptransport.NewServer(
		endpoints.PostCountEndpoint,
		decodePostCountRequest,
		postEncodeResponse,
		options...,
	))

	return router
}

func decodePostCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	validate = validator.New()
	var request countRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if err := ValidateStruct(request); err != nil {
		return nil, err
	}

	return request, nil
}

func postEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type countRequest struct {
	S string `json:"s" validate:"required,min=1,max=10"`
}

type countResponse struct {
	Amount int `json:"amount"`
}
