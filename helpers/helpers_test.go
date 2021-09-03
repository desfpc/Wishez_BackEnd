package helpers

import (
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