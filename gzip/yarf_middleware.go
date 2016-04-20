package gzip

import (
	"github.com/yarf-framework/yarf"
	"strings"
)

// Gzip middleware automatically handles gzip compressed responses to clients that accepts the encoding.
// It should be inserted at the beggining of the middleware stack so it can catch every write to the response and encode it right.
type Gzip struct {
	yarf.Middleware
}

// PreDispatch
func (m *Gzip) PreDispatch(c *yarf.Context) error {
	// Check request header
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		return nil
	}

	// Set encoding header
	c.Response.Header().Set("Content-Encoding", "gzip")

	// Wrap response writer
	c.Response = &GzipWriter{
		Writer: c.Response,
	}

	return nil
}
