package pie

import (
	"errors"
	"fmt"
	"os"

	"github.com/turingschool-examples/pie/pieql"
)

var (
	// ErrTableNotFound is returned when referencing a table that doesn't exist.
	ErrTableNotFound = errors.New("table not found")

	// ErrTableExists is returned when creating a table that already exists.
	ErrTableExists = errors.New("table already exists")

	// ErrTableNameRequired is returned when a blank table name is passed in.
	ErrTableNameRequired = errors.New("table name required")
)

// Database represents a collection of tables.
type Database struct {
	path   string
	tables map[string]*Table
}

// NewDatabase returns a new instance of Database.
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

// Open opens and initializes a database at a given file path.
func (db *Database) Open(path string) error {
	// Make a new directory.
	if err := os.Mkdir(path, 0700); err != nil {
		return err
	}

	// Set the path.
	db.path = path

	// TODO: Open meta file.

	return nil
}

func (db *Database) Close() error {
	// Unset the path.
	db.path = ""
	return nil
}

// Table returns a table by name.
func (db *Database) Table(name string) *Table {
	return db.tables[name]
}

// Tables returns a list of all tables in the database.
func (db *Database) Tables() []*Table {
	var a []*Table
	for _, t := range db.tables {
		a = append(a, t)
	}
	return a
}

// CreateTable creates a new table.
// Returns an error if name is blank or if table already exists.
func (db *Database) CreateTable(name string, columns []*Column) error {
	// Check for blank name.
	// Check for existing table with the same name.
	if name == "" {
		return ErrTableNameRequired
	} else if db.tables[name] != nil {
		return ErrTableExists
	}

	// Add table to the database.
	db.tables[name] = &Table{Name: name, Columns: columns}

	return nil
}

// DeleteTable removes an existing table by name.
// Returns an error if name is blank or table is not found.
func (db *Database) DeleteTable(name string) error {
	// TODO: Check for blank name.
	if name == "" {
		return ErrTableNameRequired
	} else if db.tables[name] == nil {
		return ErrTableNotFound
	}
	// TODO: Check that table exists.
	// TODO: Remove table from the database.
	delete(db.tables, name)

	return nil
}

// Execute executes a SELECT statement and returns the results.
func (db *Database) Execute(stmt *pieql.SelectStatement) ([][]string, error) {
	// Lookup table by name.
	t := db.Table(stmt.Source)
	if t == nil {
		return nil, ErrTableNotFound
	}

	// Iterate over all the table rows.
	var result [][]string
	for _, row := range t.Rows {
		resultRow := make([]string, len(stmt.Fields))

		// Lookup row value by field name for each field.
		for i, f := range stmt.Fields {
			// Lookup column index.
			index := t.ColumnIndex(f.Name)
			if index == -1 {
				return nil, fmt.Errorf("column not found: %s", f.Name)
			}

			// Set result cell value.
			resultRow[i] = row[index]
		}

		// Add output row to the result.
		result = append(result, resultRow)
	}

	return result, nil
}

// Table represents a tabular set of data.
type Table struct {
	Name    string
	Columns []*Column
	Rows    [][]string
}

// ColumnIndex returns the position of the column by name.
// Returns -1 if column is not found.
func (t *Table) ColumnIndex(name string) int {
	for i, c := range t.Columns {
		if c.Name == name {
			return i
		}
	}
	return -1
}

// Column represents a column in a table.
type Column struct {
	Name string
}
