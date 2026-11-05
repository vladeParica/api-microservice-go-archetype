// Command api starts the HTTP API microservice.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladeParica/api-microservice-go-archetype/internal/config"
	"github.com/vladeParica/api-microservice-go-archetype/internal/database"
	"github.com/vladeParica/api-microservice-go-archetype/internal/health"
	"github.com/vladeParica/api-microservice-go-archetype/internal/logger"
	"github.com/vladeParica/api-microservice-go-archetype/internal/postgres"
	"github.com/vladeParica/api-microservice-go-archetype/internal/server"
	"github.com/vladeParica/api-microservice-go-archetype/internal/user"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "configs/config.yaml", "path to the configuration file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}
	if version != "dev" {
		cfg.App.Version = version
	}

	log, err := logger.New(os.Stdout, cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return err
	}
	log = log.With(
		slog.String("service", cfg.App.Name),
		slog.String("version", cfg.App.Version),
		slog.String("env", cfg.App.Environment),
	)
	slog.SetDefault(log)

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer pool.Close()

	users := user.NewHandler(user.NewService(user.NewPostgresRepository(database.New(pool))))
	probes := health.NewHandler(cfg.App.Version, map[string]health.Pinger{"postgres": pool})

	mux := http.NewServeMux()
	probes.RegisterRoutes(mux)
	users.RegisterRoutes(mux)

	handler := server.Chain(mux, server.RequestLogger(log), server.Recoverer)

	return server.New(cfg.HTTP, handler, log).Run(ctx)
}
