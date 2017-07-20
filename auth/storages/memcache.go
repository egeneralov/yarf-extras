package storages

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/yarf-framework/extras/auth"
)

type memcacheStorage struct {
	client *memcache.Client
}

func Memcache(servers ...string) auth.Storage {
	ms := new(memcacheStorage)
	ms.client = memcache.New(servers...)

	return ms
}

// Get data from storage
func (ms *memcacheStorage) Get(key string) (val string, err error) {
	item, err := ms.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			err = auth.InvalidKeyError{}
		}
		return
	}

	return string(item.Value), nil
}

// Set data to storage.
func (ms *memcacheStorage) Set(key, data string, duration int) error {
	return ms.client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(data),
		Expiration: int32(duration),
	})
}

// Refresh expiration
func (ms *memcacheStorage) Refresh(key string) error {
	item, err := ms.client.Get(key)
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
func (ms *memcacheStorage) Del(key string) error {
	return ms.client.Delete(key)
}
