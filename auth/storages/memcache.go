package storages

import (
	"bytes"
	"encoding/binary"
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

func expirationKey(k string) string {
	return keyPrefix + ":expiration:" + k
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
	// Save original expiration to refresh later
	go ms.setExpiration(k, duration)

	// Now set data
	return ms.client.Set(&memcache.Item{
		Key:        key(k),
		Value:      []byte(data),
		Expiration: int32(duration),
	})
}

func (ms *memcacheStorage) getExpiration(k string) (duration int) {
	item, err := ms.client.Get(key(k))
	if err != nil {
		return
	}

	d, _ := binary.Varint(item.Value)
	duration = int(d)

	return
}

func (ms *memcacheStorage) setExpiration(k string, duration int) {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, duration)

	//binary.PutVarint(data, int64(duration))

	ms.client.Set(&memcache.Item{
		Key:        expirationKey(k),
		Value:      data.Bytes(),
		Expiration: int32(duration * 2),
	})
}

// Refresh expiration
func (ms *memcacheStorage) Refresh(k string) (err error) {
	d := ms.getExpiration(expirationKey(k))
	if d != 0 {
		go ms.client.Touch(expirationKey(k), int32(d))

		err = ms.client.Touch(key(k), int32(d))
	}

	return
}

// Delete data to storage.
func (ms *memcacheStorage) Del(k string) error {
	return ms.client.Delete(key(k))
}
