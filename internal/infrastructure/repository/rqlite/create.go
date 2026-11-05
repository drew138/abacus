package rqlite

import (
	"context"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Create(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapCreateMovementRow(movement)
	at := now()
	wr, err := r.conn.WriteOneParameterizedContext(ctx, gorqlite.ParameterizedStatement{
		Query: `INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
		        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		Arguments: []any{row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Note, row.Source, row.Edited, at, at},
	})
	if err != nil {
		return domain.Movement{}, domainError(err, "inserting a movement")
	}
	movement.ID = wr.LastInsertID
	return movement, nil
}
