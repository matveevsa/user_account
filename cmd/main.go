package main

import (
	_ "account/docs"
	"account/internal/app"
	"account/internal/config"
	"account/internal/logger"
	"context"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	lg := logger.New(cfg)

	application := app.New(&lg, cfg)

	if err := application.Run(ctx); err != nil {
		lg.Fatal().Err(err).Msg("error")
	}
}
