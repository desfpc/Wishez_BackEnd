package db

import (
	"testing"
)

func TestDb(t *testing.T) {

	//Wrong connection settings
	var wrongDbres = Db("mysql","rootz:rootz@/wishez")
	query := "SELECT count(id) cid FROM users"
	_, err := wrongDbres.Query(query)
	if err == nil {
		t.Error("Expected DB error, but not found")
	}
	wrongDbres.Close()

	//Default connection settings
	var dbres = Db("","")
	query = "SELECT count(id) cid FROM users"
	_, err = dbres.Query(query)

	if err != nil {
		t.Error("Unexpected DB errors")
	}
	dbres.Close()
}