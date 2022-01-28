package mydb

import (
	"app/utils"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

var db_pwd = utils.GetDotEnvVar("DB_PWD")
var db_user = utils.GetDotEnvVar("DB_USER")
var SQL_CONN_STR = "server=0.0.0.0; port=1499;user id=" + db_user + ";password=" + db_pwd + ";"

func GetDb() (db *sql.DB, err error) {
	db, err = sql.Open("mssql", SQL_CONN_STR)
	return db, err
}
