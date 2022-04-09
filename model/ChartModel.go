// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// Chart struct
type Chart struct {
	Hour         int `json:"hour"`
	CommitNumber int `json:"commitNumber"`
	DeployNumber int `json:"deployNumber"`
	FailNumber   int `json:"failNumber"`
	BackNumber   int `json:"backNumber"`
}

// Charts struct
type Charts []Chart

// Query chart row
func (c *Charts) Query(date string) error {
	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return err
	}
	startTimestamp := t.Unix()
	endTimestamp := startTimestamp + 86400
	rows, err := DB.Query(
		`SELECT 
		  FROM_UNIXTIME(create_time, '%H')as hour, 
		  COUNT(*) as commit, 
		  COUNT(status = 2 or null) as deploy, 
		  COUNT(status = 3 or null) as fail, 
		  COUNT(status = 4 or null) as back
		  FROM deploy 
		  WHERE create_time BETWEEN ? AND ? 
		  GROUP BY hour
		  Order BY null`,
		startTimestamp,
		endTimestamp)
	if err != nil {
		return err
	}
	for rows.Next() {
		var chart Chart
		if err := rows.Scan(&chart.Hour, &chart.CommitNumber, &chart.DeployNumber, &chart.FailNumber, &chart.BackNumber); err != nil {
			return err
		}
		*c = append(*c, chart)
	}
	return nil
}
