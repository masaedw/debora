package debora

import (
	"reflect"

	"github.com/pkg/errors"
)

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

// Row a row in a table
type Row interface {
	Columns() []Column
}

// WhereClause where clause
type WhereClause func(*Row) bool

// Queryable query
type Queryable interface {
	Where(WhereClause) Queryable
}

// Table a table
type Table interface {
	Name() string
	Get(int) (Row, error)
	Query() Queryable
	Insert(columns ...interface{}) error
}

// Database a database
type Database interface {
	Tables() []Table
	CreateTable(name string, columns ...ColumnDefinition) (Table, error)
	Get(string) (Table, error)
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

type row []*column

func (row row) Columns() []Column {
	cx := make([]Column, len(row))
	for i, r := range row {
		cx[i] = r
	}
	return cx
}

type table struct {
	name   string
	schema []ColumnDefinition
	data   []row
}

func (t *table) Name() string {
	return t.name
}

func (t *table) Get(i int) (Row, error) {
	if 0 < i && i < len(t.data) {
		return t.data[i], nil
	}
	return nil, errors.Errorf("out of index: %d", i)
}

func (t *table) Query() Queryable {
	// todo implement
	return nil
}

func (t *table) makeRow(cols []interface{}) (row, error) {
	row := make([]*column, len(t.schema))

	for i, d := range t.schema {
		switch d.Type {
		case String:
			v, ok := cols[i].(string)
			if !ok {
				return nil, errors.Errorf("requried string as cols[%d] but got %s", i, reflect.TypeOf(cols[i]))
			}
			row[i] = &column{definition: &d, stringValue: &v}
			continue

		case Integer:
			v, ok := cols[i].(int64)
			if !ok {
				return nil, errors.Errorf("requried int64 as cols[%d] but got %s", i, reflect.TypeOf(cols[i]))
			}
			row[i] = &column{definition: &d, integerValue: &v}
			continue

		default:
			return nil, errors.Errorf("cols[%d] is unsupported type: %s", i, reflect.TypeOf(cols[i]))
		}
	}
	return row, nil
}

func (t *table) Insert(cols ...interface{}) error {
	if len(t.schema) != len(cols) {
		return errors.Errorf("requried %d columns but got %d", len(t.schema), len(cols))
	}

	row, err := t.makeRow(cols)
	if err != nil {
		return errors.WithMessage(err, "can't make a row")
	}

	t.data = append(t.data, row)
	return nil
}

type database struct {
	tables []*table
}

func (d *database) Tables() []Table {
	tx := make([]Table, len(d.tables))
	for i, t := range d.tables {
		tx[i] = t
	}
	return tx
}

func (d *database) has(name string) bool {
	for _, t := range d.tables {
		if t.name == name {
			return true
		}
	}
	return false
}

func (d *database) CreateTable(name string, columns ...ColumnDefinition) (Table, error) {
	if d.has(name) {
		return nil, errors.Errorf("the %s table is already exists", name)
	}

	if len(columns) == 0 {
		return nil, errors.Errorf("no columns")
	}

	t := &table{
		name:   name,
		schema: columns,
		data:   []row{},
	}
	d.tables = append(d.tables, t)

	return t, nil
}

func (d *database) Get(name string) (Table, error) {
	for _, t := range d.tables {
		if t.name == name {
			return t, nil
		}
	}
	return nil, errors.Errorf("no such table %s", name)
}
