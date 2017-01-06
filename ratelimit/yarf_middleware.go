package ratelimit

import (
	"github.com/yarf-framework/yarf"
	"strconv"
	"time"
)

// YarfError is the custom error type compatible with Yarf's YError
type YarfError struct{}

// Implements the error interface returning the ErrorMsg value of each error.
func (e *YarfError) Error() string {
	return "Unauthorized"
}

// Code returns the error's HTTP code to be used in the response.
func (e *YarfError) Code() int {
	return 429 // Too Many Requests
}

// ID returns the error's ID for further reference.
func (e *YarfError) ID() int {
	return 429
}

// Msg returns the error's message, used to implement the Error interface.
func (e *YarfError) Msg() string {
	return "Too Many Requests"
}

// Body returns the error's content body, if needed, to be returned in the HTTP response.
func (e *YarfError) Body() string {
	return "Too Many Requests: Try again later."
}

// RateLimiter middleware provides request rate limits per IP
type RateLimiter struct {
	yarf.Middleware

	// rate limiter
	rl *RateLimit
}

// YarfMiddleware constructor receives the requests limit and a time window (in seconds) to allow. 
// Any IP that requests more than the limit within the time window will be blocked until the time window ends and a new one starts.  
func YarfMiddleware(limit, window int) *RateLimiter {
	return &RateLimiter{
		rl: New(limit, window),
	}
}

// PreDispatch performs the requests counting and handle blocks/ 
func (m *RateLimiter) PreDispatch(c *yarf.Context) error {
	// IP as key
	key := c.GetClientIP()

	// Count
	err := m.rl.Count(key)
	if err != nil {
		if _, ok := err.(RateLimitError); ok {
			return &YarfError{}
		}

		return err
	}

	// Set rate limit info on headers
	rate := m.rl.Get(key)
	c.Response.Header().Set("X-RateLimit-Limit", strconv.Itoa(rate.Limit))
	c.Response.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rate.Limit-rate.EventCount))
	c.Response.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(rate.Start.Add(time.Second*time.Duration(rate.Window)).Unix())))

	return nil
}
