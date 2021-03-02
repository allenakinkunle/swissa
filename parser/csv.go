package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVParser parses a CSV-parsable io.Reader
// It only parses files that adheres to the RFC 4180
type CSVParser struct {
	delimiter   rune
	headers     []string    // CSV file headers
	numRecords  int         // Number of records in CSV file
	headersRead bool        // Flag if header has been read
	reader      *csv.Reader // CSV Reader
}

// NewCSVParser constructs a new instance of CSVParser and returns a pointer to it
func NewCSVParser(file io.Reader) *CSVParser {

	reader := csv.NewReader(file)
	reader.Comment = '#'

	return &CSVParser{
		delimiter: ',',
		reader:    reader,
	}
}

// GetHeaders returns a slice of strings containing the file headers
func (c *CSVParser) GetHeaders() ([]string, error) {

	if c.headersRead {
		return c.headers, nil
	}

	headers, err := c.reader.Read()

	// Infer the file delimiter, since it is not comma
	if len(headers) == 1 {
		delimiter := getDelimiter(headers[0])
		c.delimiter = delimiter
		c.reader.Comma = delimiter

		headers = strings.Split(headers[0], string(delimiter))
	}

	if err != nil {
		return nil, fmt.Errorf("could not read from file, %w", err)
	}

	c.headersRead = true
	c.headers = headers

	return headers, nil
}

// GetNumRecords returns the number of records in the CSV file
func (c *CSVParser) GetNumRecords() (int, error) {

	c.GetHeaders()

	numRecords := 0

	for {
		_, err := c.reader.Read()

		switch err {
		case nil:
			numRecords++
		default:
			return numRecords, err
		}
	}
}

// getDelimiter tries to detect the delimiter of the file (rather crudely)
func getDelimiter(row string) rune {

	delimiters := []rune{'\t', ':', ';', '|'}

	for _, delimiter := range delimiters {
		split := strings.Split(row, string(delimiter))
		if len(split) != 1 {
			return delimiter
		}
	}

	return ','
}
