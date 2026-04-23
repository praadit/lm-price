package usecase

import (
	"context"

	gocache "github.com/patrickmn/go-cache"
	"github.com/praadit/lm-price/internal/domain/antaremas"
)

type AntaremasUsecase struct {
	Source antaremas.RawSource
	Cache  *gocache.Cache
}

func (u *AntaremasUsecase) FetchRaw(ctx context.Context) ([]byte, error) {
	return u.Source.Fetch(ctx)
}

func (u *AntaremasUsecase) GetBuyPrices(ctx context.Context) (antaremas.PricesResponse, error) {
	if u.Cache != nil {
		if v, ok := u.Cache.Get("antaremas:buy"); ok {
			if cast, ok := v.(antaremas.PricesResponse); ok {
				return cast, nil
			}
		}
	}
	raw, err := u.Source.Fetch(ctx)
	if err != nil {
		return antaremas.PricesResponse{}, err
	}
	parsed, err := antaremas.ParsePricesDocument(raw)
	if err != nil {
		return antaremas.PricesResponse{}, err
	}
	if u.Cache != nil {
		u.Cache.Set("antaremas:buy", parsed, gocache.DefaultExpiration)
	}
	return parsed, nil
}
