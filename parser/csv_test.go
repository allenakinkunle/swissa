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

		tmpFile := createTempCSVFile(t, records)
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		csvParser := parser.NewCSVParser(tmpFile)

		got, err := csvParser.GetHeaders()
		want := []string{"ID", "First Name", "Last Name"}

		if err != nil {
			t.Fatalf("could not read headers from CSV file %v, %v", tmpFile.Name(), err)
		}

		assertCorrectHeaders(t, got, want)

		// Get header again to make sure it returns headers are consistent
		got, _ = csvParser.GetHeaders()
		assertCorrectHeaders(t, got, want)
	})

	t.Run("empty CSV file returns no header", func(t *testing.T) {

		var records [][]string

		tmpFile := createTempCSVFile(t, records)
		defer tmpFile.Close()

		csvParser := parser.NewCSVParser(tmpFile)

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

func createTempCSVFile(t testing.TB, records [][]string) *os.File {

	t.Helper()

	tmpCSVFile, err := os.CreateTemp("", "test-*.csv")

	if err != nil {
		t.Fatalf("could not create temp CSV file %v", err)
	}

	csv.NewWriter(tmpCSVFile).WriteAll(records)

	return tmpCSVFile

}
