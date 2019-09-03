package core

import (
	"encoding/json"
	"net/http"
)

//Response struct
type Response struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// respone code
const (
	Pass            = 0
	Deny            = 1
	AccountDisabled = 10000
	IllegalRequest  = 10001
	LoginExpired    = 10086
)

//JSON response
func (r *Response) JSON(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
