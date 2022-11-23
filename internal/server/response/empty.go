// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package response

import (
	"net/http"
)

type Empty struct{}

func (Empty) Write(http.ResponseWriter, *http.Request) error { return nil }
