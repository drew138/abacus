//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_ListCategoryTotals(t *testing.T) {
	tests := []struct {
		name      string
		movements []domain.Movement
		want      map[string]map[string]int64
	}{
		{
			name: "totals each category per month, without deleted movements",
			movements: []domain.Movement{
				{Date: day(2026, 10, 1), Merchant: "D1", Amount: 10000, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceManual},
				{Date: day(2026, 10, 9), Merchant: "Ara", Amount: 5000, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceManual},
				{Date: day(2026, 10, 2), Merchant: "To Hapi", Amount: 500000, Category: domain.CategorySavings, Method: domain.MethodDebit, Source: domain.SourceManual},
				{Date: day(2026, 9, 3), Merchant: "Uber", Amount: 9000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceManual},
				{Date: day(2026, 10, 3), Merchant: "Gone", Amount: 1, Category: "Groceries", Method: domain.MethodDebit, Source: domain.SourceManual, Deleted: true},
			},
			want: map[string]map[string]int64{
				"2026-09": {"Rides and transport": 9000},
				"2026-10": {"Groceries": 15000, domain.CategorySavings: 500000},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			db.seed(t, tt.movements...)

			got, err := db.ListCategoryTotals(t.Context())

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
