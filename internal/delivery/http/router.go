package http

import (
	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/config"
	"github.com/praadit/lm-price/internal/delivery/http/handler"
	"github.com/praadit/lm-price/internal/delivery/http/middleware"
	"github.com/praadit/lm-price/internal/usecase"
)

// NewRouter wires middleware and HTTP routes.
func NewRouter(cfg config.Config, lmUC *usecase.LMUsecase, antUC *usecase.AntaremasUsecase, g24UC *usecase.Galeri24Usecase) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	rateLimitConfig := middleware.RateLimitConfig{
		BasicAuthUser:         cfg.BasicAuthUser,
		BasicAuthPass:         cfg.BasicAuthPass,
		UnauthorizedPerMinute: cfg.RateLimitUnauthorizedPerMinute,
		AuthorizedPerMinute:   cfg.RateLimitAuthorizedPerMinute,
	}

	lmH := &handler.LM{UC: lmUC, ReqTimeout: cfg.PricesTimeout}
	amH := &handler.Antaremas{UC: antUC, ReqTimeout: cfg.PricesTimeout}
	g24H := &handler.Galeri24{UC: g24UC, ReqTimeout: cfg.PricesTimeout}
	r.GET("/health", handler.Health)

	v1 := r.Group("/v1")

	prices := v1.Group("/prices")
	prices.GET("/antam",
		middleware.RateLimit(rateLimitConfig),
		lmH.GetPrices,
	)
	prices.GET("/hfgold",
		middleware.RateLimit(rateLimitConfig),
		amH.GetBuyPrices,
	)
	prices.GET("/galeri24",
		middleware.RateLimit(rateLimitConfig),
		g24H.GetAntamPrices,
	)
	return r
}
