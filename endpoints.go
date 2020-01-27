package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	PostCountEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		PostCountEndpoint:   MakePostCountEndpoint(s),
	}
}

func MakePostCountEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		amount, err := s.Count(req.S)

		return countResponse{amount}, err
	}
}




