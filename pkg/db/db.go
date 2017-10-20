// Package db provides operations on SQLite DB.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	// Registers sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Database is a database for storing information
type Database struct {
	conn *sql.DB
}

// Record type represents Database record.
type Record struct {
	UID      string
	Name     string
	Host     string
	DataType string
	Data     string
}

// New initializes and return a new Database.
func New(dbDir, dbName string) (*Database, error) {
	// Creates the dir, if not exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database dir: %v", err)
	}

	// Open the database
	c, err := sql.Open("sqlite3", filepath.Join(dbDir, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &Database{conn: c}, nil
}

// Exec executes given SQL statement.
func (db *Database) Exec(sqlStatement string) error {
	_, err := db.conn.Exec(sqlStatement)
	return err
}

// Close closes a Database connections.
func (db *Database) Close() error {
	return db.conn.Close()
}

// SelectByUIDNameTypeAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameTypeAndHost(uid, name, dataType, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, dataType, host)
}

// SelectByUIDNameAndType prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameAndType(uid, name, dataType string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND data_type = ?"
	return db.internalSelect(whereClause, uid, name, dataType)
}

// SelectByNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByNameAndHost(name, host string) ([]Record, error) {
	whereClause := "name = ? AND host = ?"
	return db.internalSelect(whereClause, name, host)
}

// SelectByUIDNameAndHost prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDNameAndHost(uid, name, host string) ([]Record, error) {
	whereClause := "uid = ? AND name = ? AND host = ?"
	return db.internalSelect(whereClause, uid, name, host)
}

// SelectByUIDAndName prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUIDAndName(uid, name string) ([]Record, error) {
	whereClause := "uid = ? AND name = ?"
	return db.internalSelect(whereClause, uid, name)
}

// SelectByUID prepares WHERE clause for SELECT query and returns array of the Record(s) found by the SQL query.
func (db *Database) SelectByUID(uid string) ([]Record, error) {
	whereClause := "uid = ?"
	return db.internalSelect(whereClause, uid)
}

// internalSelect executes SQL SELECT query with given WHERE clause and fill the result(s) to array of the Record type.
func (db *Database) internalSelect(whereClause string, params ...interface{}) ([]Record, error) {
	query := "SELECT uid, name, host, data_type, data FROM migrations"
	if len(whereClause) > 0 {
		query += " WHERE " + whereClause
	}
	rows, err := db.conn.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	selectResult := []Record{}
	for rows.Next() {
		rec := Record{}
		err = rows.Scan(&rec.UID, &rec.Name, &rec.Host, &rec.DataType, &rec.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		selectResult = append(selectResult, rec)
	}
	return selectResult, nil
}

// SelectAll executes SQL SELECT query and returns all records in the migration table.
func (db *Database) SelectAll() ([]Record, error) {
	return db.internalSelect("")
}

// Insert executes SQL INSERT query with a given record.
func (db *Database) Insert(rec []interface{}) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("could not start sql transaction: %v", err)
	}

	sqlstmt := `INSERT INTO migrations(uid, name, host, data_type, data) values (?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(sqlstmt)
	if err != nil {
		return fmt.Errorf("failed to prepare sql: %v", err)
	}

	_, err = stmt.Exec(rec...)
	if err != nil {
		return fmt.Errorf("failed to execute sql: %v", err)
	}
	return tx.Commit()
}
