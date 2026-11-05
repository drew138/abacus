package rqlite

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) SetDeleted(ctx context.Context, id int64, deleted bool) error {
	at := now()
	var deletedAt any
	if deleted {
		deletedAt = at
	}
	wr, err := r.conn.WriteOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query:     "UPDATE movements SET deleted_at = ?, updated_at = ? WHERE id = ?",
		Arguments: []any{deletedAt, at, id},
	})
	if err != nil {
		return domainError(err, fmt.Sprintf("marking movement %d", id))
	}
	if wr.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
