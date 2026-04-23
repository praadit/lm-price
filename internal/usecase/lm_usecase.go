package usecase

import (
	"context"

	gocache "github.com/patrickmn/go-cache"
	"github.com/praadit/lm-price/internal/domain/lm"
)

// LMUsecase loads and interprets upstream LM price data.
type LMUsecase struct {
	Source lm.RawSource
	Cache  *gocache.Cache
}

// FetchRaw returns the upstream document bytes.
func (u *LMUsecase) FetchRaw(ctx context.Context) ([]byte, error) {
	return u.Source.Fetch(ctx)
}

// ListPrices returns parsed rows and upstream last-update text, optionally filtered by area and/or location.
func (u *LMUsecase) ListPrices(ctx context.Context, area, location string) (lm.PricesResponse, error) {
	var doc lm.PricesResponse
	if u.Cache != nil {
		if v, ok := u.Cache.Get("lm:doc"); ok {
			if cast, ok := v.(lm.PricesResponse); ok {
				doc = cast
			}
		}
	}
	if doc.Data == nil && doc.LastUpdate.IsZero() {
		raw, err := u.Source.Fetch(ctx)
		if err != nil {
			return lm.PricesResponse{}, err
		}
		doc, err = lm.ParsePricesDocument(raw)
		if err != nil {
			return lm.PricesResponse{}, err
		}
		if u.Cache != nil {
			u.Cache.Set("lm:doc", doc, gocache.DefaultExpiration)
		}
	}

	filtered, err := lm.FilterPrices(doc.Data, area, location)
	if err != nil {
		return lm.PricesResponse{}, err
	}
	return lm.PricesResponse{
		LastUpdate: doc.LastUpdate,
		Data:       filtered,
	}, nil
}
