package parser

// Parser interface is implemented by all supported file parsers
// in the 'convert' package
type Parser interface {
	GetHeaders() ([]string, error)
}
