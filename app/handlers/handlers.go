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
)

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
		Pwd:        r.FormValue("password"),
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
