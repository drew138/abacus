package sqlite

import (
	"context"
)

func (r *Repository) MerchantRules(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT merchant_key, category FROM merchant_rules")
	if err != nil {
		return nil, domainError(err, "selecting merchant rules")
	}
	defer rows.Close()

	rules := map[string]string{}
	for rows.Next() {
		var key, category string
		if err := rows.Scan(&key, &category); err != nil {
			return nil, domainError(err, "reading a merchant rule")
		}
		rules[key] = category
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading merchant rules")
	}
	return rules, nil
}
