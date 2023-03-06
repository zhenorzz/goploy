// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	sq "github.com/Masterminds/squirrel"
)

const monitorTable = "`monitor`"

type Monitor struct {
	ID           int64  `json:"id"`
	NamespaceID  int64  `json:"namespaceId"`
	Name         string `json:"name"`
	Type         int    `json:"type"`
	Target       string `json:"target"`
	Second       int    `json:"second"`
	Times        uint16 `json:"times"`
	SilentCycle  int    `json:"silentCycle"`
	NotifyType   uint8  `json:"notifyType"`
	NotifyTarget string `json:"notifyTarget"`
	Description  string `json:"description"`
	ErrorContent string `json:"errorContent"`
	State        uint8  `json:"state"`
	InsertTime   string `json:"insertTime"`
	UpdateTime   string `json:"updateTime"`
}

type Monitors []Monitor

func (m Monitor) GetList() (Monitors, error) {
	rows, err := sq.
		Select("id, name, type, target, second, times, silent_cycle, notify_type, notify_target, description, error_content, state, insert_time, update_time").
		From(monitorTable).
		Where(sq.Eq{
			"namespace_id": m.NamespaceID,
		}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	monitors := Monitors{}
	for rows.Next() {
		var monitor Monitor

		if err := rows.Scan(
			&monitor.ID,
			&monitor.Name,
			&monitor.Type,
			&monitor.Target,
			&monitor.Second,
			&monitor.Times,
			&monitor.SilentCycle,
			&monitor.NotifyType,
			&monitor.NotifyTarget,
			&monitor.Description,
			&monitor.ErrorContent,
			&monitor.State,
			&monitor.InsertTime,
			&monitor.UpdateTime); err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (m Monitor) GetData() (Monitor, error) {
	var monitor Monitor
	err := sq.
		Select("id, name, type, target, second, times, silent_cycle, notify_type, notify_target, state").
		From(monitorTable).
		Where(sq.Eq{"id": m.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&monitor.ID, &monitor.Name, &monitor.Type, &monitor.Target, &monitor.Second, &monitor.Times, &monitor.SilentCycle, &monitor.NotifyType, &monitor.NotifyTarget, &monitor.State)
	if err != nil {
		return monitor, err
	}
	return monitor, nil
}

func (m Monitor) GetAllByState() (Monitors, error) {
	rows, err := sq.
		Select("id, name, type, target, second, times, silent_cycle, notify_type, notify_target, description").
		From(monitorTable).
		Where(sq.Eq{
			"state": m.State,
		}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	monitors := Monitors{}
	for rows.Next() {
		var monitor Monitor

		if err := rows.Scan(
			&monitor.ID,
			&monitor.Name,
			&monitor.Type,
			&monitor.Target,
			&monitor.Second,
			&monitor.Times,
			&monitor.SilentCycle,
			&monitor.NotifyType,
			&monitor.NotifyTarget,
			&monitor.Description); err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (m Monitor) AddRow() (int64, error) {
	result, err := sq.
		Insert(monitorTable).
		Columns("namespace_id", "name", "type", "target", "second", "times", "silent_cycle", "notify_type", "notify_target", "description", "error_content").
		Values(m.NamespaceID, m.Name, m.Type, m.Target, m.Second, m.Times, m.SilentCycle, m.NotifyType, m.NotifyTarget, m.Description, "").
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (m Monitor) EditRow() error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"name":          m.Name,
			"type":          m.Type,
			"target":        m.Target,
			"second":        m.Second,
			"times":         m.Times,
			"silent_cycle":  m.SilentCycle,
			"notify_type":   m.NotifyType,
			"notify_target": m.NotifyTarget,
			"description":   m.Description,
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (m Monitor) ToggleState() error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"state": m.State,
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (m Monitor) DeleteRow() error {
	_, err := sq.
		Delete(monitorTable).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (m Monitor) TurnOff(errorContent string) error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"state":         Disable,
			"error_content": errorContent,
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (m Monitor) Notify(errMsg string) (string, error) {
	var err error
	var resp *http.Response
	if m.NotifyType == NotifyWeiXin {
		type markdown struct {
			Content string `json:"content"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		content := "Monitor: <font color=\"warning\">" + m.Name + "</font>\n "
		content += "> <font color=\"warning\">can not access</font> \n "
		content += "> <font color=\"comment\">" + errMsg + "</font> \n "

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Content: content,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(m.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if m.NotifyType == NotifyDingTalk {
		type markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}
		type message struct {
			Msgtype  string   `json:"msgtype"`
			Markdown markdown `json:"markdown"`
		}
		text := "#### Monitor: " + m.Name + " can not access \n >" + errMsg

		msg := message{
			Msgtype: "markdown",
			Markdown: markdown{
				Title: m.Name,
				Text:  text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(m.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if m.NotifyType == NotifyFeiShu {
		type content struct {
			Text string `json:"text"`
		}
		type message struct {
			MsgType string  `json:"msg_type"`
			Content content `json:"content"`
		}

		text := m.Name + " can not access\n "
		text += "detail:  " + errMsg

		msg := message{
			MsgType: "text",
			Content: content{
				Text: text,
			},
		}
		b, _ := json.Marshal(msg)
		resp, err = http.Post(m.NotifyTarget, "application/json", bytes.NewBuffer(b))
	} else if m.NotifyType == NotifyCustom {
		type message struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				MonitorName string `json:"monitorName"`
				Type        int    `json:"type"`
				Target      string `json:"target"`
				Second      int    `json:"second"`
				Times       uint16 `json:"times"`
				Error       string `json:"error"`
			} `json:"data"`
		}
		code := 0
		msg := message{
			Code:    code,
			Message: m.Name + " can not access",
		}
		msg.Data.MonitorName = m.Name
		msg.Data.Type = m.Type
		msg.Data.Target = m.Target
		msg.Data.Second = m.Second
		msg.Data.Times = m.Times
		msg.Data.Error = errMsg
		b, _ := json.Marshal(msg)
		resp, err = http.Post(m.NotifyTarget, "application/json", bytes.NewBuffer(b))
	}
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		return string(responseData), err
	}
}
