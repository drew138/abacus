package sqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Find(ctx context.Context, id int64) (domain.Movement, error) {
	row, err := scanMovement(r.db.QueryRowContext(ctx, "SELECT "+movementColumns+" FROM movements WHERE id = ?", id).Scan)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("selecting movement %d", id))
	}
	return row.toDomain(), nil
}
