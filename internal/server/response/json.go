// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package response

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	// Pass = 0
	// Deny = 1
	// Error = 2, see Message for details
	// AccountDisabled  = 10000
	// IllegalRequest = 10001
	// NamespaceInvalid = 10002
	// IllegalParam = 10003
	// LoginExpired = 10086
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	Pass             = 0
	Deny             = 1
	Error            = 2
	AccountDisabled  = 10000
	IllegalRequest   = 10001
	NamespaceInvalid = 10002
	IllegalParam     = 10003
	PasswordExpired  = 10004
	LoginExpired     = 10086
)

func (j JSON) Write(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(j)
}
