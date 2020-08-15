package core

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
	"strings"
	"unicode"
)

// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate

// Trans Translator
var Trans ut.Translator

// CreateValidator create a single Validator
func CreateValidator() {

	english := en.New()
	uni := ut.New(english, english)
	Trans, _ = uni.GetTranslator("english")
	Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate, Trans)
	registerTagName()
	registerPassword()
	registerRole()
}

func registerTagName() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func registerPassword() {
	Validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		// 8到16个字符，至少包含字母、数字、特殊符号中的两种
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

	Validate.RegisterTranslation("password", Trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} policy is min:8, max:16 and at least one alpha and at least one special char!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())

		return t
	})
}

func registerRole() {
	Validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		// 8到16个字符，至少包含字母、数字、特殊符号中的两种
		role := fl.Field().String()
		for _, v := range Roles {
			if role == v {
				return true
			}
		}
		return false
	})

	Validate.RegisterTranslation("role", Trans, func(ut ut.Translator) error {
		return ut.Add("role", "{0} is invalid", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("role", fe.Field())

		return t
	})
}
