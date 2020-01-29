package main

import (
	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) MakeXlsx(s string) (path string, err error) {
	return mw.next.MakeXlsx(s)
}