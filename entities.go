package main

import "github.com/tealeg/xlsx"

type XlsxFile struct {
	File *xlsx.File
	Sheet *xlsx.Sheet
	Row *xlsx.Row
	Cell *xlsx.Cell
	Err error
}
