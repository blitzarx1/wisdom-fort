package rps

import (
	"log"
	"time"

	"blitzarx1/wisdom-fort/server/service/storage"
)

// Service tracks rps per ip.
type Service struct {
	logger *log.Logger

	storageID storage.StorageID
	storage   *storage.Service
}

func New(l *log.Logger, s *storage.Service) *Service {
	l.Println("initializing rps service")

	storageID := s.AddStore()

	// reset rps every second.
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			s.Clear(storageID)
		}
	}()

	return &Service{
		logger: l,

		storageID: storageID,
		storage:   s,
	}
}

func (s *Service) IncAndGet(ip string) uint {
	s.storage.Increment(s.storageID, ip)

	rps, err := s.storage.Get(s.storageID, ip)
	if err != nil {
		return 0
	}

	return rps
}
