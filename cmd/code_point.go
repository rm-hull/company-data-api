package cmd

import (
	"log"

	"github.com/rm-hull/company-data-api/internal"
	"github.com/rm-hull/company-data-api/internal/importer"
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

	err = internal.TransientDownload(zipFile, importer.NewCodePointImporter(db).Import)
	if err != nil {
		log.Fatalf("failed to import code points: %v", err)
	}
}
