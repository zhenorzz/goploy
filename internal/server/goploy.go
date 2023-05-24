package server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	"github.com/zhenorzz/goploy/internal/model"
	_validator "github.com/zhenorzz/goploy/internal/validator"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/url"
)

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

var decoder = schema.NewDecoder()

// Decode data from URL query and body
// if the URL query and body have the same field, the value in the body will be taken first
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

	if err := _validator.Validate.Struct(v); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(_validator.Trans))
		}
	}

	return nil
}
