package postgres

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) RestoreMovement(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "UPDATE movements SET deleted_at = NULL, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return domainError(err, fmt.Sprintf("marking movement %d", id))
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
