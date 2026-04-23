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

	out := PricesEnvelope{
		LastUpdate: timeToRFC3339String(resp.LastUpdate),
		Source:     "Antaremas",
		Data: []LocationData{
			{
				Location: "Indonesia",
				Product:  "ANTAM",
				Area:     "Indonesia",
				Prices:   make([]PriceEntry, 0, len(resp.Data)),
			},
		},
	}
	for _, r := range resp.Data {
		out.Data[0].Prices = append(out.Data[0].Prices, PriceEntry{
			Gramasi:   parseLeadingFloat(r.Size),
			BuyPrice:  r.BuyPrice,
			SellPrice: 0,
			Stock:     0,
			SoldOut:   false,
		})
	}

	c.JSON(http.StatusOK, out)
}
