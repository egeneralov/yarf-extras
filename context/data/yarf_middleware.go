package data

import (
	"github.com/yarf-framework/yarf"
)

// SetStrData middleware creates a new yarf.Context.Data object based on StrData type.
type SetStrData struct {
	yarf.Middleware
}

func (m *SetStrData) PreDispatch(c *yarf.Context) error {
	c.Data = new(StrData)

	return nil
}
