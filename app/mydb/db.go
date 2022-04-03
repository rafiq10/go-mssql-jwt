package mydb

import (
	"app/utils"
	"database/sql"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db_pwd = utils.GetDotEnvVar("DB_PWD")
var db_user = utils.GetDotEnvVar("DB_USER")
var db_uri = utils.GetDotEnvVar("DB_URI")



func GetDb() (*sql.DB, error) {
    uri := os.Getenv("DATABASE_URI")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
    return db,err
}

func RunSqlFromFile(path string) error {
	db, err := GetDb()
	if err != nil {
		return fmt.Errorf("RunSqlFromFile(%s)=\n %v\nError trying to GetDb()", path, err)
	}

	defer db.Close()
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("RunSqlFromFile(%s)=\n %v\nError trying to ioutil.ReadFile(%s)", path, err, path)
	}
	sql := string(c)
	_, err = db.Exec(sql)

	
	if err != nil {
		return fmt.Errorf("RunSqlFromFile(%s)=\n %v\nError trying to db.Exec(%s)", path, err, sql)
	}

	return nil
}

func RunSQL(sql string) error {
	db, err := GetDb()

	if err != nil {
		return fmt.Errorf("RunSQL(%s)=\n %v\nError trying to GetDb()", sql, err)
	}
	defer db.Close()
	_, err = db.Exec(sql)

	if err != nil {
		return fmt.Errorf("RunSQL(%s)=\n %v\nError trying to db.Exec(%s)", sql, err, sql)
	}
	return nil
}

func GenerateKeys() error {
	err := RunSQL("delete from auth.keys")
	if err != nil {
		return fmt.Errorf("Error generating random key: GenerateRandomString(15)=%v", err)
	}
	for i := 0; i < 10; i++ {
		key, err := utils.GenerateRandomString(15)
		if err != nil {
			return fmt.Errorf("Error generating random key: GenerateRandomString(15)=%v", err)
		}
		keyId := uuid.New().String()
		err = RunSQL("insert into auth.keys (id,key) values ('" + keyId + "','" + key + "') ")
		if err != nil {
			return fmt.Errorf("Error generating random key: GenerateRandomString(15)=%v", err)
		}
	}
	return nil

}
