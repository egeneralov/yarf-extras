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
func (m *Gzip) PreDispatch() error {
	// Check request header
	if !strings.Contains(m.Context.Request.Header.Get("Accept-Encoding"), "gzip") {
		return nil
	}

	// Set encoding header
	m.Context.Response.Header().Set("Content-Encoding", "gzip")

	// Wrap response writer
	m.Context.Response = &GzipWriter{
		Writer: m.Context.Response,
	}

	return nil
}
