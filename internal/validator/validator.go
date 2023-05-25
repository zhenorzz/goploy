package validator

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
)

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
		name := fld.Tag.Get("json")

		if name == "" {
			name = fld.Tag.Get("schema")
		}

		if name == "" {
			name = fld.Name
		}

		if name == "-" {
			return ""
		}

		return name
	})

	registerPassword()
	registerFilepath()
}

func Verify(v interface{}) error {
	if err := Validate.Struct(v); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(Trans))
		}
	}
	return nil
}
