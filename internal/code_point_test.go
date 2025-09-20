package internal

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// createTestZipCodePoint creates a temporary zip file with a single CSV file for testing
func createTestZipCodePoint(t *testing.T, numRecords int) string {
	t.Helper()
	tempFile, err := os.CreateTemp("", "test-*.zip")
	assert.NoError(t, err)
	defer func() {
		if err := tempFile.Close(); err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}()

	zipWriter := zip.NewWriter(tempFile)
	defer func() {
		assert.NoError(t, zipWriter.Close())
	}()

	f, err := zipWriter.Create("Data/CSV/test.csv")
	assert.NoError(t, err)

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	for i := range numRecords {
		record := make([]string, 4)
		record[0] = fmt.Sprintf("AB12 3CD%d", i) // Postcode
		record[1] = "1"                          // Quality
		record[2] = fmt.Sprintf("%d", 300000+i)  // Easting
		record[3] = fmt.Sprintf("%d", 700000+i)  // Northing
		// Fill other fields with dummy data if necessary, or leave empty
		assert.NoError(t, csvWriter.Write(record))
	}

	return tempFile.Name()
}

func TestFromCodePointCSV(t *testing.T) {
	headers := []string{"Postcode", "Quality", "Easting", "Northing"}
	record := []string{"AB12 3CD", "1", "300000", "700000"}

	expected := &CodePoint{
		PostCode: "AB12 3CD",
		Easting:  300000,
		Northing: 700000,
	}

	actual, err := fromCodePointCSV(record, headers)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFromCodePointCSVInvalidEasting(t *testing.T) {
	headers := []string{"Postcode", "Quality", "Easting", "Northing"}
	record := []string{"AB12 3CD", "1", "invalid", "700000"}

	_, err := fromCodePointCSV(record, headers)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strconv.Atoi: parsing \"invalid\": invalid syntax")
}

func TestFromCodePointCSVInvalidNorthing(t *testing.T) {
	headers := []string{"Postcode", "Quality", "Easting", "Northing"}
	record := []string{"AB12 3CD", "1", "300000", "invalid"}

	_, err := fromCodePointCSV(record, headers)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strconv.Atoi: parsing \"invalid\": invalid syntax")
}

func TestCodePointToTuple(t *testing.T) {
	codePoint := CodePoint{
		PostCode: "AB12 3CD",
		Easting:  300000,
		Northing: 700000,
	}

	expected := []any{
		"AB12 3CD",
		300000,
		700000,
	}

	actual := codePointToTuple(codePoint)

	assert.Equal(t, expected, actual)
}

func TestImportCodePoint(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	zipPath := createTestZipCodePoint(t, 1) // Create a zip with 1 record
	defer func() {
		assert.NoError(t, os.Remove(zipPath))
	}()

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr) // Restore default output
	}()

	mock.ExpectBegin()
	mock.ExpectPrepare(InsertCodePointSQL)
	mock.ExpectExec(InsertCodePointSQL).
		WithArgs("AB12 3CD0", 300000, 700000).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = ImportCodePoint(zipPath, db)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Contains(t, buf.String(), "Total records imported: 1")
}

func TestImportCodePointMultipleRecords(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	numRecords := 3
	zipPath := createTestZipCodePoint(t, numRecords)
	defer func() {
		assert.NoError(t, os.Remove(zipPath))
	}()

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr) // Restore default output
	}()

	mock.ExpectBegin()
	mock.ExpectPrepare(InsertCodePointSQL)
	for i := 0; i < numRecords; i++ {
		mock.ExpectExec(InsertCodePointSQL).
			WithArgs(fmt.Sprintf("AB12 3CD%d", i), 300000+i, 700000+i).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	err = ImportCodePoint(zipPath, db)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Contains(t, buf.String(), fmt.Sprintf("Total records imported: %d", numRecords))
}

func TestImportCodePointPrepareError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	zipPath := createTestZipCodePoint(t, 1)
	defer func() {
		assert.NoError(t, os.Remove(zipPath))
	}()

	mock.ExpectBegin()
	mock.ExpectPrepare(InsertCodePointSQL).WillReturnError(fmt.Errorf("prepare error"))
	mock.ExpectRollback()

	err = ImportCodePoint(zipPath, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to prepare statement: prepare error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestImportCodePointExecError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	zipPath := createTestZipCodePoint(t, 1)
	defer func() {
		assert.NoError(t, os.Remove(zipPath))
	}()

	mock.ExpectBegin()
	mock.ExpectPrepare(InsertCodePointSQL)
	mock.ExpectExec(InsertCodePointSQL).
		WithArgs("AB12 3CD0", 300000, 700000).
		WillReturnError(fmt.Errorf("exec error"))
	mock.ExpectRollback()

	err = ImportCodePoint(zipPath, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute individual insert: exec error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProcessCodePointCSV(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	zipPath := createTestZipCodePoint(t, 1) // Create a zip with 1 record
	defer func() {
		assert.NoError(t, os.Remove(zipPath))
	}()

	r, err := zip.OpenReader(zipPath)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, r.Close())
	}()

	mock.ExpectBegin()
	mock.ExpectPrepare(InsertCodePointSQL)
	mock.ExpectExec(InsertCodePointSQL).
		WithArgs("AB12 3CD0", 300000, 700000).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Find the CSV file within the zip
	var csvFile *zip.File
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && strings.HasSuffix(f.Name, ".csv") {
			csvFile = f
			break
		}
	}
	assert.NotNil(t, csvFile, "CSV file not found in zip")

	numRecords, err := processCodePointCSV(csvFile, db)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 1, numRecords)
}
