package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "SELECT id, username, password_hash FROM users WHERE username = ?",
		Arguments: []any{username},
	})
	if err != nil {
		return domain.User{}, domainError(err, "selecting a user")
	}
	if !qr.Next() {
		return domain.User{}, domain.ErrNotFound
	}
	var row findByUsernameRow
	if err := qr.Scan(&row.ID, &row.Username, &row.PasswordHash); err != nil {
		return domain.User{}, fmt.Errorf("reading a user: %w: %w", domain.ErrInternal, err)
	}
	return row.toDomain(), nil
}
