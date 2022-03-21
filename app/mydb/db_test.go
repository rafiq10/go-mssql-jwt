package mydb

import (
	"database/sql"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestGetDb(t *testing.T) {
	t.Run("testing db connection", func(t *testing.T) {
		db, err := sql.Open("mssql", SQL_CONN_STR)
		if err != nil {
			t.Errorf("Not able to connect to the dataase")
		}

		err = db.Ping()
		if err != nil {
			t.Errorf("Not Ping: %v", err.Error())
		}
		defer db.Close()
	})

}

// func TestRunSqlFromFile(t *testing.T) {
// 	t.Run("testing db initialization", func(t *testing.T) {
// 		path, err := os.Getwd()
// 		if err != nil {
// 			t.Errorf("Not able to get current path: os.Getwd()=%v", err)
// 		}

// 		path = filepath.Join(path, "init.sql")
// 		err = RunSqlFromFile(path)
// 		if err != nil {
// 			t.Errorf("Not able execute sql from file: RunSqlFromFile(%s)=%v", path, err)
// 		}
// 	})
// }

func TestRunSQL(t *testing.T) {
	err := GenerateKeys()
	if err != nil {
		t.Errorf("Error generating random key: GenerateKeys(15)=%v", err)
	}
	// err = RunSQL("use [users-db] insert into Keys values ('" + key + "') ")
	// if err != nil {
	// 	t.Errorf("Error generating random key: GenerateRandomString(15)=%v", err)
	// }
}
