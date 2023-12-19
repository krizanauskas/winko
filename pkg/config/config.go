package config

import (
	"errors"
	"os"
)

type Config struct {
	BscApiKey string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		BscApiKey: os.Getenv("BSC_API_KEY"),
	}

	if cfg.BscApiKey == "" {
		return nil, errors.New("BSC_API_KEY is not set")
	}

	return cfg, nil
}
