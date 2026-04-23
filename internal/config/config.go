package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds process-level settings.
type Config struct {
	LMURL                          string
	AntaremasURL                   string
	Galeri24URL                    string
	HTTPTimeout                    time.Duration
	PricesTimeout                  time.Duration
	CacheTTL                       time.Duration
	BasicAuthUser                  string
	BasicAuthPass                  string
	RateLimitUnauthorizedPerMinute int
	RateLimitAuthorizedPerMinute   int
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

	cacheTTL := parseIntEnv("CACHE_TTL_SECOND", 60)
	return Config{
		LMURL:                          u,
		AntaremasURL:                   a,
		Galeri24URL:                    g,
		HTTPTimeout:                    15 * time.Second,
		PricesTimeout:                  20 * time.Second,
		CacheTTL:                       time.Duration(cacheTTL) * time.Second,
		BasicAuthUser:                  strings.TrimSpace(os.Getenv("BASIC_AUTH_USER")),
		BasicAuthPass:                  os.Getenv("BASIC_AUTH_PASS"),
		RateLimitUnauthorizedPerMinute: parseIntEnv("RATE_LIMIT_UNAUTHORIZED_PER_MINUTE", 1),
		RateLimitAuthorizedPerMinute:   parseIntEnv("RATE_LIMIT_AUTHORIZED_PER_MINUTE", 100),
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

func parseIntEnv(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}
