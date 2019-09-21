package controller

import (
	"encoding/json"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"goploy/core"
)

// Controller struct
type Controller struct {
}

func verify(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if err := core.Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(core.Trans))
		}
	}
	return nil
}
