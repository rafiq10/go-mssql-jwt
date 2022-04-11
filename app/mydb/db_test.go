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
}

// insert into auth.users(tf, user_name,email,pwd, created_at,usr_role,department) values
// ('TF05069','Edu MS','edums@gmail.com','$2a$10$ybMx2eDHAOUqi65lQEyBSeeBKdQHUDdvLdn64600S.5ax26VVcJKu',current_timestamp,'tis-gf-oper','GF')
// on conflict (tf) do nothing;
