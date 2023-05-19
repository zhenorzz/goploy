// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const sftpLogTable = "`sftp_log`"

type SftpLog struct {
	ID          int64  `json:"id"`
	NamespaceID int64  `json:"namespaceID"`
	UserID      int64  `json:"userID"`
	Username    string `json:"username"`
	ServerID    int64  `json:"serverID"`
	ServerName  string `json:"serverName"`
	RemoteAddr  string `json:"remoteAddr"`
	UserAgent   string `json:"userAgent"`
	Type        string `json:"type"`
	Path        string `json:"path"`
	Reason      string `json:"reason"`
	InsertTime  string `json:"insertTime"`
	UpdateTime  string `json:"updateTime"`
}

const (
	SftpLogTypeDownload = "DOWNLOAD"
	SftpLogTypeUpload   = "UPLOAD"
	SftpLogTypeRead     = "READ"
	SftpLogTypePreview  = "PREVIEW"
	SftpLogTypeRename   = "RENAME"
	SftpLogTypeCopy     = "COPY"
	SftpLogTypeEdit     = "EDIT"
	SftpLogTypeDelete   = "DELETE"
)

type SftpLogs []SftpLog

func (sl SftpLog) GetList(page, limit uint64) (SftpLogs, error) {
	builder := sq.
		Select(
			sftpLogTable+".id",
			sftpLogTable+".namespace_id",
			sftpLogTable+".user_id",
			userTable+".name",
			sftpLogTable+".server_id",
			fmt.Sprintf("IFNULL(%s.name, '')", serverTable),
			sftpLogTable+".remote_addr",
			sftpLogTable+".user_agent",
			sftpLogTable+".type",
			sftpLogTable+".path",
			sftpLogTable+".reason",
			sftpLogTable+".insert_time",
		).
		From(sftpLogTable).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, sftpLogTable)).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.server_id", serverTable, sftpLogTable))

	if sl.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{sftpLogTable + ".namespace_id": sl.NamespaceID})
	}
	if sl.Username != "" {
		builder = builder.Where(sq.Eq{userTable + ".name": sl.Username})
	}
	if sl.ServerName != "" {
		builder = builder.Where(sq.Eq{serverTable + ".name": sl.ServerName})
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
	sftpLogs := SftpLogs{}
	for rows.Next() {
		var sftpLog SftpLog
		if err := rows.Scan(
			&sftpLog.ID,
			&sftpLog.NamespaceID,
			&sftpLog.UserID,
			&sftpLog.Username,
			&sftpLog.ServerID,
			&sftpLog.ServerName,
			&sftpLog.RemoteAddr,
			&sftpLog.UserAgent,
			&sftpLog.Type,
			&sftpLog.Path,
			&sftpLog.Reason,
			&sftpLog.InsertTime,
		); err != nil {
			return nil, err
		}
		sftpLogs = append(sftpLogs, sftpLog)
	}

	return sftpLogs, nil
}

func (sl SftpLog) GetTotal() (int64, error) {
	var total int64
	builder := sq.Select("COUNT(*) AS count").From(sftpLogTable)

	if sl.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{sftpLogTable + ".namespace_id": sl.NamespaceID})
	}
	if sl.Username != "" {
		builder = builder.
			LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, sftpLogTable)).
			Where(sq.Eq{userTable + ".name": sl.Username})
	}
	if sl.ServerName != "" {
		builder = builder.
			LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.server_id", serverTable, sftpLogTable)).
			Where(sq.Eq{serverTable + ".name": sl.ServerName})
	}
	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (sl SftpLog) AddRow() error {
	_, err := sq.
		Insert(sftpLogTable).
		Columns("namespace_id", "user_id", "server_id", "remote_addr", "user_agent", "type", "path", "reason").
		Values(sl.NamespaceID, sl.UserID, sl.ServerID, sl.RemoteAddr, sl.UserAgent, sl.Type, sl.Path, sl.Reason).
		RunWith(DB).
		Exec()
	return err
}
