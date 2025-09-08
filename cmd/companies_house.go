package cmd

import (
	"log"

	"github.com/rm-hull/company-data-api/internal"
)

func ImportCompaniesHouseZipFile(zipFile string, dbPath string) {
	db, err := internal.Connect(dbPath, internal.ReadWrite)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	err = internal.ImportCompanyData(zipFile, db)
	if err != nil {
		log.Fatalf("failed to import company data: %v", err)
	}
}
