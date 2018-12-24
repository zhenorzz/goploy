package core

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) Json(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
