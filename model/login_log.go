// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const loginLogTable = "`login_log`"

type LoginLog struct {
	ID         int64  `json:"id"`
	Account    string `json:"account"`
	RemoteAddr string `json:"remoteAddr"`
	UserAgent  string `json:"userAgent"`
	Referer    string `json:"referer"`
	Reason     string `json:"reason"`
	LoginTime  string `json:"loginTime"`
}

type LoginLogs []LoginLog

func (ll LoginLog) GetList(page, limit uint64) (LoginLogs, error) {
	builder := sq.
		Select(
			"id",
			"account",
			"remote_addr",
			"user_agent",
			"referer",
			"reason",
			"login_time",
		).
		From(loginLogTable)
	if len(ll.Account) > 0 {
		builder = builder.Where(sq.Eq{"account": ll.Account})
	}

	rows, err := builder.
		Limit(limit).
		Offset((page - 1) * limit).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	loginLogs := LoginLogs{}
	for rows.Next() {
		var loginLog LoginLog
		if err := rows.Scan(
			&loginLog.ID,
			&loginLog.Account,
			&loginLog.RemoteAddr,
			&loginLog.UserAgent,
			&loginLog.Referer,
			&loginLog.Reason,
			&loginLog.LoginTime,
		); err != nil {
			return nil, err
		}
		loginLogs = append(loginLogs, loginLog)
	}

	return loginLogs, nil
}

func (ll LoginLog) GetTotal() (int64, error) {
	var total int64
	builder := sq.Select("COUNT(*) AS count").From(loginLogTable)

	if len(ll.Account) > 0 {
		builder = builder.Where(sq.Eq{"account": ll.Account})
	}
	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (ll LoginLog) AddRow() error {
	_, err := sq.
		Insert(loginLogTable).
		Columns("account", "remote_addr", "user_agent", "referer", "reason", "login_time").
		Values(ll.Account, ll.RemoteAddr, ll.UserAgent, ll.Referer, ll.Reason, ll.LoginTime).
		RunWith(DB).
		Exec()
	return err
}
