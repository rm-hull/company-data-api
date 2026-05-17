package cmd

import (
	"log/slog"
	"os"

	"github.com/map-services/company-data-api/internal"
	"github.com/map-services/company-data-api/internal/importer"
	"github.com/map-services/company-data-api/pkg/logger"
	"github.com/rm-hull/godx"
)

func ImportCompaniesHouseZipFile(zipFile string, dbPath string) {
	logger.SetupLogger()

	godx.GitVersion()
	godx.EnvironmentVars()
	godx.UserInfo()

	db, err := internal.Connect(dbPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}()

	err = internal.TransientDownload(zipFile, importer.NewCompanyDataImporter(db).Import)
	if err != nil {
		slog.Error("failed to import company data", "error", err)
		os.Exit(1)
	}
}
