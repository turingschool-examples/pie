package pie

import (
	"errors"
	"os"
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

// CreateTable creates a new table.
// Returns an error if name is blank or if table already exists.
func (db *Database) CreateTable(name string, columns []*Column) error {
	// TODO: Check for blank name.
	// TODO: Check for existing table with the same name.

	// Create table.
	t := &Table{Name: name, Columns: columns}

	// Add table to the database.
	db.tables[name] = t

	return nil
}

// DeleteTable removes an existing table by name.
// Returns an error if name is blank or table is not found.
func (db *Database) DeleteTable(name string) error {
	// TODO: Check for blank name.
	// TODO: Check that table exists.
	// TODO: Remove table from the database.
	return nil
}

// Table represents a tabular set of data.
type Table struct {
	Name    string
	Columns []*Column
	Rows    [][]string
}

// Column represents a column in a table.
type Column struct {
	Name string
}
