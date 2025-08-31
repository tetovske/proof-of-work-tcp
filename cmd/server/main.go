package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/tetovske/proof-of-work-tcp/internal/config"
	"github.com/tetovske/proof-of-work-tcp/internal/model"
	"github.com/tetovske/proof-of-work-tcp/internal/repository/quote"
	"github.com/tetovske/proof-of-work-tcp/internal/server"
	"github.com/tetovske/proof-of-work-tcp/internal/tcp"
	"github.com/tetovske/proof-of-work-tcp/pkg/cache"
)

const cacheSize = 10

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	configPath := pflag.StringP("config", "c", "", "path to config file")
	pflag.Parse()

	cfg := config.New()

	err := cfg.ReadFromConfigAndENV(configPath)
	if err != nil {
		log.Error().Msgf("failed to read from config: %v", err)
		return
	}

	c := cache.New[*model.Quote](cacheSize)
	quotesRepo := quote.New(c)

	quotesRepo.WarmUpCache(cfg.Quotes.Data)

	tcpTransport := tcp.New(cfg, quotesRepo)
	srv := server.New(cfg, tcpTransport)

	log.Info().Msg("starting server")

	if err = srv.Run(ctx); err != nil {
		log.Error().Msgf("failed to run tcp server: %v", err)
		return
	}

	log.Info().Msg("server stopped")
}
