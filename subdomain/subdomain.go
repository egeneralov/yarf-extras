package subdomain

import (
	"github.com/yarf-framework/yarf"
	"strings"
)

// Subdomain middleware parses the Request URL
// and sets the subdomain value into "_subdomain" index on Context Data
type Subdomain struct {
	yarf.Middleware
}

// PreDispatch parses the URL and sets the value.
func (s *Subdomain) PreDispatch(c *yarf.Context) error {
	// One-liner
	c.Data.Set("_subdomain", strings.Split(c.Request.Host, ".")[0])

	return nil
}
