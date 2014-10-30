package pie_test

import (
	"testing"

	"github.com/turingschool-examples/pie"
)

// Ensure database can correctly open and close.
func TestDatabase_Open(t *testing.T) {
	db := pie.NewDatabase()
	if err := db.Open("/tmp/pie"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
