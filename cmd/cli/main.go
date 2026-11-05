package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
	_ "modernc.org/sqlite"

	"github.com/drew138/abacus/internal/application"
	"github.com/drew138/abacus/internal/infrastructure/clock/system"
	"github.com/drew138/abacus/internal/infrastructure/password/bcrypt"
	"github.com/drew138/abacus/internal/infrastructure/repository/sqlite"
	"github.com/drew138/abacus/internal/infrastructure/token/random"
)

const usage = "usage: cli set-password <username>"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err := run(os.Args[1:]); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 || args[0] != "set-password" {
		return errors.New(usage)
	}
	username := strings.TrimSpace(args[1])
	if username == "" || len(username) > 60 {
		return errors.New("the username must be 1 to 60 characters")
	}

	plain, err := readPassword()
	if err != nil {
		return fmt.Errorf("reading password: %w", err)
	}
	if len(plain) > 72 {
		return errors.New("the password must be at most 72 bytes")
	}

	ctx := context.Background()
	cfg := loadConfig()
	db, err := sql.Open("sqlite", "file:"+cfg.databasePath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_txlock=immediate")
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("reaching database: %w", err)
	}

	repo := sqlite.New(db)
	app := application.New(repo, repo, repo, repo, bcrypt.New(), random.New(), system.New(), time.Hour)
	user, err := app.SetPassword(ctx, username, plain)
	if err != nil {
		return fmt.Errorf("setting password: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Password set for user %d. Every existing session was signed out.\n", user.ID)
	return nil
}

func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && line == "" {
			return "", err
		}
		return strings.TrimRight(line, "\r\n"), nil
	}
	fmt.Fprint(os.Stderr, "New password: ")
	first, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	fmt.Fprint(os.Stderr, "Repeat it: ")
	second, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if string(first) != string(second) {
		return "", errors.New("the two passwords differ")
	}
	return string(first), nil
}
