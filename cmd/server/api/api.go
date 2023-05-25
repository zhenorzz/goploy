// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/zhenorzz/goploy/internal/validator"
)

type API struct{}

var decoder = schema.NewDecoder()

func decodeJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return validator.Verify(v)
}

func decodeQuery(data map[string][]string, v interface{}) error {
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(v, data)
	if err != nil {
		return err
	}
	return validator.Verify(v)
}
