// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package response

import (
	"net/http"
)

type Redirect struct {
	URL  string
	Code int
}

func (rdr Redirect) Write(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, rdr.URL, rdr.Code)
	return nil
}
