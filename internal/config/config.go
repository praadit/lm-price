package config

import (
	"os"
	"strings"
	"time"
)

// Config holds process-level settings.
type Config struct {
	LMURL         string
	AntaremasURL  string
	Galeri24URL   string
	HTTPTimeout   time.Duration
	PricesTimeout time.Duration
	CacheTTL      time.Duration
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
	g := os.Getenv("GALERI24_SOURCE_URL")
	if g == "" {
		g = "https://galeri24.co.id/harga-emas"
	}

	cacheTTL := parseDurationEnv("CACHE_TTL", 60*time.Second)
	return Config{
		LMURL:         u,
		AntaremasURL:  a,
		Galeri24URL:   g,
		HTTPTimeout:   15 * time.Second,
		PricesTimeout: 20 * time.Second,
		CacheTTL:      cacheTTL,
	}
}

func parseDurationEnv(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil || d < 0 {
		return def
	}
	return d
}
