package main

import (
	"errors"
)

type StringService interface {
	Count(string) (int, error)
}

type stringService struct{}

func (stringService) Count(s string) (int, error) {
	if s == "" {
		return 0, errors.New("The string can't be empty!")
	}
	return len(s), nil
}
