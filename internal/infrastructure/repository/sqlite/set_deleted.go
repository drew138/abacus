package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) SetDeleted(ctx context.Context, id int64, deleted bool) error {
	deletedAt := sql.NullString{String: now(), Valid: deleted}
	res, err := r.db.ExecContext(ctx, "UPDATE movements SET deleted_at = ?, updated_at = ? WHERE id = ?", deletedAt, now(), id)
	if err != nil {
		return domainError(err, fmt.Sprintf("marking movement %d", id))
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
