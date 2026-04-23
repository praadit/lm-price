package galeri24

import (
	"context"
	"time"
)

// PricesResponse is the JSON envelope for Galeri24 "Harga ANTAM" table rows.
type PricesResponse struct {
	LastUpdate time.Time  `json:"last_update,omitempty"`
	Data       []PriceRow `json:"data"`
}

type PriceRow struct {
	Weight       float64 `json:"weight"`
	SellPrice    int     `json:"sell_price"`
	BuybackPrice int     `json:"buyback_price"`
}

// RawSource fetches upstream Galeri24 HTML.
type RawSource interface {
	Fetch(ctx context.Context) ([]byte, error)
}
