package quotes

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/pkg/scheme"
)

// Service manages the set of available quotes.
type Service struct {
	quotes []scheme.Quote
}

func New(ctx context.Context, quotesFilePath string) (*Service, error) {
	l := logger.MustFromCtx(ctx)
	l.Println("initializing quotes service")

	data, err := os.ReadFile(quotesFilePath)
	if err != nil {
		return nil, err
	}

	l.Println("reading quotes")

	s := &Service{}
	if err := json.Unmarshal(data, &s.quotes); err != nil {
		return nil, err
	}

	if len(s.quotes) == 0 {
		return nil, errors.New("no quotes found")
	}

	l.Printf("read %d quotes\n", len(s.quotes))

	return s, nil
}

func (s *Service) GetRandom() scheme.Quote {
	return s.quotes[rand.Intn(len(s.quotes))]
}
