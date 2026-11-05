package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

func (r *Repository) CreateFirstUser(ctx context.Context, username, plain string) (domain.User, error) {
	hash, err := password.Hash(plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing the password: %w: %w", domain.ErrInternal, err)
	}

	var row userRow
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash, created_at)
		 SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users)
		 RETURNING id, username`,
		username, hash, now()).Scan(&row.ID, &row.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrRegistrationClosed
	}
	if err != nil {
		return domain.User{}, domainError(err, "creating the first user")
	}
	return row.toDomain(), nil
}
