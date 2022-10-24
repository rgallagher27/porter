package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines the input config required for the app to run
type Config struct {
	RedisAddr    string
	InputFile    string
	IgnoreErrors bool
}

// Load reads configuration from environment variables.
// This is a simplistic version of a config loader, using only env vars.
func Load() (Config, error) {
	var cfg Config
	cfg.RedisAddr = os.Getenv("REDIS_ADDRESS")
	cfg.InputFile = os.Getenv("INPUT_FILE")

	ignoreBool, err := strconv.ParseBool(os.Getenv("IGNORE_ERRORS"))
	if err != nil {
		return Config{}, fmt.Errorf("parse ignore errors: %w", err)
	}
	cfg.IgnoreErrors = ignoreBool

	return cfg, nil
}
