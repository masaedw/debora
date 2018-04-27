package debora

import (
	"testing"
)

var columnTests = []struct {
	pattern    string
	columnType ColumnType
	name       string
}{
	{"id:integer", Integer, "id"},
	{"name:string", String, "name"},
	{"default", String, "default"},
}

var columnErrors = []string {
	"id:int",
	"id:",
	"hoge:fuga",
	"hoge:s",
}

func TestParseColumnDefinition(t *testing.T) {
	for _, p := range columnTests {
		c, err := parseColumnDefinition(p.pattern)
		if err != nil {
			t.Error(err)
		}
		if c.Type != p.columnType {
			t.Errorf("the type %v expected but got %v from %s", p.columnType, c.Type, p.pattern)
		}

		if c.Name != p.name {
			t.Errorf("name %v expected but got %v from %s", p.name, c.Name, p.pattern)
		}
	}

	for _, p := range columnErrors {
		c, err := parseColumnDefinition(p)
		if err == nil {
			t.Errorf("column definition %v should be error but got %v", p, *c)
		}
	}
}
