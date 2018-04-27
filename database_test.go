package debora

import (
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	db := CreateDatabase()
	tables := db.Tables()

	if len(tables) != 0 {
		t.Error()
	}
}

func TestCreateTable(t *testing.T) {
	db := CreateDatabase()
	if _, err := db.CreateTable("users"); err == nil {
		t.Error("creating table with no column definition should be error")
	}

	if _, err := db.CreateTable("users", "id:fuga"); err == nil {
		t.Error("creating table with invalid column definition should be error")
	}

	_, err := db.CreateTable("users", "id:integer", "name:string", "email:string", "password:string")
	if err != nil {
		t.Error(err)
	}

	if _, err := db.CreateTable("users", "id:integer"); err == nil {
		t.Error("creating table with duplicated name should be error")
	}

	tx := db.Tables()
	if len(tx) != 1 {
		t.Error("unexpected table count")
	}

	if _, err := db.Get("users"); err != nil {
		t.Error("can't get users table")
	}

	if _, err := db.Get("undefined"); err == nil {
		t.Error("got undefined table unexpectedly")
	}

}
