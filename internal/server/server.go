package server

import (
	"log/slog"

	"github.com/suhrobdomoiZ/yadro-test-2026/config"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/closer"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/logger"
)

type Server struct {
	Config *config.AppConfig
	Logger *slog.Logger
	Closer *closer.Closer
}

func NewServer() (*Server, error) {
	cfg := config.NewAppConfig()
	Logger := logger.With("env", cfg.EnvType())
	logger.Setup(cfg.EnvType())

	clsr := closer.New(Logger)

	return &Server{
		Config: cfg,
		Logger: Logger,
		Closer: clsr,
	}, nil
}
