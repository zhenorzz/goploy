package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const namespaceTable = "`namespace`"

// Namespace -
type Namespace struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	UserID     int64  `json:"-"`
	Role       string `json:"role"`
	InsertTime string `json:"insertTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

// Namespaces -
type Namespaces []Namespace

// AddRow return LastInsertId
func (ns Namespace) AddRow() (int64, error) {
	result, err := sq.
		Insert(namespaceTable).
		Columns("name").
		Values(ns.Name).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow -
func (ns Namespace) EditRow() error {
	_, err := sq.
		Update(namespaceTable).
		Set("name", ns.Name).
		Where(sq.Eq{"id": ns.ID}).
		RunWith(DB).
		Exec()
	return err
}

// GetAllByUserID -
func (ns Namespace) GetAllByUserID() (Namespaces, error) {
	rows, err := sq.
		Select("namespace.id, namespace.name, role").
		From(namespaceTable).
		Join(namespaceUserTable + " ON namespace_user.namespace_id = namespace.id").
		Where(sq.Eq{"user_id": ns.UserID}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	namespaces := Namespaces{}
	for rows.Next() {
		var namespace Namespace
		if err := rows.Scan(&namespace.ID, &namespace.Name, &namespace.Role); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}
	return namespaces, nil
}

// GetListByUserID -
func (ns Namespace) GetListByUserID(pagination Pagination) (Namespaces, error) {
	rows, err := sq.
		Select("namespace.id, namespace.name, namespace.insert_time, namespace.update_time").
		From(namespaceTable).
		Join(namespaceUserTable + " ON namespace_user.namespace_id = namespace.id").
		Where(sq.Eq{
			"user_id": ns.UserID,
			"role":    []string{"admin", "manager"},
		}).
		OrderBy("namespace.id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	namespaces := Namespaces{}
	for rows.Next() {
		var namespace Namespace

		if err := rows.Scan(&namespace.ID, &namespace.Name, &namespace.InsertTime, &namespace.UpdateTime); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

// GetTotalByUserID -
func (ns Namespace) GetTotalByUserID() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(namespaceUserTable).
		Where(sq.Eq{
			"user_id": ns.UserID,
			"role":    []string{"admin", "manager"},
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
func (ns Namespace) GetData() (Namespace, error) {
	var namespace Namespace
	err := sq.
		Select("name, insert_time, update_time").
		From(namespaceTable).
		Where(sq.Eq{"id": ns.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&namespace.Name, &namespace.InsertTime, &namespace.UpdateTime)
	if err != nil {
		return namespace, errors.New("数据查询失败")
	}
	return namespace, nil
}
