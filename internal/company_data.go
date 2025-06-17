package internal

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type CompanyData struct {
	CompanyName                       string     `json:"company_name"`
	CompanyNumber                     string     `json:"company_number"`
	RegAddressCareOf                  string     `json:"reg_address_care_of,omitempty"`
	RegAddressPOBox                   string     `json:"reg_address_po_box,omitempty"`
	RegAddressAddressLine1            string     `json:"reg_address_address_line1"`
	RegAddressAddressLine2            string     `json:"reg_address_address_line2,omitempty"`
	RegAddressPostTown                string     `json:"reg_address_post_town"`
	RegAddressCounty                  string     `json:"reg_address_county,omitempty"`
	RegAddressCountry                 string     `json:"reg_address_country"`
	RegAddressPostCode                string     `json:"reg_address_post_code"`
	CompanyCategory                   string     `json:"company_category"`
	CompanyStatus                     string     `json:"company_status"`
	CountryOfOrigin                   string     `json:"country_of_origin"`
	DissolutionDate                   *time.Time `json:"dissolution_date,omitempty"`
	IncorporationDate                 *time.Time `json:"incorporation_date"`
	AccountsAccountRefDay             int        `json:"accounts_account_ref_day"`
	AccountsAccountRefMonth           int        `json:"accounts_account_ref_month"`
	AccountsNextDueDate               *time.Time `json:"accounts_next_due_date"`
	AccountsLastMadeUpDate            *time.Time `json:"accounts_last_made_up_date"`
	AccountsAccountCategory           string     `json:"accounts_account_category"`
	ReturnsNextDueDate                *time.Time `json:"returns_next_due_date"`
	ReturnsLastMadeUpDate             *time.Time `json:"returns_last_made_up_date"`
	MortgagesNumMortCharges           int        `json:"mortgages_num_mort_charges"`
	MortgagesNumMortOutstanding       int        `json:"mortgages_num_mort_outstanding"`
	MortgagesNumMortPartSatisfied     int        `json:"mortgages_num_mort_part_satisfied"`
	MortgagesNumMortSatisfied         int        `json:"mortgages_num_mort_satisfied"`
	SICCodeSicText_1                  string     `json:"sic_code_sic_text_1"`
	SICCodeSicText_2                  string     `json:"sic_code_sic_text_2,omitempty"`
	SICCodeSicText_3                  string     `json:"sic_code_sic_text_3,omitempty"`
	SICCodeSicText_4                  string     `json:"sic_code_sic_text_4,omitempty"`
	LimitedPartnershipsNumGenPartners int        `json:"limited_partnerships_num_gen_partners"`
	LimitedPartnershipsNumLimPartners int        `json:"limited_partnerships_num_lim_partners"`
	URI                               string     `json:"uri"`
	ConfStmtNextDueDate               *time.Time `json:"conf_stmt_next_due_date,omitempty"`
	ConfStmtLastMadeUpDate            *time.Time `json:"conf_stmt_last_made_up_date,omitempty"`
}

func fromCompanyDataCSV(record []string, headers []string) (CompanyData, error) {
	if len(record) < len(headers) {
		return CompanyData{}, fmt.Errorf("record has fewer fields than headers: %d vs %d", len(record), len(headers))
	}

	dissolutionDate, err := parseDate(record[13])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid DissolutionDate: %w", err)
	}
	incorporationDate, err := parseDate(record[14])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid IncorporationDate: %w", err)
	}
	confStmtNextDueDate, err := parseDate(record[53])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid ConfStmtNextDueDate: %w", err)
	}
	confStmtLastMadeUpDate, err := parseDate(record[54])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid ConfStmtLastMadeUpDate: %w", err)
	}

	accountsNextDueDate, err := parseDate(record[17])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid AccountsNextDueDate: %w", err)
	}
	accountsLastMadeUpDate, err := parseDate(record[18])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid AccountsLastMadeUpDate: %w", err)
	}
	returnsNextDueDate, err := parseDate(record[20])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid ReturnsNextDueDate: %w", err)
	}
	returnsLastMadeUpDate, err := parseDate(record[21])
	if err != nil {
		return CompanyData{}, fmt.Errorf("invalid ReturnsLastMadeUpDate: %w", err)
	}

	company := CompanyData{
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
		MortgagesNumMortCharges:           parseInt(record[22]),
		MortgagesNumMortOutstanding:       parseInt(record[23]),
		MortgagesNumMortPartSatisfied:     parseInt(record[24]),
		MortgagesNumMortSatisfied:         parseInt(record[25]),
		SICCodeSicText_1:                  record[26],
		SICCodeSicText_2:                  record[27],
		SICCodeSicText_3:                  record[28],
		SICCodeSicText_4:                  record[29],
		LimitedPartnershipsNumGenPartners: parseInt(record[30]),
		LimitedPartnershipsNumLimPartners: parseInt(record[31]),
		URI:                               record[32],
		ConfStmtNextDueDate:               confStmtNextDueDate,
		ConfStmtLastMadeUpDate:            confStmtLastMadeUpDate,
	}

	return company, nil
}

func companyDataToTuple(companyData CompanyData) []any {
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
		companyData.MortgagesNumMortCharges,
		companyData.MortgagesNumMortOutstanding,
		companyData.MortgagesNumMortPartSatisfied,
		companyData.MortgagesNumMortSatisfied,
		companyData.SICCodeSicText_1,
		companyData.SICCodeSicText_2,
		companyData.SICCodeSicText_3,
		companyData.SICCodeSicText_4,
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
	defer stmt.Close()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if err := processCompanyDataCSV(f, stmt); err != nil {
			return fmt.Errorf("failed to process CSV data: %w", err)
		}
	}

	return nil
}

func processCompanyDataCSV(f *zip.File, stmt *sql.Stmt) error {
	rc, err := f.Open()
	defer rc.Close()
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s in zip: %w", f.Name, err)
	}

	for result := range parseCSV(rc, true, fromCompanyDataCSV) {

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
