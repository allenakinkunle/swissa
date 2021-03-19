package converter_test

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/allenakinkunle/swissa/converter"
)

func TestGetHeaders(t *testing.T) {

	t.Run("headers from a CSV file", func(t *testing.T) {

		records := [][]string{
			{"ID", "First Name", "Last Name"},
			{"1", "James", "Bond"},
		}

		csvConverter, clean := createCSVConverterFromFile(t, records, ',')
		defer clean()

		got, err := csvConverter.GetHeaders()
		want := []string{"ID", "First Name", "Last Name"}

		assertNoError(t, err, "could not read headers from CSV file")
		assertCorrectHeaders(t, got, want)

		// Get header again to make sure it returns headers are consistent
		got, _ = csvConverter.GetHeaders()
		assertCorrectHeaders(t, got, want)
	})

	t.Run("headers from file with a different delimiters other than comma", func(t *testing.T) {

		records := [][]string{
			{"ID", "First Name", "Last Name"},
			{"1", "James", "Bond"},
		}

		tests := []struct {
			name      string
			delimiter rune
		}{
			{"tab", '\t'},
			{"colon", ':'},
			{"semicolon", ';'},
			{"pipe", '|'},
		}

		want := []string{"ID", "First Name", "Last Name"}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {

				csvConverter, clean := createCSVConverterFromFile(t, records, test.delimiter)
				defer clean()

				got, err := csvConverter.GetHeaders()

				assertNoError(t, err, "could not read headers from CSV file")
				assertCorrectHeaders(t, got, want)
			})
		}
	})

	t.Run("empty CSV file returns no header", func(t *testing.T) {

		csvConverter, clean := createCSVConverterFromFile(t, nil, ',')
		defer clean()

		got, err := csvConverter.GetHeaders()

		assertCorrectHeaders(t, got, nil)
		assertNoError(t, err, "")
	})
}

func TestNumRecords(t *testing.T) {

	records := [][]string{
		{"ID", "First Name", "Last Name"},
		{"1", "James", "Bond"},
		{"2", "Akinkunle", "Allen"},
		{"# This is a comment", "James", "Bond"},
	}

	csvConverter, clean := createCSVConverterFromFile(t, records, ',')
	defer clean()

	got, err := csvConverter.GetNumRecords()
	want := 2

	if err != io.EOF && err != nil {
		assertNoError(t, err, "could not read CSV file")
	}

	if got != want {
		t.Errorf("incorrect number of records, got %v but want %v", got, want)
	}
}

func TestConvert(t *testing.T) {

	t.Run("convert to JSON", func(t *testing.T) {

		records := [][]string{
			{"ID", "First Name", "Last Name"},
			{"1", "James", "Bond"},
			{"2", "Akinkunle", "Allen"},
		}

		csvConverter, clean := createCSVConverterFromFile(t, records, ',')
		defer clean()

		// Create a temporary JSON file
		tmpJSONFile, err := os.CreateTemp("", "test-*.json")
		defer tmpJSONFile.Close()
		defer os.Remove(tmpJSONFile.Name())

		assertNoError(t, err, "could not create temp JSON file")

		got, err := csvConverter.Convert(converter.FormatJSON, tmpJSONFile)

		assertNoError(t, err, "could not write to file")

		if got != 2 {
			t.Errorf("incorrect number of records converted, got %d but want %d", got, 2)
		}
	})
}

func assertNoError(t testing.TB, err error, message string) {

	t.Helper()

	if err != nil {
		t.Errorf("%s, %v", message, err)
	}
}

func assertCorrectHeaders(t testing.TB, got, want []string) {

	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect headers, got %v but want %v", got, want)
	}
}

func createTempCSVFile(t testing.TB, records [][]string, delimiter rune) *os.File {

	t.Helper()

	tmpCSVFile, err := os.CreateTemp("", "test-*.csv")

	if err != nil {
		t.Errorf("could not create temp CSV file %v", err)
	}

	writer := csv.NewWriter(tmpCSVFile)
	writer.Comma = delimiter
	writer.WriteAll(records)

	tmpCSVFile.Seek(0, 0)

	return tmpCSVFile
}

func createCSVConverterFromFile(t testing.TB, records [][]string, delimiter rune) (*converter.CSVConverter, func()) {

	t.Helper()

	tmpFile := createTempCSVFile(t, records, delimiter)

	cleanUp := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	csvConverter := converter.NewCSVConverter(tmpFile)

	return csvConverter, cleanUp
}
