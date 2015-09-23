package logger

import (
	"net/http"
)

// LoggerWriter will replace (wrap) the http.ResponseWriter to log all content written to the response.
type LoggerWriter struct {
	StatusCode int
	Writer     http.ResponseWriter
}

// Header is a wrapper for http.ResponseWriter.Header()
func (lw *LoggerWriter) Header() http.Header {
	return lw.Writer.Header()
}

// WriteHeader is a wrapper for http.ResponseWriter.WriteHeader()
// It saves the status code to be returned so we can log it.
func (lw *LoggerWriter) WriteHeader(code int) {
	lw.StatusCode = code

	lw.Writer.WriteHeader(code)
}

// Write is a wrapper for
func (lw *LoggerWriter) Write(content []byte) (int, error) {
	return lw.Writer.Write(content)
}
