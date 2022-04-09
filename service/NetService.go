// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package service

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Gnet -
type Gnet struct {
	URL string
}

func (gnet Gnet) Ping() error {
	urlParse, err := url.Parse(gnet.URL)
	if err != nil {
		return err
	}
	switch urlParse.Scheme {
	case "tcp":
		conn, err := net.DialTimeout("tcp", urlParse.Host, 5*time.Second)
		if err != nil {
			return err
		}
		_ = conn.Close()
	case "http", "https":
		if resp, err := http.Get(urlParse.String()); err != nil {
			return err
		} else if resp.StatusCode != 200 {
			return errors.New("Unexpected response status code: " + strconv.Itoa(resp.StatusCode))
		}
	default:
		return errors.New("URL scheme " + urlParse.Scheme + " does not supported")
	}
	return nil
}
