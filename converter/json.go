package converter

import (
	"encoding/json"
	"io"
)

// JSONConverter parses a JSON-parsable io.Reader
// It only parses files that adhere to the RFC 7159 standard
type JSONConverter struct {
	headers     []string
	numRecords  int
	headersRead bool
	decoder     *json.Decoder
}

// NewJSONConverter constructs a new instance of JSONConverter and returns a pointer to it
func NewJSONConverter(file io.Reader) *JSONConverter {
	return &JSONConverter{
		decoder:    json.NewDecoder(file),
		numRecords: -1,
	}
}

// GetHeaders returns a slice of strings containing the file headers
// func (j *JSONConverter) GetHeaders() ([]string, error) {

// 	// if j.headersRead {
// 	// 	return j.headers, nil
// 	// }
// }

// GetNumRecords returns the number of records in the JSON file
// func (j *JSONConverter) GetNumRecords() (int, error) {

// 	if j.numRecords != 1 {
// 		return j.numRecords, nil
// 	}

// }
