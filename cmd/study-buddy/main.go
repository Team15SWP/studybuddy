package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"study_buddy/internal/config"
	"study_buddy/internal/db"
	"study_buddy/internal/server"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting study_buddy", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	pgPool, err := db.NewPgPool(ctx, &cfg.PGConfig)
	if err != nil {
		log.Error(fmt.Sprintf("db.NewPgPool: %v", err))
		return
	}
	log.Info("pgPool is created successfully")

	source := server.NewServer(cfg, log, pgPool)
	source.Run(ctx)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
