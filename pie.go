package pie

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/turingschool-examples/pie/pieql"
)

var (
	// ErrNotOpen is returned using the database while it's closed.
	ErrNotOpen = errors.New("database not open")

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
	// Set the path.
	db.path = path

	// Make a new directory.
	if err := os.MkdirAll(db.path, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(db.dataPath(), 0700); err != nil {
		return err
	}

	// Open meta file.
	if err := db.load(); err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() error {
	db.path = ""
	db.tables = make(map[string]*Table)
	return nil
}

// Path returns the root path of the database.
func (db *Database) Path() string { return db.path }

func (db *Database) dataPath() string {
	if db.path == "" {
		return ""
	}
	return filepath.Join(db.path, "data")
}

// load reads the metadata from disk.
func (db *Database) load() error {
	// Open the meta file.
	f, err := os.Open(filepath.Join(db.path, "meta"))
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Unmarshal the meta file.
	if err := json.NewDecoder(f).Decode(&db); err != nil {
		return err
	}

	return nil
}

// save persists the metadata to disk.
func (db *Database) save() error {
	if db.path == "" {
		return nil
	}

	// Open file for writing.
	f, err := os.Create(filepath.Join(db.path, "meta"))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Marshal metadata to file.
	if err := json.NewEncoder(f).Encode(db); err != nil {
		return err
	}

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

	return db.save()
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

	return db.save()
}

// TableRows retrieves the rows for a table from disk.
func (db *Database) TableRows(name string) ([][]string, error) {
	// Open data file for reading.
	f, err := os.Open(filepath.Join(db.dataPath(), name))
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	// Decode rows from disk.
	var rows [][]string
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		return nil, err
	}

	return rows, nil
}

// SetTableRows sets the rows on a table and saves the rows to disk.
func (db *Database) SetTableRows(name string, rows [][]string) error {
	// Verify database is open.
	if db.path == "" {
		return ErrNotOpen
	}

	// Open data file for writing.
	f, err := os.Create(filepath.Join(db.dataPath(), name))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Encode rows to disk.
	if err := json.NewEncoder(f).Encode(rows); err != nil {
		return err
	}

	return nil
}

// Execute executes a SELECT statement and returns the results.
func (db *Database) Execute(stmt *pieql.SelectStatement) ([][]string, error) {
	// Lookup table by name.
	t := db.Table(stmt.Source)
	if t == nil {
		return nil, ErrTableNotFound
	}

	// Expand out SELECT ALL.
	if len(stmt.Fields) > 0 && stmt.Fields[0].Name == "*" {
		stmt.Fields = nil
		for _, c := range t.Columns {
			stmt.Fields = append(stmt.Fields, &pieql.Field{Name: c.Name})
		}
	}

	// Retrieve rows for table.
	rows, err := db.TableRows(stmt.Source)
	if err != nil {
		return nil, err
	}

	// Iterate over all the table rows.
	var result [][]string
	for _, row := range rows {
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

// MarshalJSON encodes the database metadata as JSON.
func (db *Database) MarshalJSON() ([]byte, error) {
	var dm databaseJSONMarshaler
	for _, t := range db.tables {
		tm := &tableJSONMarshaler{
			Name:    t.Name,
			Columns: t.Columns,
		}
		dm.Tables = append(dm.Tables, tm)
	}
	return json.Marshal(dm)
}

// UnmarshalJSON decodes the JSON as database metadata.
func (db *Database) UnmarshalJSON(data []byte) error {
	var dm databaseJSONMarshaler
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}

	// Copy marshaled data to internal types.
	db.tables = make(map[string]*Table)
	for _, tm := range dm.Tables {
		t := &Table{
			Name:    tm.Name,
			Columns: tm.Columns,
		}
		db.tables[t.Name] = t
	}

	return nil
}

type databaseJSONMarshaler struct {
	Tables []*tableJSONMarshaler `json:"tables"`
}

// Table represents a tabular set of data.
type Table struct {
	Name    string
	Columns []*Column
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
	Name string `json:"name"`
}

type tableJSONMarshaler struct {
	Name    string    `json:"name"`
	Columns []*Column `json:"columns"`
}
