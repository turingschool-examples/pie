package pieql

import (
	"io"
)

// Parser represents a PieQL parser.
type Parser struct {
	s *Scanner
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses the next statement from the underlying reader.
func (p *Parser) Parse() (*SelectStatement, error) {
	panic("aaaaaaah!!!!")
}
