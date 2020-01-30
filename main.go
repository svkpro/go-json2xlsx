package main

import (
	"json2xls/config"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	/* TODO: should be implemented with aws uploading.
	ssl := false

	// Initialize minio client object.
	minioClient, err := minio.New("0.0.0.0:19002", "MinioAccessKey", "MinioSecretKey", ssl)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = minioClient.MakeBucket("mybucket", "us-east-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully created mybucket.")*/

	settings := config.New()

	logger := log.NewLogfmtLogger(os.Stderr)
	service := NewService()
	service = loggingMiddleware{logger, service}

	handler := MakeHTTPHandler(service, logger)

	logger.Log("err", http.ListenAndServe(settings.HttpPort, handler))
}
