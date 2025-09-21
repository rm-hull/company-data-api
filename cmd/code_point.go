package cmd

import (
	"log"

	"github.com/rm-hull/company-data-api/internal"
	"github.com/rm-hull/company-data-api/internal/importer"
)

func ImportCodepointZipFile(zipFile string, dbPath string) {
	db, err := internal.Connect(dbPath, internal.ReadWrite)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	err = importer.NewCodePointImporter(db).Import(zipFile)
	if err != nil {
		log.Fatalf("failed to import code points: %v", err)
	}
}
