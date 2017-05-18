package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"net/http"
)

// Internal Storage object
var store Storage

// Init storage using the internal package version.
// Can be overwritten by calling RegisterStorage(s Storage) at any time.
func init() {
	if store == nil {
		store = &authStorage{
			store: make(map[string]authToken),
		}
	}
}

// RegisterStorage replaces the default storage engine by a custom one.
// Replacing the storage means all data stored previously will be lost, so it should be done during initialization.
// Takes a Storage interface parameter and isn't safe for concurrent access.
func RegisterStorage(s Storage) {
	store = s
}

// generateToken creates a new random token long enough to avoid collisions.
func generateToken() string {
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

// NewToken creates and stores a new token on the local storage. It handles the token's uniqueness.
// It associates a token to the provided data so it can be identified and returned by the ValidateToken method.
// It should be used from a Login method after a successful authentication.
func NewToken(data string, d int) string {
	t := generateToken()

	// Check non-existence of the new token
	_, check := store.Get(t)
	for check == nil {
		t = generateToken()
		_, check = store.Get(t)
	}

	// Save data
	store.Set(t, data, d)

	return t
}

// GetToken tries to retrieve the token from the request object.
// It looks for the value of a request cookie named "Auth" first,
// and then for the value of a request header named "Auth" to retrieve the first value found:
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
	return store.Get(token)
}

// RefreshToken resets the timer of the token to extend its valid status.
// It sets the same duration time as when it was created, but starting now.
func RefreshToken(token string) {
	store.Refresh(token)
}

// DeleteToken removes the token data from the storage.
func DeleteToken(token string) {
	store.Del(token)
}
