package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var s *authStorage

type authStorage struct {
	// Data store
	store map[string]authToken

	// Sync Mutex
	sync.RWMutex
}

type authToken struct {
	data       string    // Any data you want to save for this token.
	duration   int       // Seconds. Stored to be used by the RefreshToken function.
	expiration time.Time // Expiration time calculated after duration
}

func initStorage() {
	if s == nil {
		s = new(authStorage)
		s.store = make(map[string]authToken)

		go garbageCollector()
	}
}

func garbageCollector() {
	// Run every 5 minutes.
	t := time.NewTicker(5 * time.Minute)

	for _ = range t.C {
		// Cancel when storage not present
		if s == nil {
			t.Stop()
			return
		}

		// Check for expired storage entries.
		now := time.Now()

		s.Lock()
		for token, data := range s.store {
			if data.expiration.After(now) {
				delete(s.store, token)
			}
		}
		s.Unlock()
	}
}

func generateToken() string {
	initStorage()

	// Generate random bytes
	b := make([]byte, 256)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	// SHA512 -> string
	h := sha512.Sum512(b)

	return fmt.Sprintf("%x", h) // Encode the right UTF-8 bytes.
}

// NewToken creates and stores a new token on the local storage.
// It associates a token to an id so it can be identified and returned by the ValidateToken method.
// It should be used from a Login method after a successful authentication.
func NewToken(id string, d int) string {
	initStorage()

	t := generateToken()

	// Check non-existence of a new token
	s.Lock()
	defer s.Unlock()

	if _, ok := s.store[t]; ok {
		for ok {
			t = generateToken()
			_, ok = s.store[t]
		}
	}

	// Calculate expiration
	exp := time.Now().Add(time.Duration(d) * time.Second)

	// Save token and id
	s.store[t] = authToken{data: id, duration: d, expiration: exp}

	return t
}

// GetToken tries to retrieve the token from the request object.
// If the token is not found, returns an empty string.
func GetToken(r *http.Request) string {
	var token string

	// First try the cookie method
	cookie, err := r.Cookie("Auth")
	if err != nil {
		// Try header
		token = r.Header.Get("Auth")
	} else {
		token = cookie.Value
	}

	return token
}

// ValidateToken checks if a token is valid and returns the data contained on it.
// Otherwise it will return an error status together with an empty string.
func ValidateToken(token string) (string, error) {
	initStorage()

	s.RLock()
	defer s.RUnlock()

	if data, ok := s.store[token]; ok {
		if data.expiration.After(time.Now()) {
			return data.data, nil
		}
	}

	return "", errors.New("Invalid token")
}

// RefreshToken resets the timer of the token to extend its valid status.
// It sets the same expiration time as when it was created, but starting now.
func RefreshToken(token string) {
	initStorage()

	s.Lock()
	defer s.Unlock()

	if t, ok := s.store[token]; ok {
		if t.expiration.After(time.Now()) {
			t.expiration = time.Now().Add(time.Duration(s.store[token].duration) * time.Second)
			s.store[token] = t
		}
	}
}

// DeleteToken removes the token data from the storage.
func DeleteToken(token string) {
	initStorage()

	s.Lock()
	defer s.Unlock()

	delete(s.store, token)
}
