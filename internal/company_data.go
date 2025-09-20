package internal

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rm-hull/company-data-api/models"
)

func fromCompanyDataCSV(record []string, headers []string) (*models.CompanyData, error) {
	if len(record) < len(headers) {
		return nil, fmt.Errorf("record has fewer fields than headers: %d vs %d", len(record), len(headers))
	}

	var err error
	parseDateField := func(index int, fieldName string) *time.Time {
		if err != nil {
			return nil
		}
		var date *time.Time
		date, err = parseDate(record[index])
		if err != nil {
			err = fmt.Errorf("invalid %s: %w", fieldName, err)
		}
		return date
	}

	parseIntField := func(index int, fieldName string) int {
		if err != nil {
			return 0
		}
		var val int
		val, err = parseInt(record[index])
		if err != nil {
			err = fmt.Errorf("invalid %s: %w", fieldName, err)
		}
		return val
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
		DissolutionDate:                   parseDateField(13, "DissolutionDate"),
		IncorporationDate:                 parseDateField(14, "IncorporationDate"),
		AccountsAccountRefDay:             parseIntField(15, "AccountsAccountRefDay"),
		AccountsAccountRefMonth:           parseIntField(16, "AccountsAccountRefMonth"),
		AccountsNextDueDate:               parseDateField(17, "AccountsNextDueDate"),
		AccountsLastMadeUpDate:            parseDateField(18, "AccountsLastMadeUpDate"),
		AccountsAccountCategory:           record[19],
		ReturnsNextDueDate:                parseDateField(20, "ReturnsNextDueDate"),
		ReturnsLastMadeUpDate:             parseDateField(21, "ReturnsLastMadeUpDate"),
		MortgagesNumCharges:               parseIntField(22, "MortgagesNumCharges"),
		MortgagesNumOutstanding:           parseIntField(23, "MortgagesNumOutstanding"),
		MortgagesNumPartSatisfied:         parseIntField(24, "MortgagesNumPartSatisfied"),
		MortgagesNumSatisfied:             parseIntField(25, "MortgagesNumSatisfied"),
		SICCode1:                          record[26],
		SICCode2:                          record[27],
		SICCode3:                          record[28],
		SICCode4:                          record[29],
		LimitedPartnershipsNumGenPartners: parseIntField(30, "LimitedPartnershipsNumGenPartners"),
		LimitedPartnershipsNumLimPartners: parseIntField(31, "LimitedPartnershipsNumLimPartners"),
		URI:                               record[32],
		ConfStmtNextDueDate:               parseDateField(53, "ConfStmtNextDueDate"),
		ConfStmtLastMadeUpDate:            parseDateField(54, "ConfStmtLastMadeUpDate"),
	}

	if err != nil {
		return nil, err
	}

	return &company, nil
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

const batchSize = 5000

func ImportCompanyData(zipPath string, db *sql.DB) error {
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
		if err := processCompanyDataCSV(f, db); err != nil {
			return fmt.Errorf("failed to process CSV data: %w", err)
		}
	}

	return nil
}

func processCompanyDataCSV(f *zip.File, db *sql.DB) error {
	r, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s in zip: %w", f.Name, err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing embedded zip file: %v", err)
		}
	}()

	var (
		batch   []models.CompanyData
		lineNum int
	)

	for result := range parseCSV(r, true, fromCompanyDataCSV) {
		lineNum = result.LineNum
		if result.Error != nil {
			return fmt.Errorf("error parsing line %d: %w", lineNum, result.Error)
		}

		batch = append(batch, *result.Value)

		if len(batch) >= batchSize {
			if err := insertCompanyDataBatch(db, batch, lineNum); err != nil {
				return fmt.Errorf("failed to insert company data batch at line %d: %w", lineNum, err)
			}
			batch = batch[:0] // Clear the buffer, retaining capacity
		}
	}

	// Insert any remaining records in the buffer
	if len(batch) > 0 {
		if err := insertCompanyDataBatch(db, batch, lineNum); err != nil {
			return fmt.Errorf("failed to insert final company data batch at line %d: %w", lineNum, err)
		}
	}

	return nil
}

func insertCompanyDataBatch(db *sql.DB, batch []models.CompanyData, lastLineNum int) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("error rolling back transaction: %v", rbErr)
			}
		}
	}()

	stmt, err := tx.Prepare(InsertCompanyDataSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	for _, companyData := range batch {
		_, err = stmt.Exec(companyDataToTuple(companyData)...)
		if err != nil {
			return fmt.Errorf("failed to execute individual insert: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Inserted %d records...", lastLineNum)
	return nil
}
