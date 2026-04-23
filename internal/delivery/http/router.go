package http

import (
	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/config"
	"github.com/praadit/lm-price/internal/delivery/http/handler"
	"github.com/praadit/lm-price/internal/usecase"
)

// NewRouter wires middleware and HTTP routes.
func NewRouter(cfg config.Config, uc *usecase.LMUsecase) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	lmH := &handler.LM{UC: uc, ReqTimeout: cfg.PricesTimeout}
	r.GET("/health", handler.Health)

	v1 := r.Group("/v1")
	v1.GET("/prices", lmH.GetPrices)
	return r
}
