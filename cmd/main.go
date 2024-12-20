package main

import (
	"DatsNewWay/client"
	"DatsNewWay/config"
	"DatsNewWay/entity"
	"context"
	"github.com/rs/zerolog/log"
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

}
