package server

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/zhenorzz/goploy/internal/model"
	"github.com/zhenorzz/goploy/internal/validator"
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
	return validator.Verify(v)
}
