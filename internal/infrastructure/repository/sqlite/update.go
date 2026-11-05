package sqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Update(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapMovementRow(movement)
	res, err := r.db.ExecContext(ctx,
		`UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, note = ?, edited = ?, updated_at = ?
		 WHERE id = ?`,
		row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Note, row.Edited, now(), row.ID)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("updating movement %d", row.ID))
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.Movement{}, domain.ErrNotFound
	}
	return r.Find(ctx, row.ID)
}
