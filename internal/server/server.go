package server

import (
	"log/slog"

	"github.com/suhrobdomoiZ/yadro-test-2026/config"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/closer"
)

type Server struct {
	Config *config.AppConfig
	Logger *slog.Logger
	Closer *closer.Closer
}

func NewServer(cfg *config.AppConfig, appLogger *slog.Logger, closer *closer.Closer) *Server {
	return &Server{
		Config: cfg,
		Logger: appLogger,
		Closer: closer,
	}
}

func (s *Server) Start() {}
