package http

import (
	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/config"
	"github.com/praadit/lm-price/internal/delivery/http/handler"
	"github.com/praadit/lm-price/internal/usecase"
)

// NewRouter wires middleware and HTTP routes.
func NewRouter(cfg config.Config, lmUC *usecase.LMUsecase, antUC *usecase.AntaremasUsecase) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	lmH := &handler.LM{UC: lmUC, ReqTimeout: cfg.PricesTimeout}
	amH := &handler.Antaremas{UC: antUC, ReqTimeout: cfg.PricesTimeout}
	r.GET("/health", handler.Health)

	v1 := r.Group("/v1")

	prices := v1.Group("/prices")
	prices.GET("/antam", lmH.GetPrices)
	prices.GET("/hf", amH.GetBuyPrices)
	return r
}
