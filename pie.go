package pie

import (
	"os"
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
