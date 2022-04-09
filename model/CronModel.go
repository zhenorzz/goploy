// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import sq "github.com/Masterminds/squirrel"

const cronTable = "`cron`"

type Cron struct {
	ID          int64  `json:"id"`
	ServerID    int64  `json:"serverId"`
	Expression  string `json:"expression"`
	Command     string `json:"command"`
	SingleMode  uint8  `json:"singleMode"`
	LogLevel    uint8  `json:"logLevel"`
	Description string `json:"description"`
	State       uint8  `json:"state"`
	Creator     string `json:"creator"`
	Editor      string `json:"editor"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type Crons []Cron

func (c Cron) GetList() (Crons, error) {
	rows, err := sq.
		Select(
			"id",
			"server_id",
			"expression",
			"command",
			"single_mode",
			"log_level",
			"description",
			"state",
			"creator",
			"editor",
			"insert_time",
			"update_time",
		).
		From(cronTable).
		Where(sq.Eq{"server_id": c.ServerID}).
		Where(sq.Eq{"state": Enable}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	crons := Crons{}
	for rows.Next() {
		var cron Cron
		if err := rows.Scan(
			&cron.ID,
			&cron.ServerID,
			&cron.Expression,
			&cron.Command,
			&cron.SingleMode,
			&cron.LogLevel,
			&cron.Description,
			&cron.State,
			&cron.Creator,
			&cron.Editor,
			&cron.InsertTime,
			&cron.UpdateTime,
		); err != nil {
			return nil, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}

func (c Cron) AddRow() (int64, error) {
	result, err := sq.
		Insert(cronTable).
		Columns(
			"server_id",
			"expression",
			"command",
			"single_mode",
			"log_level",
			"description",
			"creator",
		).
		Values(
			c.ServerID,
			c.Expression,
			c.Command,
			c.SingleMode,
			c.LogLevel,
			c.Description,
			c.Creator,
		).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (c Cron) EditRow() error {
	_, err := sq.
		Update(cronTable).
		SetMap(sq.Eq{
			"expression":  c.Expression,
			"command":     c.Command,
			"single_mode": c.SingleMode,
			"log_level":   c.LogLevel,
			"description": c.Description,
			"editor":      c.Editor,
		}).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (c Cron) RemoveRow() error {
	_, err := sq.
		Update(cronTable).
		SetMap(sq.Eq{
			"state": Disable,
		}).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}
