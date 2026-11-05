package postgres

import (
	"context"

	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) GetSession(ctx context.Context, raw string) (int64, error) {
	var userID int64
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > now()", token.Hash(raw)).Scan(&userID)
	if err != nil {
		return 0, domainError(err, "selecting a session")
	}
	return userID, nil
}
