package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) FindSession(ctx context.Context, raw string) (int64, error) {
	qr, err := r.conn.QueryOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "SELECT user_id FROM sessions WHERE token_hash = ? AND expires_at > ?",
		Arguments: []any{token.Hash(raw), now()},
	})
	if err != nil {
		return 0, domainError(err, "selecting a session")
	}
	if !qr.Next() {
		return 0, domain.ErrNotFound
	}
	var userID int64
	if err := qr.Scan(&userID); err != nil {
		return 0, fmt.Errorf("reading a session: %w: %w", domain.ErrInternal, err)
	}
	return userID, nil
}
