package pieql_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie/pieql"
)

// Ensure the parser can parse PieQL into an AST.
func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		q    string
		stmt *pieql.SelectStatement
	}{
		// 0. Simple SELECT statement.
		{
			q: `SELECT fname FROM tbl`,
			stmt: &pieql.SelectStatement{
				Fields: pieql.Fields{
					&pieql.Field{Name: "fname"},
				},
				Source: "tbl",
			},
		},
	}

	// Parse querystring into AST.
	for i, tt := range tests {
		// Parse the query.
		p := pieql.NewParser(strings.NewReader(tt.q))
		stmt, err := p.Parse()
		if err != nil {
			t.Errorf("%d. %q: error: %s", i, tt.q, err)
			continue
		}

		// Ensure AST matches.
		if !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q: stmt mismatch:\n\n%#v", stmt)
			continue
		}
	}
}
