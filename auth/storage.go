package auth

import (
    "sync"
    "sync/atomic"
    "time"
)

// Storage interface is used to register any custom storage system for auth module. 
type Storage interface {
    // Get returns the data for a given key or an error if the key isn't valid.
    Get(key string) (string, error)
    
    // Set stores the data for a key for a given duration in seconds.
    // Returns error if it fails. 
    Set(key, data string, duration int) error
    
    // Refresh extends the expiration of a key by the same time it had when it was created. 
    Refresh(key string) error
    
    // Del removes the data and invalidates a key.
    // Returns error if it fails.
    Del(key string) error
}

// authToken is the storage unit used for auth module. 
type authToken struct {
	data       string    // Any data you want to save for this token.
	duration   int       // Seconds. Stored to be used by the RefreshToken function.
	expiration time.Time // Expiration time calculated after duration
}

// authStorage is the internal implementation for Storage interface. 
// It uses a in-memory map to store auth data.
type authStorage struct {
	// Data store
	store map[string]authToken
    
    // Garbage collector running? 
    gcFlag int64
    
	// Sync Mutex
	sync.RWMutex
}

// authStorage's garbage collector
func (as *authStorage) gc() {
    // Set running flag
    atomic.StoreInt64(&as.gcFlag, 1)
    
	// Run every minute.
	t := time.NewTicker(1 * time.Minute)

	for _ = range t.C {
		// Cancel when storage not present
		if as.store == nil {
			t.Stop()
			return
		}

		// Check for expired storage entries.
		now := time.Now()
        
        // Read lock
		as.RLock()
		for key, data := range as.store {
			if now.After(data.expiration) {
			    // Switch form read lock to full lock to delete expired data.
			    as.RUnlock()
			    as.Lock()
				
				// Delete
				as.Del(key)
				
				// Switch lock back.
				as.Unlock()
				as.RLock()
			}
		}
		as.RUnlock()
	}
}

// Get data from storage
func (as *authStorage) Get(key string) (string, error) {
    as.RLock()
	defer as.RUnlock()

    // Return data if available
	if data, ok := as.store[key]; ok {
		if data.expiration.After(time.Now()) {
			return data.data, nil
		}
	}
	
	// Key not found
	return "", InvalidKeyError{}
}

// Set data to storage.
func (as *authStorage) Set(key, data string, duration int) error {
	// Calculate expiration time
	exp := time.Now().Add(time.Duration(duration) * time.Second)
	
	as.Lock()
	defer as.Unlock()
	
	// Save data
	as.store[key] = authToken{data: data, duration: duration, expiration: exp}
	
	// Init GC if not running yet. 
	// Write lock comes handy here.
	if atomic.LoadInt64(&as.gcFlag) == 0 {
	    go as.gc()
	}
	
	return nil
}

// Refresh expiration
func (as *authStorage) Refresh(key string) error {
    as.Lock()
	defer as.Unlock()
	
	if data, ok := as.store[key]; ok {
	    // Validate expiration. Expired tokens can't be refreshed.
		if data.expiration.After(time.Now()) {
			data.expiration = time.Now().Add(time.Duration(data.duration) * time.Second)
			as.store[key] = data
		}
	}
	
	return nil
}

// Delete data to storage.
func (as *authStorage) Del(key string) error {
	as.Lock()
	defer as.Unlock()
	
	delete(as.store, key)
	
	return nil
}
