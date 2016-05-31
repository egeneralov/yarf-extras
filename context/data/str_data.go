package data

import (
	"net/url"
	"sync"
)

// StrData implements the yarf.ContextData interface to be used as a simple string storage.
// All interface{} values passed through this methods are treated as strings.
type StrData struct {
	data url.Values

	// Sync Mutex
	sync.RWMutex
}

// Get retrieves a data item by it's key name.
func (sd *StrData) Get(key string) (interface{}, error) {
	sd.Lock()
	defer sd.Unlock()

	if sd.data == nil {
		sd.data = url.Values{}
	}

	return sd.data.Get(key), nil
}

// Set saves a data item under a key name.
func (sd *StrData) Set(key string, data interface{}) error {
	sd.Lock()
	defer sd.Unlock()

	if sd.data == nil {
		sd.data = url.Values{}
	}

	sd.data.Set(key, data.(string))

	return nil
}

// Del removes the data item and key name for a given key.
func (sd *StrData) Del(key string) error {
	sd.Lock()
	defer sd.Unlock()

	if sd.data == nil {
		sd.data = url.Values{}
	}

	sd.data.Del(key)

	return nil
}
