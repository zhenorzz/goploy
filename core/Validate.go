package core

import (
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)
// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate
// Trans Translator
var Trans ut.Translator

// CreateValidator create a single Validator
func CreateValidator() {

	en := en.New()
	uni := ut.New(en, en)
	Trans, _ = uni.GetTranslator("en")
	Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate, Trans)

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
		return ut.Add("password", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())

		return t
	})
}