// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const serverMonitorTable = "`server_monitor`"

// ServerMonitor -
type ServerMonitor struct {
	ID           int64  `json:"id"`
	ServerID     int64  `json:"serverId"`
	Item         string `json:"item"`
	Formula      string `json:"formula"`
	Operator     string `json:"operator"`
	Value        string `json:"value"`
	GroupCycle   int    `json:"groupCycle"`
	LastCycle    int    `json:"lastCycle"`
	SilentCycle  int    `json:"silentCycle"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	NotifyType   uint8  `json:"notifyType"`
	NotifyTarget string `json:"notifyTarget"`
	InsertTime   string `json:"insertTime"`
	UpdateTime   string `json:"updateTime"`
}

// ServerMonitors -
type ServerMonitors []ServerMonitor

// GetAll -
func (sm ServerMonitor) GetAll() (ServerMonitors, error) {
	rows, err := sq.
		Select(
			"id",
			"server_id",
			"item",
			"formula",
			"operator",
			"value",
			"group_cycle",
			"last_cycle",
			"silent_cycle",
			"start_time",
			"end_time",
			"notify_type",
			"notify_target",
			"insert_time",
			"update_time",
		).
		From(serverMonitorTable).
		Where(sq.Eq{"server_id": sm.ServerID}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	serverMonitors := ServerMonitors{}
	for rows.Next() {
		var serverMonitor ServerMonitor
		if err := rows.Scan(
			&serverMonitor.ID,
			&serverMonitor.ServerID,
			&serverMonitor.Item,
			&serverMonitor.Formula,
			&serverMonitor.Operator,
			&serverMonitor.Value,
			&serverMonitor.GroupCycle,
			&serverMonitor.LastCycle,
			&serverMonitor.SilentCycle,
			&serverMonitor.StartTime,
			&serverMonitor.EndTime,
			&serverMonitor.NotifyType,
			&serverMonitor.NotifyTarget,
			&serverMonitor.InsertTime,
			&serverMonitor.UpdateTime,
		); err != nil {
			return nil, err
		}
		serverMonitors = append(serverMonitors, serverMonitor)
	}
	return serverMonitors, nil
}

// GetAllModBy -
func (sm ServerMonitor) GetAllModBy(number int, time string) (ServerMonitors, error) {
	rows, err := sq.
		Select(
			"id",
			"server_id",
			"item",
			"formula",
			"operator",
			"value",
			"group_cycle",
			"last_cycle",
			"silent_cycle",
			"start_time",
			"end_time",
			"notify_type",
			"notify_target",
		).
		From(serverMonitorTable).
		Where("? % `group_cycle` = 0", number).
		Where(sq.GtOrEq{"end_time": time}).
		Where(sq.LtOrEq{"start_time": time}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	serverMonitors := ServerMonitors{}
	for rows.Next() {
		var serverMonitor ServerMonitor
		if err := rows.Scan(
			&serverMonitor.ID,
			&serverMonitor.ServerID,
			&serverMonitor.Item,
			&serverMonitor.Formula,
			&serverMonitor.Operator,
			&serverMonitor.Value,
			&serverMonitor.GroupCycle,
			&serverMonitor.LastCycle,
			&serverMonitor.SilentCycle,
			&serverMonitor.StartTime,
			&serverMonitor.EndTime,
			&serverMonitor.NotifyType,
			&serverMonitor.NotifyTarget,
		); err != nil {
			return nil, err
		}
		serverMonitors = append(serverMonitors, serverMonitor)
	}
	return serverMonitors, nil
}

// AddRow return LastInsertId
func (sm ServerMonitor) AddRow() (int64, error) {
	result, err := sq.
		Insert(serverMonitorTable).
		Columns(
			"server_id",
			"item",
			"formula",
			"operator",
			"value",
			"group_cycle",
			"last_cycle",
			"silent_cycle",
			"start_time",
			"end_time",
			"notify_type",
			"notify_target",
		).
		Values(
			sm.ServerID,
			sm.Item,
			sm.Formula,
			sm.Operator,
			sm.Value,
			sm.GroupCycle,
			sm.LastCycle,
			sm.SilentCycle,
			sm.StartTime,
			sm.EndTime,
			sm.NotifyType,
			sm.NotifyTarget,
		).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (sm ServerMonitor) EditRow() error {
	_, err := sq.
		Update(serverMonitorTable).
		SetMap(sq.Eq{
			"item":          sm.Item,
			"formula":       sm.Formula,
			"operator":      sm.Operator,
			"value":         sm.Value,
			"group_cycle":   sm.GroupCycle,
			"last_cycle":    sm.LastCycle,
			"silent_cycle":  sm.SilentCycle,
			"start_time":    sm.StartTime,
			"end_time":      sm.EndTime,
			"notify_type":   sm.NotifyType,
			"notify_target": sm.NotifyTarget,
		}).
		Where(sq.Eq{"id": sm.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRow -
func (sm ServerMonitor) DeleteRow() error {
	_, err := sq.
		Delete(serverMonitorTable).
		Where(sq.Eq{"id": sm.ID}).
		RunWith(DB).
		Exec()
	return err
}
