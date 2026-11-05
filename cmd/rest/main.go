package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/application"
	"github.com/drew138/abacus/internal/infrastructure/api/rest"
	"github.com/drew138/abacus/internal/infrastructure/repository/rqlite"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	conn, err := gorqlite.Open(cfg.databaseURL)
	if err != nil {
		return fmt.Errorf("reaching database: %w", err)
	}
	defer conn.Close()
	if err := conn.SetExecutionWithTransaction(true); err != nil {
		return fmt.Errorf("configuring database: %w", err)
	}

	app := application.New(rqlite.New(conn))
	router := chi.NewRouter()
	if _, err := rest.New(router, app, cfg.importToken, cfg.secureCookie); err != nil {
		return fmt.Errorf("building web interface: %w", err)
	}
	server := &http.Server{Addr: cfg.httpAddr, Handler: router}

	shutdownErr := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdownErr <- server.Shutdown(context.Background())
	}()

	slog.InfoContext(ctx, "service started", "addr", cfg.httpAddr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving http: %w", err)
	}
	if err := <-shutdownErr; err != nil {
		return fmt.Errorf("shutting down http: %w", err)
	}
	return nil
}
