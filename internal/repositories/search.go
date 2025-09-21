package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rm-hull/company-data-api/internal"
	"github.com/rm-hull/company-data-api/internal/models"
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
	stmt, err := db.Prepare(internal.SearchSQL)
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
