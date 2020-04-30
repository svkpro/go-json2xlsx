package main

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func ValidateStruct(st interface{}) error {
	err := validate.Struct(st)

	if validate.Struct(st) != nil {
		return err
	}

	return nil
}
