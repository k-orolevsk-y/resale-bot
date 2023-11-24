package bot

import (
	"errors"
	"sync"
	"time"
)

type (
	Storage interface {
		Add(id string, data interface{}) error
		Get(id string) (interface{}, error)
		Delete(id string) error
	}

	memStorage struct {
		mx   sync.RWMutex
		data map[string]item
	}

	item struct {
		Value      interface{}
		Expiration int64
	}
)

func newMemStorage() *memStorage {
	storage := &memStorage{data: make(map[string]item)}
	go storage.GC()

	return storage
}

func (storage *memStorage) Add(id string, data interface{}) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.data[id] = item{Value: data, Expiration: time.Now().Add(60 * 60 * 24).UnixNano()}
	return nil
}

func (storage *memStorage) Get(id string) (interface{}, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	data, ok := storage.data[id]
	if !ok {
		return nil, errors.New("invalid id")
	}

	return data.Value, nil
}

func (storage *memStorage) Delete(id string) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	_, ok := storage.data[id]
	if !ok {
		return errors.New("invalid id")
	}

	delete(storage.data, id)
	return nil
}

func (storage *memStorage) GC() {
	for {
		<-time.After(time.Minute * time.Duration(10))

		if storage.data == nil {
			return
		}

		storage.mx.Lock()
		if keys := storage.expiredKeys(); len(keys) != 0 {
			storage.clearItems(keys)
		}
		storage.mx.Unlock()

	}
}

func (storage *memStorage) expiredKeys() (keys []string) {
	for k, i := range storage.data {
		if time.Now().UnixNano() > i.Expiration {
			keys = append(keys, k)
		}
	}

	return
}

func (storage *memStorage) clearItems(keys []string) {
	for _, k := range keys {
		delete(storage.data, k)
	}
}
