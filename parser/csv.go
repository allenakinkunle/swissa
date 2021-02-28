package parser

import (
	"encoding/csv"
	"fmt"
	"io"
)

// CSVParser parses a CSV-parsable io.Reader
type CSVParser struct {
	headers     []string      // CSV file headers
	headersRead bool          // Flag if header has been read
	file        io.ReadSeeker // Reader to parse
	reader      *csv.Reader   // CSV Reader
}

// NewCSVParser constructs a new instance of CSVParser and returns a pointer to it
func NewCSVParser(file io.ReadSeeker) *CSVParser {
	return &CSVParser{
		file:   file,
		reader: csv.NewReader(file),
	}
}

// GetHeaders returns a slice of strings containing the file headers
// It returns an error if an error occurs while parsing file
func (c *CSVParser) GetHeaders() ([]string, error) {
	if c.headersRead {
		return c.headers, nil
	}

	c.file.Seek(0, 0)

	headers, err := c.reader.Read()

	if err != nil {
		return nil, fmt.Errorf("could not read from file, %w", err)
	}

	c.headersRead = true
	c.headers = headers

	return headers, nil
}
