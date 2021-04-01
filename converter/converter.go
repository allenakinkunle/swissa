package converter

import "io"

// Conversion formats supported
var SupportedFormats = []string{"csv", "json"}

// Converter interface is implemented by all supported file converters
// in the 'converter' package
type Converter interface {
	GetHeaders() ([]string, error)
	Convert(toFormat string, writer io.Writer) (int, error)
}
