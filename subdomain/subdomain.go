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
	err := c.Data.Set("_subdomain", strings.Split(c.Request.Host, ".")[0])
	if err != nil {
	    c.Render("SUBDOMAIN SET ERROR: " + err.Error())
	}
	
	str, err := c.Data.Get("_subdomain")
	if err != nil {
	    c.Render("SUBDOMAIN ERROR: " + err.Error())
	}
	
	c.Render("SUBDOMAIN: " + str.(string))

	return nil
}
