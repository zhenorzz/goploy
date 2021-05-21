package model

import (
	sq "github.com/Masterminds/squirrel"
)

const crontabTable = "`crontab`"

// Crontab -
type Crontab struct {
	ID          int64  `json:"id"`
	NamespaceID int64  `json:"namespaceId"`
	Command     string `json:"command"`
	CommandMD5  string `json:"commandMD5"`
	Creator     string `json:"creator"`
	CreatorID   int64  `json:"creatorId"`
	Editor      string `json:"editor"`
	EditorID    int64  `json:"editorId"`
	InsertTime  string `json:"insertTime"`
	UpdateTime  string `json:"updateTime"`
}

// Crontabs -
type Crontabs []Crontab

// GetList -
func (c Crontab) GetList(pagination Pagination) (Crontabs, error) {
	builder := sq.
		Select("id, command, creator, creator_id, editor, editor_id, insert_time, update_time").
		From(crontabTable).
		Where(sq.Eq{"namespace_id": c.NamespaceID})
	if len(c.Command) > 0 {
		builder = builder.Where(sq.Like{"command": "%" + c.Command + "%"})
	}
	rows, err := builder.
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	crontabs := Crontabs{}
	for rows.Next() {
		var crontab Crontab
		if err := rows.Scan(&crontab.ID, &crontab.Command, &crontab.Creator, &crontab.CreatorID, &crontab.Editor, &crontab.EditorID, &crontab.InsertTime, &crontab.UpdateTime); err != nil {
			return nil, err
		}
		crontabs = append(crontabs, crontab)
	}
	return crontabs, nil
}

// GetTotal -
func (c Crontab) GetTotal() (int64, error) {
	var total int64
	builder := sq.
		Select("COUNT(*) AS count").
		From(crontabTable).
		Where(sq.Eq{"namespace_id": c.NamespaceID})
	if len(c.Command) > 0 {
		builder = builder.Where(sq.Like{"command": "%" + c.Command + "%"})
	}

	err := builder.
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetAllInCommandMD5 return all row in command md5
func (c Crontab) GetAllInCommandMD5(commandMD5s []string) (Crontabs, error) {
	rows, err := sq.
		Select("id, command, command_md5").
		From(crontabTable).
		Where(sq.Eq{"command_md5": commandMD5s}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	crontabs := Crontabs{}
	for rows.Next() {
		var crontab Crontab
		if err := rows.Scan(&crontab.ID, &crontab.Command, &crontab.CommandMD5); err != nil {
			return nil, err
		}
		crontabs = append(crontabs, crontab)
	}
	return crontabs, nil
}

// GetData return Crontab
func (c Crontab) GetData() (Crontab, error) {
	var crontab Crontab
	err := sq.
		Select("id, command, creator, creator_id, editor, editor_id").
		From(crontabTable).
		Where(sq.Eq{"id": c.ID}).
		OrderBy("id DESC").
		RunWith(DB).
		QueryRow().
		Scan(&crontab.ID, &crontab.Command, &crontab.Creator, &crontab.CreatorID, &crontab.Editor, &crontab.EditorID)
	if err != nil {
		return crontab, err
	}
	return crontab, nil
}

// AddRow return LastInsertId
func (c Crontab) AddRow() (int64, error) {
	result, err := sq.
		Insert(crontabTable).
		Columns("namespace_id", "command", "command_md5", "creator", "creator_id").
		Values(c.NamespaceID, c.Command, sq.Expr("md5(?)", c.Command), c.Creator, c.CreatorID).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// AddRowsInCommand add many rows in command
func (c Crontab) AddRowsInCommand(commands []string) error {
	builder := sq.
		Insert(crontabTable).
		Columns("command", "command_md5", "creator", "creator_id")
	for _, command := range commands {
		builder = builder.Values(command, sq.Expr("md5(?)", command), c.Creator, c.CreatorID)
	}
	_, err := builder.RunWith(DB).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

// EditRow -
func (c Crontab) EditRow() error {
	_, err := sq.
		Update(crontabTable).
		SetMap(sq.Eq{
			"command":     c.Command,
			"command_md5": sq.Expr("md5(?)", c.Command),
			"editor":      c.Editor,
			"editor_id":   c.EditorID,
		}).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}

// DeleteRow -
func (c Crontab) DeleteRow() error {
	_, err := sq.
		Delete(crontabTable).
		Where(sq.Eq{"id": c.ID}).
		RunWith(DB).
		Exec()
	return err
}
