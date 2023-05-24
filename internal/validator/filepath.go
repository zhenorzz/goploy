package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/zhenorzz/goploy/internal/pkg"
	"gopkg.in/go-playground/validator.v9"
)

func registerFilepath() {
	_ = Validate.RegisterValidation("filepath", func(fl validator.FieldLevel) bool {
		return pkg.IsFilePath(fl.Field().String())
	})

	_ = Validate.RegisterTranslation("filepath", Trans, func(ut ut.Translator) error {
		return ut.Add("filepath", "{0} policy is start with slash and can not end with slash", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filepath", fe.Field())

		return t
	})
}
