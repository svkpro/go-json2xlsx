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
)

type Service interface {
	GetXlsx(string) (string, error)
	MakeXlsx(data XlsxPayloadData) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (sr *service) GetXlsx(f string) (string, error) {
	path := fmt.Sprintf(RepositoryPath, f)

	return path, nil
}

func (sr *service) MakeXlsx(data XlsxPayloadData) (string, error) {
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
	for _, value := range data.Headers {
		xf.Cell = xf.Row.AddCell()
		xf.Cell.Value = value
	}

	for _, row := range data.Rows {
		xf.Row = xf.Sheet.AddRow()
		for _, value := range row {
			xf.Cell = xf.Row.AddCell()
			xf.Cell.Value = string(value)
		}
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

	s3l, err := s3.Upload(fr, fileName)
	if xf.Err != nil {
		return "", err
	}

	return s3l, xf.Err
}
