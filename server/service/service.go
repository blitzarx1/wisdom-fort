package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"blitzarx1/wisdom-fort/server/service/quotes"
)

type (
	difficulty uint8
	solution   uint64
)

const (
	quotesFilePath = "server/quotes.json"
)

// Service encapsulates logic of handling requests from the client.
// It defines the difficulty of the challenge for given client and
// check correctness of the results.
type Service struct {
	logger *log.Logger

	quotesService *quotes.Service

	lock             sync.RWMutex
	activeChallenges map[Token]difficulty
}

func New(logger *log.Logger) (*Service, error) {
	logger.Println("initializing service")

	quotesService, err := quotes.New(NewLogger(logger, "quotes"), quotesFilePath)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger: logger,

		quotesService: quotesService,

		lock:             sync.RWMutex{},
		activeChallenges: make(map[Token]difficulty),
	}, nil
}

func (s *Service) GenerateToken(ip string) Token {
	return newToken(ip)
}

func (s *Service) GenerateChallenge(ip string, t Token) ([]byte, *Error) {
	reqLogger := NewLogger(s.logger, string(t))
	reqLogger.Println("handling challenge request")

	difficulty, ok := s.activeChallenges[t]
	if !ok {
		difficulty = s.computeChallenge(t)
	}

	payload := payloadChallenge{Target: difficulty}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, NewError(ErrGeneric, err)
	}

	return data, nil
}

func (s *Service) CheckSolution(ip string, t Token, payload []byte) ([]byte, *Error) {
	reqLogger := NewLogger(s.logger, string(t))
	reqLogger.Println("handling solution request")

	if payload == nil {
		return nil, NewError(ErrInvalidPayloadFormat, errors.New("empty payload"))
	}

	var reqPayload payloadRequestSolution
	err := json.Unmarshal(payload, &reqPayload)
	if err != nil {
		return nil, NewError(ErrInvalidPayloadFormat, err)
	}

	correct, checkSolErr := s.checkSolution(t, reqPayload.Solution)
	if checkSolErr != nil {
		return nil, checkSolErr
	}

	if !correct {
		return nil, NewError(ErrInvalidSolution, fmt.Errorf("solution is invalid: %d", reqPayload.Solution))
	}

	quote := s.quotesService.GetRandom()
	respPayload := payloadResponseSolution{Quote: quote}
	data, err := json.Marshal(respPayload)
	if err != nil {
		return nil, NewError(ErrGeneric, err)
	}

	return data, nil
}

func (s *Service) computeChallenge(token Token) difficulty {
	s.lock.Lock()
	defer s.lock.Unlock()

	// TODO: check rps and choose corresponding difficulty
	diff := 1

	s.activeChallenges[token] = difficulty(diff)
	return difficulty(diff)
}

// checkSolution validates the solution provided by the client.
//
// If the solution is correct, the corresponding challenge is removed from active challenges.
//
// The function returns a boolean indicating whether the solution is correct, and an error if something went wrong.
func (s *Service) checkSolution(t Token, sol solution) (bool, *Error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	diff, ok := s.activeChallenges[t]
	if !ok {
		return false, NewError(ErrNoActiveChallenge, fmt.Errorf("active challenges not found for token: %s", t))
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
