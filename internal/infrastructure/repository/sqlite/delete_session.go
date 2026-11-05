package sqlite

import (
	"context"
)

func (r *Repository) DeleteSession(ctx context.Context, tokenHash string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash); err != nil {
		return domainError(err, "deleting a session")
	}
	return nil
}
