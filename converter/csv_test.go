package converter_test

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/allenakinkunle/swissa/converter"
)

func TestGetHeaders(t *testing.T) {

	t.Run("headers from a CSV file", func(t *testing.T) {

		const csvString = `ID,First Name,Last Name
		1,James,Bond`

		csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))

		got, err := csvConverter.GetHeaders()
		want := []string{"ID", "First Name", "Last Name"}

		assertNoError(t, err, "could not read headers from CSV file")
		assertCorrectHeaders(t, got, want)

		// Get header again to make sure it returns headers are consistent
		got, _ = csvConverter.GetHeaders()
		assertCorrectHeaders(t, got, want)
	})

	t.Run("headers from file with a different delimiters other than comma", func(t *testing.T) {

		tests := []struct {
			csvString string
			delimiter string
		}{
			{`ID	First Name	Last Name
			1	James	Bond`, "tab"},

			{`ID:First Name:Last Name
			1:James:Bond`, "colon"},

			{`ID;First Name;Last Name
			1;James;Bond`, "semicolon"},

			{`ID|First Name|Last Name
			1|James|Bond`, "pipe"},
		}

		want := []string{"ID", "First Name", "Last Name"}

		for _, test := range tests {
			t.Run(test.csvString, func(t *testing.T) {

				csvConverter := converter.NewCSVConverter(strings.NewReader(test.csvString))

				got, err := csvConverter.GetHeaders()

				assertNoError(t, err, "could not read headers from CSV file")
				assertCorrectHeaders(t, got, want)
			})
		}
	})

	t.Run("empty CSV file returns no header", func(t *testing.T) {

		const csvString = ``
		csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))

		got, err := csvConverter.GetHeaders()

		assertCorrectHeaders(t, got, nil)
		assertNoError(t, err, "")
	})
}

func TestNumRecords(t *testing.T) {

	const csvString = `ID,First Name,Last Name
		1,James,Bond
		2,Akinkunle,Allen`

	csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))

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

		const csvString = `ID,First Name,Last Name
		1,James,Bond
		2,Akinkunle,Allen`

		csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))

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
