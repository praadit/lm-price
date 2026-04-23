package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/domain/lm"
	"github.com/praadit/lm-price/internal/usecase"
)

// LM exposes HTTP handlers for LM prices.
type LM struct {
	UC         *usecase.LMUsecase
	ReqTimeout time.Duration
}

// GetPrices handles GET /prices (optional ?raw=1, ?area=, ?location=).
func (h *LM) GetPrices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.ReqTimeout)
	defer cancel()

	if c.Query("raw") == "1" {
		raw, err := h.UC.FetchRaw(ctx)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", raw)
		return
	}

	payload, err := h.UC.ListPrices(ctx, c.Query("area"), c.Query("location"))
	if err != nil {
		var qv *lm.QueryValidationError
		if errors.As(err, &qv) {
			c.JSON(http.StatusBadRequest, qv)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payload)
}
