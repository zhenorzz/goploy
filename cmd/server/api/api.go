// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/schema"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
	"strings"
	"unicode"
)

type API struct{}

// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate

// Trans Translator
var Trans ut.Translator

func init() {
	english := en.New()
	uni := ut.New(english, english)
	Trans, _ = uni.GetTranslator("english")
	Validate = validator.New()
	_ = enTranslations.RegisterDefaultTranslations(Validate, Trans)
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	_ = Validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 || len(password) > 16 {
			return false
		}
		var (
			hasLetter  = false
			hasNumber  = false
			hasSpecial = false
		)

		for _, char := range password {
			switch {
			case unicode.IsLetter(char):
				hasLetter = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		if hasLetter && hasNumber {
			return true
		} else if hasLetter && hasSpecial {
			return true
		} else if hasNumber && hasSpecial {
			return true
		} else {
			return false
		}
	})

	_ = Validate.RegisterTranslation("password", Trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} policy is min:8, max:16 and at least one alpha and at least one special char!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())

		return t
	})
}

var decoder = schema.NewDecoder()

func decodeJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if err := Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(Trans))
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
	if err := Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(Trans))
		}
	}
	return nil
}
