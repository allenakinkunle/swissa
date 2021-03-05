package converter

import "io"

// Conversion formats supported
const (
	FormatJSON = "json"
)

// Converter interface is implemented by all supported file converters
// in the 'converter' package
type Converter interface {
	GetHeaders() ([]string, error)
	GetNumRecords() (int, error)
	Convert(toFormat string, writer io.Writer) (int, error)
}
