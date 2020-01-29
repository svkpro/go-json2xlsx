package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	PostMakeXlsxEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		PostMakeXlsxEndpoint:   MakePostMakeXlsxEndpoint(s),
	}
}

func MakePostMakeXlsxEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(xlsxRequest)
		fp, err:= s.MakeXlsx(req.Data)

		return xlsxResponse{fp}, err
	}
}




