package antaremas

import (
	"context"
	"time"
)

// PricesResponse is the JSON envelope for Antaremas buy-price table rows.
type PricesResponse struct {
	LastUpdate time.Time  `json:"last_update,omitempty"`
	Data       []PriceRow `json:"data"`
}

type PriceRow struct {
	Size     string `json:"size"`
	BuyPrice int    `json:"buy_price"`
}

// RawSource fetches upstream Antaremas HTML.
type RawSource interface {
	Fetch(ctx context.Context) ([]byte, error)
}
