package storages

import (
	"github.com/yarf-framework/extras/auth"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	mc := Memcache("localhost:11211")
	auth.RegisterStorage(mc)
}

func TestNewToken(t *testing.T) {
	token := auth.NewToken("someid", 5)
	if token == "" {
		t.Error("No token received after NewToken ")
	}
}

func TestGetToken(t *testing.T) {
	token := auth.NewToken("someid", 5)

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Auth", token)

	token2 := auth.GetToken(r)
	if token != token2 {
		t.Error("Token missmatch")
	}
}

func TestValidateToken(t *testing.T) {
	id := "someid"
	token := auth.NewToken(id, 5)

	data, err := auth.ValidateToken(token)
	if err != nil {
		t.Error(err.Error())
	}

	if data != id {
		t.Error("Token data missmatch")
	}
}

func TestDeleteToken(t *testing.T) {
	id := "someid"
	token := auth.NewToken(id, 5)

	auth.DeleteToken(token)

	data, err := auth.ValidateToken(token)
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
			token := auth.NewToken(strconv.Itoa(int(time.Now().UnixNano())), 1)
			auth.ValidateToken(token)
			auth.RefreshToken(token)
			auth.DeleteToken(token)

			time.Sleep(5 * time.Second)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
