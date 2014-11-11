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
}
