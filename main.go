package main

import (
	"archive/zip"
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

	zipPath := filepath.Join("data", "BasicCompanyDataAsOneFile-2025-06-01.zip")
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("Failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		defer rc.Close()
		if err != nil {
			log.Fatalf("Failed to open file in zip: %v", err)
		}

		for result := range internal.ParseCSV(rc, internal.FromCSV) {

			if result.Error != nil {
				log.Fatalf("Error parsing line %d: %v", result.LineNum, result.Error)
			}

			companyData := result.Value

			err := internal.InsertCompanyData(db, companyData)
			if err != nil {
				log.Fatalf("failed to insert company data for line %d: %v", result.LineNum, err)
			}

			if result.LineNum%379 == 0 {
				log.Printf("Inserted %d records...", result.LineNum)
			}
		}
	}
}
