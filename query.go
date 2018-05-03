package debora

import (
	"github.com/pkg/errors"
)

type whereQuery struct {
	table     TableLike
	condition WhereClause
}

// From create query
func From(t TableLike) Queryable {
	return &whereQuery{
		table:     t,
		condition: func(Row) bool { return true },
	}
}

func (w *whereQuery) Schema() []ColumnDefinition {
	return w.table.Schema()
}

func (w *whereQuery) ForEach(proc Iterator) {
	w.ForEach(func(row Row) {
		if w.condition(row) {
			proc(row)
		}
	})
}

func (w *whereQuery) Where(cond WhereClause) Queryable {
	return &whereQuery{
		table:     w,
		condition: cond,
	}
}

func (w *whereQuery) Select(sel ...ColumnSelector) Queryable {
	// TODO
	return nil
}

func (w *whereQuery) Slice() ([]Row, error) {
	rows := make([]Row, 0)
	w.ForEach(func(r Row) {
		rows = append(rows, r)
	})
	return rows, nil
}

func (w *whereQuery) First() (Row, error) {
	slice, err := w.Slice()
	if err != nil {
		return nil, errors.Wrapf(err, "can't get slice")
	}

	if len(slice) > 0 {
		return slice[0], nil
	}

	return nil, errors.Errorf("no contents")
}
