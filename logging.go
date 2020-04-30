package main

import (
	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GetXlsx(s string) (path string, err error) {
	return mw.next.GetXlsx(s)
}

func (mw loggingMiddleware) MakeXlsx(data XlsxRequestData) (path string, err error) {
	return mw.next.MakeXlsx(data)
}

func (mw loggingMiddleware) DeleteXlsx(s string) (err error) {
	return mw.next.DeleteXlsx(s)
}
