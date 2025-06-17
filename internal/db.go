package internal

import (
	"database/sql"
	_ "embed"
)

//go:embed migration.sql
var migrationSQL string

func CreateDB(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	return err
}

func InsertCompanyData(db *sql.DB, companyData CompanyData) error {
	_, err := db.Exec(`
		INSERT OR REPLACE INTO company_data (
			company_name,
			company_number,
			reg_address_care_of,
			reg_address_po_box,
			reg_address_address_line_1,
			reg_address_address_line_2,
			reg_address_post_town,
			reg_address_county,
			reg_address_country,
			reg_address_post_code,
			company_category,
			company_status,
			country_of_origin,
			dissolution_date,
			incorporation_date,
			accounts_account_ref_day,
			accounts_account_ref_month,
			accounts_next_due_date,
			accounts_last_made_up_date,
			accounts_account_category,
			returns_next_due_date,
			returns_last_made_up_date,
			mortgages_num_charges,
			mortgages_num_outstanding,
			mortgages_num_part_satisfied,
			mortgages_num_satisfied,
			sic_code_1,
			sic_code_2,
			sic_code_3,
			sic_code_4,
			limited_partnerships_num_gen_partners,
			limited_partnerships_num_lim_partners,
			uri,
			conf_stmt_next_due_date,
			conf_stmt_last_made_up_date
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
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
	)
	return err
}
