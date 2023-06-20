package challenges

import (
	"context"
	"time"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/server/internal/service/rps"
	"blitzarx1/wisdom-fort/server/internal/service/storage"
	"blitzarx1/wisdom-fort/server/internal/token"
)

// Service tracks challenges for client, validates solutions and computes difficulty.
// Challenge has a ttl afteer which it expires.
type Service struct {
	storageID storage.StorageID

	diffMult uint8

	storageService *storage.Service
	rpsService     *rps.Service
}

func New(
	ctx context.Context,
	diffMult uint8,
	ttlSeconds uint,
	storageService *storage.Service,
	rpsService *rps.Service,
) *Service {
	l := logger.MustFromCtx(ctx)
	l.Println("initializing challenges service")

	return &Service{
		storageID: storageService.AddStorageWithTTL(
			logger.WithCtx(ctx, l, "addStorage"),
			time.Duration(ttlSeconds)*time.Second,
		),

		diffMult: diffMult,

		storageService: storageService,
		rpsService:     rpsService,
	}
}

func (s *Service) ComputeChallenge(t token.Token) uint8 {
	rps := s.rpsService.Get(t.IP())

	diff := uint8(rps) * s.diffMult

	s.storageService.Set(s.storageID, string(t), uint(diff))
	return diff
}

func (s *Service) Challenge(key string) (uint8, error) {
	challenge, err := s.storageService.Get(s.storageID, key)
	if err != nil {
		return 0, err
	}

	return uint8(challenge), nil
}

// CheckSolution validates the solution provided by the client.
//
// If the solution is correct, the corresponding challenge is removed from active challenges.
//
// The function returns a boolean indicating whether the solution is correct, and an error if something went wrong.
func (s *Service) CheckSolution(t token.Token, sol uint64) (bool, error) {
	diff, err := s.storageService.Get(s.storageID, string(t))
	if err != nil {
		return false, err
	}

	// generate hash of solution with the token
	hash := generateHash(string(t), sol)

	// check if the hash meets the difficulty requirement
	isCorrect := checkHash(hash, uint8(diff))

	if isCorrect {
		s.storageService.Delete(s.storageID, string(t))
	}

	return isCorrect, nil
}
