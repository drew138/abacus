package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/drew138/abacus/internal/domain"
)

const dateLayout = "2006-01-02"

var constraintErrors = map[string]error{
	"users_single_key":   domain.ErrRegistrationClosed,
	"users_username_key": domain.ErrConflict,
	"movements_ref_key":  domain.ErrConflict,
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func merchantKey(merchant string) string {
	return strings.ToLower(strings.Join(strings.Fields(merchant), " "))
}

func domainError(err error, doing string) error {
	var pgErr *pgconn.PgError
	var netErr net.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrNotFound
	case errors.As(err, &pgErr):
		return postgresError(pgErr, err, doing)
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &netErr):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	}
	return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
}

func postgresError(pgErr *pgconn.PgError, err error, doing string) error {
	if mapped, ok := constraintErrors[pgErr.ConstraintName]; ok {
		return fmt.Errorf("%s: %w", doing, mapped)
	}
	switch {
	case pgErr.Code == "23505":
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrConflict, err)
	case strings.HasPrefix(pgErr.Code, "23"), strings.HasPrefix(pgErr.Code, "22"):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInvalid, err)
	case strings.HasPrefix(pgErr.Code, "08"), strings.HasPrefix(pgErr.Code, "53"), strings.HasPrefix(pgErr.Code, "57P"):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	}
	return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
}
