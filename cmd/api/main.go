package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/praadit/lm-price/internal/config"
	deliveryhttp "github.com/praadit/lm-price/internal/delivery/http"
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
	r := deliveryhttp.NewRouter(cfg, uc)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
