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

// GetData add user information to u *User
func (u User) GetData() (User, error) {
	var user User
	db := NewDB()
	err := db.QueryRow("SELECT account, name, role_id FROM user WHERE id = ? AND state = ?", u.ID, 1).Scan(&user.Account, &user.Name, &user.RoleID)
	if err != nil {
		return user, errors.New("数据查询失败")
	}
	return user, nil
}

// GetList get many user row
func (u Users) GetList(pagination *Pagination) (Users, error) {
	db := NewDB()
	rows, err := db.Query(
		"SELECT id, account, name, mobile, role_id, create_time, update_time FROM user ORDER BY id DESC LIMIT ?, ?",
		(pagination.Page-1)*pagination.Rows,
		pagination.Rows)
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.RoleID, &user.CreateTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	err = db.QueryRow(`SELECT COUNT(*) AS count FROM user`).Scan(&pagination.Total)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAll user row
func (u User) GetAll() (Users, error) {
	db := NewDB()
	rows, err := db.Query("SELECT id, account, name, mobile, create_time, update_time FROM user ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.CreateTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// AddRow add one row to table user and add id to u.ID
func (u User) AddRow() (uint32, error) {
	db := NewDB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) AS count FROM user WHERE account = ?", u.Account).Scan(&count)
	fmt.Println(count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("账号已存在")
	}

	if u.Password == "" {
		u.Password = u.Account + "!@#"
	}
	password := []byte(u.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return 0, err
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
	return uint32(id), err
}

// EditRow edit one row to table server
func (u User) EditRow() error {
	var err error
	db := NewDB()
	if u.Password == "" {
		_, err = db.Exec(
			`UPDATE user SET 
			  name = ?,
			  mobile = ?,
			  role_id = ?
			WHERE
			 id = ?`,
			u.Name,
			u.Mobile,
			u.RoleID,
			u.ID,
		)
	} else {
		password := []byte(u.Password)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec(
			`UPDATE user SET 
			  name = ?,
			  mobile = ?,
			  role_id = ?,
			  password = ?
			WHERE
			 id = ?`,
			u.Name,
			u.Mobile,
			u.RoleID,
			hashedPassword,
			u.ID,
		)
	}

	return err
}

// Vaildate if user exists
func (u User) Vaildate() error {
	var hashPassword string
	db := NewDB()
	err := db.QueryRow("SELECT password FROM user WHERE id = ?", u.ID).Scan(&hashPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(u.Password))
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// VaildateByAccount if user exists
func (u User) VaildateByAccount() (User, error) {
	var user User
	var hashPassword string
	db := NewDB()
	err := db.QueryRow("SELECT id, password, name FROM user WHERE account = ?", u.Account).Scan(&user.ID, &hashPassword, &user.Name)
	if err != nil {
		return user, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(u.Password))
	if err != nil {
		return user, errors.New("密码错误")
	}
	return user, nil
}

// UpdatePassword return err
func (u User) UpdatePassword() error {
	password := []byte(u.Password)
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
