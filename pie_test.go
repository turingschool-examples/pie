package pie_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie"
	"github.com/turingschool-examples/pie/pieql"
)

// Ensure the database can create a table.
func TestDatabase_CreateTable(t *testing.T) {
	db := pie.NewDatabase()

	// Create table named "foo" with two columns.
	columns := []*pie.Column{
		&pie.Column{Name: "first_name"},
		&pie.Column{Name: "last_name"},
	}
	if err := db.CreateTable("foo", columns); err != nil {
		t.Fatal(err)
	}

	// Retrieve table and verify it's correct.
	if table := db.Table("foo"); table == nil {
		t.Fatal("expected table")
	} else if table.Name != "foo" {
		t.Fatalf("unexpected name: %q", table.Name)
	} else if len(table.Columns) != 2 {
		t.Fatalf("unexpected column count: %d", len(table.Columns))
	} else if table.Columns[0].Name != "first_name" {
		t.Fatalf("unexpected column(0) name: %s", table.Columns[0].Name)
	} else if table.Columns[1].Name != "last_name" {
		t.Fatalf("unexpected column(1) name: %s", table.Columns[1].Name)
	}
}

// Ensure the database returns an error when creating a table without a name.
func TestDatabase_CreateTable_ErrTableNameRequired(t *testing.T) {
	db := pie.NewDatabase()
	if err := db.CreateTable("", nil); err != pie.ErrTableNameRequired {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Ensure the database returns an error when creating a duplicate table.
func TestDatabase_CreateTable_ErrTableExists(t *testing.T) {
	db := pie.NewDatabase()
	db.CreateTable("foo", nil)
	if err := db.CreateTable("foo", nil); err != pie.ErrTableExists {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Ensure the database can delete a table.
func TestDatabase_DeleteTable(t *testing.T) {
	db := pie.NewDatabase()

	// Create the table and verify it exists.
	if err := db.CreateTable("foo", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if db.Table("foo") == nil {
		t.Fatal("table not actually created")
	}

	// Delete the table and verify it's gone.
	if err := db.DeleteTable("foo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if db.Table("foo") != nil {
		t.Fatal("table not actually delete")
	}
}

// Ensure the database returns an error when deleting a table without a name.
func TestDatabase_DeleteTable_ErrTableNameRequired(t *testing.T) {
	db := pie.NewDatabase()
	if err := db.DeleteTable(""); err != pie.ErrTableNameRequired {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Ensure the database returns an error when deleting a table that doesn't exist.
func TestDatabase_DeleteTable_ErrTableNotFound(t *testing.T) {
	db := pie.NewDatabase()
	if err := db.DeleteTable("no_such_table"); err != pie.ErrTableNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Ensure the database can execute a selection query.
func TestDatabase_Execute(t *testing.T) {
	// Create database and seed table.
	db := pie.NewDatabase()
	db.CreateTable("foo", []*pie.Column{{Name: "fname"}, {Name: "lname"}})
	db.Table("foo").Rows = [][]string{
		{"susy", "que"},
		{"bob", "smith"},
	}

	// Parse PieQL statement.
	stmt, err := pieql.NewParser(strings.NewReader(`SELECT lname, fname FROM foo`)).Parse()
	if err != nil {
		t.Fatal(err)
	}

	// Execute statement.
	res, err := db.Execute(stmt)
	if err != nil {
		t.Fatal(err)
	}

	// Verify results.
	if len(res) != 2 {
		t.Fatalf("result len mismatch: %d", len(res))
	} else if !reflect.DeepEqual(res[0], []string{"que", "susy"}) {
		t.Fatalf("row(0) mismatch: %#v", res[0])
	} else if !reflect.DeepEqual(res[1], []string{"smith", "bob"}) {
		t.Fatalf("row(1) mismatch: %#v", res[1])
	}
}
