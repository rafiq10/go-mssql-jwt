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
