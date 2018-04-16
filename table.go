package debora

import (
	"reflect"

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
