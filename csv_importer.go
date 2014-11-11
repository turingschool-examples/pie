package pie

import (
	"encoding/csv"
)

// CSVImporter creates a table by importing from a CSV reader into a database.
type CSVImporter struct{}

// NewCSVImporter returns a new instance of CSVImporter.
func NewCSVImporter() *CSVImporter {
	return &CSVImporter{}
}

// Import creates a new table in the database from data in the CSV reader.
func (i *CSVImporter) Import(db *Database, name string, r *csv.Reader) error {
	// TODO: Read CSV headers.
	// TODO: Create columns from headers.
	// TODO: Create table in database.

	panic("not yet implemented: CSVImporter.Import()")
}
