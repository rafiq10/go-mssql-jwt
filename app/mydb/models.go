package mydb

import (
	"time"
	jwt "github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	jwt.StandardClaims
	sessionID int64
}    

func (u *UserClaims) Valid error {
	u.VerifyExpiresAt(time.Now().Unix())
}