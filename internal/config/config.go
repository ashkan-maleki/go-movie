package config

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func IsValid(config any) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return errors.New(fmt.Sprintf("Missing required attributes %v\n", err))
	}
	return nil
}
