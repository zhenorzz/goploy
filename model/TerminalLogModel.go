// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const terminalLogTable = "`terminal_log`"

type TerminalLog struct {
	ID          int64  `json:"id"`
	NamespaceID int64  `json:"namespaceID"`
	UserID      int64  `json:"userID"`
	Username    string `json:"username"`
	ServerID    int64  `json:"serverID"`
	ServerName  string `json:"serverName"`
	RemoteAddr  string `json:"remoteAddr"`
	UserAgent   string `json:"userAgent"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	InsertTime  string `json:"insertTime"`
	UpdateTime  string `json:"updateTime"`
}

type TerminalLogs []TerminalLog

func (tl TerminalLog) GetData() (TerminalLog, error) {
	var terminalLog TerminalLog
	err := sq.
		Select(
			"id",
			"namespace_id",
			"user_id",
			"server_id",
			"remote_addr",
			"user_agent",
			"start_time",
			"end_time",
		).
		From(terminalLogTable).
		Where(sq.Eq{"id": tl.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&terminalLog.ID,
			&terminalLog.NamespaceID,
			&terminalLog.UserID,
			&terminalLog.ServerID,
			&terminalLog.RemoteAddr,
			&terminalLog.UserAgent,
			&terminalLog.StartTime,
			&terminalLog.EndTime,
		)
	if err != nil {
		return terminalLog, err
	}
	return terminalLog, nil
}

func (tl TerminalLog) GetList(page, limit uint64) (TerminalLogs, error) {
	builder := sq.
		Select(
			terminalLogTable+".id",
			terminalLogTable+".namespace_id",
			terminalLogTable+".user_id",
			userTable+".name",
			terminalLogTable+".server_id",
			fmt.Sprintf("IFNULL(%s.name, '')", serverTable),
			terminalLogTable+".remote_addr",
			terminalLogTable+".user_agent",
			terminalLogTable+".start_time",
			terminalLogTable+".end_time",
		).
		From(terminalLogTable).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, terminalLogTable)).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.server_id", serverTable, terminalLogTable))

	if tl.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{terminalLogTable + ".namespace_id": tl.NamespaceID})
	}
	if tl.Username != "" {
		builder = builder.Where(sq.Eq{userTable + ".name": tl.Username})
	}
	if tl.ServerName != "" {
		builder = builder.Where(sq.Eq{serverTable + ".name": tl.ServerName})
	}

	rows, err := builder.Limit(limit).Offset((page - 1) * limit).OrderBy("id DESC").
		RunWith(DB).Query()
	if err != nil {
		return nil, err
	}
	terminalLogs := TerminalLogs{}
	for rows.Next() {
		var terminalLog TerminalLog
		if err := rows.Scan(
			&terminalLog.ID,
			&terminalLog.NamespaceID,
			&terminalLog.UserID,
			&terminalLog.Username,
			&terminalLog.ServerID,
			&terminalLog.ServerName,
			&terminalLog.RemoteAddr,
			&terminalLog.UserAgent,
			&terminalLog.StartTime,
			&terminalLog.EndTime,
		); err != nil {
			return nil, err
		}
		terminalLogs = append(terminalLogs, terminalLog)
	}

	return terminalLogs, nil
}

func (tl TerminalLog) GetTotal() (int64, error) {
	var total int64
	builder := sq.Select("COUNT(*) AS count").From(terminalLogTable)

	if tl.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{terminalLogTable + ".namespace_id": tl.NamespaceID})
	}
	if tl.Username != "" {
		builder = builder.
			LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, terminalLogTable)).
			Where(sq.Eq{userTable + ".name": tl.Username})
	}
	if tl.ServerName != "" {
		builder = builder.
			LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.server_id", serverTable, terminalLogTable)).
			Where(sq.Eq{serverTable + ".name": tl.ServerName})
	}
	err := builder.RunWith(DB).QueryRow().Scan(&total)

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (tl TerminalLog) AddRow() (int64, error) {
	result, err := sq.
		Insert(terminalLogTable).
		Columns("namespace_id", "user_id", "server_id", "remote_addr", "user_agent", "start_time").
		Values(tl.NamespaceID, tl.UserID, tl.ServerID, tl.RemoteAddr, tl.UserAgent, tl.StartTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (tl TerminalLog) EditRow() error {
	_, err := sq.
		Update(terminalLogTable).
		SetMap(sq.Eq{
			"end_time": tl.EndTime,
		}).
		Where(sq.Eq{"id": tl.ID}).
		RunWith(DB).
		Exec()
	return err
}
