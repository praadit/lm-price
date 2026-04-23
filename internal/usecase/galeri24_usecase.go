package usecase

import (
	"context"

	gocache "github.com/patrickmn/go-cache"
	"github.com/praadit/lm-price/internal/domain/galeri24"
)

type Galeri24Usecase struct {
	Source galeri24.RawSource
	Cache  *gocache.Cache
}

func (u *Galeri24Usecase) FetchRaw(ctx context.Context) ([]byte, error) {
	return u.Source.Fetch(ctx)
}

func (u *Galeri24Usecase) GetAntamPrices(ctx context.Context) (galeri24.PricesResponse, error) {
	if u.Cache != nil {
		if v, ok := u.Cache.Get("galeri24:antam"); ok {
			if cast, ok := v.(galeri24.PricesResponse); ok {
				return cast, nil
			}
		}
	}
	raw, err := u.Source.Fetch(ctx)
	if err != nil {
		return galeri24.PricesResponse{}, err
	}
	parsed, err := galeri24.ParseAntamPricesDocument(raw)
	if err != nil {
		return galeri24.PricesResponse{}, err
	}
	if u.Cache != nil {
		u.Cache.Set("galeri24:antam", parsed, gocache.DefaultExpiration)
	}
	return parsed, nil
}
