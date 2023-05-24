// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	_validator "github.com/zhenorzz/goploy/internal/validator"
	"gopkg.in/go-playground/validator.v9"
)

type API struct{}

var decoder = schema.NewDecoder()

func decodeJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if err := _validator.Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(_validator.Trans))
		}
	}
	return nil
}

func decodeQuery(data map[string][]string, v interface{}) error {
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(v, data)
	if err != nil {
		return err
	}
	if err := _validator.Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(_validator.Trans))
		}
	}
	return nil
}
