package helpers

import (
	"database/sql"
	"errors"
	types "github.com/desfpc/Wishez_Type"
	"testing"
)

func TestAuthErrorAnswer(t *testing.T) {
	var errors types.Errors
	var code int

	errors, code = AuthErrorAnswer(true, true)
	if code != 401 {
		t.Error("Expected 401, got ", code)
	}
	if errors[0] != "Authorization Required" {
		t.Error("Expected Authorization Required error, got ", errors[0])
	}

	errors, code = AuthErrorAnswer(false, true)
	if code != 401 {
		t.Error("Expected 401, got ", code)
	}
	if errors[0] != "Access Token is Expired" {
		t.Error("Expected Access Token is Expired error, got ", errors[0])
	}

	errors, code = AuthErrorAnswer(false, false)
	if code != 200 {
		t.Error("Expected 200, got ", code)
	}
	if len(errors) > 0 {
		t.Error("Expected 0 errors, got ", len(errors))
	}
}

func TestNoRouteErrorAnswer(t *testing.T) {
	var errors types.Errors
	var code int

	errors, code = NoRouteErrorAnswer()
	if code != 404 {
		t.Error("Expected 404, got ", code)
	}
	if errors[0] != "Entity and/or action not found" {
		t.Error("Expected Entity and/or action not found error, got ", errors[0])
	}
}

func TestMakeStringFromIntSQL(t *testing.T) {
	var sqlNullInt = sql.NullInt64{Int64: 1000, Valid: true}
	var str = MakeStringFromIntSQL(sqlNullInt)
	if str != "1000" {
		t.Error("Expected 1000, got ", str)
	}

	sqlNullInt = sql.NullInt64{Int64: 1000, Valid: false}
	str = MakeStringFromIntSQL(sqlNullInt)
	if str != "" {
		t.Error("Expected '', got ", str)
	}
}

func TestMakeStringFromSQL(t *testing.T) {
	var sqlNullString = sql.NullString{String: "Test String", Valid: true}
	var str = MakeStringFromSQL(sqlNullString)

	if str != "Test String" {
		t.Error("Expected 'Test String', got ", str)
	}

	sqlNullString = sql.NullString{String: "Test String", Valid: false}
	str = MakeStringFromSQL(sqlNullString)
	if str != "" {
		t.Error("Expected '', got ", str)
	}
}

func TestCheckErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	CheckErr(nil)

	var err = errors.New("this is error message")
	CheckErr(err)
}

func TestIsEmailValid(t *testing.T) {
	if !IsEmailValid("desfpc@gmail.com") {
		t.Errorf("Error in valid email")
	}
	if IsEmailValid("desfpcgmail.com") {
		t.Errorf("Valid wrong email")
	}
	if IsEmailValid("d@m") {
		t.Errorf("Valid wrong email")
	}
	if IsEmailValid("dm") {
		t.Errorf("Valid wrong email")
	}
}