package rqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) Import(ctx context.Context, movements []domain.Movement) (added, updated, kept int, err error) {
	if len(movements) == 0 {
		return 0, 0, 0, nil
	}
	refs := make([]any, len(movements))
	for i, m := range movements {
		refs[i] = m.Ref
	}
	queries, err := r.conn.QueryParameterizedContext(ctx, []gorqlite.ParameterizedStatement{
		{Query: "SELECT merchant_key, category FROM merchant_rules"},
		{Query: "SELECT ref, edited, coalesce(deleted_at, '') FROM movements WHERE ref IN (?" + strings.Repeat(", ?", len(refs)-1) + ")", Arguments: refs},
	})
	if err != nil {
		return 0, 0, 0, domainError(err, "selecting rules and imported movements")
	}
	rules := map[string]string{}
	for queries[0].Next() {
		var key, category string
		if err := queries[0].Scan(&key, &category); err != nil {
			return 0, 0, 0, fmt.Errorf("reading a merchant rule: %w: %w", domain.ErrInternal, err)
		}
		rules[key] = category
	}
	existing := map[string]importExistingRow{}
	for queries[1].Next() {
		var row importExistingRow
		if err := queries[1].Scan(&row.Ref, &row.Edited, &row.DeletedAt); err != nil {
			return 0, 0, 0, fmt.Errorf("reading an imported movement: %w: %w", domain.ErrInternal, err)
		}
		existing[row.Ref] = row
	}

	at := now()
	var statements []gorqlite.ParameterizedStatement
	var inserts []bool
	for _, m := range movements {
		row := mapImportMovementRow(m, rules)
		current, found := existing[row.Ref]
		switch {
		case !found:
			statements = append(statements, gorqlite.ParameterizedStatement{
				Query: `INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
				        VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, 0, ?, ?)`,
				Arguments: []any{row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Source, at, at},
			})
			inserts = append(inserts, true)
		case current.Edited || current.DeletedAt != "":
			kept++
		default:
			statements = append(statements, gorqlite.ParameterizedStatement{
				Query: `UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, updated_at = ?
				        WHERE ref = ? AND (occurred_on <> ? OR merchant <> ? OR amount <> ? OR category <> ? OR method <> ?)`,
				Arguments: []any{row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, at,
					row.Ref, row.OccurredOn, row.Merchant, row.Amount, row.Category, row.Method},
			})
			inserts = append(inserts, false)
		}
	}
	if len(statements) == 0 {
		return 0, 0, kept, nil
	}

	results, err := r.conn.WriteParameterizedContext(ctx, statements)
	if err != nil {
		return 0, 0, 0, domainError(err, fmt.Sprintf("writing %d imported movements", len(statements)))
	}
	for i, result := range results {
		switch {
		case inserts[i]:
			added++
		case result.RowsAffected > 0:
			updated++
		default:
			kept++
		}
	}
	return added, updated, kept, nil
}
