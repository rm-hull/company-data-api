package internal

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"

	"github.com/rm-hull/company-data-api/models"
)

func fromCompanyDataCSV(record []string, headers []string) (models.CompanyData, error) {
	if len(record) < len(headers) {
		return models.CompanyData{}, fmt.Errorf("record has fewer fields than headers: %d vs %d", len(record), len(headers))
	}

	dissolutionDate, err := parseDate(record[13])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid DissolutionDate: %w", err)
	}
	incorporationDate, err := parseDate(record[14])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid IncorporationDate: %w", err)
	}
	confStmtNextDueDate, err := parseDate(record[53])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid ConfStmtNextDueDate: %w", err)
	}
	confStmtLastMadeUpDate, err := parseDate(record[54])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid ConfStmtLastMadeUpDate: %w", err)
	}

	accountsNextDueDate, err := parseDate(record[17])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid AccountsNextDueDate: %w", err)
	}
	accountsLastMadeUpDate, err := parseDate(record[18])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid AccountsLastMadeUpDate: %w", err)
	}
	returnsNextDueDate, err := parseDate(record[20])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid ReturnsNextDueDate: %w", err)
	}
	returnsLastMadeUpDate, err := parseDate(record[21])
	if err != nil {
		return models.CompanyData{}, fmt.Errorf("invalid ReturnsLastMadeUpDate: %w", err)
	}

	company := models.CompanyData{
		CompanyName:                       record[0],
		CompanyNumber:                     record[1],
		RegAddressCareOf:                  record[2],
		RegAddressPOBox:                   record[3],
		RegAddressAddressLine1:            record[4],
		RegAddressAddressLine2:            record[5],
		RegAddressPostTown:                record[6],
		RegAddressCounty:                  record[7],
		RegAddressCountry:                 record[8],
		RegAddressPostCode:                record[9],
		CompanyCategory:                   record[10],
		CompanyStatus:                     record[11],
		CountryOfOrigin:                   record[12],
		DissolutionDate:                   dissolutionDate,
		IncorporationDate:                 incorporationDate,
		AccountsAccountRefDay:             parseInt(record[15]),
		AccountsAccountRefMonth:           parseInt(record[16]),
		AccountsNextDueDate:               accountsNextDueDate,
		AccountsLastMadeUpDate:            accountsLastMadeUpDate,
		AccountsAccountCategory:           record[19],
		ReturnsNextDueDate:                returnsNextDueDate,
		ReturnsLastMadeUpDate:             returnsLastMadeUpDate,
		MortgagesNumCharges:               parseInt(record[22]),
		MortgagesNumOutstanding:           parseInt(record[23]),
		MortgagesNumPartSatisfied:         parseInt(record[24]),
		MortgagesNumSatisfied:             parseInt(record[25]),
		SICCode1:                          record[26],
		SICCode2:                          record[27],
		SICCode3:                          record[28],
		SICCode4:                          record[29],
		LimitedPartnershipsNumGenPartners: parseInt(record[30]),
		LimitedPartnershipsNumLimPartners: parseInt(record[31]),
		URI:                               record[32],
		ConfStmtNextDueDate:               confStmtNextDueDate,
		ConfStmtLastMadeUpDate:            confStmtLastMadeUpDate,
	}

	return company, nil
}

func companyDataToTuple(companyData models.CompanyData) []any {
	return []any{
		companyData.CompanyName,
		companyData.CompanyNumber,
		companyData.RegAddressCareOf,
		companyData.RegAddressPOBox,
		companyData.RegAddressAddressLine1,
		companyData.RegAddressAddressLine2,
		companyData.RegAddressPostTown,
		companyData.RegAddressCounty,
		companyData.RegAddressCountry,
		companyData.RegAddressPostCode,
		companyData.CompanyCategory,
		companyData.CompanyStatus,
		companyData.CountryOfOrigin,
		companyData.DissolutionDate,
		companyData.IncorporationDate,
		companyData.AccountsAccountRefDay,
		companyData.AccountsAccountRefMonth,
		companyData.AccountsNextDueDate,
		companyData.AccountsLastMadeUpDate,
		companyData.AccountsAccountCategory,
		companyData.ReturnsNextDueDate,
		companyData.ReturnsLastMadeUpDate,
		companyData.MortgagesNumCharges,
		companyData.MortgagesNumOutstanding,
		companyData.MortgagesNumPartSatisfied,
		companyData.MortgagesNumSatisfied,
		companyData.SICCode1,
		companyData.SICCode2,
		companyData.SICCode3,
		companyData.SICCode4,
		companyData.LimitedPartnershipsNumGenPartners,
		companyData.LimitedPartnershipsNumLimPartners,
		companyData.URI,
		companyData.ConfStmtNextDueDate,
		companyData.ConfStmtLastMadeUpDate,
	}
}

func ImportCompanyData(zipPath string, db *sql.DB) error {

	stmt, err := db.Prepare(InsertCompanyDataSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("error closing statement: %v", err)
		}
	}()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing zip file: %v", err)
		}
	}()

	for _, f := range r.File {
		if err := processCompanyDataCSV(f, stmt); err != nil {
			return fmt.Errorf("failed to process CSV data: %w", err)
		}
	}

	return nil
}

func processCompanyDataCSV(f *zip.File, stmt *sql.Stmt) error {
	r, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s in zip: %w", f.Name, err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing embedded zip file: %v", err)
		}
	}()

	for result := range parseCSV(r, true, fromCompanyDataCSV) {

		if result.Error != nil {
			return fmt.Errorf("error parsing line %d: %w", result.LineNum, result.Error)
		}

		companyData := result.Value

		_, err := stmt.Exec(companyDataToTuple(companyData)...)
		if err != nil {
			return fmt.Errorf("failed to insert company data for line %d: %w", result.LineNum, err)
		}

		if result.LineNum%379 == 0 {
			log.Printf("Inserted %d records...", result.LineNum)
		}
	}

	return nil
}
