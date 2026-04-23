package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/usecase"
)

type Galeri24 struct {
	UC         *usecase.Galeri24Usecase
	ReqTimeout time.Duration
}

// GetAntamPrices handles GET /galeri24/antam (optional ?raw=1).
func (h *Galeri24) GetAntamPrices(c *gin.Context) {
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

	resp, err := h.UC.GetAntamPrices(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	out := PricesEnvelope{
		LastUpdate: timeToRFC3339String(resp.LastUpdate),
		Source:     "Galeri 24",
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
			Gramasi:   r.Weight,
			BuyPrice:  r.SellPrice,
			SellPrice: r.BuybackPrice,
			Stock:     0,
			SoldOut:   false,
		})
	}

	c.JSON(http.StatusOK, out)
}
