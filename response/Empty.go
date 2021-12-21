package response

import (
	"net/http"
)

type Empty struct{}

//JSON response
func (Empty) Write(http.ResponseWriter) error { return nil }
