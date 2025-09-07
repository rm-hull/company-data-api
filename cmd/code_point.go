package cmd

import (
	"company-data-api/internal"
	"log"
)

func ImportCodepointZipFile(zipFile string, dbPath string) {
	db, err := internal.Connect(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	err = internal.ImportCodePoint(zipFile, db)
	if err != nil {
		log.Fatalf("failed to import code points: %v", err)
	}
}
