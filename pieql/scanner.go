package pieql

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token int

const (
	//Special Tokens
	ILLEGAL Token = iota
	EOF
	WS

	//Literals
	IDENT

	// Misc
	ASTERISK
	COMMA

	// Keywords
	SELECT
	FROM
)

// Scanner represents a lexical scanner for PieQL
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns an instance of Scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the enxt token and position from the reader
// Also returns the literal text read for strings
func (s *Scanner) Scan() (tok Token, lit string) {
	// read the next rune
	ch := s.read()

	// if whitespace then consume all contiguous whitespace
	// if we see a letter then put back into reader and consume
	// as an ident or reserved word
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	// otherwise read the individual character
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character in
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer
	// Non-whitespace characters and EOF will exit the loop
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character in
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer
	// Non-ident characters and EOF will exit the loop
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword, return that keyword
	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	}

	// Otherwise return as a regular identifier
	return IDENT, buf.String()
}

// reads the next rune from the reader
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread pleaces the previously read rune back onto the reader
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

type reader struct {
	r   io.RuneScanner
	i   int
	n   int
	buf [3]struct {
		ch rune
	}
}

var eof = rune(0)

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isLetter(ch rune) bool     { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }
func isDigit(ch rune) bool      { return (ch >= '0' && ch <= '9') }
