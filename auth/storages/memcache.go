package storages

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/yarf-framework/extras/auth"
)

var (
	keyPrefix = "github.com/yarf-framework/extras/auth/storage/memcache:"
)

type memcacheStorage struct {
	client *memcache.Client
}

func Memcache(servers ...string) auth.Storage {
	ms := new(memcacheStorage)
	ms.client = memcache.New(servers...)

	return ms
}

func key(k string) string {
	return keyPrefix + k
}

// Get data from storage
func (ms *memcacheStorage) Get(k string) (val string, err error) {
	item, err := ms.client.Get(key(k))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = auth.InvalidKeyError{}
		}
		return
	}

	return string(item.Value), nil
}

// Set data to storage.
func (ms *memcacheStorage) Set(k, data string, duration int) error {
	return ms.client.Set(&memcache.Item{
		Key:        key(k),
		Value:      []byte(data),
		Expiration: int32(duration),
	})
}

// Refresh expiration
func (ms *memcacheStorage) Refresh(k string) error {
	item, err := ms.client.Get(key(k))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = auth.InvalidKeyError{}
		}
		return err
	}

	return ms.client.Set(&memcache.Item{
		Key:        item.Key,
		Value:      item.Value,
		Expiration: item.Expiration,
	})
}

// Delete data to storage.
func (ms *memcacheStorage) Del(k string) error {
	return ms.client.Delete(key(k))
}
