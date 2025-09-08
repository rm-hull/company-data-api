package internal

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:embed migration.sql
var migrationSQL string

type Mode int

const (
	ReadOnly Mode = iota
	ReadWrite
)

func CreateDB(db *sql.DB) error {
	_, err := db.Exec(migrationSQL)
	return err
}

const InsertCompanyDataSQL = `
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
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

const InsertCodePointSQL = `
		INSERT OR REPLACE INTO code_point (
			post_code,
			easting,
			northing
		) VALUES (?,?,?)`

func Connect(dbPath string, mode Mode) (*sql.DB, error) {
	dsn := dbPath
	if mode == ReadOnly {
		dsn = dsn + "?mode=ro"
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Printf("connected to database: %s", dsn)

	if mode == ReadWrite {
		err = CreateDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}
	return db, nil
}
