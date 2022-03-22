package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token has expired")
	}
	if u.VerifyAudience("tis-gf", true) {
		return fmt.Errorf("Invalid token audience")
	}
	if u.VerifyIssuer("tis-gf-api", true) {
		return fmt.Errorf("Invalid token issuer")
	}
	if u.SessionID == 0 {
		return fmt.Errorf("Invalid session ID")
	}
	return nil
}

func (u *user) ValidatePasswordHash(pwdhash string) bool {
	return u.pwd == pwdhash
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
