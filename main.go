package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	s := NewService()
	s = loggingMiddleware{logger, s}

	handler := MakeHTTPHandler(s, logger)

	logger.Log("err", http.ListenAndServe(":9009", handler))
}
