package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/config"
	deliveryhttp "github.com/praadit/lm-price/internal/delivery/http"
	"github.com/praadit/lm-price/internal/repository/antaremasremote"
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

	amSrc := antaremasremote.NewRawSource(cfg.AntaremasURL, cfg.HTTPTimeout)
	amUC := &usecase.AntaremasUsecase{Source: amSrc}

	r := deliveryhttp.NewRouter(cfg, uc, amUC)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
