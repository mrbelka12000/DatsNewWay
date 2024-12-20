package config

import (
	"errors"
	"os"
)

type Config struct {
	Token string
}

func LoadConfig() (*Config, error) {
	token := os.Getenv("token")
	if token == "" {
		return nil, errors.New("token environment variable not set")
	}

	return &Config{
		Token: token,
	}, nil
}
