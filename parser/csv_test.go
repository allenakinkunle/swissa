package parser_test

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"

	"github.com/allenakinkunle/swissa/parser"
)

func TestGetHeaders(t *testing.T) {

	t.Run("headers from a CSV file", func(t *testing.T) {

		records := [][]string{
			{"ID", "First Name", "Last Name"},
			{"1", "James", "Bond"},
		}

		csvParser, clean := createCSVParserFromFile(t, records, ',')
		defer clean()

		got, err := csvParser.GetHeaders()
		want := []string{"ID", "First Name", "Last Name"}

		if err != nil {
			t.Errorf("could not read headers from CSV file, %v", err)
		}

		assertCorrectHeaders(t, got, want)

		// Get header again to make sure it returns headers are consistent
		got, _ = csvParser.GetHeaders()
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

				csvParser, clean := createCSVParserFromFile(t, records, ',')
				defer clean()

				got, err := csvParser.GetHeaders()

				if err != nil {
					t.Errorf("could not read headers from CSV file %v", err)
				}

				assertCorrectHeaders(t, got, want)
			})
		}
	})

	t.Run("empty CSV file returns no header", func(t *testing.T) {

		csvParser, clean := createCSVParserFromFile(t, nil, ',')
		defer clean()

		_, err := csvParser.GetHeaders()

		if err == nil {
			t.Errorf("expected an error, but did not get one, %v", err)
		}
	})
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

func createCSVParserFromFile(t testing.TB, records [][]string, delimiter rune) (*parser.CSVParser, func()) {

	t.Helper()

	tmpFile := createTempCSVFile(t, records, delimiter)

	cleanUp := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	csvParser := parser.NewCSVParser(tmpFile)

	return csvParser, cleanUp
}
