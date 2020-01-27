package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Count(s string) (n int, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"result", n,
			"errors", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	n, err = mw.next.Count(s)

	return
}
