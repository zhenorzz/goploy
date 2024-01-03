// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package response

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

type File struct {
	Filename    string
	Disposition string // attachment | inline
}

func (f File) Write(w http.ResponseWriter, _ *http.Request) error {
	file, err := os.Open(f.Filename)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if f.Disposition == "attachment" {
		w.Header().Set("Content-Disposition", "attachment;filename="+fileStat.Name())
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		w.Header().Set("Content-Disposition", "inline;filename="+fileStat.Name())
	}

	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	return nil
}
