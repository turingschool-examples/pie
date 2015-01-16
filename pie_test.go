package pie_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/turingschool-examples/pie"
	"github.com/turingschool-examples/pie/pieql"
)

// Ensure the database can create a table.
func TestDatabase_CreateTable(t *testing.T) {
	path := tempfile()

	db := pie.NewDatabase()
	if err := db.Open(path); err != nil {
		t.Fatalf("open: %s", err)
	}
	defer db.Close()

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

	// Close the database.
	if err := db.Close(); err != nil {
		t.Fatalf("close: %s", err)
	} else if db.Table("foo") != nil {
		t.Fatalf("table foo still exists")
	}

	// Reopen the database and make sure shit is still there.
	if err := db.Open(path); err != nil {
		t.Fatalf("reopen: %s", err)
	} else if db.Table("foo") == nil {
		t.Fatalf("expected foo table after reopen")
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
	db := OpenDatabase()
	defer db.Close()
	db.CreateTable("foo", []*pie.Column{{Name: "fname"}, {Name: "lname"}})
	db.SetTableRows("foo", [][]string{
		{"susy", "que"},
		{"bob", "smith"},
	})

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

// Ensure the database can marshal metadata to JSON.
func TestDatabase_MarshalJSON(t *testing.T) {
	// Create a database with two tables.
	db := pie.NewDatabase()
	db.CreateTable("foo", []*pie.Column{{Name: "fname"}, {Name: "lname"}})
	db.CreateTable("bar", []*pie.Column{{Name: "age"}})

	// Add data to one table.
	db.SetTableRows("foo", [][]string{{"bob", "smith"}})

	// Marshal database into JSON.
	if b, err := json.Marshal(db); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if string(b) != `{"tables":[{"name":"foo","columns":[{"name":"fname"},{"name":"lname"}]},{"name":"bar","columns":[{"name":"age"}]}]}` {
		t.Fatalf("unexpected bytes: %s", b)
	}
}

// Ensure the database can unmarshal JSON to metadata.
func TestDatabase_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"tables":[{"name":"foo","columns":[{"name":"fname"},{"name":"lname"}]},{"name":"bar","columns":[{"name":"age"}]}]}`)

	// Unmarshal JSON into database.
	db := pie.NewDatabase()
	if err := json.Unmarshal(data, &db); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if len(db.Tables()) != 2 {
		t.Fatalf("table count mismatch: %d", len(db.Tables()))
	} else if db.Table("foo") == nil {
		t.Fatalf("table foo not found")
	}
}

// Database is a test wrapper for pie.Database.
type Database struct {
	*pie.Database
}

// OpenDatabase returns a new, opened instance of Database.
func OpenDatabase() *Database {
	db := pie.NewDatabase()
	if err := db.Open(tempfile()); err != nil {
		panic(err.Error())
	}
	return &Database{db}
}

// Close closes the database and removes the underlying data.
func (db *Database) Close() {
	defer os.RemoveAll(db.Path())
	db.Database.Close()
}

// tempfile returns the path to a non-existent temporary file.
func tempfile() string {
	f, _ := ioutil.TempFile("", "pie-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}
