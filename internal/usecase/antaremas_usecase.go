package usecase

import (
	"context"

	"github.com/praadit/lm-price/internal/domain/antaremas"
)

type AntaremasUsecase struct {
	Source antaremas.RawSource
}

func (u *AntaremasUsecase) FetchRaw(ctx context.Context) ([]byte, error) {
	return u.Source.Fetch(ctx)
}

func (u *AntaremasUsecase) GetBuyPrices(ctx context.Context) (antaremas.PricesResponse, error) {
	raw, err := u.Source.Fetch(ctx)
	if err != nil {
		return antaremas.PricesResponse{}, err
	}
	return antaremas.ParsePricesDocument(raw)
}
