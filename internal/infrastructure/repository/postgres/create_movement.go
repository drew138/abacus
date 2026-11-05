package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapCreateMovementRow(movement)
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		row.Ref, row.OccurredOn.Format(dateLayout), row.Merchant, merchantKey(row.Merchant), row.Amount, row.Category, row.Method, row.Note, row.Source, row.Edited).Scan(&row.ID)
	if err != nil {
		return domain.Movement{}, domainError(err, "inserting a movement")
	}
	return row.toDomain(), nil
}
