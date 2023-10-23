// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const cronLogTable = "`cron_log`"

type CronLog struct {
	ID         int64  `json:"id"`
	CronID     int64  `json:"cronId"`
	ServerID   int64  `json:"serverId"`
	ExecCode   int    `json:"execCode"`
	Message    string `json:"message"`
	ReportTime string `json:"reportTime"`
	InsertTime string `json:"insertTime"`
}

type CronLogs []CronLog

func (cl CronLog) GetList(page, limit uint64) (CronLogs, error) {
	rows, err := sq.
		Select(
			"id",
			"cron_id",
			"server_id",
			"exec_code",
			"message",
			"report_time",
		).
		From(cronLogTable).
		Where(sq.Eq{"server_id": cl.ServerID}).
		Where(sq.Eq{"cron_id": cl.CronID}).
		Limit(limit).
		Offset((page - 1) * limit).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	cronLogs := CronLogs{}
	for rows.Next() {
		var cronLog CronLog
		if err := rows.Scan(
			&cronLog.ID,
			&cronLog.CronID,
			&cronLog.ServerID,
			&cronLog.ExecCode,
			&cronLog.Message,
			&cronLog.ReportTime,
		); err != nil {
			return nil, err
		}
		cronLogs = append(cronLogs, cronLog)
	}
	return cronLogs, nil
}

// AddRow -
func (cl CronLog) AddRow() error {
	_, err := sq.
		Insert(cronLogTable).
		Columns("cron_id", "server_id", "exec_code", "message", "report_time").
		Values(cl.CronID, cl.ServerID, cl.ExecCode, cl.Message, cl.ReportTime).
		RunWith(DB).
		Exec()
	return err
}
