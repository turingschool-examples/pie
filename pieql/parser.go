package pieql

import (
	"fmt"
	"io"
)

// Parser represents a PieQL parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses the next statement from the underlying reader.
func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// Parse fields.
	fields, err := p.parseFields()
	if err != nil {
		return nil, err
	}
	stmt.Fields = fields

	// Parse the source.
	source, err := p.parseSource()
	if err != nil {
		return nil, err
	}
	stmt.Source = source

	return stmt, nil
}

// parseFields parses one to all fields.
func (p *Parser) parseFields() (Fields, error) {
	var fields Fields

	// Expect to see the "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != MUL {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		fields = append(fields, &Field{Name: lit})

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	return fields, nil
}

// parseSource parses the source table for the query.
func (p *Parser) parseSource() (string, error) {
	// Expect to see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return "", fmt.Errorf("found %q, expected FROM", lit)
	}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return "", fmt.Errorf("found %q, expected table name", lit)
	}

	return lit, nil
}

// scan returns the next token from the scanner.
// If a token been unscanned, read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we need to unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}
