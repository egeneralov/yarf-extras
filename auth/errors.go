package auth

import ()

// InvalidKeyError indicates that a key isn't present or that has expired so the data isn't available.
type InvalidKeyError struct{}

func (err InvalidKeyError) Error() string {
	return "Invalid key"
}
