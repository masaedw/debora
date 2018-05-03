package debora

import (
	"testing"
)

func TestFromNil(t *testing.T) {
	query := From(nil)
	schema := query.Schema()
	if len(schema) != 0 {
		t.Fatal("schema is unexpectedly not empty")
	}

	slice, err := query.Slice()
	if err != nil {
		t.Fatal(err)
	}

	if len(slice) != 0 {
		t.Fatal("slice is unexpectedly not empty")
	}
}

func TestFromTable(t *testing.T) {
	db := CreateDatabase()
	table, _ := db.CreateTable("users", "id:integer", "name:string")
	tests := []struct {
		id   int64
		name string
	}{
		{1, "user1"},
		{2, "user2"},
		{3, "user3"},
	}
	for _, t := range tests {
		table.Insert(t.id, t.name)
	}
	testColumns := 2

	query := From(table)
	schema := query.Schema()
	if len(schema) != testColumns {
		t.Fatal("schema is not same as source table")
	}

	sourceSchema := query.Schema()
	for i, c := range sourceSchema {
		if c.ColumnType() != schema[i].ColumnType() || c.Name() != schema[i].Name() {
			t.Errorf("expected %v but got %v", c, schema[i])
		}
	}

	slice, err := query.Slice()
	if err != nil {
		t.Fatal(err)
	}

	if len(slice) != len(tests) {
		t.Error("the len of Slice() is not same as the source table's")
	}

	for i, r := range slice {
		if len(r) != testColumns {
			t.Errorf("invalid len of slice[%d]", i)
		}

		if tests[i].id != *r[0].IntegerValue() {
			t.Errorf("[%d].id: expected %d but got %d", i, tests[i].id, *r[0].IntegerValue())
		}

		if tests[i].name != *r[1].StringValue() {
			t.Errorf("[%d].name: expected %s but got %s", i, tests[i].name, *r[1].StringValue())
		}
	}
}
