package validator

import (
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	"unicode"
)

func registerPassword() {
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
