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
	headersRead bool        // Flag if header has been read
	reader      *csv.Reader // CSV Reader
}

// NewCSVParser constructs a new instance of CSVParser and returns a pointer to it
func NewCSVParser(file io.Reader) *CSVParser {
	return &CSVParser{
		delimiter: ',',
		reader:    csv.NewReader(file),
	}
}

// GetHeaders returns a slice of strings containing the file headers
// It returns an error if an error occurs while parsing file
func (c *CSVParser) GetHeaders() ([]string, error) {
	if c.headersRead {
		return c.headers, nil
	}

	headers, err := c.reader.Read()

	// Infer the file delimiter, since it is not comma
	if len(headers) == 1 {
		delimiter := getDelimiter(headers[0])
		c.reader.Comma = delimiter
		c.delimiter = delimiter
		headers = strings.Split(headers[0], string(delimiter))
	}

	if err != nil {
		return nil, fmt.Errorf("could not read from file, %w", err)
	}

	c.headersRead = true
	c.headers = headers

	return headers, nil
}

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
