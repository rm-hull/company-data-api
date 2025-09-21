package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	results := ParseCSV(reader, true, fromFunc)

	expected := []Result[testData]{
		{Value: testData{Name: "John Doe", Age: 30}, LineNum: 1},
		{Value: testData{Name: "Jane Doe", Age: 25}, LineNum: 2},
	}

	var actual []Result[testData]
	for result := range results {
		actual = append(actual, result)
	}

	require.Len(t, actual, len(expected))
	for i := range expected {
		require.NoError(t, actual[i].Error)
		assert.Equal(t, expected[i].Value, actual[i].Value)
		assert.Equal(t, expected[i].LineNum, actual[i].LineNum)
	}
}

func TestParseCSVWithoutHeader(t *testing.T) {
	csvData := `"John Doe",30
"Jane Doe",25
`
	reader := strings.NewReader(csvData)
	results := ParseCSV(reader, false, fromFunc)

	expected := []Result[testData]{
		{Value: testData{Name: "John Doe", Age: 30}, LineNum: 1},
		{Value: testData{Name: "Jane Doe", Age: 25}, LineNum: 2},
	}

	var actual []Result[testData]
	for result := range results {
		actual = append(actual, result)
	}

	require.Len(t, actual, len(expected))
	for i := range expected {
		require.NoError(t, actual[i].Error)
		assert.Equal(t, expected[i].Value, actual[i].Value)
		assert.Equal(t, expected[i].LineNum, actual[i].LineNum)
	}
}

func TestParseCSVMalformed(t *testing.T) {
	csvData := `name,age
"John Doe",30
"Jane Doe",25,extra
`
	reader := strings.NewReader(csvData)
	results := ParseCSV(reader, true, fromFunc)

	var resultsSlice []Result[testData]
	for result := range results {
		resultsSlice = append(resultsSlice, result)
	}
	require.Len(t, resultsSlice, 2, "expected two results")
	assert.NoError(t, resultsSlice[0].Error, "first result should not have an error")
	assert.Error(t, resultsSlice[1].Error, "second result should have an error")
}

func TestParseCSVEmpty(t *testing.T) {
	cases := map[string]struct {
		withHeader        bool
		expectedNumResult int
		expectError       bool
	}{
		"with header":    {true, 1, true},
		"without header": {false, 0, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			reader := strings.NewReader("")
			results := ParseCSV(reader, tc.withHeader, fromFunc)

			var resultsSlice []Result[testData]
			for result := range results {
				resultsSlice = append(resultsSlice, result)
			}

			require.Len(t, resultsSlice, tc.expectedNumResult)
			if tc.expectError {
				assert.Error(t, resultsSlice[0].Error)
			}
		})
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
	results := ParseCSV(reader, true, fromFuncErr)

	var resultsSlice []Result[testData]
	for result := range results {
		resultsSlice = append(resultsSlice, result)
	}
	require.Len(t, resultsSlice, 1, "expected one result")
	assert.Error(t, resultsSlice[0].Error, "expected an error from fromFunc")
}
