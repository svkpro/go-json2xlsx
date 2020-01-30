package application

import (
	"io"
	"log"
)

type DocumentUploader interface {
	Upload(data io.Reader, key string) (uuid string, err error)
	Delete(key string) error
	Ping() error
	SignedRetrievalURL(key string, originalFileName string, signedTTL int64) (url string, err error)
}

type DocumentStorageService struct {
	fileUploader DocumentUploader
	logger       log.Logger
}

func NewDocumentStorageService(du DocumentUploader, logger log.Logger) *DocumentStorageService {
	dss := &DocumentStorageService{
		fileUploader: du,
		logger: logger,
	}

	return dss
}