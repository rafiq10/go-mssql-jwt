package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	err := godotenv.Load("../.env", ".env")
	if err != nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("GetDotEnvVar(%s) error: %s", key, err.Error())
		}
	}
	return os.Getenv(key)
}
