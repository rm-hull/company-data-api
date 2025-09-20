package internal

import (
	"testing"
	"time"

	"github.com/rm-hull/company-data-api/models"
	"github.com/stretchr/testify/assert"
)

func TestFromCompanyDataCSV(t *testing.T) {
	headers := []string{
		"CompanyName", "CompanyNumber", "RegAddress.CareOf", "RegAddress.POBox",
		"RegAddress.AddressLine1", "RegAddress.AddressLine2", "RegAddress.PostTown",
		"RegAddress.County", "RegAddress.Country", "RegAddress.PostCode",
		"CompanyCategory", "CompanyStatus", "CountryOfOrigin", "DissolutionDate",
		"IncorporationDate", "Accounts.AccountRefDay", "Accounts.AccountRefMonth",
		"Accounts.NextDueDate", "Accounts.LastMadeUpDate", "Accounts.AccountCategory",
		"Returns.NextDueDate", "Returns.LastMadeUpDate", "Mortgages.NumCharges",
		"Mortgages.NumOutstanding", "Mortgages.NumPartSatisfied", "Mortgages.NumSatisfied",
		"SICCode.SicText_1", "SICCode.SicText_2", "SICCode.SicText_3", "SICCode.SicText_4",
		"LimitedPartnerships.NumGenPartners", "LimitedPartnerships.NumLimPartners", "URI",
		"ConfStmtNextDueDate", "ConfStmtLastMadeUpDate",
	}

	record := make([]string, 55)
	record[0] = "company"
	record[1] = "123456"
	record[4] = "address1"
	record[5] = "address2"
	record[6] = "posttown"
	record[7] = "county"
	record[8] = "country"
	record[9] = "postcode"
	record[10] = "category"
	record[11] = "status"
	record[12] = "origin"
	record[13] = "01/01/2025"
	record[14] = "01/01/2024"
	record[15] = "1"
	record[16] = "1"
	record[17] = "01/01/2025"
	record[18] = "01/01/2024"
	record[19] = "category"
	record[20] = "01/01/2025"
	record[21] = "01/01/2024"
	record[22] = "1"
	record[23] = "1"
	record[24] = "1"
	record[25] = "1"
	record[26] = "sic1"
	record[27] = "sic2"
	record[28] = "sic3"
	record[29] = "sic4"
	record[30] = "1"
	record[31] = "1"
	record[32] = "uri"
	record[53] = "01/01/2025"
	record[54] = "01/01/2024"

	dissolutionDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	incorporationDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	accountsNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	accountsLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	returnsNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	returnsLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confStmtNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	confStmtLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	expected := models.CompanyData{
		CompanyName:             "company",
		CompanyNumber:           "123456",
		RegAddressAddressLine1:  "address1",
		RegAddressAddressLine2:  "address2",
		RegAddressPostTown:      "posttown",
		RegAddressCounty:        "county",
		RegAddressCountry:       "country",
		RegAddressPostCode:      "postcode",
		CompanyCategory:         "category",
		CompanyStatus:           "status",
		CountryOfOrigin:         "origin",
		DissolutionDate:         &dissolutionDate,
		IncorporationDate:       &incorporationDate,
		AccountsAccountRefDay:   1,
		AccountsAccountRefMonth: 1,
		AccountsNextDueDate:     &accountsNextDueDate,
		AccountsLastMadeUpDate:  &accountsLastMadeUpDate,
		AccountsAccountCategory: "category",
		ReturnsNextDueDate:      &returnsNextDueDate,
		ReturnsLastMadeUpDate:   &returnsLastMadeUpDate,
		MortgagesNumCharges:     1,
		MortgagesNumOutstanding: 1,
		MortgagesNumPartSatisfied: 1,
		MortgagesNumSatisfied:   1,
		SICCode1:                "sic1",
		SICCode2:                "sic2",
		SICCode3:                "sic3",
		SICCode4:                "sic4",
		LimitedPartnershipsNumGenPartners: 1,
		LimitedPartnershipsNumLimPartners: 1,
		URI:                     "uri",
		ConfStmtNextDueDate:     &confStmtNextDueDate,
		ConfStmtLastMadeUpDate:  &confStmtLastMadeUpDate,
	}

	actual, err := fromCompanyDataCSV(record, headers)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFromCompanyDataCSVShortRecord(t *testing.T) {
	headers := []string{"CompanyName", "CompanyNumber"}
	record := []string{"company"}

	_, err := fromCompanyDataCSV(record, headers)

	assert.Error(t, err)
}

func TestFromCompanyDataCSVInvalidDate(t *testing.T) {
	headers := []string{
		"CompanyName", "CompanyNumber", "RegAddress.CareOf", "RegAddress.POBox",
		"RegAddress.AddressLine1", "RegAddress.AddressLine2", "RegAddress.PostTown",
		"RegAddress.County", "RegAddress.Country", "RegAddress.PostCode",
		"CompanyCategory", "CompanyStatus", "CountryOfOrigin", "DissolutionDate",
	}
	record := make([]string, 55)
	record[13] = "invalid-date"

	_, err := fromCompanyDataCSV(record, headers)

	assert.Error(t, err)
}

func TestCompanyDataToTuple(t *testing.T) {
	dissolutionDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	incorporationDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	accountsNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	accountsLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	returnsNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	returnsLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confStmtNextDueDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	confStmtLastMadeUpDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	companyData := models.CompanyData{
		CompanyName:             "company",
		CompanyNumber:           "123456",
		RegAddressAddressLine1:  "address1",
		RegAddressAddressLine2:  "address2",
		RegAddressPostTown:      "posttown",
		RegAddressCounty:        "county",
		RegAddressCountry:       "country",
		RegAddressPostCode:      "postcode",
		CompanyCategory:         "category",
		CompanyStatus:           "status",
		CountryOfOrigin:         "origin",
		DissolutionDate:         &dissolutionDate,
		IncorporationDate:       &incorporationDate,
		AccountsAccountRefDay:   1,
		AccountsAccountRefMonth: 1,
		AccountsNextDueDate:     &accountsNextDueDate,
		AccountsLastMadeUpDate:  &accountsLastMadeUpDate,
		AccountsAccountCategory: "category",
		ReturnsNextDueDate:      &returnsNextDueDate,
		ReturnsLastMadeUpDate:   &returnsLastMadeUpDate,
		MortgagesNumCharges:     1,
		MortgagesNumOutstanding: 1,
		MortgagesNumPartSatisfied: 1,
		MortgagesNumSatisfied:   1,
		SICCode1:                "sic1",
		SICCode2:                "sic2",
		SICCode3:                "sic3",
		SICCode4:                "sic4",
		LimitedPartnershipsNumGenPartners: 1,
		LimitedPartnershipsNumLimPartners: 1,
		URI:                     "uri",
		ConfStmtNextDueDate:     &confStmtNextDueDate,
		ConfStmtLastMadeUpDate:  &confStmtLastMadeUpDate,
	}

	expected := []any{
		"company", "123456", "", "", "address1", "address2", "posttown", "county",
		"country", "postcode", "category", "status", "origin",
		&dissolutionDate,
		&incorporationDate,
		1, 1, &accountsNextDueDate,
		&accountsLastMadeUpDate,
		"category", &returnsNextDueDate,
		&returnsLastMadeUpDate,
		1, 1, 1, 1, "sic1", "sic2", "sic3", "sic4", 1, 1, "uri",
		&confStmtNextDueDate,
		&confStmtLastMadeUpDate,
	}

	actual := companyDataToTuple(companyData)

	assert.Equal(t, expected, actual)
}
