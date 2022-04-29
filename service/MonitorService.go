// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/model"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Monitor struct {
	Type    int
	Items   []string
	Timeout time.Duration
	Process string
	Script  string
}

func NewMonitorFromTarget(t int, target string) (Monitor, error) {
	var m Monitor
	if err := json.Unmarshal([]byte(target), &m); err != nil {
		return m, err
	}
	m.Type = t
	return m, nil
}

func (m Monitor) Check() error {
	var err error
	switch m.Type {
	case 1:
		err = m.CheckSite()
	case 2:
		err = m.CheckPort()
	case 3:
		err = m.CheckHostAlive()
	case 4:
		m.Script = fmt.Sprintf("ps -ef|grep -v grep|grep %s", m.Process)
		fallthrough
	case 5:
		err = m.CheckScript()
	default:
		err = errors.New("type error")
	}
	return err
}

func (m Monitor) CheckSite() error {
	for _, url := range m.Items {
		client := http.Client{
			Timeout: m.Timeout * time.Second,
		}
		if resp, err := client.Get(url); err != nil {
			return err
		} else if 200 < resp.StatusCode || resp.StatusCode >= 400 {
			return errors.New("Unexpected response status code: " + strconv.Itoa(resp.StatusCode))
		}
	}
	return nil
}

func (m Monitor) CheckPort() error {
	for _, address := range m.Items {
		conn, err := net.DialTimeout("tcp", address, m.Timeout*time.Second)
		if err != nil {
			return err
		}
		_ = conn.Close()
	}
	return nil
}

func (m Monitor) CheckHostAlive() error {
	var stdout, stderr bytes.Buffer
	for _, addr := range m.Items {
		var arg []string
		if runtime.GOOS == "windows" {
			arg = append(arg, "-n", "1", "-w", strconv.Itoa(int(m.Timeout*1000)), addr)
		} else {
			arg = append(arg, "-c", "1", "-W", strconv.Itoa(int(m.Timeout)), addr)
		}

		stdout.Reset()
		stderr.Reset()
		cmd := exec.Command("ping", arg...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.New(err.Error() + ", detail: " + stderr.String())
		}
	}

	return nil
}

func (m Monitor) CheckScript() error {
	for _, serverIDStr := range m.Items {
		serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
		if err != nil {
			return err
		}
		server, err := (model.Server{ID: serverID}).GetData()
		if err != nil {
			return err
		} else if server.State == model.Disable {
			continue
		}
		client, err := server.ToSSHConfig().SetTimeout(m.Timeout * time.Second).Dial()
		if err != nil {
			return err
		}
		session, err := client.NewSession()
		if err != nil {
			return err
		}
		var stdout, stderr bytes.Buffer
		session.Stdout = &stdout
		session.Stderr = &stderr
		if err := session.Run(m.Script); err != nil {
			return errors.New(err.Error() + ", stdout: " + stdout.String() + ", stderr: " + stderr.String())
		}
	}
	return nil
}
