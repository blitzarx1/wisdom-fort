package challenges

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Service struct {
	logger *log.Logger

	lock             sync.RWMutex
	activeChallenges map[string]uint8
}

func New(logger *log.Logger) *Service {
	logger.Println("initializing challenges service")

	return &Service{
		logger: logger,

		lock:             sync.RWMutex{},
		activeChallenges: make(map[string]uint8),
	}
}

func (s *Service) ComputeChallenge(token string) uint8 {
	s.lock.Lock()
	defer s.lock.Unlock()

	// TODO: check rps and choose corresponding difficulty
	var diff uint8 = 1

	s.activeChallenges[token] = diff
	return diff
}

func (s *Service) Challenge(key string) (uint8, error) {
	challenge, ok := s.activeChallenges[key]
	if !ok {
		return 0, errors.New("no challenge found")
	}

	return challenge, nil
}

// CheckSolution validates the solution provided by the client.
//
// If the solution is correct, the corresponding challenge is removed from active challenges.
//
// The function returns a boolean indicating whether the solution is correct, and an error if something went wrong.
func (s *Service) CheckSolution(t string, sol uint64) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	diff, ok := s.activeChallenges[t]
	if !ok {
		return false, fmt.Errorf("active challenges not found for token: %s", t)
	}

	// generate hash of solution with the token
	hash := generateHash(t, sol)

	// check if the hash meets the difficulty requirement
	isCorrect := checkHash(hash, diff)

	if isCorrect {
		delete(s.activeChallenges, t)
	}

	return isCorrect, nil
}
