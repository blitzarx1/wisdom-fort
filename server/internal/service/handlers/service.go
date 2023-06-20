package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/pkg/scheme"
	wfErrors "blitzarx1/wisdom-fort/server/internal/errors"
	"blitzarx1/wisdom-fort/server/internal/service/challenges"
	"blitzarx1/wisdom-fort/server/internal/service/quotes"
	"blitzarx1/wisdom-fort/server/internal/service/rps"
	"blitzarx1/wisdom-fort/server/internal/service/storage"
	"blitzarx1/wisdom-fort/server/internal/token"
)

// Service encapsulates logic of handling requests from the client.
// It defines the difficulty of the challenge for given client and
// check correctness of the results.
type Service struct {
	quotesService     *quotes.Service
	rpsService        *rps.Service
	challengesService *challenges.Service
}

func New(
	ctx context.Context,
	rpsService *rps.Service,
	storageService *storage.Service,
	quotesService *quotes.Service,
	challengesService *challenges.Service,
) (*Service, error) {
	logger.MustFromCtx(ctx).Println("initializing service")

	return &Service{
		quotesService:     quotesService,
		challengesService: challengesService,
		rpsService:        rpsService,
	}, nil
}

func (s *Service) GenerateToken(ip string) token.Token {
	return token.New(ip)
}

// GenerateChallenge generates a challenge for the given token.
func (s *Service) GenerateChallenge(ctx context.Context, t token.Token) ([]byte, *wfErrors.Error) {
	l := logger.MustFromCtx(ctx)
	l.Println("handling challenge request")

	var diff uint8
	var err error
	challengeKey := string(t)
	if diff, err = s.challengesService.Challenge(challengeKey); err != nil {
		diff = s.challengesService.ComputeChallenge(t)
	}

	payload := scheme.PayloadResponseChallenge{Target: diff}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, wfErrors.NewError(wfErrors.ErrGeneric, err)
	}

	return data, nil
}

// CheckSolution checks correctness of the solution for the given token.
func (s *Service) CheckSolution(ctx context.Context, t token.Token, payload []byte) ([]byte, *wfErrors.Error) {
	l := logger.MustFromCtx(ctx)
	l.Println("handling solution request")

	if payload == nil {
		return nil, wfErrors.NewError(wfErrors.ErrInvalidPayloadFormat, errors.New("empty payload"))
	}

	var reqPayload scheme.PayloadRequestSolution
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
	respPayload := scheme.PayloadResponseSolution{Quote: quote}
	data, err := json.Marshal(respPayload)
	if err != nil {
		return nil, wfErrors.NewError(wfErrors.ErrGeneric, err)
	}

	return data, nil
}
