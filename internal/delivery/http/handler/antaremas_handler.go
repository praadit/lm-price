package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/usecase"
)

type Antaremas struct {
	UC         *usecase.AntaremasUsecase
	ReqTimeout time.Duration
}

// GetBuyPrices handles GET /antaremas (optional ?raw=1).
func (h *Antaremas) GetBuyPrices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.ReqTimeout)
	defer cancel()

	if c.Query("raw") == "1" {
		raw, err := h.UC.FetchRaw(ctx)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", raw)
		return
	}

	resp, err := h.UC.GetBuyPrices(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
