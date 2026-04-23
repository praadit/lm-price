package usecase

import (
	"context"

	"github.com/praadit/lm-price/internal/domain/lm"
)

// LMUsecase loads and interprets upstream LM price data.
type LMUsecase struct {
	Source lm.RawSource
}

// FetchRaw returns the upstream document bytes.
func (u *LMUsecase) FetchRaw(ctx context.Context) ([]byte, error) {
	return u.Source.Fetch(ctx)
}

// ListPrices returns parsed rows, optionally filtered by area and/or location.
func (u *LMUsecase) ListPrices(ctx context.Context, area, location string) ([]lm.LocationPrices, error) {
	raw, err := u.Source.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := lm.ParsePrices(raw)
	if err != nil {
		return nil, err
	}
	return lm.FilterPrices(rows, area, location)
}
