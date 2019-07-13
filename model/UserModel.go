package model

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User mysql table user
type User struct {
	ID         uint32 `json:"id"`
	Account    string `json:"account"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	RoleID     uint32 `json:"roleId"`
	State      uint8  `json:"state"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

// Users many user
type Users []User

// QueryRow add user information to u *User
func (u *User) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT account, name, role_id FROM user WHERE id = ? AND state = ?", u.ID, 1).Scan(&u.Account, &u.Name, &u.RoleID)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}

// Query user row
func (u *Users) Query(pagination *Pagination) error {
	db := NewDB()
	rows, err := db.Query(
		"SELECT id, account, name, mobile, create_time, update_time FROM user ORDER BY id DESC LIMIT ?, ?",
		(pagination.Page-1)*pagination.Rows,
		pagination.Rows)
	if err != nil {
		return err
	}
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.CreateTime, &user.UpdateTime); err != nil {
			return err
		}
		*u = append(*u, user)
	}
	err = db.QueryRow(`SELECT COUNT(*) AS count FROM user`).Scan(&pagination.Total)
	if err != nil {
		return err
	}
	return nil
}

// QueryAll user row
func (u *Users) QueryAll() error {
	db := NewDB()
	rows, err := db.Query("SELECT id, account, name, mobile, create_time, update_time FROM user ORDER BY id DESC")
	if err != nil {
		return err
	}
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.CreateTime, &user.UpdateTime); err != nil {
			return err
		}
		*u = append(*u, user)
	}
	return nil
}

// AddRow add one row to table user and add id to u.ID
func (u *User) AddRow() error {
	db := NewDB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) AS count FROM user WHERE account = ?", u.Account).Scan(&count)
	fmt.Println(count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("账号已存在")
	}
	password := []byte(u.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	result, err := db.Exec(
		"INSERT INTO user (account, password, name, mobile, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)",
		u.Account,
		string(hashedPassword),
		u.Name,
		u.Mobile,
		u.CreateTime,
		u.UpdateTime,
	)
	id, err := result.LastInsertId()
	u.ID = uint32(id)
	return err
}

// Vaildate if user exists
func (u *User) Vaildate() error {
	var hashPassword string
	db := NewDB()
	err := db.QueryRow("SELECT id, password, name FROM user WHERE account = ?", u.Account).Scan(&u.ID, &hashPassword, &u.Name)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(u.Password))
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// UpdatePassword return err
func (u *User) UpdatePassword(newPassword string) error {
	password := []byte(newPassword)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	db := NewDB()
	_, err = db.Exec(
		"UPDATE user SET password = ? where id = ?",
		string(hashedPassword),
		u.ID,
	)
	return err
}
