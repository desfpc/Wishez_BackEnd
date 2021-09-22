package db

import (
	"strconv"
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

func TestCheckCount(t *testing.T) {
	var dbres = Db("","")
	query := "SELECT count(id) cnt FROM users"
	rows, err := dbres.Query(query)

	if err != nil {
		t.Error("Unexpected DB errors")
	}
	var count = CheckCount(rows)

	query = "SELECT count(id) cnt FROM users"
	qCountRes, err := dbres.Query(query)

	var qCount int
	for qCountRes.Next() {
		_ = qCountRes.Scan(&qCount)
	}

	if count != qCount {
		t.Error("Count " + strconv.Itoa(count) + " not equal " + strconv.Itoa(qCount))
	}
}

func TestEscape(t *testing.T) {
	var asciiCode int32
	asciiCode = 0
	testString := "Test string 0 ' \n \r \032 \\ 0 \" " + string(rune(asciiCode))
	testString2 := "Test string 0 \\' \\n \\r \\Z \\\\ 0 \\\" \\0"
	escapedString := Escape(testString)
	if escapedString != testString2 {
		t.Error("Wrong escaped string: " + escapedString + ". Needed: " + testString2 + ".")
	}
}