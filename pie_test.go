package pie_test

import (
	"testing"

	"github.com/turingschool-examples/pie"
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
