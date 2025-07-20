package cmd

import (
	"company-data-api/internal"
	"database/sql"
	"log"
	"path/filepath"
)

func ImportCompaniesHouseZipFile(path string) {
	dbPath := filepath.Join("data", "companies_data.db")

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

	// err = internal.ImportCodePoint("./data/codepo_gb.zip", db)
	// if err != nil {
	// 	log.Fatalf("failed to import code points: %v", err)
	// }

	err = internal.ImportCompanyData(path, db)
	if err != nil {
		log.Fatalf("failed to import company data: %v", err)
	}
}
