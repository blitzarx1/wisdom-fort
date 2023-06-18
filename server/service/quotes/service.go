package quotes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
)

type Service struct {
	logger *log.Logger
	quotes []Quote
}

func New(logger *log.Logger, quotesFilePath string) (*Service, error) {
	logger.Println("initializing quotes service")

	data, err := os.ReadFile(quotesFilePath)
	if err != nil {
		return nil, err
	}

	logger.Println("reading quotes")

	s := &Service{logger: logger}
	if err := json.Unmarshal(data, &s.quotes); err != nil {
		return nil, err
	}

	if len(s.quotes) == 0 {
		return nil, errors.New("no quotes found")
	}

	logger.Printf("read %d quotes\n", len(s.quotes))

	return s, nil
}

func (s *Service) GetRandom() Quote {
	return s.quotes[rand.Intn(len(s.quotes))]
}
