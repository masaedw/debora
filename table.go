package debora

import (
	"github.com/pkg/errors"
)

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

type table struct {
	name   string
	schema []*ColumnDefinition
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
