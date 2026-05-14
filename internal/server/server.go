package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/config"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/handlers"
	middleware "github.com/suhrobdomoiZ/yadro-test-2026/internal/midleware"
	"github.com/suhrobdomoiZ/yadro-test-2026/pkg/closer"
)

type Server struct {
	Config  *config.AppConfig
	Logger  *slog.Logger
	Closer  *closer.Closer
	Handler *handlers.AppHandler
}

func NewServer(
	cfg *config.AppConfig,
	appLogger *slog.Logger,
	closer *closer.Closer,
	pool *pgxpool.Pool,
) *Server {
	return &Server{
		Config:  cfg,
		Logger:  appLogger,
		Closer:  closer,
		Handler: handlers.NewAppHandlerHandler(appLogger, pool),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/parse/", s.Handler.ParseHandler.ServeHTTP)

	handler := middleware.LoggingMiddleware(mux)

	srv := &http.Server{
		Addr:         ":" + s.Config.HTTPPort(),
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.TimeoutConfig.ServerReadTimeout()) * time.Second,
		WriteTimeout: time.Duration(s.Config.TimeoutConfig.ServerWriteTimeout()) * time.Second,
		IdleTimeout:  time.Duration(s.Config.TimeoutConfig.ServerIdleTimeout()) * time.Second,
	}

	s.Closer.Add("http server", srv.Shutdown)

	errCh := make(chan error, 1)

	go func() {
		s.Logger.Info("server.Start: Starting server on port " + s.Config.HTTPPort())

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error("server.Start: listen failed", "error", err)

			errCh <- fmt.Errorf("server.Start: %w", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		s.Logger.Error("server.Start: error occurred", "error", err)

		return fmt.Errorf("server.Start: %w", err)
	case sig := <-sigCh:
		s.Logger.Info("server.Start: received signal", "signal", sig.String())

		shutdownCtx, cancel := context.WithTimeout(
			ctx,
			time.Duration(s.Config.TimeoutConfig.ShutdownCtxTimeout())*time.Second,
		)

		defer cancel()

		err := s.Closer.Close(shutdownCtx)
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("graceful shutdown: %w", err)
			}

			s.Logger.Warn("shutdown timed out, forcing close")
		}
	}

	s.Logger.Info("server stopped")

	return nil
}
