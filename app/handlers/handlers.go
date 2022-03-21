package handlers

import "net/http"

type user struct {
	tf         string
	userName   string
	email      string
	salt       string
	pwd        string
	createDate string
	role       int
	department string
}

func (u *user) ValidatePasswordHash(pwdhash string) bool {
	return u.pwd == pwdhash
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
