// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package monitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy/internal/model"
	"golang.org/x/crypto/ssh"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type ScriptError struct {
	Message  string
	ServerID int64
}

func (se ScriptError) Error() string {
	return se.Message
}

func (se ScriptError) Server() int64 {
	return se.ServerID
}

type Script struct {
	ServerID int64
	Content  string
}

func (s Script) IsValid() bool {
	return s.Content != ""
}

type Monitor struct {
	Type          int
	Items         []string
	Timeout       time.Duration
	Process       string
	Script        string
	FailScript    Script
	SuccessScript Script
}

type Option func(*Monitor)

func WithSuccessScript(serverID int64, content string) Option {
	return func(c *Monitor) {
		c.SuccessScript = Script{
			ServerID: serverID,
			Content:  content,
		}
	}
}

func WithFailScript(serverID int64, content string) Option {
	return func(c *Monitor) {
		c.FailScript = Script{
			ServerID: serverID,
			Content:  content,
		}
	}
}

func NewMonitorFromTarget(t int, target string, options ...Option) Monitor {
	m := Monitor{}
	if err := json.Unmarshal([]byte(target), &m); err != nil {
		return m
	}
	m.Type = t

	for _, option := range options {
		option(&m)
	}

	return m
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
		server, client, session, err := NewServerSession(serverID, m.Timeout*time.Second)
		if err != nil {
			return err
		}
		var stdout, stderr bytes.Buffer
		session.Stdout = &stdout
		session.Stderr = &stderr
		err = session.Run(server.ReplaceVars(m.Script))
		_ = client.Close()
		_ = session.Close()
		if err != nil {
			return ScriptError{Message: err.Error() + ", stdout: " + stdout.String() + ", stderr: " + stderr.String(), ServerID: serverID}
		}
	}
	return nil
}
func (m Monitor) RunFailScript(serverID int64) error {
	if m.FailScript.IsValid() {
		sId := m.FailScript.ServerID
		if sId == 0 && serverID != -1 {
			sId = serverID
		}

		if sId != -1 {
			server, client, session, err := NewServerSession(sId, m.Timeout*time.Second)
			if err != nil {
				return err
			}

			err = session.Run(server.ReplaceVars(m.FailScript.Content))
			_ = client.Close()
			_ = session.Close()
			return err
		} else {
			cmd := exec.Command(m.FailScript.Content)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("err: %s, detail: %s", err, string(output))
			}
		}

	}
	return nil
}
func (m Monitor) RunSuccessScript(serverID int64) error {
	if m.SuccessScript.IsValid() {
		sId := m.SuccessScript.ServerID

		if sId == 0 && serverID != -1 {
			sId = serverID
		}

		if sId != -1 {
			server, client, session, err := NewServerSession(sId, m.Timeout*time.Second)
			if err != nil {
				return err
			}

			err = session.Run(server.ReplaceVars(m.SuccessScript.Content))
			_ = client.Close()
			_ = session.Close()
			return err
		} else {
			cmd := exec.Command(m.SuccessScript.Content)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("err: %s, detail: %s", err, string(output))
			}
		}

	}
	return nil
}

func NewServerSession(serverID int64, timeout time.Duration) (model.Server, *ssh.Client, *ssh.Session, error) {
	server, err := (model.Server{ID: serverID}).GetData()
	if err != nil {
		return server, nil, nil, err
	} else if server.State == model.Disable {
		return server, nil, nil, errors.New("Server Disable [" + server.Name + "]")
	}
	client, err := server.ToSSHConfig().SetTimeout(timeout).Dial()
	if err != nil {
		return server, nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
	}
	return server, client, session, err
}
