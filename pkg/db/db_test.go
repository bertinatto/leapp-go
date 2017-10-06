package db

import (
	"os"
	"testing"
)

const testDbPath string = ":memory:"

var (
	testConn    *Database
	testConnErr error
)

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	teardown()
	os.Exit(c)
}

func setup() {
	testConn = connect(testDbPath)
	testConn.CreateTable()
}

func teardown() {
	testConn.Close()
}

func TestConnect(t *testing.T) {
	if testConn == nil {
		t.Errorf("Error connection is nil")
	}
	if testConnErr != nil {
		t.Errorf("Connection to the DB failed: %v", testConnErr)
	}
}

func TestExec(t *testing.T) {
	succStmt := "SELECT * from migrations;"
	err := testConn.Exec(succStmt)
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
	_, err := testConn.SelectAll()
	checkErr(t, "SelectAll", err)

	//SelectByUID
	_, err = testConn.SelectByUID(rec["UID"])
	checkErr(t, "SelectByUID", err)
	//SelectByUIDAndName
	_, err = testConn.SelectByUIDAndName(rec["UID"], rec["name"])
	checkErr(t, "SelectByUIDAndName", err)
	//SelectByUIDNameAndHost
	_, err = testConn.SelectByUIDNameAndHost(rec["UID"], rec["name"], rec["host"])
	checkErr(t, "SelectByUIDNameAndHost", err)
	//SelectByNameAndHost
	_, err = testConn.SelectByNameAndHost(rec["name"], rec["host"])
	checkErr(t, "SelectByNameAndHost", err)
	//SelectByUIDNameAndType
	_, err = testConn.SelectByUIDNameAndType(rec["UID"], rec["name"], rec["dataType"])
	checkErr(t, "SelectByUIDNameAndType", err)
	//SelectByUIDNameTypeAndHost
	_, err = testConn.SelectByUIDNameTypeAndHost(rec["UID"], rec["name"], rec["dataType"], rec["host"])
	checkErr(t, "SelectByUIDNameTypeAndHost", err)
}

func checkErr(t *testing.T, funcName string, err error) {
	if err != nil {
		t.Errorf("%v failed: %v", funcName, err)
	}
}

func TestInsert(t *testing.T) {
	rec := []interface{}{"testUID", "test", "testHost.io", "testType", "testData"}
	err := testConn.Insert(rec)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}
}
