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
	// Read CSV headers.
	record, err := r.Read()
	if err != nil {
		return err
	}

	// Create columns from headers.
	var columns []*Column
	for _, name := range record {
		columns = append(columns, &Column{Name: name})
	}

	// Create table in database.
	if err := db.CreateTable(name, columns); err != nil {
		return err
	}

	// Read remaining rows into Table.Rows.
	rows, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Write table rows to disk.
	if err := db.SetTableRows(name, rows); err != nil {
		return err
	}

	return nil
}
