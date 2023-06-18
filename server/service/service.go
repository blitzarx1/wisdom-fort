package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"blitzarx1/wisdom-fort/server/service/challenges"
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

	quotesService     *quotes.Service
	challengesService *challenges.Service
}

func New(logger *log.Logger) (*Service, error) {
	logger.Println("initializing service")

	quotesService, err := quotes.New(NewLogger(logger, "quotes"), quotesFilePath)
	if err != nil {
		return nil, err
	}

	challengesService := challenges.New(NewLogger(logger, "challenges"))
	return &Service{
		logger: logger,

		quotesService:     quotesService,
		challengesService: challengesService,
	}, nil
}

func (s *Service) GenerateToken(ip string) Token {
	return newToken(ip)
}

func (s *Service) GenerateChallenge(ip string, t Token) ([]byte, *Error) {
	reqLogger := NewLogger(s.logger, string(t))
	reqLogger.Println("handling challenge request")

	var diff uint8
	var err error
	challengeKey := string(t)
	if diff, err = s.challengesService.Challenge(challengeKey); err != nil {
		diff = s.challengesService.ComputeChallenge(challengeKey)
	}

	payload := payloadChallenge{Target: difficulty(diff)}
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

	correct, checkSolErr := s.challengesService.CheckSolution(string(t), uint64(reqPayload.Solution))
	if checkSolErr != nil {
		return nil, NewError(ErrInvalidSolution, checkSolErr)
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
