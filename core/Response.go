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

//JSON response
func (r *Response) JSON(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
