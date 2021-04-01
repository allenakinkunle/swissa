package converter

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// CSVConverter parses a CSV-parsable io.Reader
// It only parses files that adheres to the RFC 4180
type CSVConverter struct {
	headers     []string    // CSV file headers
	numRecords  int         // Number of records in CSV file
	headersRead bool        // Flag if header has been read
	reader      *csv.Reader // CSV Reader
}

// NewCSVConverter constructs a new instance of CSVConverter and returns a pointer to it
func NewCSVConverter(file io.Reader) *CSVConverter {

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.LazyQuotes = true

	return &CSVConverter{
		reader:     reader,
		numRecords: -1,
	}
}

// GetHeaders returns a slice of strings containing the file headers
func (c *CSVConverter) GetHeaders() ([]string, error) {

	if c.headersRead {
		return c.headers, nil
	}

	headers, err := c.reader.Read()

	if headers == nil {
		return nil, errors.New("cannot read from an empty file")
	}

	if err != io.EOF && err != nil {
		return nil, fmt.Errorf("error getting headers, %w", err)
	}

	// Infer the file delimiter, since it is not comma
	if len(headers) == 1 {
		delimiter := getDelimiter(headers[0])
		headers = strings.Split(headers[0], string(delimiter))

		c.reader.Comma = delimiter
		c.reader.FieldsPerRecord = len(headers)
	}

	headers = cleanRecord(headers)
	c.headersRead = true
	c.headers = headers

	return headers, nil
}

// GetNumRecords returns the number of records in the CSV file
func (c *CSVConverter) GetNumRecords() (int, error) {

	if c.numRecords != -1 {
		return c.numRecords, nil
	}

	// Skip header
	_, err := c.GetHeaders()

	if err != nil {
		return 0, fmt.Errorf("could not get number of records in CSV file, %w", err)
	}

	numRecords := 0

	for {
		_, err := c.reader.Read()

		switch err {
		case nil:
			numRecords++
		case io.EOF:
			c.numRecords = numRecords
			return c.numRecords, nil
		default:
			return 0, fmt.Errorf("could not get number of records in CSV file, %w", err)
		}
	}
}

// Convert converts the CSV file into the specified formats and
// writes it to the provided io.Writer
func (c *CSVConverter) Convert(toFormat string, writer io.Writer) (int, error) {

	switch toFormat {
	case "json":
		return c.convertToJSON(writer)
	default:
		return 0, errors.New("unsupported conversion file type")
	}
}

func (c *CSVConverter) convertToJSON(writer io.Writer) (int, error) {

	headers, err := c.GetHeaders()

	if err != nil {
		return 0, err
	}

	writer.Write([]byte{'['})

	numRecordsConverted, err := c.buildJSON(headers, nil, 0, writer)

	if err != nil {
		return 0, err
	}

	writer.Write([]byte{']'})

	return numRecordsConverted, err

}

// buildJSON recursively builds the JSON array from CSV records and write them to provided writer
func (c *CSVConverter) buildJSON(headers, record []string, numRecordsConverted int, writer io.Writer) (int, error) {

	var err error

	if numRecordsConverted == 0 {
		record, err = c.reader.Read()
		if err != io.EOF && err != nil {
			return 0, err
		}

		record = cleanRecord(record)
	}

	dictRecord := map[string]string{}

	for i, header := range headers {
		dictRecord[header] = record[i]
	}

	jsonRecord, err := json.MarshalIndent(&dictRecord, "", "  ")

	if err != nil {
		return 0, err
	}

	_, err = writer.Write(jsonRecord)
	if err != nil {
		return 0, err
	}

	numRecordsConverted++

	// Convert records to JSON if there are more to process
	record, err = c.reader.Read()

	if err == nil {
		writer.Write([]byte{',', '\n'})
		record = cleanRecord(record)
		numRecordsConverted, err = c.buildJSON(headers, record, numRecordsConverted, writer)
	}
	if err == io.EOF {
		return numRecordsConverted, nil
	}

	return numRecordsConverted, err
}

func cleanRecord(record []string) []string {
	clean_record := make([]string, len(record))
	for indx, rec := range record {
		clean := strings.TrimFunc(rec, func(r rune) bool {
			return r == '"' || r == '\'' || unicode.IsSpace(r)
		})
		clean_record[indx] = clean
	}
	return clean_record
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
