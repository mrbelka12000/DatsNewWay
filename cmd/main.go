package main

import (
	"DatsNewWay/client"
	"DatsNewWay/config"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	cli := client.NewClient(cfg.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create client")
	}

	_ = cli

	if err = start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("BELKA MONSTER")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
