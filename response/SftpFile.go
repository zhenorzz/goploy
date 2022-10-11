// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package response

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"strconv"
)

type SftpFile struct {
	Client      *ssh.Client
	Filename    string
	Disposition string // attachment | inline
}

func (sf SftpFile) Write(w http.ResponseWriter, _ *http.Request) error {
	defer sf.Client.Close()

	sftpClient, err := sftp.NewClient(sf.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(sf.Filename)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileStat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if sf.Disposition == "attachment" {
		w.Header().Set("Content-Disposition", "attachment;filename="+fileStat.Name())
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		w.Header().Set("Content-Disposition", "inline;filename="+fileStat.Name())
	}

	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	_, err = io.Copy(w, srcFile)

	return err
}
