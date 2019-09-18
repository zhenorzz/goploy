package controller

import (
	"encoding/json"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

// Controller struct
type Controller struct {
}

// Init DB open
func Init() {
	validate = validator.New()
}

func verify(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s %s", err.Field(), err.Tag())
		}
	}
	return nil
}
