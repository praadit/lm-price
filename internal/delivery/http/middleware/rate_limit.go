package middleware

import (
	"crypto/subtle"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	BasicAuthUser string
	BasicAuthPass string

	UnauthorizedPerMinute int
	AuthorizedPerMinute   int
}

type limiterStore struct {
	mu sync.Mutex
	m  map[string]*rate.Limiter
}

func (s *limiterStore) get(key string, limit rate.Limit, burst int) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.m == nil {
		s.m = make(map[string]*rate.Limiter)
	}
	if l, ok := s.m[key]; ok {
		return l
	}
	l := rate.NewLimiter(limit, burst)
	s.m[key] = l
	return l
}

func RateLimit(cfg RateLimitConfig) gin.HandlerFunc {
	var unauth limiterStore
	var auth limiterStore

	unauthPerMin := cfg.UnauthorizedPerMinute
	if unauthPerMin <= 0 {
		unauthPerMin = 1
	}
	authPerMin := cfg.AuthorizedPerMinute
	if authPerMin <= 0 {
		authPerMin = 100
	}

	unauthLimit := rate.Every(time.Minute / time.Duration(unauthPerMin))
	authLimit := rate.Every(time.Minute / time.Duration(authPerMin))

	return func(c *gin.Context) {
		ip := c.ClientIP()
		authorized := basicAuthOK(c, cfg.BasicAuthUser, cfg.BasicAuthPass)

		var l *rate.Limiter
		if authorized {
			l = auth.get("ip:"+ip, authLimit, authPerMin)
		} else {
			l = unauth.get("ip:"+ip, unauthLimit, 1)
		}

		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func basicAuthOK(c *gin.Context, user, pass string) bool {
	if user == "" || pass == "" {
		return false
	}
	u, p, ok := c.Request.BasicAuth()
	if !ok {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(u), []byte(user)) != 1 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(p), []byte(pass)) != 1 {
		return false
	}
	return true
}
