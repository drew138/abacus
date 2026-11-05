package rqlite

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

const (
	dateLayout = "2006-01-02"

	codeBusy              = 5
	codeLocked            = 6
	codeConstraintCheck   = 275
	codeConstraintFK      = 787
	codeConstraintNotNull = 1299
	codeConstraintPK      = 1555
	codeConstraintUnique  = 2067
)

var constraintErrors = map[string]error{
	"users.username": domain.ErrConflict,
	"movements.ref":  domain.ErrConflict,
}

type Repository struct {
	conn *gorqlite.Connection
}

func New(conn *gorqlite.Connection) *Repository {
	return &Repository{conn: conn}
}

func merchantKey(merchant string) string {
	return strings.ToLower(strings.Join(strings.Fields(merchant), " "))
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func domainError(err error, doing string) error {
	var statement *gorqlite.StatementError
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &netErr):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	case errors.As(err, &statement) && statement.SQLite != nil:
		return sqliteError(statement, doing)
	case errors.As(err, &statement):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
	}
	return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
}

func sqliteError(err *gorqlite.StatementError, doing string) error {
	switch err.SQLite.ExtendedCode {
	case codeConstraintUnique, codeConstraintPK:
		if mapped, ok := constraintErrors[uniqueTarget(err.Message)]; ok {
			return fmt.Errorf("%s: %w", doing, mapped)
		}
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrConflict, err)
	case codeConstraintFK, codeConstraintNotNull, codeConstraintCheck:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInvalid, err)
	}
	switch err.SQLite.Code {
	case codeBusy, codeLocked:
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
