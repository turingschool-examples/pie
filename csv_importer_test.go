package pie_test

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie"
)

// Ensure the importer can create a table with columns based on the CSV header.
func TestCSVImporter_Import(t *testing.T) {
	db := pie.NewDatabase()
	i := pie.NewCSVImporter()

	// Create incoming data.
	data := strings.TrimSpace(`
first_name,last_name,company
susy,que,acme
bob,smith,ford
`)

	// Import CSV data.
	if err := i.Import(db, "my_peeps", csv.NewReader(strings.NewReader(data))); err != nil {
		t.Fatal(err)
	}
	if tbl := db.Table("my_peeps"); tbl == nil {
		t.Fatal("table expected")
	} else if len(tbl.Columns) != 3 {
		t.Fatalf("unexpected table count: %d", len(tbl.Columns))
	} else if tbl.Columns[0].Name != "first_name" {
		t.Fatal("unexpected table name(0)")
	} else if tbl.Columns[1].Name != "last_name" {
		t.Fatal("unexpected table name(0)")
	} else if tbl.Columns[2].Name != "company" {
		t.Fatal("unexpected table name(0)")
	} else if len(tbl.Rows) != 2 {
		t.Fatalf("unexpected row count: %d", len(tbl.Rows))
	}
}
