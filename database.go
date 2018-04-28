package debora

import (
	"github.com/pkg/errors"
)

// Database a database
type Database interface {
	Tables() []Table
	CreateTable(name string, column ...string) (Table, error)
	Get(string) (Table, error)
}

type database struct {
	tables []*table
}

// CreateDatabase create a database
func CreateDatabase() Database {
	return &database{
		tables: []*table{},
	}
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

func parseColumns(columns []string) ([]*simpleDefinition, error) {
	cx := make([]*simpleDefinition, len(columns))
	for i, column := range columns {
		c, err := parseColumnDefinition(column)
		if err != nil {
			return nil, errors.Wrapf(err, "can't parse columns[%d]", i)
		}
		cx[i] = c
	}
	return cx, nil
}

func (d *database) CreateTable(name string, columns ...string) (Table, error) {
	if d.has(name) {
		return nil, errors.Errorf("the %s table is already exists", name)
	}

	if len(columns) == 0 {
		return nil, errors.Errorf("no columns")
	}

	schema, err := parseColumns(columns)
	if err != nil {
		return nil, err
	}

	t := &table{
		name:   name,
		schema: schema,
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
