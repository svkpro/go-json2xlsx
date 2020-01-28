package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var s Service
	s = service{}
	s = loggingMiddleware{logger, s}

	handler := MakeHTTPHandler(s, logger)

	logger.Log("err", http.ListenAndServe(":9090", handler))
}
