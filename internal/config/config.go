package config

import (
	"os"
	"time"
)

// Config holds process-level settings.
type Config struct {
	LMURL        string
	HTTPTimeout  time.Duration
	PricesTimeout time.Duration
}

// Load reads configuration from the environment with sensible defaults.
func Load() Config {
	u := os.Getenv("LM_SOURCE_URL")
	if u == "" {
		u = "https://emasantam.id/content/lm.txt"
	}
	return Config{
		LMURL:         u,
		HTTPTimeout:   15 * time.Second,
		PricesTimeout: 20 * time.Second,
	}
}
