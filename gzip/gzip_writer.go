package gzip

import (
	"compress/gzip"
	"net/http"
)

// GzipWriter will replace (wrap) the http.ResponseWriter to gzip all content written to the response.
// It implements the http.ResponseWriter interface.
type GzipWriter struct {
	Writer http.ResponseWriter
}

// Header is a wrapper for http.ResponseWriter.Header()
func (gzw *GzipWriter) Header() http.Header {
	return gzw.Writer.Header()
}

// WriteHeader is a wrapper for http.ResponseWriter.WriteHeader()
func (gzw *GzipWriter) WriteHeader(code int) {
	gzw.Writer.WriteHeader(code)
}

// Write compress the content received and writes it to the client through the http.ResponseWriter
func (gzw *GzipWriter) Write(content []byte) (int, error) {
	// Create writer
	gz := gzip.NewWriter(gzw.Writer)
	defer gz.Close()

	// Write gzip bytes
	return gz.Write(content)
}
