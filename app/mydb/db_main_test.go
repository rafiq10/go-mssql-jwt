package mydb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	path, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Not able to get current path: os.Getwd()=%v", err))
	}

	path = filepath.Join(path, "init_db.sql")
	err = RunSqlFromFile(path)
	if err != nil {
		panic(fmt.Errorf("Not able to RunSqlFromFile(%s)=%v", path, err))
	}
	exitVal := m.Run()
	// teardown()
	os.Exit(exitVal)
}

func teardown() {
	path, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("teardown() error:Not able to get current path: os.Getwd()=%v", err))
	}
	path = filepath.Join(path, "clear_db.sql")
	err = RunSqlFromFile(path)
	if err != nil {
		panic(fmt.Errorf("teardown() error: Not able to RunSqlFromFile(%s)=%v", path, err))
	}
}
