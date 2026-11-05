package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) CreateFirstUser(ctx context.Context, username, passwordHash string) (domain.User, error) {
	var row userRow
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash, created_at)
		 SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users)
		 RETURNING id, username`,
		username, passwordHash, now()).Scan(&row.ID, &row.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrRegistrationClosed
	}
	if err != nil {
		return domain.User{}, domainError(err, "creating the first user")
	}
	return row.toDomain(), nil
}
