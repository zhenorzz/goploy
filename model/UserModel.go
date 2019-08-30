package model

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User mysql table user
type User struct {
	ID             int64  `json:"id"`
	Account        string `json:"account"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	Mobile         string `json:"mobile"`
	Role           string `json:"role"`
	ManageGroupStr string `json:"manageGroupStr"`
	State          uint8  `json:"state"`
	CreateTime     int64  `json:"createTime"`
	UpdateTime     int64  `json:"updateTime"`
	LastLoginTime  int64  `json:"lastLoginTime"`
}

// Users many user
type Users []User

// GetData add user information to u *User
func (u User) GetData() (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT 
			id, 
			account, 
			password,
			name, 
			mobile, 
			role, 
			manage_group_str,
			state,
			create_time, 
			update_time  
		FROM 
			user 
		WHERE 
			id = ?`, u.ID).Scan(
		&user.ID,
		&user.Account,
		&user.Password,
		&user.Name,
		&user.Mobile,
		&user.Role,
		&user.ManageGroupStr,
		&user.State,
		&user.CreateTime,
		&user.UpdateTime)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetDataByAccount get user information
func (u User) GetDataByAccount() (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT 
			id, 
			account, 
			password,
			name, 
			mobile, 
			role, 
			manage_group_str,
			state,
			create_time, 
			update_time  
		FROM 
			user 
		WHERE 
		account = ?`, u.Account).Scan(
		&user.ID,
		&user.Account,
		&user.Password,
		&user.Name,
		&user.Mobile,
		&user.Role,
		&user.ManageGroupStr,
		&user.State,
		&user.CreateTime,
		&user.UpdateTime)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetList get many user row
func (u Users) GetList(pagination *Pagination) (Users, error) {
	rows, err := DB.Query(
		"SELECT id, account, name, mobile, role, manage_group_str, create_time, update_time FROM user WHERE state = 1 ORDER BY id DESC LIMIT ?, ?",
		(pagination.Page-1)*pagination.Rows,
		pagination.Rows)
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.Role, &user.ManageGroupStr, &user.CreateTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	err = DB.QueryRow(`SELECT COUNT(*) AS count FROM user`).Scan(&pagination.Total)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAll user row
func (u User) GetAll() (Users, error) {
	rows, err := DB.Query("SELECT id, account, name, mobile, create_time, update_time FROM user WHERE state = 1 ORDER BY id DESC")
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
func (u User) AddRow() (int64, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) AS count FROM user WHERE account = ?", u.Account).Scan(&count)
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
	result, err := DB.Exec(
		"INSERT INTO user (account, password, name, mobile, role, manage_group_str, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		u.Account,
		string(hashedPassword),
		u.Name,
		u.Mobile,
		u.Role,
		u.ManageGroupStr,
		u.CreateTime,
		u.UpdateTime,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table server
func (u User) EditRow() error {
	var err error
	if u.Password == "" {
		_, err = DB.Exec(
			`UPDATE user SET 
			  name = ?,
			  mobile = ?,
			  role = ?,
			  manage_group_str = ?
			WHERE
			 id = ?`,
			u.Name,
			u.Mobile,
			u.Role,
			u.ManageGroupStr,
			u.ID,
		)
	} else {
		password := []byte(u.Password)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = DB.Exec(
			`UPDATE user SET 
			  name = ?,
			  mobile = ?,
			  role = ?,
			  manage_group_str = ?,
			  password = ?
			WHERE
			 id = ?`,
			u.Name,
			u.Mobile,
			u.Role,
			u.ManageGroupStr,
			hashedPassword,
			u.ID,
		)
	}

	return err
}

// RemoveRow User
func (u User) RemoveRow() error {
	_, err := DB.Exec(
		`UPDATE user SET 
		  state = 0,
		  update_time = ?
		WHERE
		 id = ?`,
		u.UpdateTime,
		u.ID,
	)
	return err
}

// UpdatePassword return err
func (u User) UpdatePassword() error {
	password := []byte(u.Password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = DB.Exec(
		"UPDATE user SET password = ? where id = ?",
		string(hashedPassword),
		u.ID,
	)
	return err
}

// UpdateLastLoginTime return err
func (u User) UpdateLastLoginTime() error {
	_, err := DB.Exec(
		"UPDATE user SET last_login_time = ? where id = ?",
		u.LastLoginTime,
		u.ID,
	)
	return err
}

// Vaildate if user exists
func (u User) Vaildate(inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// CreateToken create token
func (u User) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   u.ID,
		"name": u.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"nbf":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))

	//Sign and get the complete encoded token as string
	return tokenString, err
}
