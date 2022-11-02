// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import sq "github.com/Masterminds/squirrel"

const crontabTable = "`crontab`"

type Crontab struct {
	ID          int64  `json:"id"`
	ServerID    int64  `json:"serverId"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	Save        uint8  `json:"save"`
	State       uint8  `json:"state"`
	InsertTime  string `json:"insertTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

type Crontabs []Crontab

func (c Crontab) GetList() (Crontabs, error) {
	rows, err := sq.
		Select(
			"id",
			"server_id",
			"expression",
			"description",
			"save",
			"state",
			"insert_time",
			"update_time",
		).
		From(crontabTable).
		Where(sq.Eq{"server_id": c.ServerID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	crons := Crontabs{}
	for rows.Next() {
		var cron Crontab
		if err := rows.Scan(
			&cron.ID,
			&cron.ServerID,
			&cron.Expression,
			&cron.Description,
			&cron.Save,
			&cron.State,
			&cron.InsertTime,
			&cron.UpdateTime,
		); err != nil {
			return nil, err
		}
		crons = append(crons, cron)
	}
	return crons, nil
}

func (c Crontab) AddRow() (int64, error) {
	result, err := sq.
		Insert(crontabTable).
		Columns(
			"server_id",
			"expression",
			"description",
			"save",
			"state",
		).
		Values(
			c.ServerID,
			c.Expression,
			c.Description,
			c.Save,
			c.State,
		).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (c Crontab) EditRow() error {
	_, err := sq.
		Update(crontabTable).
		SetMap(sq.Eq{
			"expression":  c.Expression,
			"description": c.Description,
			"save":        c.Save,
			"state":       c.State,
		}).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}

func (c Crontab) RemoveRow() error {
	_, err := sq.
		Delete(crontabTable).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}
