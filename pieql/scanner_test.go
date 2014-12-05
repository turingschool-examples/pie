package pieql_test

import (
	"strings"
	"testing"

	"github.com/turingschool-examples/pie/pieql"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok pieql.Token
		lit string
	}{
		// Special Tokens (whitespace)
		{s: `!`, tok: pieql.ILLEGAL, lit: "!"},
		{s: ``, tok: pieql.EOF, lit: ""},
		{s: ` `, tok: pieql.WS, lit: " "},
		{s: ` X`, tok: pieql.WS, lit: " "},
		{s: "\n \t", tok: pieql.WS, lit: "\n \t"},

		// Misc
		{s: `*`, tok: pieql.MUL, lit: `*`},
		{s: `,`, tok: pieql.COMMA, lit: `,`},

		// Identifiers
		{s: `foo`, tok: pieql.IDENT, lit: `foo`},
		{s: `foo_20 `, tok: pieql.IDENT, lit: `foo_20`},

		// Keywords
		{s: `SELECT`, tok: pieql.SELECT, lit: `SELECT`},
		{s: `FROM`, tok: pieql.FROM, lit: `FROM`},
		{s: `from`, tok: pieql.FROM, lit: `from`},
	}

	for i, tt := range tests {
		s := pieql.NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
