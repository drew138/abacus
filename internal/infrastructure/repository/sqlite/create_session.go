package sqlite

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) CreateSession(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		tokenHash, userID, expiresAt.UTC().Format(time.RFC3339)); err != nil {
		return domainError(err, fmt.Sprintf("inserting a session for user %d", userID))
	}
	return nil
}
