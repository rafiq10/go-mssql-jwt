package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	TF         string `json: "tf"`
	User_Name  string `json: "user_name"`
	Email      string `json: "email"`
	Pwd        string `json: "pwd"`
	CreatedAt  int64  `json: "created_at"`
	UsrRole    string `json: "usr_role"`
	Department string `json: "department"`
	SessionId  string `json: "session_id"`
}

func (u *User) Save(db *sql.DB) (*User, error) {
	h, err := hashPassword(u.Pwd)
	if err != nil {
		return nil, fmt.Errorf("hashPassword(u.Pwd)=%w", err)
	}

	mySQl := `insert into auth.users 
	(tf,user_name, email, pwd, created_at, usr_role, department,sid) 
	values 
	('` + u.TF + `','` + u.User_Name + `','` + u.Email + `','` + string(h) + `',Now(),'` + u.UsrRole + `','` + u.Department + `','` + u.SessionId + `') on conflict (tf) do nothing;`

	_, err = db.Exec(mySQl)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) GetByTF(db *sql.DB, TF string) (*User, error) {
	mySQl := `select tf,user_name, email, pwd, usr_role, department,session_id from auth.users where TF = '` + TF + `';`
	return getUserFromQuery(db, mySQl)
}

func (u *User) GetBySessionID(db *sql.DB, sID string) (*User, error) {
	mySQl := `select tf,user_name, email, pwd, usr_role, department,session_id from auth.users where session_id = '` + sID + `';`
	usr, err := getUserFromQuery(db, mySQl)
	if err != nil {
		log.Printf("GetBySessionID err: %v \n", err)
		return nil, err
	}

	return usr, nil

}
func (u *User) UpdateSID(db *sql.DB) error {

	session_id := uuid.New().String()

	mySQl := `update auth.users set session_id = '` + session_id + `' where TF = '` + u.TF + `';`
	_, err := db.Exec(mySQl)
	if err != nil {
		fmt.Errorf("(%v) UpdateSID( %s, %s)=%w", u, u.TF, session_id, err)
	}
	u.SessionId = session_id
	return nil
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

func getUserFromQuery(db *sql.DB, mySQl string) (*User, error) {
	var u *User

	rows, err := db.Query(mySQl)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if rows.Err() != nil {
			fmt.Errorf("Error iterating rows")
		}
		var TF, User_Name, Email, Pwd, UsrRole, Department, SessionId string
		err = rows.Scan(&TF, &User_Name, &Email, &Pwd, &UsrRole, &Department, &SessionId)

		if err != nil {
			return &User{}, fmt.Errorf("error scanning the row in db: %v", err)
		} else {
			u = &User{
				TF:         TF,
				User_Name:  User_Name,
				Email:      Email,
				Pwd:        Pwd,
				UsrRole:    UsrRole,
				Department: Department,
				SessionId:  SessionId,
			}
			return u, nil
		}

	}
	return nil, nil
}
