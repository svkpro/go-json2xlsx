package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"json2xls/application/aws"
	"json2xls/config"
	"os"
	"time"
)

const (
	RepositoryPath = "data/%s"
	FileNameLayout = "%s.xlsx"
	DateTimeLayout = "2006-01-02_15:04:05"
)

type Service interface {
	GetXlsx(string) (string, error)
	MakeXlsx(string) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (sr *service) GetXlsx(f string) (string, error) {
	path := fmt.Sprintf(RepositoryPath, f)

	return path, nil
}

func (sr *service) MakeXlsx(s string) (string, error) {
	c := config.New()
	s3 := aws.S3FileUploader{
		URL:        c.AwsURL,
		BucketName: c.AwsBucketName,
		Region:     c.AwsRegion,
		AccessKey:  c.AwsAccessKey,
		SecretKey:  c.AwsSecretKey,
		DisableSSL: c.AwsDisableSSL,
	}
	err := s3.Ping()
	if s3.Ping() != nil {
		return "", err
	}
	xf := XlsxFile{}
	xf.File = xlsx.NewFile()
	xf.Sheet, xf.Err = xf.File.AddSheet("Sheet1")
	if xf.Err != nil {
		return "", err
	}
	xf.Row = xf.Sheet.AddRow()
	xf.Cell = xf.Row.AddCell()
	xf.Cell.Value = s
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

	s3l, err := s3.Upload(fr, fileName)
	if xf.Err != nil {
		return "", err
	}

	return s3l, xf.Err
}
