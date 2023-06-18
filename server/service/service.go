package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"blitzarx1/wisdom-fort/server/logger"
	"blitzarx1/wisdom-fort/server/service/challenges"
	"blitzarx1/wisdom-fort/server/service/quotes"
	"blitzarx1/wisdom-fort/server/service/rps"
	"blitzarx1/wisdom-fort/server/service/storage"
	"blitzarx1/wisdom-fort/server/token"
	wfErrors "blitzarx1/wisdom-fort/server/errors"
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
	rpsService        *rps.Service
	challengesService *challenges.Service
}

func New(l *log.Logger, rpsService *rps.Service, storageService *storage.Service) (*Service, error) {
	l.Println("initializing service")

	quotesService, err := quotes.New(logger.NewLogger(l, "quotes"), quotesFilePath)
	if err != nil {
		return nil, err
	}

	challengesService := challenges.New(logger.NewLogger(l, "challenges"), storageService, rpsService)
	return &Service{
		logger: l,

		quotesService:     quotesService,
		challengesService: challengesService,
		rpsService:        rpsService,
	}, nil
}

func (s *Service) GenerateToken(ip string) token.Token {
	return token.New(ip)
}

// GenerateChallenge generates a challenge for the given token.
func (s *Service) GenerateChallenge(t token.Token) ([]byte, *wfErrors.Error) {
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
		return nil, wfErrors.NewError(wfErrors.ErrGeneric, err)
	}

	return data, nil
}

// CheckSolution checks correctness of the solution for the given token.
func (s *Service) CheckSolution(t token.Token, payload []byte) ([]byte, *wfErrors.Error) {
	reqLogger := logger.NewLogger(s.logger, string(t))
	reqLogger.Println("handling solution request")

	if payload == nil {
		return nil, wfErrors.NewError(wfErrors.ErrInvalidPayloadFormat, errors.New("empty payload"))
	}

	var reqPayload payloadRequestSolution
	err := json.Unmarshal(payload, &reqPayload)
	if err != nil {
		return nil, wfErrors.NewError(wfErrors.ErrInvalidPayloadFormat, err)
	}

	correct, checkSolErr := s.challengesService.CheckSolution(t, uint64(reqPayload.Solution))
	if checkSolErr != nil {
		return nil, wfErrors.NewError(wfErrors.ErrInvalidSolution, checkSolErr)
	}

	if !correct {
		return nil, wfErrors.NewError(wfErrors.ErrInvalidSolution, fmt.Errorf("solution is invalid: %d", reqPayload.Solution))
	}

	quote := s.quotesService.GetRandom()
	respPayload := payloadResponseSolution{Quote: quote}
	data, err := json.Marshal(respPayload)
	if err != nil {
		return nil, wfErrors.NewError(wfErrors.ErrGeneric, err)
	}

	return data, nil
}
