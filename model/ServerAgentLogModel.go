// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const serverAgentLogTable = "`server_agent_log`"

// ServerAgentLog -
type ServerAgentLog struct {
	ServerID   int64  `json:"serverId,omitempty"`
	Type       int    `json:"type,omitempty"`
	Item       string `json:"item"`
	Value      string `json:"value"`
	ReportTime string `json:"reportTime"`
}

// ServerAgentLogs -
type ServerAgentLogs []ServerAgentLog

// GetListBetweenTime -
func (sal ServerAgentLog) GetListBetweenTime(low, high string) (ServerAgentLogs, error) {
	rows, err := sq.
		Select(
			"item",
			"value",
			"report_time",
		).
		From(serverAgentLogTable).
		Where(sq.Eq{"server_id": sal.ServerID}).
		Where(sq.Eq{"type": sal.Type}).
		Where(sq.Expr("report_time BETWEEN ? AND ?", low, high)).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	serverAgentLogs := ServerAgentLogs{}
	for rows.Next() {
		var serverAgentLog ServerAgentLog
		if err := rows.Scan(
			&serverAgentLog.Item,
			&serverAgentLog.Value,
			&serverAgentLog.ReportTime,
		); err != nil {
			return nil, err
		}
		serverAgentLogs = append(serverAgentLogs, serverAgentLog)
	}
	return serverAgentLogs, nil
}

// GetCycleValue -
func (sal ServerAgentLog) GetCycleValue(groupCycle int, formula string) (value string, err error) {
	var builder sq.SelectBuilder
	switch formula {
	case "min":
		builder = sq.
			Select("min(value) as value")
	case "max":
		builder = sq.
			Select("max(value) as value")
	default:
		builder = sq.
			Select("avg(value) as value")
	}

	err = builder.
		From(serverAgentLogTable).
		Where(sq.Eq{"server_id": sal.ServerID}).
		Where(sq.Eq{"item": sal.Item}).
		Where("now() - interval ? minute", groupCycle).
		RunWith(DB).
		QueryRow().
		Scan(&value)
	return
}

// AddRow return LastInsertId
func (sal ServerAgentLog) AddRow() error {
	_, err := sq.
		Insert(serverAgentLogTable).
		Columns("server_id", "type", "item", "value", "report_time").
		Values(sal.ServerID, sal.Type, sal.Item, sal.Value, sal.ReportTime).
		RunWith(DB).
		Exec()
	if err != nil {
		return err
	}
	return err
}
