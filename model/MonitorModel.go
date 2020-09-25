package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const monitorTable = "`monitor`"

// Monitor -
type Monitor struct {
	ID           int64  `json:"id"`
	NamespaceID  int64  `json:"namespaceId"`
	Name         string `json:"name"`
	Domain       string `json:"domain"`
	Port         int    `json:"port"`
	Second       int    `json:"second"`
	Times        uint16 `json:"times"`
	NotifyType   uint8  `json:"notifyType"`
	NotifyTarget string `json:"notifyTarget"`
	NotifyTimes  uint16 `json:"notifyTimes"`
	Description  string `json:"description"`
	ErrorContent string `json:"errorContent"`
	State        uint8  `json:"state"`
	InsertTime   string `json:"insertTime"`
	UpdateTime   string `json:"updateTime"`
}

// Monitors -
type Monitors []Monitor

// GetList -
func (m Monitor) GetList(pagination Pagination) (Monitors, error) {
	rows, err := sq.
		Select("id, name, domain, port, second, times, notify_type, notify_target, notify_times, description, error_content, state, insert_time, update_time").
		From(monitorTable).
		Where(sq.Eq{
			"namespace_id": m.NamespaceID,
		}).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	monitors := Monitors{}
	for rows.Next() {
		var monitor Monitor

		if err := rows.Scan(
			&monitor.ID,
			&monitor.Name,
			&monitor.Domain,
			&monitor.Port,
			&monitor.Second,
			&monitor.Times,
			&monitor.NotifyType,
			&monitor.NotifyTarget,
			&monitor.NotifyTimes,
			&monitor.Description,
			&monitor.ErrorContent,
			&monitor.State,
			&monitor.InsertTime,
			&monitor.UpdateTime); err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

// GetTotal -
func (m Monitor) GetTotal() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(monitorTable).
		Where(sq.Eq{
			"namespace_id": m.NamespaceID,
		}).
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetData -
func (m Monitor) GetData() (Monitor, error) {
	var monitor Monitor
	err := sq.
		Select("id, name, domain, ip, port, second, times, notify_type, notify_target, notify_times, state").
		From(monitorTable).
		Where(sq.Eq{"id": m.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&monitor.ID, &monitor.Name, &monitor.Domain, &monitor.Port, &monitor.Second, &monitor.Times, &monitor.NotifyType, &monitor.NotifyTarget, &monitor.NotifyTimes, &monitor.State)
	if err != nil {
		return monitor, errors.New("数据查询失败")
	}
	return monitor, nil
}

// GetAllByState -
func (m Monitor) GetAllByState() (Monitors, error) {
	rows, err := sq.
		Select("id, name, domain, port, second, times, notify_type, notify_target, notify_times, description").
		From(monitorTable).
		Where(sq.Eq{
			"state": m.State,
		}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	monitors := Monitors{}
	for rows.Next() {
		var monitor Monitor

		if err := rows.Scan(
			&monitor.ID,
			&monitor.Name,
			&monitor.Domain,
			&monitor.Port,
			&monitor.Second,
			&monitor.Times,
			&monitor.NotifyType,
			&monitor.NotifyTarget,
			&monitor.NotifyTimes,
			&monitor.Description); err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

// AddRow return LastInsertId
func (m Monitor) AddRow() (int64, error) {
	result, err := sq.
		Insert(monitorTable).
		Columns("namespace_id", "name", "domain", "port", "second", "times", "notify_type", "notify_target", "notify_times", "description", "error_content").
		Values(m.NamespaceID, m.Name, m.Domain, m.Port, m.Second, m.Times, m.NotifyType, m.NotifyTarget, m.NotifyTimes, m.Description, "").
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (m Monitor) EditRow() error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"name":          m.Name,
			"domain":        m.Domain,
			"port":          m.Port,
			"second":        m.Second,
			"times":         m.Times,
			"notify_type":   m.NotifyType,
			"notify_target": m.NotifyTarget,
			"notify_times":  m.NotifyTimes,
			"description":   m.Description,
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

// ToggleState -
func (m Monitor) ToggleState() error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"state": sq.Expr("!state"),
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRow -
func (m Monitor) DeleteRow() error {
	_, err := sq.
		Delete(monitorTable).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}

// TurnOff -
func (m Monitor) TurnOff(errorContent string) error {
	_, err := sq.
		Update(monitorTable).
		SetMap(sq.Eq{
			"state": Disable,
			"error_content": errorContent,
		}).
		Where(sq.Eq{"id": m.ID}).
		RunWith(DB).
		Exec()
	return err
}