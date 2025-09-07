package repositories

import (
	"company-data-api/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	LEFT = iota
	BOTTOM
	RIGHT
	TOP
)

type SearchRepository interface {
	Find(bbox []float64, processRow func(cd *models.CompanyDataWithLocation)) error
	LastUpdated() *time.Time
}

type SqliteDbRepository struct {
	findStmt    *sql.Stmt
	lastUpdated *time.Time
}

func NewSqliteDbRepository(db *sql.DB) (SearchRepository, error) {
	findStmt, err := prepareStatement(db)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}

	lastUpdated, err := getLastUpdated(db)
	if err != nil {
		return nil, fmt.Errorf("error querying server: %w", err)
	}

	return &SqliteDbRepository{findStmt: findStmt, lastUpdated: lastUpdated}, nil
}

func prepareStatement(db *sql.DB) (*sql.Stmt, error) {
	stmt, err := db.Prepare(`
		SELECT
			cd.company_name, cd.company_number, cd.reg_address_care_of, cd.reg_address_po_box,
			cd.reg_address_address_line_1, cd.reg_address_address_line_2, cd.reg_address_post_town,
			cd.reg_address_county, cd.reg_address_country, cd.reg_address_post_code,
			cd.company_category, cd.company_status, cd.country_of_origin, cd.dissolution_date,
			cd.incorporation_date, cd.accounts_account_ref_day, cd.accounts_account_ref_month,
			cd.accounts_next_due_date, cd.accounts_last_made_up_date, cd.accounts_account_category,
			cd.returns_next_due_date, cd.returns_last_made_up_date, cd.mortgages_num_charges,
			cd.mortgages_num_outstanding, cd.mortgages_num_part_satisfied, cd.mortgages_num_satisfied,
			cd.sic_code_1, cd.sic_code_2, cd.sic_code_3, cd.sic_code_4,
			cd.limited_partnerships_num_gen_partners, cd.limited_partnerships_num_lim_partners,
			cd.uri, cd.conf_stmt_next_due_date, cd.conf_stmt_last_made_up_date,
			cp.easting, cp.northing
		FROM code_point cp
		INNER JOIN company_data cd ON cp.post_code = cd.reg_address_post_code
		WHERE cp.easting BETWEEN ? AND ?
		AND cp.northing BETWEEN ? AND ?
	`)
	if err != nil {
		return nil, fmt.Errorf("error preparing SQL statement: %w", err)
	}
	return stmt, nil
}

func (repo *SqliteDbRepository) Find(bbox []float64, rowProcessor func(companyData *models.CompanyDataWithLocation)) error {

	// In bbox: [LEFT, BOTTOM, RIGHT, TOP]
	rows, err := repo.findStmt.Query(
		bbox[LEFT],   // = min easting
		bbox[RIGHT],  // = max easting
		bbox[BOTTOM], // = min northing
		bbox[TOP],    // = max northing
	)
	if err != nil {
		return fmt.Errorf("error querying database: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("error closing rows: %v\n", err)
		}
	}()

	var cd models.CompanyDataWithLocation

	for rows.Next() {
		if err := rows.Scan(
			&cd.CompanyName,
			&cd.CompanyNumber,
			&cd.RegAddressCareOf,
			&cd.RegAddressPOBox,
			&cd.RegAddressAddressLine1,
			&cd.RegAddressAddressLine2,
			&cd.RegAddressPostTown,
			&cd.RegAddressCounty,
			&cd.RegAddressCountry,
			&cd.RegAddressPostCode,
			&cd.CompanyCategory,
			&cd.CompanyStatus,
			&cd.CountryOfOrigin,
			&cd.DissolutionDate,
			&cd.IncorporationDate,
			&cd.AccountsAccountRefDay,
			&cd.AccountsAccountRefMonth,
			&cd.AccountsNextDueDate,
			&cd.AccountsLastMadeUpDate,
			&cd.AccountsAccountCategory,
			&cd.ReturnsNextDueDate,
			&cd.ReturnsLastMadeUpDate,
			&cd.MortgagesNumCharges,
			&cd.MortgagesNumOutstanding,
			&cd.MortgagesNumPartSatisfied,
			&cd.MortgagesNumSatisfied,
			&cd.SICCode1,
			&cd.SICCode2,
			&cd.SICCode3,
			&cd.SICCode4,
			&cd.LimitedPartnershipsNumGenPartners,
			&cd.LimitedPartnershipsNumLimPartners,
			&cd.URI,
			&cd.ConfStmtNextDueDate,
			&cd.ConfStmtLastMadeUpDate,
			&cd.Easting,
			&cd.Northing,
		); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		rowProcessor(&cd)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %w", err)
	}

	return nil
}

func (repo *SqliteDbRepository) LastUpdated() *time.Time {
	return repo.lastUpdated
}

func getLastUpdated(db *sql.DB) (*time.Time, error) {
	var lastUpdateStr sql.NullString
	row := db.QueryRow(`SELECT MAX(incorporation_date) FROM company_data`)
	if err := row.Scan(&lastUpdateStr); err != nil {
		return nil, fmt.Errorf("failed to determine last update: %w", err)
	}

	if !lastUpdateStr.Valid {
		// Table is empty or all incorporation_date values are NULL.
		return nil, nil
	}

	rfc3339str := strings.Replace(lastUpdateStr.String, " ", "T", 1)
	lastUpdate, err := time.Parse(time.RFC3339, rfc3339str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update string '%s': %w", lastUpdateStr.String, err)
	}

	return &lastUpdate, nil
}
