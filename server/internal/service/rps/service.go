package rps

import (
	"context"
	"time"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/server/internal/service/storage"
)

// Service tracks rps per ip.
type Service struct {
	storageID storage.ID
	storage   *storage.Service
}

func New(ctx context.Context, s *storage.Service) *Service {
	l := logger.MustFromCtx(ctx)
	l.Println("initializing rps service")

	storageID := s.AddStorageWithTTL(logger.WithCtx(ctx, l, "addStorage"), time.Second)

	return &Service{
		storageID: storageID,
		storage:   s,
	}
}

func (s *Service) Inc(ip string) {
	s.storage.Increment(s.storageID, ip)
}

func (s *Service) Get(ip string) uint {
	rps, err := s.storage.Get(s.storageID, ip)
	if err != nil {
		return 0
	}

	return rps
}
