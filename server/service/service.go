package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"blitzarx1/wisdom-fort/server/logger"
	"blitzarx1/wisdom-fort/server/service/challenges"
	"blitzarx1/wisdom-fort/server/service/quotes"
	"blitzarx1/wisdom-fort/server/service/storage"
	"blitzarx1/wisdom-fort/server/token"
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

func New(l *log.Logger) (*Service, error) {
	l.Println("initializing service")

	quotesService, err := quotes.New(logger.NewLogger(l, "quotes"), quotesFilePath)
	if err != nil {
		return nil, err
	}

	storageService := storage.New(logger.NewLogger(l, "storage"))

	challengesService := challenges.New(logger.NewLogger(l, "challenges"), storageService)
	return &Service{
		logger: l,

		quotesService:     quotesService,
		challengesService: challengesService,
	}, nil
}

func (s *Service) GenerateToken(ip string) token.Token {
	return token.New(ip)
}

func (s *Service) GenerateChallenge(t token.Token) ([]byte, *Error) {
	reqLogger := logger.NewLogger(s.logger, string(t))
	reqLogger.Println("handling challenge request")

	var diff uint8
	var err error
	challengeKey := string(t)
	if diff, err = s.challengesService.Challenge(challengeKey); err != nil {
		diff = s.challengesService.ComputeChallenge(t)
	}

	payload := payloadChallenge{Target: diff}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, NewError(ErrGeneric, err)
	}

	return data, nil
}

func (s *Service) CheckSolution(t token.Token, payload []byte) ([]byte, *Error) {
	reqLogger := logger.NewLogger(s.logger, string(t))
	reqLogger.Println("handling solution request")

	if payload == nil {
		return nil, NewError(ErrInvalidPayloadFormat, errors.New("empty payload"))
	}

	var reqPayload payloadRequestSolution
	err := json.Unmarshal(payload, &reqPayload)
	if err != nil {
		return nil, NewError(ErrInvalidPayloadFormat, err)
	}

	correct, checkSolErr := s.challengesService.CheckSolution(t, uint64(reqPayload.Solution))
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
