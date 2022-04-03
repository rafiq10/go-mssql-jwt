package mydb

import (
	"testing"
	// "os"
	// "database/sql"
)

func TestGetDb(t *testing.T) {
	t.Run("testing db connection", func(t *testing.T) {
		db, err := GetDb()
		err = db.Ping()
		if err != nil {
			t.Errorf("Not Ping: %v", err.Error())
		}
		defer db.Close()
	})

}

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