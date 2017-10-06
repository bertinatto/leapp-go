//Package db provides operations on SQLite DB
package db

import (
	"database/sql"
	"os"
	// blank import to register sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const dbDir string = "/var/lib/leapp/"
const dbName string = "leapp.db"

//Database type is sql.DB
type Database sql.DB

//Record type represents DB record description
type Record struct {
	UID, Name, Host, DataType, Data string
}

//EnsureDbPath checks wether the leapp directory exists, if not it tries to create it
//temporary function, will be solved in packaging later
func EnsureDbPath() {
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func (db *Database) sql() *sql.DB {
	return (*sql.DB)(db)
}

//connect represents internal function for DB connection to specific DB (i.e. for test, we use in-memory DB)
func connect(dbPath string) *Database {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	db := (*Database)(conn)
	return db
}

//Connect connects to DB
func Connect() *Database {
	conn := connect(dbDir + dbName)
	return conn
}

//Close closes the connection to the database
func (db *Database) Close() error {
	return db.sql().Close()
}

//CreateTable is temporary function and will be solved in packaging later
func (db *Database) CreateTable() {
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
	if err := db.Exec(sql); err != nil {
		panic(err)
	}
}

//Exec executes given SQL statement
func (db *Database) Exec(sqlStatement string) error {
	_, err := db.sql().Exec(sqlStatement)
	return err
}

//SelectByUIDNameTypeAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByUIDNameTypeAndHost(uid, name, dataType, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, dataType, host)
}

//SelectByUIDNameAndType prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByUIDNameAndType(uid, name, dataType string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ?"
	return db.internalSelect(whereClause, uid, name, dataType)
}

//SelectByNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByNameAndHost(name, host string) ([]Record, error) {
	whereClause := "name = ? AND host = ?"
	return db.internalSelect(whereClause, name, host)
}

//SelectByUIDNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByUIDNameAndHost(uid, name, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, host)
}

//SelectByUIDAndName prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByUIDAndName(uid, name string) ([]Record, error) {
	whereClause := "uid = ? AND name = ?"
	return db.internalSelect(whereClause, uid, name)
}

//SelectByUID prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query
func (db *Database) SelectByUID(uid string) ([]Record, error) {
	whereClause := "uid = ?"
	return db.internalSelect(whereClause, uid)
}

//internalSelect executes SQL SELECT query with given WHERE clause and fill the result(s) to array of the Record type
func (db *Database) internalSelect(whereClause string, params ...interface{}) ([]Record, error) {
	query := "SELECT uid, name, host, data_type, data FROM migrations"
	if len(whereClause) > 0 {
		query += " WHERE " + whereClause
	}
	rows, err := db.sql().Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	selectResult := []Record{}
	for rows.Next() {
		rec := Record{}
		err = rows.Scan(&rec.UID, &rec.Name, &rec.Host, &rec.DataType, &rec.Data)
		if err != nil {
			return nil, err
		}
		selectResult = append(selectResult, rec)
	}
	return selectResult, nil
}

//SelectAll executes SQL SELECT query and returns all records in the migration table
func (db *Database) SelectAll() ([]Record, error) {
	return db.internalSelect("")
}

//Insert executes SQL INSERT query with a given record
func (db *Database) Insert(rec []interface{}) error {
	tx, err := db.sql().Begin()
	sqlstmt := `INSERT INTO migrations(uid, name, host, data_type, data) values (?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(sqlstmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rec...)
	if err != nil {
		return err
	}
	return tx.Commit()
}
