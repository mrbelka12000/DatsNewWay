package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Token string `env:"token,required"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := new(Config)

	if err = env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
