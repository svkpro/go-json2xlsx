package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
)

const (
	RepositoryPath = "repository/%s"
	DateTimeLayout = "2006-01-02 15:04:05"
)

type Service interface {
	GetXlsx(string) (string, error)
	MakeXlsx(string) (string, error)
}

type service struct{}

func NewService() Service {
	return &service {}
}

func (sr *service) GetXlsx(f string)(string, error) {
	path := fmt.Sprintf(RepositoryPath, f)

	return path, nil
}

func (sr *service) MakeXlsx(s string)(string, error) {
	xf := XlsxFile{}
	xf.File = xlsx.NewFile()
	xf.Sheet, xf.Err = xf.File.AddSheet("Sheet1")
	if xf.Err != nil {
		fmt.Printf(xf.Err.Error())
	}
	xf.Row = xf.Sheet.AddRow()
	xf.Cell = xf.Row.AddCell()
	xf.Cell.Value = s

	fileName := time.Now().Format(DateTimeLayout) + ".xlsx"
	xf.Err = xf.File.Save(fmt.Sprintf(RepositoryPath, fileName))
	if xf.Err != nil {
		fmt.Printf(xf.Err.Error())
	}

	return fileName, xf.Err
}
