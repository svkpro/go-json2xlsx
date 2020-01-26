package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte (http.StatusText(http.StatusOK)))
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
