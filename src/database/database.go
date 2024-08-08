package database

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "data.db"
)

// initDB initializes the database connection
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// setupSchema creates tables if they do not exist
func SetupSchema(db *sql.DB) error {
	schemaFile, err := os.Open("schema.sql")
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	schema, err := io.ReadAll(schemaFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
