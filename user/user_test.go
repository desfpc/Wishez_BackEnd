package user

import (
	types "github.com/desfpc/Wishez_Type"
	"testing"
)

func TestMakeToken(t *testing.T) {
	user := GetUserFromBD("1")
	token := MakeToken("access", user)
	dToken := deconcatToken(token)
	if dToken.Head != "{\"alg\":\"HS256\",\"typ\":\"JWT\"}" {
		t.Error("Wrong token head: " + dToken.Head)
	}
	if !checkToken(dToken) {
		t.Error("Wrong token: " + dToken.Body)
	}

	rtoken := MakeToken("refresh", user)
	if !CheckUserToken(rtoken) {
		t.Error("Wrong token: " + rtoken)
	}
}

func TestGetAuthorization(t *testing.T) {
	user := GetUserFromBD("1")
	token := MakeToken("access", user)
	systemUser, authorizeError, expireError := GetAuthorization(token, "access")
	if user != systemUser {
		t.Error("Wrong token user: " + systemUser.Email)
	}
	if authorizeError {
		t.Error("Wrong token authorize error")
	}
	if expireError {
		t.Error("Wrong token expire error")
	}

	wrongToken := "eklmldkmfvldkmfvlkdfvkdl;fkvmldkfmvkldmfklvmdklfmvkldmfklvmdklfmkvlkmdlflkdfv="
	_, authorizeError, _ = GetAuthorization(wrongToken, "access")
	if !authorizeError {
		t.Error("No token authorize error, but it's is")
	}

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJ1c2VyX2lkIjoiMSIsImV4cCI6IjE2MzI0MTE3NTciLCJraW5kIjoiYWNjZXNzIn16WnlVd3EyeEZZeXYzMDZGWFpTMkhNYlJyd2w5KzNIQWdDOTZmRFJZK0Q4PQ=="
	_, _, expireError = GetAuthorization(expiredToken, "access")
	if !expireError {
		t.Error("No token expire error, but it's is")
	}

	wrongKindToken := MakeToken("access", user)
	_, authorizeError, _ = GetAuthorization(wrongKindToken, "refresh")
	if !authorizeError {
		t.Error("No token authorize error when wrong token kind")
	}

	emptyExpToken := makeTokenFromStrings("{\"alg\":\"HS256\",\"typ\":\"JWT\"}", "{\"user_id\":\"1\",\"kind\":\"access\"}")
	_, authorizeError, expireError = GetAuthorization(emptyExpToken, "access")
	if !expireError {
		t.Error("No token expire error, but it's is")
	}

	wrongSignatureToken := makeTokenFromStringsVsSignature("{\"alg\":\"HS256\",\"typ\":\"JWT\"}", "{\"user_id\":\"1\",\"kind\":\"access\"}", "DFgbfgbffgb423as")
	_, authorizeError, _ = GetAuthorization(wrongSignatureToken, "access")
	if !authorizeError {
		t.Error("No token authorize error when wrong signature")
	}

}

func TestGetUserByID(t *testing.T) {
	var request = types.JsonRequest{
		Entity: "user",
		Id:     "",
		Action: "getById",
		Params: make(map[string]string),
	}
	request.Params["id"] = "1"
	_, err := getUserByID(request)
	if len(err) > 0 {
		t.Error("Errors when getting user request by ID")
	}

	request.Params["id"] = "-1"
	_, err = getUserByID(request)
	if len(err) == 0 {
		t.Error("No errors when getting user request by wrong ID")
	}

	request.Params = nil
	_, err = getUserByID(request)
	if len(err) == 0 {
		t.Error("No errors when getting user by wrong request")
	}

}