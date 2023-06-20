package server

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config is a struct that holds all the configuration for the server
type Config struct {
	// Port to listen on
	Port uint `env:"PORT"`

	// RPSLimitUnauth is ip rps limit for requests without valid token
	RPSLimitUnauth uint `env:"RPS_LIMIT_UNAUTH"`

	// DiffMult is difficulty multiplier for challenges. If set to 1 the
	// difficulty is equal to the client IPs RPS. 0 makes
	// challenges trivial. Recommended value is 1.
	DiffMult uint8 `env:"DIFF_MULT"`

	// ChallengeTTLSeconds is expiration time for challenge in seconds.
	// When the time is passed the challenge is considered invalid and
	// the client needs to request a new one.
	ChallengeTTLSeconds uint `env:"CHALLENGE_TTL_SECONDS"`
}

func NewConfigFromEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
