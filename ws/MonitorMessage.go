// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package ws

type MonitorMessage struct {
	MonitorID    int64  `json:"monitorId"`
	State        uint8  `json:"state"`
	ErrorContent string `json:"errorContent"`
}

const (
	MonitorTurnOff = 0
)

func (projectMessage MonitorMessage) canSendTo(*Client) error {
	return nil
}
