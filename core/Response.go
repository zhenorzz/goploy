package core

import (
	"encoding/json"
	"net/http"
)

//Response struct
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// response code
const (
	Pass             = 0
	Deny             = 1
	Error            = 2
	AccountDisabled  = 10000
	IllegalRequest   = 10001
	NamespaceInvalid = 10002
	LoginExpired     = 10086
)

//JSON response
func (r *Response) JSON(w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(r); err != nil {
		Log(ERROR, err.Error())
	}
}
