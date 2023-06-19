package storage

import "log"

type StorageID int

// Service manages multiple key-value stores.
//
// It designed to be uses as a singleton. Method AddStore is not concurrent and suppossed
// to be used during initialization. It is safe to use other methods concurrently
// due to used storage implementation.
//
// It has a dependency on the keyvalStore interface. This allows us to easily
// switch the storage provider without having to change the service. It can be usefull
// when we need out service to scale and we need to use a storage like Redis.
type Service struct {
	logger *log.Logger

	stores []keyvalStore
}

func New(logger *log.Logger) *Service {
	logger.Println("initializing storage service")

	return &Service{
		logger: logger,
		stores: make([]keyvalStore, 0),
	}
}

func (s *Service) AddStore() StorageID {
	s.logger.Println("adding new storage")

	s.stores = append(s.stores, newStorage())
	return StorageID(len(s.stores) - 1)
}

func (s *Service) Set(id StorageID, key string, value uint) error {
	return s.stores[id].Set(key, value)
}

func (s *Service) Get(id StorageID, key string) (uint, error) {
	return s.stores[id].Get(key)
}

func (s *Service) Delete(id StorageID, key string) error {
	return s.stores[id].Delete(key)
}

func (s *Service) Increment(id StorageID, key string) {
	s.stores[id].Increment(key)
}

func (s *Service) Clear(id StorageID) {
	s.stores[id] = newStorage()
}
