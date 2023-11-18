package bot

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type (
	CallbackStorage interface {
		Add(interface{}) (uuid.UUID, error)
		Get(uuid.UUID) (interface{}, error)
		Delete(uuid.UUID) error
	}

	callbackMemStorage struct {
		mx   sync.Mutex
		data map[uuid.UUID]interface{}
	}
)

func newCallbackMemStorage() *callbackMemStorage {
	return &callbackMemStorage{data: make(map[uuid.UUID]interface{})}
}

func (storage *callbackMemStorage) Add(data interface{}) (uuid.UUID, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	id := uuid.New()
	storage.data[id] = data

	return id, nil
}

func (storage *callbackMemStorage) Get(id uuid.UUID) (interface{}, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	data, ok := storage.data[id]
	if !ok {
		return nil, errors.New("invalid id")
	}

	return data, nil
}

func (storage *callbackMemStorage) Delete(id uuid.UUID) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	_, ok := storage.data[id]
	if !ok {
		return errors.New("invalid id")
	}

	delete(storage.data, id)
	return nil
}
