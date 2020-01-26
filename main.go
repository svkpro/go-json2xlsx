package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

type StringService interface {
	Count(string) (int, error)
}

type stringService struct{}

func (stringService) Count(s string) (int, error) {
	if s == "" {
		return 0, errors.New("The string can't be empty!")
	}
	return len(s), nil
}

type countRequest struct {
	String string `json:"string"`
}

type countResponse struct {
	Amount int `json:"amount"`
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		amount, err := svc.Count(req.String)

		return countResponse{amount}, err
	}
}

func main() {
	svc := stringService{}
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
		opts...,
	)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
