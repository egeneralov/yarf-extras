package ratelimit

import (
	"sync"
	"time"
)

// RateLimitError indicates that the limit has been reached for the actual window.
type RateLimitError struct{}

func (e RateLimitError) Error() string {
	return "Rate limit exceeded"
}

// Rate represents a single client count for a given event.
// Its used internally by RateLimit to count different keys.
type Rate struct {
	// Amount of events allowed in a given window
	Limit int64

	// Time window in seconds to limit
	Window int64

	// Actual count
	EventCount int64

	// Start time
	Start time.Time

	// Sync Mutex
	sync.RWMutex
}

// Count checks for the actual limit/window and resets the window (Start) and Count when corresponding.
// If the limit has been reached for the actual window it returns a RateLimitError error.
func (r *Rate) Count() error {
	r.Lock()
	defer r.Unlock()

	// Reset window
	if time.Now().After(r.Start.Add(time.Second * time.Duration(r.Window))) {
		r.Start = time.Now()
		r.EventCount = 0
	}

	// Count
	r.EventCount++

	// Block
	if r.EventCount >= r.Limit {
		return RateLimitError{}
	}

	// Continue
	return nil
}

type RateLimit struct {
	// Amount of events allowed in a given window
	Limit int64

	// Time window in seconds to limit
	Window int64

	// Count storage
	counter map[string]*Rate

	// Sync Mutex
	sync.RWMutex

	// Garbage collection flag
	gcActive bool
}

func New(limit, window int64) *RateLimit {
	return &RateLimit{
		Limit:    limit,
		Window:   window,
		counter:  make(map[string]*Rate),
		gcActive: false,
	}
}

func (rl *RateLimit) Get(key string) *Rate {
	rl.RLock()
	defer rl.RUnlock()

	if _, ok := rl.counter[key]; !ok {
		// Default new Rate
		return &Rate{
			Limit:      rl.Limit,
			Window:     rl.Window,
			EventCount: 0,
			Start:      time.Now(),
		}
	}

	return rl.counter[key]
}

// Count checks on a given key, for the actual limit/window and resets the window (Start) and Count when corresponding.
// If the limit has been reached for the actual window it returns a RateLimitError error.
func (rl *RateLimit) Count(key string) error {
	rl.Lock()
	defer rl.Unlock()

	// Init counter
	if rl.counter == nil {
		rl.counter = make(map[string]*Rate)
	}

	// Init garbage collector
	if !rl.gcActive {
		rl.gcActive = true
		go rl.gc()
	}

	// Init rate count for the given key
	if _, ok := rl.counter[key]; !ok {
		rl.counter[key] = &Rate{
			Limit:      rl.Limit,
			Window:     rl.Window,
			EventCount: 0,
			Start:      time.Now(),
		}
	}

	return rl.counter[key].Count()
}

func (rl *RateLimit) gc() {
	// 10 times the window + 1
	t := time.NewTicker(time.Duration((rl.Window+1)*10) * time.Second)
	defer func(t *time.Ticker, rl *RateLimit) {
		t.Stop()
		rl.gcActive = false
	}(t, rl)

	for _ = range t.C {
		// Cancel when storage not present
		if rl.counter == nil {
			return
		}

		// Check for expired entries.
		now := time.Now()

		// Write lock
		rl.Lock()
		for key, rate := range rl.counter {
			// Expired 2 windows ago.
			if now.After(rate.Start.Add(time.Second * time.Duration(rate.Window*2))) {
				delete(rl.counter, key)
			}
		}
		rl.Unlock()
	}
}
