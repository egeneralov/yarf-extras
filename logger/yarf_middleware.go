package logger

import (
	"github.com/yarf-framework/yarf"
	"log"
)

// Logger middleware it's a simple log module that uses the default golang's log package.
// The log output writer can be defined by default with the log.SetOutput(w io.Writer) function.
// For more complex environments where a default logger can't be used across the system,
// a custom solution to replace this should be implemented.
type Logger struct {
	yarf.Middleware
}

// PreDispatch wraps the http.ResponseWriter with a new LoggerWritter
// so we can log information about the response.
func (l *Logger) PreDispatch(c *yarf.Context) error {
	c.Response = &LoggerWriter{
		Writer: c.Response,
	}

	return nil
}

func (l *Logger) End(c *yarf.Context) error {
	// If nobody sets the status code, it's a 200
	var code int
	if _, ok := c.Response.(*LoggerWriter); ok {
		code = c.Response.(*LoggerWriter).StatusCode
	}

	if code == 0 {
		code = 200
	}

	log.Printf(
		"| %s | %s | %d | %s",
		c.GetClientIP(),
		c.Request.Method,
		code,
		c.Request.URL.String(),
	)

	return nil
}
