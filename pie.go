package pie

// Database represents a collection of tables.
type Database struct {
	tables map[string]*Table
}

// NewDatabase returns a new instance of Database.
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
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
