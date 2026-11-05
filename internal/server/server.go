// Package server runs the HTTP server and provides cross-cutting middleware.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vladeParica/api-microservice-go-archetype/internal/config"
)

// Server is an HTTP server with graceful shutdown.
type Server struct {
	srv *http.Server
	cfg config.HTTP
	log *slog.Logger
}

// New returns a Server that serves h using the given configuration.
func New(cfg config.HTTP, h http.Handler, log *slog.Logger) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		srv: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           h,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
		},
	}
}

// Run starts the server and blocks until ctx is canceled or the server fails.
// When ctx is canceled, in-flight requests are given cfg.ShutdownTimeout to
// complete before the server is closed.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("http server listening", "addr", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
	}

	s.log.Info("shutting down http server", "timeout", s.cfg.ShutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	return nil
}
