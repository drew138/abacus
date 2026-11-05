package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) CreateSession(ctx context.Context, userID int64, lifetime time.Duration) (string, time.Time, error) {
	raw, err := token.Create()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("creating a session token: %w: %w", domain.ErrInternal, err)
	}
	var expiresAt time.Time
	err = r.db.QueryRowContext(ctx,
		"INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, now() + $3 * interval '1 second') RETURNING expires_at",
		token.Hash(raw), userID, lifetime.Seconds()).Scan(&expiresAt)
	if err != nil {
		return "", time.Time{}, domainError(err, fmt.Sprintf("inserting a session for user %d", userID))
	}
	return raw, expiresAt, nil
}
