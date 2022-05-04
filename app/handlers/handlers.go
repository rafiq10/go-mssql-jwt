package handlers

import (
	"app/errhdl"
	"app/models"

	"app/mydb"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	var usr *models.User
	var l *LoginData

	c, err := r.Cookie("sessionID")
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionID",
			Value: "",
		}
		l = &LoginData{
			LoggedIn: false,
			UsrName:  "",
			ErrMsg:   "Not logged in",
			ErrCode:  http.StatusUnauthorized,
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(wd + "/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	if c.Value != "" {
		s, err := parseToken(c.Value)
		if err != nil {
			log.Printf("Index-->parseToken(%s)=%v", c.Value, err)
		}

		if s != "" {
			db, err := mydb.GetDb()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			u := &models.User{}
			usr, err = u.GetBySessionID(db, s)
			if err != nil {
				JSONHandleError(w, err)
			}

			l = &LoginData{
				LoggedIn: true,
				UsrName:  usr.TF,
				ErrMsg:   "",
				ErrCode:  http.StatusOK,
			}
		}
	}

	err = tmpl.Execute(w, l)
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

	err = usr.UpdateSID(db)
	if err != nil {
		JSONHandleError(w, err)
		log.Fatalf("Register --> u.UpdateSID(db)=%v", err)
		return
	}
	err, token := createToken(u.SessionId)
	if err != nil {
		JSONHandleError(w, err)
		log.Fatalf("Register --> createToken(u.SessionId)=%v", err)
		return
	}

	c := &http.Cookie{
		Name:  "sessionID",
		Value: token,
	}
	http.SetCookie(w, c)

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
		tmpl.Execute(w, d)
		return
		// log.Fatalf("bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(pwd))=%v", err)
	}

	err = u.UpdateSID(db)
	if err != nil {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		tmpl.Execute(w, d)
		log.Fatalf("u.UpdateSID(db)=%v", err)
	}
	// u, err = u.GetByTF(db, tf)
	err, token := createToken(u.SessionId)
	if err != nil {
		d, tmpl = getLoginData(http.StatusUnauthorized, "")
		tmpl.Execute(w, d)
		log.Fatalf("createToken(SessionId)=%v", err)
	}

	c := &http.Cookie{
		Name:  "sessionID",
		Value: token,
	}
	http.SetCookie(w, c)

	d, tmpl = getLoginData(http.StatusOK, u.User_Name)
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

func createToken(sid string) (error, string) {
	id := rand.Intn(10)

	key, err := getKeyById(id)
	if err != nil {
		return err, ""
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(sid))

	// to base64
	signedHMAC := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// signiture | original sessionID
	return nil, signedHMAC + "|" + sid + "|" + strconv.Itoa(id)
}

func parseToken(signedStr string) (string, error) {

	xs := strings.SplitN(signedStr, "|", 3)
	if len(xs) != 3 {
		return "", fmt.Errorf("invalid signed string")
	}
	b64 := xs[0]
	sessionId := xs[1]
	keyId, _ := strconv.Atoi(xs[2])

	key, err := getKeyById(keyId)
	if err != nil {
		return "", fmt.Errorf("parseToen -> getKeyById(%s) = %w", key, err)
	}

	xb, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("parseToen -> DecodeString(b64) = %w", err)
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(sessionId))

	ok := hmac.Equal(xb, h.Sum(nil))
	if !ok {
		return "", fmt.Errorf("Could not parse token")
	}
	return sessionId, nil
}

func getKeyById(id int) (string, error) {
	var k, i string
	db, err := mydb.GetDb()
	if err != nil {
		return "", fmt.Errorf("mydb.GetDb() error in createToken(): %v", err)
	}
	rows, err := db.Query("select key, id from auth.keys offset  " + strconv.Itoa(id) + " limit 1")
	if err != nil {
		return "", fmt.Errorf("db.Query(select key, id from auth.keys offset  %d)=%v", id, err)
	}

	rows.Next()
	rows.Scan(&k, &i)
	if err != nil {
		return "", fmt.Errorf("error scanning the row in db: %v", err)
	} else {

		return k, nil
	}
}
