// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const operationLogTable = "`operation_log`"

type OperationLog struct {
	ID           int64  `json:"id"`
	NamespaceID  int64  `json:"namespaceID"`
	UserID       int64  `json:"userID"`
	Username     string `json:"username"`
	Router       string `json:"router"`
	API          string `json:"api"`
	RequestTime  string `json:"requestTime"`
	RequestData  string `json:"requestData"`
	ResponseTime string `json:"responseTime"`
	ResponseData string `json:"responseData"`
}

type OperationLogs []OperationLog

func (ol OperationLog) GetList(page, limit uint64) (OperationLogs, error) {
	builder := sq.
		Select(
			operationLogTable+".id",
			operationLogTable+".namespace_id",
			operationLogTable+".user_id",
			userTable+".name",
			operationLogTable+".router",
			operationLogTable+".api",
			operationLogTable+".request_time",
			operationLogTable+".request_data",
			operationLogTable+".response_time",
			operationLogTable+".response_data",
		).
		From(operationLogTable).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, operationLogTable))

	if ol.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{operationLogTable + ".namespace_id": ol.NamespaceID})
	}

	if ol.Router != "" {
		builder = builder.Where(sq.Eq{operationLogTable + ".router": ol.Router})
	}

	if ol.API != "" {
		builder = builder.Where(sq.Eq{operationLogTable + ".api": ol.API})
	}

	if ol.Username != "" {
		builder = builder.Where(sq.Eq{userTable + ".name": ol.Username})
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
	opLogs := OperationLogs{}
	for rows.Next() {
		var opLog OperationLog
		if err := rows.Scan(
			&opLog.ID,
			&opLog.NamespaceID,
			&opLog.UserID,
			&opLog.Username,
			&opLog.Router,
			&opLog.API,
			&opLog.RequestTime,
			&opLog.RequestData,
			&opLog.ResponseTime,
			&opLog.ResponseData,
		); err != nil {
			return nil, err
		}
		opLogs = append(opLogs, opLog)
	}

	return opLogs, nil
}

func (ol OperationLog) GetTotal() (int64, error) {
	var total int64
	builder := sq.Select("COUNT(*) AS count").From(operationLogTable)

	if ol.NamespaceID > 0 {
		builder = builder.Where(sq.Eq{operationLogTable + ".namespace_id": ol.NamespaceID})
	}

	if ol.Router != "" {
		builder = builder.Where(sq.Eq{operationLogTable + ".router": ol.Router})
	}

	if ol.API != "" {
		builder = builder.Where(sq.Eq{operationLogTable + ".api": ol.API})
	}

	if ol.Username != "" {
		builder = builder.
			LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.user_id", userTable, operationLogTable)).
			Where(sq.Eq{userTable + ".name": ol.Username})
	}

	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (ol OperationLog) AddRow() error {
	_, err := sq.
		Insert(operationLogTable).
		Columns("user_id", "namespace_id", "router", "api", "request_time", "request_data", "response_time", "response_data").
		Values(ol.UserID, ol.NamespaceID, ol.Router, ol.API, ol.RequestTime, ol.RequestData, ol.ResponseTime, ol.ResponseData).
		RunWith(DB).
		Exec()
	return err
}
