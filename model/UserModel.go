package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User mysql table user
type User struct {
	ID         uint32
	Account    string
	Password   string
	Name       string
	Role       string
	State      uint8
	CreateTime uint32
	UpdateTime uint32
}

// QueryRow add user information to u *User
func (u *User) QueryRow() error {
	db := NewDB()
	err := db.QueryRow("SELECT name, role FROM user WHERE id = ? and state = ?", u.ID, 1).Scan(&u.Name, &u.Role)
	if err != nil {
		return errors.New("数据查询失败")
	}
	return nil
}

// Query user row
func (u *User) Query() {
	db := NewDB()
	err := db.QueryRow("SELECT id, account FROM user").Scan(&u.ID, &u.Account)
	if err != nil {
		log.Fatal(err)
	}
}

// Vaildate if user exists
func (u *User) Vaildate() error {

	// password := []byte(u.Password)

	// // Hashing the password with the default cost of 10
	// hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(hashedPassword))

	// Comparing the password with the hash

	var hashPassword string
	db := NewDB()
	err := db.QueryRow("SELECT id, password FROM user WHERE account = ?", u.Account).Scan(&u.ID, &hashPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(u.Password))
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// func (u *User) QueryMany() {
// 	db := NewDB()
// 	rows, err := db.Query("SELECT * FROM ceshi")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for rows.Next() {
// 		var ceshi Ceshi

// 		if err := rows.Scan(&ceshi.id, &ceshi.name); err != nil {
// 			log.Fatal(err)
// 		}

// 		*ceshis = append(*ceshis, ceshi)
// 	}
// }
