package storage

import (
	"fmt"
	"sync"
)

type storage struct {
	lock sync.RWMutex
	data map[string]uint
}

func newStorage() *storage {
	return &storage{
		data: make(map[string]uint),
	}
}

func (s *storage) Set(key string, value uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = value
	return nil
}

func (s *storage) Delete(key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}

func (s *storage) Get(key string) (uint, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return 0, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (s *storage) Increment(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] += 1
}
