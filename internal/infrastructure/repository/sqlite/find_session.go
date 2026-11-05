package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) FindSession(ctx context.Context, tokenHash string) (int64, time.Time, error) {
	var userID int64
	var expires string
	err := r.db.QueryRowContext(ctx, "SELECT user_id, expires_at FROM sessions WHERE token_hash = ?", tokenHash).Scan(&userID, &expires)
	if err != nil {
		return 0, time.Time{}, domainError(err, "selecting a session")
	}
	expiresAt, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("reading the expiry of a session: %w: %w", domain.ErrInternal, err)
	}
	return userID, expiresAt, nil
}
