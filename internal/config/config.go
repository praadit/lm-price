package config

import (
	"os"
	"time"
)

// Config holds process-level settings.
type Config struct {
	LMURL         string
	AntaremasURL  string
	HTTPTimeout   time.Duration
	PricesTimeout time.Duration
}

// Load reads configuration from the environment with sensible defaults.
func Load() Config {
	u := os.Getenv("LM_SOURCE_URL")
	if u == "" {
		u = "https://emasantam.id/content/lm.txt"
	}
	a := os.Getenv("ANTAREMAS_SOURCE_URL")
	if a == "" {
		a = "https://antaremas.com/harga-emas/"
	}
	return Config{
		LMURL:         u,
		AntaremasURL:  a,
		HTTPTimeout:   15 * time.Second,
		PricesTimeout: 20 * time.Second,
	}
}
