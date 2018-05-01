package debora

import (
	"github.com/pkg/errors"
)

// WhereClause where clause
type WhereClause func(Row) bool

// ColumnSelector column selector
type ColumnSelector func(TableLike) ColumnDefinition

// Iterator iterator procedure
type Iterator func(Row)

// TableLike a table like object
type TableLike interface {
	Schema() []ColumnDefinition
	ForEach(Iterator)
}

// Queryable query
type Queryable interface {
	TableLike
	Where(WhereClause) Queryable
	Select(...ColumnSelector) Queryable
	Slice() ([]Row, error)
	First() (Row, error)
}

// Table a table
type Table interface {
	TableLike
	Name() string
	Get(int) (Row, error)
	Insert(columns ...interface{}) error
}

type table struct {
	name   string
	schema []*simpleDefinition
	data   []row
}

func (t *table) Name() string {
	return t.name
}

func (t *table) Get(i int) (Row, error) {
	if 0 <= i && i < len(t.data) {
		return t.data[i].Columns(), nil
	}
	return nil, errors.Errorf("out of index: %d", i)
}

func (t *table) makeRow(cols []interface{}) (row, error) {
	row := make([]*column, len(t.schema))

	for i, d := range t.schema {
		c, err := d.makeColumn(cols[i])
		if err != nil {
			return nil, errors.Wrapf(err, "cols[%d] is invalid", i)
		}
		row[i] = c
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

func (t *table) Schema() []ColumnDefinition {
	dx := make([]ColumnDefinition, len(t.schema))
	for i, d := range t.schema {
		dx[i] = d
	}
	return dx
}

func (t *table) ForEach(proc Iterator) {
	for _, row := range t.data {
		proc(row.Columns())
	}
}
