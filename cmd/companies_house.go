package cmd

import (
	"log"

	"github.com/rm-hull/company-data-api/internal"
	"github.com/rm-hull/company-data-api/internal/importer"
	"github.com/rm-hull/godx"
)

func ImportCompaniesHouseZipFile(zipFile string, dbPath string) {

	godx.GitVersion()
	godx.EnvironmentVars()
	godx.UserInfo()

	db, err := internal.Connect(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	err = internal.TransientDownload(zipFile, importer.NewCompanyDataImporter(db).Import)
	if err != nil {
		log.Fatalf("failed to import company data: %v", err)
	}
}
