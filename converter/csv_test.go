package converter_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/allenakinkunle/swissa/converter"
)

func TestConvert(t *testing.T) {

	t.Run("convert to JSON", func(t *testing.T) {

		tests := []struct {
			csvString string
			delimiter string
		}{
			{"ID,First_Name,Last_Name\n1,James,Bond\n", "comma"},
			{"ID\tFirst_Name\tLast_Name\n1\tJames\tBond\n", "tab"},
			{"ID:First_Name:Last_Name\n1:James:Bond\n", "colon"},
			{"ID;First_Name;Last_Name\n1;James;Bond\n", "semicolon"},
			{"ID|First_Name|Last_Name\n1|James|Bond\n", "pipe"},
		}

		want := 1

		for _, test := range tests {
			t.Run(test.delimiter, func(t *testing.T) {
				csvConverter := converter.NewCSVConverter(strings.NewReader(test.csvString))

				tmpJSONFile, err := os.CreateTemp("", "test-*.json")
				assertNoError(t, err)

				got, err := csvConverter.Convert("json", tmpJSONFile)
				assertNoError(t, err)

				if got != want {
					t.Errorf("incorrect number of records converted, got %d but want %d", got, want)
				}

				tmpJSONFile.Close()
				os.Remove(tmpJSONFile.Name())
			})
		}
	})

	t.Run("convert empty file to JSON raises error", func(t *testing.T) {
		const csvString = ""
		csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))
		_, err := csvConverter.Convert("json", os.Stdout)
		assertError(t, err)
	})

	t.Run("convert file with only header to JSON raises error", func(t *testing.T) {
		const csvString = "ID,First_Name,Last_Name\n"
		csvConverter := converter.NewCSVConverter(strings.NewReader(csvString))
		_, err := csvConverter.Convert("json", os.Stdout)
		assertError(t, err)
	})
}

func assertNoError(t testing.TB, err error) {

	t.Helper()

	if err != nil {
		t.Errorf("%v", err)
	}
}

func assertError(t testing.TB, err error) {

	t.Helper()

	if err == nil {
		t.Errorf("%v", err)
	}
}

func assertCorrectHeaders(t testing.TB, got, want []string) {

	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect headers, got %v but want %v", got, want)
	}
}
