package utils

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("GetDotEnvVar(%s) error: %s", key, err.Error())
		}
	}
	return os.Getenv(key)
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
