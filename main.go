package main

import (
	"companies-house-api/internal"
	"database/sql"
	_ "embed"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := filepath.Join("data", "companies_house.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Printf("connected to database: %s\n", dbPath)

	err = internal.CreateDB(db)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	err = internal.ImportCodePoint("./data/codepo_gb.zip", db)
	if err != nil {
		log.Fatalf("failed to import code points: %v", err)
	}

	// err = internal.ImportCompanyData("./data/BasicCompanyDataAsOneFile-2025-06-01.zip", db)
	// if err != nil {
	// 	log.Fatalf("failed to import company data: %v", err)
	// }
}
