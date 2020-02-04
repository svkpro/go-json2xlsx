package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-kit/kit/transport/http"
	"net/http"
)

type Endpoints struct {
	GetXlsxEndpoint      endpoint.Endpoint
	PostMakeXlsxEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetXlsxEndpoint:      MakeGetXlsxEndpoint(s),
		PostMakeXlsxEndpoint: MakePostMakeXlsxEndpoint(s),
	}
}

func MakeGetXlsxEndpoint(sr Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(map[string]string)
		if !ok {
			panic(http.StatusText(http.StatusBadRequest))
		}
		file := r["file"]
		path, err := sr.GetXlsx(file)

		return path, err
	}
}

func MakePostMakeXlsxEndpoint(sr Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(XlsxRequest)
		if !ok {
			panic(http.StatusText(http.StatusBadRequest))
		}
		fp, err := sr.MakeXlsx(r.Data)

		return xlsxResponse{fp}, err
	}
}
