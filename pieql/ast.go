package pieql

// SelectStatement represents a statement for retrieving data.
type SelectStatement struct {
	Fields Fields
	Source string
}

// Fields represents a list of fields.
type Fields []*Field

// Field represents a column to be selected.
type Field struct {
	Name string
}
