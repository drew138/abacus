package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	driver "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/drew138/abacus/internal/domain"
)

const dateLayout = "2006-01-02"

var constraintErrors = map[string]error{
	"users.username": domain.ErrConflict,
	"movements.ref":  domain.ErrConflict,
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

const movementColumns = "id, ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, deleted_at"

func scanMovement(scan func(...any) error) (movementRow, error) {
	var r movementRow
	err := scan(&r.ID, &r.Ref, &r.OccurredOn, &r.Merchant, &r.MerchantKey, &r.Amount, &r.Category, &r.Method, &r.Note, &r.Source, &r.Edited, &r.DeletedAt)
	return r, err
}

func domainError(err error, doing string) error {
	var sqliteErr *driver.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrNotFound
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	case errors.As(err, &sqliteErr):
		return sqliteError(sqliteErr, doing)
	default:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
	}
}

func sqliteError(err *driver.Error, doing string) error {
	switch err.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
		if mapped, ok := constraintErrors[uniqueTarget(err.Error())]; ok {
			return fmt.Errorf("%s: %w", doing, mapped)
		}
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrConflict, err)
	case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY, sqlite3.SQLITE_CONSTRAINT_NOTNULL, sqlite3.SQLITE_CONSTRAINT_CHECK:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInvalid, err)
	}
	switch err.Code() & 0xff {
	case sqlite3.SQLITE_BUSY, sqlite3.SQLITE_LOCKED, sqlite3.SQLITE_CANTOPEN, sqlite3.SQLITE_IOERR, sqlite3.SQLITE_FULL:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	}
	return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
}

func uniqueTarget(message string) string {
	_, after, ok := strings.Cut(message, "constraint failed: ")
	if !ok {
		return ""
	}
	target, _, _ := strings.Cut(after, " ")
	return strings.TrimSuffix(target, ",")
}
