package user

import (
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
		t.Error("Token authorize error")
	}
	if expireError {
		t.Error("Token expire error")
	}

	wrongToken := "eklmldkmfvldkmfvlkdfvkdl;fkvmldkfmvkldmfklvmdklfmvkldmfklvmdklfmkvlkmdlflkdfv="
	_, authorizeError, _ = GetAuthorization(wrongToken, "access")
	if !authorizeError {
		t.Error("No token authorize error, but it's is")
	}

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJ1c2VyX2lkIjoiNSIsImV4cCI6IjE2MjkyNjk4NDkiLCJraW5kIjoiYWNjZXNzIn2Bc3BHSHZp8Jgch/Cfb8E9uA6vnKYayiBa1t7iVN40VQ=="
	_, _, expireError = GetAuthorization(expiredToken, "access")
	if !expireError {
		t.Error("No token expire error, but it's is")
	}

	wrongSignatureToken := ""
}