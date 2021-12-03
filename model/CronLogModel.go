package model

import (
	sq "github.com/Masterminds/squirrel"
)

const cronLogTable = "`cron_log`"

// CronLog -
type CronLog struct {
	ID         int64  `json:"id"`
	CronID     int64  `json:"cronId"`
	ServerID   int64  `json:"serverId"`
	ExecCode   int    `json:"execCode"`
	Message    string `json:"message"`
	ReportTime string `json:"reportTime"`
	InsertTime string `json:"insertTime"`
}

// CronLogs -
type CronLogs []CronLog

// GetList -
func (cl CronLog) GetList(pagination Pagination) (CronLogs, error) {
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
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
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
