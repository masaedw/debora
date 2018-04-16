package debora

// ColumnType a type of a column
type ColumnType int

const (
	// String string
	String ColumnType = iota
	// Integer 64bit integer
	Integer
)

// ColumnDefinition a column definition
type ColumnDefinition struct {
	Type ColumnType
	Name string
}

// Column a column
type Column interface {
	Type() ColumnType
	Name() string
	StringValue() *string
	IntegerValue() *int64
}

type column struct {
	definition   *ColumnDefinition
	integerValue *int64
	stringValue  *string
}

func (c *column) Type() ColumnType {
	return c.definition.Type
}

func (c *column) Name() string {
	return c.definition.Name
}

func (c *column) StringValue() *string {
	return c.stringValue
}

func (c *column) IntegerValue() *int64 {
	return c.integerValue
}
