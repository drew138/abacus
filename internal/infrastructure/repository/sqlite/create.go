package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Create(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapCreateMovementRow(movement)
	at := now()
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Note, row.Source, row.Edited, at, at).Scan(&row.ID)
	if err != nil {
		return domain.Movement{}, domainError(err, "inserting a movement")
	}
	return row.toDomain(), nil
}
