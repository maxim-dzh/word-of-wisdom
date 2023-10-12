package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/exp/slog"

	"github.com/maxim-dzh/word-of-wisdom/internal/config"
	serverpkg "github.com/maxim-dzh/word-of-wisdom/internal/server"
	"github.com/maxim-dzh/word-of-wisdom/internal/service"
	"github.com/maxim-dzh/word-of-wisdom/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var (
		cfg    config.Config
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	)
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}
	// setup dependencies
	words := []string{
		"The pen that writes your life story must be held in your own hand.",
		"Life is a daring adventure or it is nothing at all.",
		"When life gives you a hundred reasons to cry, show life that you have a thousand reasons to smile.",
		"You get in life what you have the courage to ask for.",
		"The meaning of life is to find your gift. The purpose of life is to give it away.",
		"Too many of us are not living our dreams because we are living our fears.",
		"The purpose of life is not to fight against evil and misfortune; it is to unveil magnificence.",
	}
	server := serverpkg.NewServer(
		uuid.New().String(),
		cfg.ServerAddr,
		cfg.ChallengeComplexity,
		cfg.ChallengeTimeout,
		cfg.ReadTimeout,
		storage.NewStorage(),
		service.NewService(words),
		logger,
	)

	// run the server
	listenErrors := make(chan error, 1)
	go func() {
		err := server.Serve(ctx)
		if err != nil {
			listenErrors <- err
		}
	}()
	logger.Info("server started")
	select {
	case <-ctx.Done():
	case err := <-listenErrors:
		logger.Error("server failed", "error", err)
		return
	}
	logger.Info("server stopped")
}
