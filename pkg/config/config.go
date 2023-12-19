package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	BscApiKey    string
	Timeout      int64
	BscApiUrl    string
	OpenAiApiKey string
	TgBotApiKEy  string
}

func LoadConfig() (*Config, error) {
	timeout, err := strconv.ParseInt(os.Getenv("BSC_API_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeout: %w", err)
	}
	cfg := &Config{
		BscApiKey:    os.Getenv("BSC_API_KEY"),
		BscApiUrl:    os.Getenv("BSC_API_URL"),
		OpenAiApiKey: os.Getenv("OPENAI_API_KEY"),
		Timeout:      timeout,
		TgBotApiKEy:  os.Getenv("TG_BOT_API_KEY"),
	}

	if cfg.BscApiKey == "" {
		return nil, errors.New("BSC_API_KEY is not set")
	}

	return cfg, nil
}
