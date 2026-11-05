package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/drew138/abacus/migration"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := loadConfig()
	db, err := sql.Open("sqlite", "file:"+cfg.databasePath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	source, err := iofs.New(migration.Files, ".")
	if err != nil {
		return fmt.Errorf("reading migrations: %w", err)
	}
	target, err := migratesqlite.WithInstance(db, &migratesqlite.Config{})
	if err != nil {
		return fmt.Errorf("preparing database: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "sqlite", target)
	if err != nil {
		return fmt.Errorf("preparing migrations: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}
	version, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("reading version: %w", err)
	}
	slog.Info("database up to date", "version", version)
	return nil
}
