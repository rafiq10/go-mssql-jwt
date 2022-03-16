package mydb

import (
	"app/utils"
	"database/sql"
	"fmt"
	"io/ioutil"

	"crypto/rand"
	"math/big"

	_ "github.com/denisenkom/go-mssqldb"
)

var db_pwd = utils.GetDotEnvVar("DB_PWD")
var db_user = utils.GetDotEnvVar("DB_USER")
var SQL_CONN_STR = "server=0.0.0.0; port=1499;user id=" + db_user + ";password=" + db_pwd + ";"

func GetDb() (db *sql.DB, err error) {
	db, err = sql.Open("mssql", SQL_CONN_STR)
	return db, err
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

	fmt.Errorf("SQL: %s", sql)
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
