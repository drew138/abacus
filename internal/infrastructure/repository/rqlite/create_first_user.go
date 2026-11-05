package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

func (r *Repository) CreateFirstUser(ctx context.Context, username, plain string) (domain.User, error) {
	hash, err := password.Hash(plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing the password: %w: %w", domain.ErrInternal, err)
	}
	wr, err := r.conn.WriteOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query: `INSERT INTO users (username, password_hash, created_at)
		        SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users)`,
		Arguments: []any{username, hash, now()},
	})
	if err != nil {
		return domain.User{}, domainError(err, "creating the first user")
	}
	if wr.RowsAffected == 0 {
		return domain.User{}, domain.ErrRegistrationClosed
	}
	return domain.User{ID: wr.LastInsertID, Username: username}, nil
}
