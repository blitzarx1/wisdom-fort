package challenges

import (
	"log"

	"blitzarx1/wisdom-fort/server/service/rps"
	"blitzarx1/wisdom-fort/server/service/storage"
	"blitzarx1/wisdom-fort/server/token"
)

// Service tracks challenges for client, validates solutions and computes difficulty.
type Service struct {
	logger *log.Logger

	storageID storage.StorageID
	storage   *storage.Service

	rpsService *rps.Service
}

func New(l *log.Logger, storageService *storage.Service, rpsService *rps.Service) *Service {
	l.Println("initializing challenges service")

	return &Service{
		logger: l,

		storageID: storageService.AddStore(),
		storage:   storageService,

		rpsService: rpsService,
	}
}

func (s *Service) ComputeChallenge(t token.Token) uint8 {
	rps := s.rpsService.Get(t.IP())

	diff := uint8(rps)

	s.storage.Set(s.storageID, string(t), uint(diff))
	return diff
}

func (s *Service) Challenge(key string) (uint8, error) {
	challenge, err := s.storage.Get(s.storageID, key)
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
	diff, err := s.storage.Get(s.storageID, string(t))
	if err != nil {
		return false, err
	}

	// generate hash of solution with the token
	hash := generateHash(string(t), sol)

	// check if the hash meets the difficulty requirement
	isCorrect := checkHash(hash, uint8(diff))

	if isCorrect {
		s.storage.Delete(s.storageID, string(t))
	}

	return isCorrect, nil
}
