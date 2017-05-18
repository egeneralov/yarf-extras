package auth

import (
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
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

func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5000; i++ {
		// Count goroutines
		wg.Add(1)

		go func() {
			// Decrement goroutines count
			defer wg.Done()

			// Go full token process
			token := NewToken(strconv.Itoa(int(time.Now().UnixNano())), 1)
			ValidateToken(token)
			RefreshToken(token)
			DeleteToken(token)

			time.Sleep(5 * time.Second)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func BenchmarkNewToken(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewToken("data", 60)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	token := NewToken(strconv.Itoa(int(time.Now().UnixNano())), 3600)

	for n := 0; n < b.N; n++ {
		ValidateToken(token)
	}
}

func BenchmarkRefreshToken(b *testing.B) {
	token := NewToken(strconv.Itoa(int(time.Now().UnixNano())), 3600)

	for n := 0; n < b.N; n++ {
		RefreshToken(token)
	}
}
