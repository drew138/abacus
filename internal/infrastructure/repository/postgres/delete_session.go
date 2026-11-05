package postgres

import (
	"context"

	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) DeleteSession(ctx context.Context, raw string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = $1", token.Hash(raw)); err != nil {
		return domainError(err, "deleting a session")
	}
	return nil
}
