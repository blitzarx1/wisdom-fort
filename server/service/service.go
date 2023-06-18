package service

import (
	"encoding/json"
	"log"
	"sync"

	"blitzarx1/wisdom-fort/server/service/quotes"
)

type Difficulty uint8

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
	activeChallenges map[Token]Difficulty
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
		activeChallenges: make(map[Token]Difficulty),
	}, nil
}

func (s *Service) GenerateToken(ip string) Token {
	return newToken(ip)
}

func (s *Service) HandleChallenge(ip string, t Token) ([]byte, *Error) {
	reqLogger := NewLogger(s.logger, string(t))
	reqLogger.Println("handling challenge request")

	difficulty, ok := s.activeChallenges[t]
	if !ok {
		difficulty = s.computeChallenge(t)
	}

	payload := PayloadChallenge{Target: uint8(difficulty)}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, NewError(ErrGeneric, err)
	}

	return data, nil
}

func (s *Service) computeChallenge(token Token) Difficulty {
	s.lock.Lock()
	defer s.lock.Unlock()

	// TODO: check rps and choose corresponding difficulty
	difficulty := 1

	return Difficulty(difficulty)
}

func (s *Service) HandleSolution(ip string, t Token, payload []byte) ([]byte, *Error) {
	reqLogger := NewLogger(s.logger, string(t))
	reqLogger.Println("handling solution request")

	return nil, nil
}
