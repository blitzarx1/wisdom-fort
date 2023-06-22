package storage

import (
	"context"
	"time"

	"blitzarx1/wisdom-fort/pkg/logger"
)

type entry struct {
	id  ID
	key string
}

type ID int

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
	stores []keyvalStore

	withTTL    map[ID]time.Duration
	expiration map[time.Time]entry
}

func New(ctx context.Context) *Service {
	l := logger.MustFromCtx(ctx)
	l.Println("initializing storage service")

	s := &Service{
		stores:     make([]keyvalStore, 0),
		withTTL:    make(map[ID]time.Duration),
		expiration: make(map[time.Time]entry),
	}

	// clear expired keys with check interval of 1 second.
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for t, e := range s.expiration {
					if time.Now().After(t) {
						s.Delete(e.id, e.key)
						delete(s.expiration, t)
					}
				}
			}
		}
	}()

	return s
}

func (s *Service) AddStorageWithTTL(ctx context.Context, ttl time.Duration) ID {
	logger.MustFromCtx(ctx).Println("adding new storage with ttl")

	id := s.addStore()
	s.withTTL[id] = ttl

	return id
}

func (s *Service) Set(id ID, key string, value uint) {
	if ttl, ok := s.withTTL[id]; ok {
		s.expiration[time.Now().Add(ttl)] = entry{id: id, key: key}
	}
	s.stores[id].Set(key, value)
}

func (s *Service) Increment(id ID, key string) {
	if ttl, ok := s.withTTL[id]; ok {
		s.expiration[time.Now().Add(ttl)] = entry{id: id, key: key}
	}
	s.stores[id].Increment(key)
}

func (s *Service) Get(id ID, key string) (uint, error) {
	return s.stores[id].Get(key)
}

func (s *Service) Delete(id ID, key string) {
	s.stores[id].Delete(key)
}

func (s *Service) addStore() ID {
	s.stores = append(s.stores, newStorage())
	return ID(len(s.stores) - 1)
}
