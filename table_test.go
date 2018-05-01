package debora

import (
	"testing"
)

func TestInsert(t *testing.T) {
	db := CreateDatabase()
	users, _ := db.CreateTable("users", "id:integer", "name:string")

	data := []struct {
		id   int64
		name string
	}{
		{1, "user1"},
		{2, "user2"},
		{3, "user3"},
	}

	for _, u := range data {
		if err := users.Insert(u.id, u.name); err != nil {
			t.Error(err, "failed to insert", u.name)
		}
	}

	for i, expected := range data {
		actual, err := users.Get(i)
		if err != nil {
			t.Error(err, "failed to get", expected.name)
		}

		if *actual[0].IntegerValue() != expected.id {
			t.Errorf("expected %d but got %d", expected.id, *actual[0].IntegerValue())
		}
		if *actual[1].StringValue() != expected.name {
			t.Errorf("expected %s but got %s", expected.name, *actual[1].StringValue())
		}
	}
}
