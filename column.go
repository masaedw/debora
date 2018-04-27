package debora

import (
	"reflect"
	"regexp"

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

var columnPattern = regexp.MustCompile(`^([^:]+)(:(string|integer))?$`)

func parseColumnDefinition(column string) (*ColumnDefinition, error) {
	m := columnPattern.FindStringSubmatch(column)
	if m == nil {
		return nil, errors.Errorf("invalid column name: %v", column)
	}

	c := &ColumnDefinition{
		Name: m[1],
		Type: String,
	}

	if m[3] == "integer" {
		c.Type = Integer
	}

	return c, nil
}

func (d *ColumnDefinition) makeColumn(col interface{}) (*column, error) {
	switch d.Type {
	case String:
		v, ok := col.(string)
		if !ok {
			return nil, errors.Errorf("requried string but got %s", reflect.TypeOf(col))
		}
		return &column{definition: d, stringValue: &v}, nil

	case Integer:
		v, ok := col.(int64)
		if !ok {
			return nil, errors.Errorf("requried int64 but got %s", reflect.TypeOf(col))
		}
		return &column{definition: d, integerValue: &v}, nil

	default:
		return nil, errors.Errorf("col is unsupported type: %s", reflect.TypeOf(col))
	}
}
