package auth

import (
	"github.com/yarf-framework/yarf"
)

// Auth middleware performs auth on pre-dispatch after a token expected on the request.
// It also provides methods to generate and validate the tokens, that can be used by clients to perform authentication and authorization.
type Auth struct {
	yarf.Middleware
}

// PreDispatch checks if a token has been sent on the request, either by cookie or Auth header.
// If the token is invalid or non-present, it will return an error to stop execution of the following resources.
// If a token is valid, it returns its data on the "Auth" index of the yarf.Context.Data object.
func (a *Auth) PreDispatch(c *yarf.Context) error {
	token := GetToken(c.Request)
	
	data, err := ValidateToken(token)
	if err != nil {
		return new(UnauthorizedError)
	}

	c.Data.Set("_authData", data)
	c.Data.Set("_authToken", token)

	// Refresh token expiration on every request.
	RefreshToken(token)

	return nil
}
