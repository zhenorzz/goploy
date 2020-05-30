package model

import (
	"errors"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const userTable = "`user`"

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
	InsertTime     string `json:"insertTime"`
	UpdateTime     string `json:"updateTime"`
	LastLoginTime  string `json:"lastLoginTime"`
}

// Users many user
type Users []User

// GetData add user information to u *User
func (u User) GetData() (User, error) {
	var user User
	err := sq.
		Select("id, account, password, name, mobile, role, manage_group_str, state, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"id": u.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&user.ID, &user.Account, &user.Password, &user.Name, &user.Mobile, &user.Role, &user.ManageGroupStr, &user.State, &user.InsertTime, &user.UpdateTime)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetDataByAccount get user information
func (u User) GetDataByAccount() (User, error) {
	var user User
	err := sq.
		Select("id, account, password, name, mobile, role, manage_group_str, state, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"account": u.Account}).
		RunWith(DB).
		QueryRow().
		Scan(&user.ID, &user.Account, &user.Password, &user.Name, &user.Mobile, &user.Role, &user.ManageGroupStr, &user.State, &user.InsertTime, &user.UpdateTime)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetList get many user row
func (u Users) GetList(pagination Pagination) (Users, Pagination, error) {
	rows, err := sq.
		Select("id, account, name, mobile, role, manage_group_str, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"state": Enable}).
		OrderBy("id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, pagination, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.Role, &user.ManageGroupStr, &user.InsertTime, &user.UpdateTime); err != nil {
			return users, pagination, err
		}
		users = append(users, user)
	}
	err = sq.
		Select("COUNT(*) AS count").
		From(userTable).
		Where(sq.Eq{"state": Enable}).
		RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return users, pagination, nil
}

// GetAll user row
func (u User) GetAll() (Users, error) {
	rows, err := sq.
		Select("id, account, name, mobile, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"state": Enable}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.InsertTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetAll user row
func (u User) GetAllByRole() (Users, error) {
	rows, err := sq.
		Select("id, account, name, mobile, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"role": u.Role}).
		Where(sq.Eq{"state": Enable}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.InsertTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetAll user row
func (u User) GetCanBindProjectUser() (Users, error) {
	rows, err := sq.
		Select("id, account, name, mobile, insert_time, update_time").
		From(userTable).
		Where(sq.Eq{"role": []string{"group-manager", "member"}}).
		Where(sq.Eq{"state": Enable}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Mobile, &user.InsertTime, &user.UpdateTime); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// AddRow add one row to table user and add id to u.ID
func (u User) AddRow() (int64, error) {
	if u.Password == "" {
		u.Password = u.Account + "!@#"
	}
	password := []byte(u.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	result, err := sq.
		Insert(userTable).
		Columns("account", "password", "name", "mobile", "role", "manage_group_str").
		Values(u.Account, string(hashedPassword), u.Name, u.Mobile, u.Role, u.ManageGroupStr).
		RunWith(DB).
		Exec()

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table server
func (u User) EditRow() error {
	builder := sq.
		Update(userTable).
		SetMap(sq.Eq{
			"name":             u.Name,
			"mobile":           u.Mobile,
			"role":             u.Role,
			"manage_group_str": u.ManageGroupStr,
		}).
		Where(sq.Eq{"id": u.ID})
	if u.Password != "" {
		password := []byte(u.Password)
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		builder = builder.Set("password", hashedPassword)
	}
	_, err := builder.RunWith(DB).Exec()
	return err
}

// RemoveRow User
func (u User) RemoveRow() error {
	_, err := sq.
		Update(userTable).
		SetMap(sq.Eq{
			"state": Disable,
		}).
		Where(sq.Eq{"id": u.ID}).
		RunWith(DB).
		Exec()
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
	_, err = sq.
		Update(userTable).
		SetMap(sq.Eq{
			"password": string(hashedPassword),
		}).
		Where(sq.Eq{"id": u.ID}).
		RunWith(DB).
		Exec()
	return err
}

// UpdateLastLoginTime return err
func (u User) UpdateLastLoginTime() error {
	_, err := sq.
		Update(userTable).
		SetMap(sq.Eq{
			"last_login_time": u.LastLoginTime,
		}).
		Where(sq.Eq{"id": u.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Validate if user exists
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
