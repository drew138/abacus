//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_ListMonths(t *testing.T) {
	tests := []struct {
		name      string
		movements []domain.Movement
		want      []domain.MonthSummary
	}{
		{
			name: "sums spending and savings per month, ignoring deleted movements",
			movements: []domain.Movement{
				{Date: day(2026, 9, 30), Merchant: "D1", Amount: 10000, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceManual},
				{Date: day(2026, 10, 1), Merchant: "Uber", Amount: 9000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceManual},
				{Date: day(2026, 10, 2), Merchant: "To Hapi", Amount: 500000, Category: domain.CategorySavings, Method: domain.MethodDebit, Source: domain.SourceManual},
				{Date: day(2026, 10, 3), Merchant: "Gone", Amount: 1, Category: "Other", Method: domain.MethodDebit, Source: domain.SourceManual, Deleted: true},
			},
			want: []domain.MonthSummary{
				{Month: "2026-09", Spent: 10000, Movements: 1},
				{Month: "2026-10", Spent: 9000, Saved: 500000, Movements: 1},
			},
		},
		{
			name: "returns no months for an empty database",
			want: []domain.MonthSummary{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			db.seed(t, tt.movements...)

			got, err := db.ListMonths(t.Context())

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
