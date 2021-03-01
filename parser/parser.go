package parser

// Parser interface is implemented by all supported file parsers
// in the 'parser' package
type Parser interface {
	GetHeaders() ([]string, error)
}
