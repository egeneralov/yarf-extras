package auth

import (
	"net/http"
	"testing"
)

func TestNewToken(t *testing.T) {
	token := NewToken("someid", 5)
	if token == "" {
		t.Error("No token received after NewToken ")
	}
}

func TestGetToken(t *testing.T) {
	token := NewToken("someid", 5)

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Auth", token)

	token2 := GetToken(r)
	if token != token2 {
		t.Error("Token missmatch")
	}
}

func TestValidateToken(t *testing.T) {
	id := "someid"
	token := NewToken(id, 5)

	data, err := ValidateToken(token)
	if err != nil {
		t.Error(err.Error())
	}

	if data != id {
		t.Error("Token data missmatch")
	}
}

func TestDeleteToken(t *testing.T) {
	id := "someid"
	token := NewToken(id, 5)

	DeleteToken(token)

	data, err := ValidateToken(token)
	if err == nil {
		t.Error("Token still valid after delete: " + data)
	}
}
