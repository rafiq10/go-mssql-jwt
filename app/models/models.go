package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TF         string `json: "tf"`
	User_Name  string `json: "user_name"`
	Email      string `json: "email"`
	Salt       string `json: "salt"`
	Pwd        string `json: "pwd"`
	CreatedAt  int64  `json: "created_at"`
	UsrRole    string `json: "usr_role"`
	Department string `json: "department"`
}

func (u *User) Save(db *sql.DB) (*User, error) {
	h, err := hashPassword(u.Pwd)
	if err != nil {
		return nil, fmt.Errorf("hashPassword(u.Pwd)=%w", err)
	}

	fmt.Printf("User interface: %v", u)
	mySQl := `insert into auth.users 
	(tf,user_name, email, salt, pwd, created_at, usr_role, department) 
	values 
	('` + u.TF + `','` + u.User_Name + `','` + u.Email + `','` + u.Salt + `','` + string(h) + `',Now(),'` + u.UsrRole + `','` + u.Department + `') on conflict (tf) do nothing;`

	fmt.Println(mySQl)
	_, err = db.Exec(mySQl)

	if err != nil {
		fmt.Printf("Error saving user: %w", err)
		return nil, err
	}
	fmt.Println("User saved")
	return u, nil
}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password with bcrypt: %w", err)
	}
	return hash, nil
}

func comparePassword(password string, hashedPwd []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(password))
	if err != nil {
		return fmt.Errorf("Error comparing pwd: %w", err)
	}
	return nil

}
