package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type testData struct {
	Name string
	Age  int
}

func fromFunc(data []string, headers []string) (testData, error) {
	if len(data) != 2 {
		return testData{}, fmt.Errorf("expected 2 fields, got %d", len(data))
	}
	age, err := strconv.Atoi(data[1])
	if err != nil {
		return testData{}, fmt.Errorf("failed to parse age: %w", err)
	}
	return testData{
		Name: data[0],
		Age:  age,
	}, nil
}

func TestParseCSVWithHeader(t *testing.T) {
	csvData := `name,age
"John Doe",30
"Jane Doe",25
`
	reader := strings.NewReader(csvData)
	results := parseCSV(reader, true, fromFunc)

	expected := []Result[testData]{
		{Value: testData{Name: "John Doe", Age: 30}, LineNum: 1},
		{Value: testData{Name: "Jane Doe", Age: 25}, LineNum: 2},
	}

	i := 0
	for result := range results {
		if result.Error != nil {
			t.Errorf("unexpected error: %v", result.Error)
		}
		if result.Value != expected[i].Value {
			t.Errorf("expected value %v, got %v", expected[i].Value, result.Value)
		}
		if result.LineNum != expected[i].LineNum {
			t.Errorf("expected line number %d, got %d", expected[i].LineNum, result.LineNum)
		}
		i++
	}
}

func TestParseCSVWithoutHeader(t *testing.T) {
	csvData := `"John Doe",30
"Jane Doe",25
`
	reader := strings.NewReader(csvData)
	results := parseCSV(reader, false, fromFunc)

	expected := []Result[testData]{
		{Value: testData{Name: "John Doe", Age: 30}, LineNum: 1},
		{Value: testData{Name: "Jane Doe", Age: 25}, LineNum: 2},
	}

	i := 0
	for result := range results {
		if result.Error != nil {
			t.Errorf("unexpected error: %v", result.Error)
		}
		if result.Value != expected[i].Value {
			t.Errorf("expected value %v, got %v", expected[i].Value, result.Value)
		}
		if result.LineNum != expected[i].LineNum {
			t.Errorf("expected line number %d, got %d", expected[i].LineNum, result.LineNum)
		}
		i++
	}
}

func TestParseCSVMalformed(t *testing.T) {
	csvData := `name,age
"John Doe",30
"Jane Doe",25,extra
`
	reader := strings.NewReader(csvData)
	results := parseCSV(reader, true, fromFunc)

	i := 0
	for result := range results {
		if i == 1 && result.Error == nil {
			t.Errorf("expected error on line 2, got nil")
		}
		i++
	}
}

func TestParseCSVEmpty(t *testing.T) {
	csvData := ``
	reader := strings.NewReader(csvData)
	results := parseCSV(reader, true, fromFunc)

	for result := range results {
		if result.Error == nil {
			t.Errorf("expected error, got nil")
		}
	}
}

func TestParseCSVFromFuncError(t *testing.T) {
	csvData := `name,age
"John Doe",30
`
	reader := strings.NewReader(csvData)
	fromFuncErr := func(data []string, headers []string) (testData, error) {
		return testData{}, errors.New("test error")
	}
	results := parseCSV(reader, true, fromFuncErr)

	for result := range results {
		if result.Error == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
