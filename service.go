package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"json2xls/application/aws"
	"os"
	"time"
)

const (
	RepositoryPath = "data/%s"
	FileNameLayout = "%s.xlsx"
	DateTimeLayout = "2006-01-02_15:04:05"
	signedTTL      = 60
)

type Service interface {
	GetXlsx(string) (string, error)
	MakeXlsx(data XlsxRequestData) (string, error)
	DeleteXlsx(string) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (sr *service) GetXlsx(f string) (string, error) {
	s3 := aws.New()
	err := s3.Ping()
	if s3.Ping() != nil {
		return "", err
	}

	s3l, err := s3.SignedRetrievalURL(f, f, signedTTL)

	return s3l, nil
}

func (sr *service) MakeXlsx(data XlsxRequestData) (string, error) {
	s3 := aws.New()
	err := s3.Ping()
	if s3.Ping() != nil {
		return "", err
	}
	xf := XlsxFile{}
	xf.File = xlsx.NewFile()
	xf.Sheet, xf.Err = xf.File.AddSheet(data.Sheet)
	if xf.Err != nil {
		return "", err
	}

	xf.Row = xf.Sheet.AddRow()
	xf.Row.WriteSlice(&data.Headers, -1)

	for _, row := range data.Rows {
		xf.Row = xf.Sheet.AddRow()
		xf.Row.Cells = row.Cells
	}

	fileName := fmt.Sprintf(FileNameLayout, time.Now().Format(DateTimeLayout))
	path := fmt.Sprintf(RepositoryPath, fileName)
	xf.Err = xf.File.Save(path)
	defer func() {
		err = os.Remove(path)
	}()
	if xf.Err != nil {
		return "", err
	}
	fr, err := os.Open(fmt.Sprintf(RepositoryPath, fileName))
	if xf.Err != nil {
		return "", err
	}

	s3l, err := s3.Upload(fr, fileName, signedTTL)
	if xf.Err != nil {
		return "", err
	}

	return s3l, xf.Err
}

func (sr *service) DeleteXlsx(f string) error {
	s3 := aws.New()
	err := s3.Ping()
	if s3.Ping() != nil {
		return err
	}

	return s3.Delete(f)
}
