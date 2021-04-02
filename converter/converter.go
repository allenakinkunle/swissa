package converter

import "io"

// Conversion formats supported
const (
	FormatJSON = "json"
	FormatCSV  = "csv"
)

var SupportedFormats = []string{FormatCSV, FormatJSON}

// Converter interface is implemented by all supported file converters
// in the 'converter' package
type Converter interface {
	Convert(toFormat string, writer io.Writer) (int, error)
}
