package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"DatsNewWay/algo"
	"DatsNewWay/client"
	"DatsNewWay/config"
	"DatsNewWay/entity"
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

	if err = start(ctx, cli); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func start(ctx context.Context, cl *client.Client) error {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	var (
		resp entity.Response
		err  error
	)

	resp, err = cl.Get(ctx, entity.Payload{})
	if err != nil {
		log.Err(err).Msg("failed to create client")
		return err
	}

	for {
		select {
		case <-ticker.C:

			start := time.Now()
			payload := algo.GetNextDirection(resp)

			resp, err = cl.Get(ctx, payload)
			if err != nil {
				log.Err(err).Msg("failed to create client")
				continue
			}

			log.Info().Msg(fmt.Sprintf("successfully send data, spent: %v", time.Since(start).Seconds()))
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
