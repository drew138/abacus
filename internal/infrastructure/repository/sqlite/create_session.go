package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) CreateSession(ctx context.Context, userID int64, lifetime time.Duration) (string, time.Time, error) {
	raw, err := token.New()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("creating a session token: %w: %w", domain.ErrInternal, err)
	}
	expiresAt := time.Now().Add(lifetime).UTC().Truncate(time.Second)
	if _, err := r.db.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		token.Hash(raw), userID, expiresAt.Format(time.RFC3339)); err != nil {
		return "", time.Time{}, domainError(err, fmt.Sprintf("inserting a session for user %d", userID))
	}
	return raw, expiresAt, nil
}
