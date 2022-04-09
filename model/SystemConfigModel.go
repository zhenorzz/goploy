// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	sq "github.com/Masterminds/squirrel"
)

const systemConfigTable = "`system_config`"

// SystemConfig -
type SystemConfig struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetDataByKey -
func (sc SystemConfig) GetDataByKey() (SystemConfig, error) {
	var systemConfig SystemConfig
	err := sq.
		Select("id, `key`, value").
		From(systemConfigTable).
		Where(sq.Eq{"`key`": sc.Key}).
		RunWith(DB).
		QueryRow().
		Scan(&systemConfig.ID, &systemConfig.Key, &systemConfig.Value)
	if err != nil {
		return systemConfig, err
	}
	return systemConfig, nil
}

// EditRowByKey -
func (sc SystemConfig) EditRowByKey() error {
	builder := sq.
		Update(systemConfigTable).
		SetMap(sq.Eq{
			"value": sc.Value,
		}).
		Where(sq.Eq{"`key`": sc.Key})
	_, err := builder.RunWith(DB).Exec()
	return err
}
