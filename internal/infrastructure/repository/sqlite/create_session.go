package sqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) CreateSession(ctx context.Context, session domain.Session) error {
	row := mapSessionRow(session)
	if _, err := r.db.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		row.TokenHash, row.UserID, row.ExpiresAt); err != nil {
		return domainError(err, fmt.Sprintf("inserting a session for user %d", row.UserID))
	}
	return nil
}
