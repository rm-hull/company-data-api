package importer

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/rm-hull/company-data-api/internal"
)

type CodePoint struct {
	PostCode string `json:"post_code"`
	Easting  int    `json:"easting"`
	Northing int    `json:"northing"`
}

func fromCodePointCSV(record []string, headers []string) (*CodePoint, error) {
	easting, err := parseInt(record[2])
	if err != nil {
		return nil, err
	}
	northing, err := parseInt(record[3])
	if err != nil {
		return nil, err
	}

	return &CodePoint{
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

type codePointImporter struct {
	batchSize int
	db        *sql.DB
}

func NewCodePointImporter(db *sql.DB) *codePointImporter {
	return &codePointImporter{
		batchSize: 5000,
		db:        db,
	}
}

func (importer *codePointImporter) Import(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing zip file: %v", err)
		}
	}()

	totalRecordsImported := 0
	for _, f := range r.File {
		if f.FileInfo().IsDir() || !strings.HasPrefix(f.Name, "Data/CSV/") {
			continue
		}
		recordsInFile, err := importer.processCSV(f)
		if err != nil {
			return fmt.Errorf("failed to process CSV data: %w", err)
		}
		log.Printf("Processed file: %s (inserted %d records)", f.Name, recordsInFile)
		totalRecordsImported += recordsInFile
	}
	log.Printf("Total records imported: %d", totalRecordsImported)
	return nil
}

func (importer *codePointImporter) processCSV(f *zip.File) (int, error) {
	r, err := f.Open()
	if err != nil {
		return 0, fmt.Errorf("failed to open embedded file %s in zip: %w", f.Name, err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("error closing embedded zip file: %v", err)
		}
	}()

	batch := make([]CodePoint, 0, importer.batchSize)
	lineNum := 0

	const batchSize = 5000
	for result := range internal.ParseCSV(r, false, fromCodePointCSV) {
		lineNum = result.LineNum
		if result.Error != nil {
			return 0, fmt.Errorf("error parsing line %d: %w", lineNum, result.Error)
		}

		batch = append(batch, *result.Value)

		if len(batch) >= importer.batchSize {
			if err := importer.insertBatch(batch); err != nil {
				return 0, fmt.Errorf("failed to insert batch: %w", err)
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := importer.insertBatch(batch); err != nil {
			return 0, fmt.Errorf("failed to insert batch: %w", err)
		}
	}
	return lineNum, nil
}

func (importer *codePointImporter) insertBatch(batch []CodePoint) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := importer.db.Begin()
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

	stmt, err := tx.Prepare(internal.InsertCodePointSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	for _, codePoint := range batch {
		_, err = stmt.Exec(codePointToTuple(codePoint)...)
		if err != nil {
			return fmt.Errorf("failed to execute individual insert: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
