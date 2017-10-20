package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDB    *Database
	testDBErr error

	tmpFileName = randFileName("leapp_")
)

func randFileName(prefix string) string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes)
}

// createTable is a temporary function and will be solved in packaging later.
func createTable(db *Database) error {
	sql := `
	  CREATE TABLE IF NOT EXISTS migrations(
		uid TEXT,
		name TEXT,
		host TEXT,
		data_type TEXT,
		data TEXT,
		PRIMARY KEY(uid, name, host, data_type)
	);
	`
	_, err := db.conn.Exec(sql)
	return err
}

func setup() {
	testDB, testDBErr = New(os.TempDir(), tmpFileName)
	if err := createTable(testDB); err != nil {
		fmt.Fprintf(os.Stderr, "Could not create table\n")
		os.Exit(1)
	}
}

func teardown() {
	testDB.Close()
	if err := os.Remove(filepath.Join(os.TempDir(), tmpFileName)); err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete database file: %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	teardown()
	os.Exit(c)
}

func TestConnect(t *testing.T) {
	if testDB == nil {
		t.Errorf("Error connection is nil")
	}
	if testDBErr != nil {
		t.Errorf("Connection to the DB failed: %v", testDBErr)
	}
}

func TestExec(t *testing.T) {
	succStmt := "SELECT * from migrations;"
	err := testDB.Exec(succStmt)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}
}

func TestSelects(t *testing.T) {
	rec := map[string]string{
		"UID":      "testUID",
		"name":     "test",
		"host":     "testHost.io",
		"dataType": "testType",
		"data":     "testData"}
	//SelectAll()
	_, err := testDB.SelectAll()
	checkErr(t, "SelectAll", err)

	//SelectByUID
	_, err = testDB.SelectByUID(rec["UID"])
	checkErr(t, "SelectByUID", err)
	//SelectByUIDAndName
	_, err = testDB.SelectByUIDAndName(rec["UID"], rec["name"])
	checkErr(t, "SelectByUIDAndName", err)
	//SelectByUIDNameAndHost
	_, err = testDB.SelectByUIDNameAndHost(rec["UID"], rec["name"], rec["host"])
	checkErr(t, "SelectByUIDNameAndHost", err)
	//SelectByNameAndHost
	_, err = testDB.SelectByNameAndHost(rec["name"], rec["host"])
	checkErr(t, "SelectByNameAndHost", err)
	//SelectByUIDNameAndType
	_, err = testDB.SelectByUIDNameAndType(rec["UID"], rec["name"], rec["dataType"])
	checkErr(t, "SelectByUIDNameAndType", err)
	//SelectByUIDNameTypeAndHost
	_, err = testDB.SelectByUIDNameTypeAndHost(rec["UID"], rec["name"], rec["dataType"], rec["host"])
	checkErr(t, "SelectByUIDNameTypeAndHost", err)
}

func checkErr(t *testing.T, funcName string, err error) {
	if err != nil {
		t.Errorf("%v failed: %v", funcName, err)
	}
}

func TestInsert(t *testing.T) {
	rec := []interface{}{"testUID", "test", "testHost.io", "testType", "testData"}
	err := testDB.Insert(rec)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}
}
