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
	err := DB.QueryRow(`
		SELECT 
			id, 
			account, 
			password,
			name, 
			mobile, 
			role_id, 
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
		&user.RoleID,
		&user.State,
		&user.CreateTime,
		&user.UpdateTime)
	if err != nil {
		return user, errors.New("数据查询失败")
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
			role_id, 
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
		&user.RoleID,
		&user.State,
		&user.CreateTime,
		&user.UpdateTime)
	if err != nil {
		return user, errors.New("数据查询失败")
	}
	return user, nil
}

// GetList get many user row
func (u Users) GetList(pagination *Pagination) (Users, error) {
	rows, err := DB.Query(
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
	err = DB.QueryRow(`SELECT COUNT(*) AS count FROM user`).Scan(&pagination.Total)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAll user row
func (u User) GetAll() (Users, error) {
	rows, err := DB.Query("SELECT id, account, name, mobile, create_time, update_time FROM user ORDER BY id DESC")
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
		"INSERT INTO user (account, password, name, mobile, role_id, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.Account,
		string(hashedPassword),
		u.Name,
		u.Mobile,
		u.RoleID,
		u.CreateTime,
		u.UpdateTime,
	)
	id, err := result.LastInsertId()
	return uint32(id), err
}

// EditRow edit one row to table server
func (u User) EditRow() error {
	var err error
	if u.Password == "" {
		_, err = DB.Exec(
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
		_, err = DB.Exec(
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
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
		"nbf":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))

	//Sign and get the complete encoded token as string
	return tokenString, err
}
