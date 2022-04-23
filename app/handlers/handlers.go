package handlers

import (
	"app/errhdl"
	"app/models"
	"app/mydb"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	LoggedIn bool
	UsrName  string
	ErrMsg   string
	ErrCode  int
}

func Index(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.ParseFiles(wd + "/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		JSONHandleError(w, err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	db, err := mydb.GetDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	u := &models.User{
		TF:         r.FormValue("tf"),
		User_Name:  r.FormValue("user_name"),
		Email:      r.FormValue("email"),
		Salt:       "",
		Pwd:        r.FormValue("pwd"),
		CreatedAt:  time.Now().Unix(),
		UsrRole:    r.FormValue("user_role"),
		Department: r.FormValue("department"),
	}

	usr, err := u.Save(db)
	if err != nil {
		JSONHandleError(w, err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.ParseFiles(wd + "/templates/register.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, usr)
}

func LogIn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	var d *LoginData
	var tmpl *template.Template

	if r.Method != http.MethodPost {
		d, tmpl = getLoginData(http.StatusMethodNotAllowed, "")
		tmpl.Execute(w, d)
		log.Fatalf("r.Method != http.MethodPost")
	}

	pwd := r.FormValue("pwd")
	if pwd == "" {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		tmpl.Execute(w, d)
		log.Fatalf("pwd empty")
	}

	tf := r.FormValue("tf")
	if tf == "" {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		tmpl.Execute(w, d)
		log.Fatalf("ptf empty")
	}

	db, err := mydb.GetDb()
	if err != nil {
		d, tmpl = getLoginData(http.StatusNotFound, "")
		tmpl.Execute(w, d)
		log.Fatalf("mydb.GetDb()=%v", err)
	}
	u := &models.User{}
	u, err = u.GetByTF(db, tf)

	if err != nil {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		tmpl.Execute(w, d)
		log.Fatalf("u.GetByTF(db, tf)=%v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(pwd))
	if err != nil {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		fmt.Printf("data: %v", d)
		tmpl.Execute(w, d)
		return
		// log.Fatalf("bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(pwd))=%v", err)
	}

	d, tmpl = getLoginData(http.StatusOK, u.User_Name)
	fmt.Printf("d: %v", d)
	tmpl.Execute(w, d)
}

func JSONHandleError(w http.ResponseWriter, err error) {
	var apiErr errhdl.ApiError
	if errors.As(err, &apiErr) {
		status, msg := apiErr.ApiErr()
		JSONError(w, status, msg)
	} else {
		JSONError(w, http.StatusInternalServerError, "internal error")
	}
}
func JSONError(w http.ResponseWriter, s int, errMsg string) error {
	type myErr struct {
		Status int    `json:"status"`
		Msg    string `json:"error_message"`
	}

	o, e := json.Marshal(&myErr{Status: s, Msg: errMsg})
	if e != nil {
		return fmt.Errorf("unable to marshal error object: json.Marshal(&myErr{status: %d, msg: %s})=%w", s, errMsg, e)
	}
	w.Write(o)
	return nil
}

func getLoginData(statusCode int, userName string) (*LoginData, *template.Template) {
	l := true
	if statusCode != http.StatusOK {
		l = false
	}
	d := &LoginData{
		UsrName:  userName,
		LoggedIn: l,
		ErrMsg:   http.StatusText(statusCode),
		ErrCode:  statusCode,
	}
	wd, err := os.Getwd()
	if err != nil {
		d = &LoginData{
			UsrName:  userName,
			LoggedIn: false,
			ErrMsg:   http.StatusText(http.StatusInternalServerError),
			ErrCode:  http.StatusInternalServerError,
		}
	}

	tmpl, err := template.ParseFiles(wd + "/templates/login.html")
	if err != nil {
		d = &LoginData{
			UsrName:  userName,
			LoggedIn: d.LoggedIn,
			ErrMsg:   http.StatusText(http.StatusInternalServerError),
			ErrCode:  http.StatusInternalServerError,
		}
	} else {
		d = &LoginData{
			UsrName:  userName,
			LoggedIn: d.LoggedIn,
			ErrMsg:   http.StatusText(http.StatusOK),
			ErrCode:  http.StatusOK,
		}
	}

	return d, tmpl
}
