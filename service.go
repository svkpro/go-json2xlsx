package main

import (
	"errors"
	"fmt"
)

const EmptyParameterError  =  "The parameter %s can't be empty!"

type Service interface {
	Count(string) (int, error)
}

type service struct{}

func (service) Count(s string) (int, error) {
	if s == "" {
		return 0, errors.New(fmt.Sprintf(EmptyParameterError, "s"))
	}
	return len(s), nil
}
