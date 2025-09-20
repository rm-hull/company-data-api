package internal

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:embed sql/migration.sql
var migrationSQL string

//go:embed sql/insert_code_point.sql
var InsertCodePointSQL string

//go:embed sql/insert_company_data.sql
var InsertCompanyDataSQL string

type Mode int

const (
	ReadOnly Mode = iota
	ReadWrite
)

func CreateDB(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	return err
}

func Connect(dbPath string, mode Mode) (*sql.DB, error) {
	dsn := dbPath + "?_journal_mode=WAL&_busy_timeout=5000"
	if mode == ReadOnly {
		dsn = dsn + "&mode=ro"
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Printf("connected to database: %s", dsn)

	if mode == ReadWrite {
		err = CreateDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}
	return db, nil
}
