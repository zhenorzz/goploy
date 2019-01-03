package model

import (
	"fmt"
	"time"
)

// Index struct
type Index struct {
}

// Query chart row
func (index *Index) Query(date string) error {
	type Chart struct {
		Hour  string
		Count int
	}
	type Charts []Chart
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return err
	}
	startTimestamp := t.Unix()
	endTimestamp := startTimestamp + 86400
	db := NewDB()
	rows, err := db.Query(
		"SELECT FROM_UNIXTIME(create_time, '%H')as hour, COUNT(status = 1) as count FROM deploy WHERE create_time BETWEEN ? AND ? GROUP BY hour",
		startTimestamp,
		endTimestamp)
	if err != nil {
		return err
	}
	charts := new(Charts)
	for rows.Next() {
		var chart Chart

		if err := rows.Scan(&chart.Hour, &chart.Count); err != nil {
			return err
		}
		*charts = append(*charts, chart)
	}
	fmt.Println(charts)
	return nil
}
