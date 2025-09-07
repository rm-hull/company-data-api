package main

import (
	"company-data-api/cmd"

	"github.com/spf13/cobra"
)

func main() {
	var err error
	var dbPath string
	var port int
	var debug bool
	var companiesHouseZipFile string
	var codepointZipFile string

	rootCmd := &cobra.Command{
		Use:  "company-data",
		Long: `Company Data API & data importers`,
	}

	apiServerCmd := &cobra.Command{
		Use:   "api-server [--db <path>] [--port <port>] [--debug]",
		Short: "Start HTTP API server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.ApiServer(dbPath, port, debug)
		},
	}
	apiServerCmd.Flags().IntVar(&port, "port", 8080, "Port to run HTTP server on")
	apiServerCmd.Flags().BoolVar(&debug, "debug", false, "Enable debugging (pprof) - WARING: do not enable in production")

	processCompaniesHouseZipCmd := &cobra.Command{
		Use:   "import-companies-house [--zip-file <path>] [--db <path>]",
		Short: "Import Companies House ZIP file",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.ImportCompaniesHouseZipFile(companiesHouseZipFile, dbPath)
		},
	}
	processCompaniesHouseZipCmd.Flags().StringVar(&companiesHouseZipFile, "zip-file", "./data/BasicCompanyDataAsOneFile-2025-07-01.zip", "Path to Companies House .zip file")

	processCodepointZipCmd := &cobra.Command{
		Use:   "import-code-point [--zip-file <path>] [--db <path>]",
		Short: "Import Codepoint ZIP file",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.ImportCodepointZipFile(codepointZipFile, dbPath)
		},
	}
	processCodepointZipCmd.Flags().StringVar(&codepointZipFile, "zip-file", "./data/codepo_gb.zip", "Path to Codepoint .zip file")

	rootCmd.AddCommand(apiServerCmd)
	rootCmd.AddCommand(processCompaniesHouseZipCmd)
	rootCmd.AddCommand(processCodepointZipCmd)
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "./data/companies_data.db", "Path to Companies data SQLite database")

	if err = rootCmd.Execute(); err != nil {
		panic(err)
	}
}
