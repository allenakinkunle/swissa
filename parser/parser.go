package parser

import "io"

// Conversion formats supported
const (
	FormatJSON = "json"
)

// Parser interface is implemented by all supported file parsers
// in the 'parser' package
type Parser interface {
	GetHeaders() ([]string, error)
	GetNumRecords() (int, error)
	Convert(toFormat string, writer io.Writer) (int, error)
}
