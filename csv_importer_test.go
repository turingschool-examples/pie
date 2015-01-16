package pie_test

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie"
)

// Ensure the importer can create a table with columns based on the CSV header.
func TestCSVImporter_Import(t *testing.T) {
	db := OpenDatabase()
	defer db.Close()
	i := pie.NewCSVImporter()

	// Create incoming data.
	data := strings.TrimSpace(`
first_name,last_name,company
susy,que,acme
bob,smith,ford
`)

	// Import CSV data.
	if err := i.Import(db.Database, "my_peeps", csv.NewReader(strings.NewReader(data))); err != nil {
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
	}

	if rows, err := db.TableRows("my_peeps"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if len(rows) != 2 {
		t.Fatalf("unexpected row count: %d", len(rows))
	} else if !reflect.DeepEqual(rows[0], []string{"susy", "que", "acme"}) {
		t.Fatalf("unexpected row(0): %#v", rows[0])
	} else if !reflect.DeepEqual(rows[1], []string{"bob", "smith", "ford"}) {
		t.Fatalf("unexpected row(0): %#v", rows[1])
	}
}
