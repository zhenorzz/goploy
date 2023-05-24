// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/schema"
	"github.com/vearutop/statigz"
	"github.com/zhenorzz/goploy/config"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/pkg"
	"github.com/zhenorzz/goploy/web"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"unicode"
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

	_ = Validate.RegisterValidation("filepath", func(fl validator.FieldLevel) bool {
		filepath := fl.Field().String()
		if pkg.IsFilePath(filepath) {
			return true
		}
		return false
	})

	_ = Validate.RegisterTranslation("filepath", Trans, func(ut ut.Translator) error {
		return ut.Add("filepath", "{0} policy is start with slash and can not end with slash", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filepath", fe.Field())

		return t
	})
}

var decoder = schema.NewDecoder()

type Goploy struct {
	UserInfo  model.User
	Namespace struct {
		ID            int64
		PermissionIDs map[int64]struct{}
	}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	URLQuery       url.Values
	Body           []byte
}

type Response interface {
	Write(http.ResponseWriter, *http.Request) error
}

type Server struct {
	http.Server
	*Router
}

func (srv *Server) Spin() {
	if config.Toml.Env == "production" {
		subFS, err := fs.Sub(web.Dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/assets/", statigz.FileServer(subFS.(fs.ReadDirFS)))
		http.Handle("/favicon.ico", statigz.FileServer(subFS.(fs.ReadDirFS)))
	}
	http.Handle("/", srv.Router)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func (g *Goploy) Decode(v interface{}) error {
	if len(g.URLQuery) > 0 {
		decoder.IgnoreUnknownKeys(true)
		err := decoder.Decode(v, g.URLQuery)
		if err != nil {
			return err
		}
	}

	if len(g.Body) > 0 {
		err := json.Unmarshal(g.Body, v)
		if err != nil {
			return err
		}
	}

	if err := Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(Trans))
		}
	}

	return nil
}
