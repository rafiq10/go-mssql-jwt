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
	fmt.Printf("user in Save: %v \n", u)
	h, err := hashPassword(u.Pwd)
	if err != nil {
		return nil, fmt.Errorf("hashPassword(u.Pwd)=%w", err)
	}

	mySQl := `insert into auth.users 
	(tf,user_name, email, salt, pwd, created_at, usr_role, department) 
	values 
	('` + u.TF + `','` + u.User_Name + `','` + u.Email + `','` + u.Salt + `','` + string(h) + `',Now(),'` + u.UsrRole + `','` + u.Department + `') on conflict (tf) do nothing;`

	_, err = db.Exec(mySQl)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) GetByTF(db *sql.DB, TF string) (*User, error) {

	mySQl := `select tf,user_name, email, salt, pwd, usr_role, department from auth.users where TF = '` + TF + `';`

	rows, err := db.Query(mySQl)
	if err != nil {
		fmt.Println(fmt.Errorf("error retrieving user: %v", err))
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&u.TF, &u.User_Name, &u.Email, &u.Salt, &u.Pwd, &u.UsrRole, &u.Department)
	}
	if err != nil {
		fmt.Println(fmt.Errorf("error scanning the row in db: %v", err))
		return &User{}, fmt.Errorf("error scanning the row in db: %v", err)
	} else {
		return u, nil
	}

}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password with bcrypt: %w", err)
	}
	return hash, nil
}

func comparePassword(password string, hashedPwd []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(password))
	if err != nil {
		return fmt.Errorf("error comparing pwd: %w", err)
	}
	return nil

}
