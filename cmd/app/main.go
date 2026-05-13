package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/config"
	"github.com/suhrobdomoiZ/yadro-test-2026/migrations"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/closer"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewAppConfig()
	appLogger := logger.With("env", cfg.EnvType())
	logger.Setup(cfg.EnvType())

	appCloser := closer.New(appLogger)
	appCloser.AddFunc("context", cancel)

	pool, err := pgxpool.New(ctx, cfg.DbConfig.DSN())
	if err != nil {
		appLogger.Error("main: failed to create pool", "error", err)
		os.Exit(1)
	}

	appCloser.AddFunc("pool", pool.Close)

	for attempt := range 10 {
		err = pool.Ping(ctx)
		if err == nil {
			appLogger.Info("main: connected to db", "host", cfg.DbConfig.DBHost())

			break
		}

		appLogger.Warn("main: db is not ready", "attempt", attempt+1, "error", err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		appLogger.Error("main: db ping failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("main: connected to db", "host", cfg.DbConfig.DBHost())

	err = migrations.Up(ctx, pool, appLogger)
	if err != nil {
		appLogger.Error("main: migration failed", "error", err)
		os.Exit(1)
	}
	// Server := server.NewServer(cfg, appLogger, appCloser)
}
