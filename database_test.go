package debora

import (
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	db := CreateDatabase()
	tables := db.Tables()

	if len(tables) != 0 {
		t.Fatal()
	}
}