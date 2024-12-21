package main

import (
	"DatsNewWay/algo_a_section"
	"context"
	"time"

	"github.com/rs/zerolog/log"

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
	ticker := time.NewTicker(900 * time.Millisecond)
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

			payload := algo_a_section.GetNextDirection(resp)

			resp, err = cl.Get(ctx, payload)
			if err != nil {
				log.Err(err).Msg("failed to create client")
				continue
			}

			log.Info().Msg("successfully send data")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
