package internal

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type CodePoint struct {
	PostCode string `json:"post_code"`
	Easting  int    `json:"easting"`
	Northing int    `json:"northing"`
}

func fromCodePointCSV(record []string, headers []string) (CodePoint, error) {

	easting := parseInt(record[2])
	northing := parseInt(record[3])

	return CodePoint{
		PostCode: record[0],
		Easting:  easting,
		Northing: northing,
	}, nil
}

func codePointToTuple(codePoint CodePoint) []any {
	return []any{
		codePoint.PostCode,
		codePoint.Easting,
		codePoint.Northing,
	}
}

func ImportCodePoint(zipPath string, db *sql.DB) error {
	stmt, err := db.Prepare(InsertCodePointSQL)
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
		if f.FileInfo().IsDir() || !strings.HasPrefix(f.Name, "Data/CSV/") {
			continue
		}
		log.Printf("Checking file: %s", f.Name)

		if err := processCodePointCSV(f, stmt); err != nil {
			return fmt.Errorf("failed to process CSV data: %w", err)
		}
	}
	return nil
}

func processCodePointCSV(f *zip.File, stmt *sql.Stmt) error {
	r, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s in zip: %w", f.Name, err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing embedded zip file: %v", err)
		}
	}()

	for result := range parseCSV(r, false, fromCodePointCSV) {
		if result.Error != nil {
			return fmt.Errorf("error parsing line %d: %w", result.LineNum, result.Error)
		}

		_, err := stmt.Exec(codePointToTuple(result.Value)...)
		if err != nil {
			return fmt.Errorf("failed to insert code point for line %d: %w", result.LineNum, err)
		}

		if result.LineNum%379 == 0 {
			log.Printf("Inserted %d records...", result.LineNum)
		}
	}

	return nil
}
