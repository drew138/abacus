package postgres

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) UpdateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapUpdateMovementRow(movement)
	var saved updateMovementRow
	err := r.db.QueryRowContext(ctx,
		`UPDATE movements SET occurred_on = $1, merchant = $2, merchant_key = $3, amount = $4, category = $5, method = $6, note = $7, edited = $8, updated_at = now()
		 WHERE id = $9 RETURNING id, ref, occurred_on, merchant, amount, category, method, note, source, edited, deleted_at`,
		row.OccurredOn.Format(dateLayout), row.Merchant, merchantKey(row.Merchant), row.Amount, row.Category, row.Method, row.Note, row.Edited, row.ID).
		Scan(&saved.ID, &saved.Ref, &saved.OccurredOn, &saved.Merchant, &saved.Amount, &saved.Category, &saved.Method, &saved.Note, &saved.Source, &saved.Edited, &saved.DeletedAt)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("updating movement %d", row.ID))
	}
	return saved.toDomain(), nil
}
