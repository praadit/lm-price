package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	gocache "github.com/patrickmn/go-cache"
	"github.com/praadit/lm-price/internal/config"
	deliveryhttp "github.com/praadit/lm-price/internal/delivery/http"
	"github.com/praadit/lm-price/internal/repository/antaremasremote"
	"github.com/praadit/lm-price/internal/repository/galeri24remote"
	"github.com/praadit/lm-price/internal/repository/lmremote"
	"github.com/praadit/lm-price/internal/usecase"
)

func main() {
	cfg := config.Load()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	src := lmremote.NewRawSource(cfg.LMURL, cfg.HTTPTimeout)
	uc := &usecase.LMUsecase{Source: src}
	if cfg.CacheTTL > 0 {
		uc.Cache = gocache.New(cfg.CacheTTL, 2*cfg.CacheTTL)
	}

	amSrc := antaremasremote.NewRawSource(cfg.AntaremasURL, cfg.HTTPTimeout)
	amUC := &usecase.AntaremasUsecase{Source: amSrc}
	if cfg.CacheTTL > 0 {
		amUC.Cache = gocache.New(cfg.CacheTTL, 2*cfg.CacheTTL)
	}

	g24Src := galeri24remote.NewRawSource(cfg.Galeri24URL, cfg.HTTPTimeout)
	g24UC := &usecase.Galeri24Usecase{Source: g24Src}
	if cfg.CacheTTL > 0 {
		g24UC.Cache = gocache.New(cfg.CacheTTL, 2*cfg.CacheTTL)
	}

	r := deliveryhttp.NewRouter(cfg, uc, amUC, g24UC)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
