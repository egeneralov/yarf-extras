package tpl

import (
	"io/ioutil"
	"sync"
)

// Tpl centralizes all package base functionality into a single type designed to be composited by others.
type Tpl struct {
	// Custom template path from where to render
	TplPath string

	// Custom template name used to build the complete path
	TplName string

	// cache
	cache map[string]string

	// Sync Mutex
	sync.RWMutex
}

// Cached reads the content of the path and returns the file's content or returns an error when fails.
// When successful, it saves the content on the cache and next time reads from there.
func (t *Tpl) Cached(path string) (string, error) {
	t.Lock()
	defer t.Unlock()

	if t.cache == nil {
		t.cache = make(map[string]string)
	}

	if _, ok := t.cache[path]; !ok {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		// Debug. Avoid caching
		return string(data), nil

		t.cache[path] = string(data)
	}

	return t.cache[path], nil
}
