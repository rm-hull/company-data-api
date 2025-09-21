package models

import "time"

type CompanyData struct {
	CompanyName                       string     `json:"company_name"`
	CompanyNumber                     string     `json:"company_number"`
	RegAddressCareOf                  string     `json:"reg_address_care_of,omitempty"`
	RegAddressPOBox                   string     `json:"reg_address_po_box,omitempty"`
	RegAddressAddressLine1            string     `json:"reg_address_address_line_1"`
	RegAddressAddressLine2            string     `json:"reg_address_address_line_2,omitempty"`
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
	AccountsNextDueDate               *time.Time `json:"accounts_next_due_date,omitempty"`
	AccountsLastMadeUpDate            *time.Time `json:"accounts_last_made_up_date,omitempty"`
	AccountsAccountCategory           string     `json:"accounts_account_category"`
	ReturnsNextDueDate                *time.Time `json:"returns_next_due_date,omitempty"`
	ReturnsLastMadeUpDate             *time.Time `json:"returns_last_made_up_date,omitempty"`
	MortgagesNumCharges               int        `json:"mortgages_num_charges"`
	MortgagesNumOutstanding           int        `json:"mortgages_num_outstanding"`
	MortgagesNumPartSatisfied         int        `json:"mortgages_num_part_satisfied"`
	MortgagesNumSatisfied             int        `json:"mortgages_num_satisfied"`
	SICCode1                          string     `json:"sic_code_1"`
	SICCode2                          string     `json:"sic_code_2,omitempty"`
	SICCode3                          string     `json:"sic_code_3,omitempty"`
	SICCode4                          string     `json:"sic_code_4,omitempty"`
	LimitedPartnershipsNumGenPartners int        `json:"limited_partnerships_num_gen_partners"`
	LimitedPartnershipsNumLimPartners int        `json:"limited_partnerships_num_lim_partners"`
	URI                               string     `json:"uri"`
	ConfStmtNextDueDate               *time.Time `json:"conf_stmt_next_due_date,omitempty"`
	ConfStmtLastMadeUpDate            *time.Time `json:"conf_stmt_last_made_up_date,omitempty"`
}

type CompanyDataWithLocation struct {
	CompanyData
	Easting  int `json:"easting"`
	Northing int `json:"northing"`
}
