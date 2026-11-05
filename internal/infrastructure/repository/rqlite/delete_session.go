package rqlite

import (
	"context"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/pkg/token"
)

func (r *Repository) DeleteSession(ctx context.Context, raw string) error {
	if _, err := r.conn.WriteOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "DELETE FROM sessions WHERE token_hash = ?",
		Arguments: []any{token.Hash(raw)},
	}); err != nil {
		return domainError(err, "deleting a session")
	}
	return nil
}
