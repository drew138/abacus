//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_Recategorise(t *testing.T) {
	tests := []struct {
		name       string
		merchant   string
		category   string
		want       int
		wantErr    error
		wantGroups map[string]string
	}{
		{
			name:       "moves every movement of the merchant, matching case and spaces, but not savings",
			merchant:   "  UBER ",
			category:   "Other",
			want:       2,
			wantGroups: map[string]string{"Uber": "Other", "uber": "Other", "Uber savings": domain.CategorySavings, "DiDi": "Rides and transport"},
		},
		{
			name:     "rejects a category that does not exist",
			merchant: "Uber",
			category: "Nope",
			wantErr:  domain.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			db.seed(t,
				domain.Movement{Date: day(2026, 10, 1), Merchant: "Uber", Amount: 9000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceStatement},
				domain.Movement{Date: day(2026, 10, 2), Merchant: "uber", Amount: 8000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceStatement},
				domain.Movement{Date: day(2026, 10, 3), Merchant: "Uber savings", Amount: 1000, Category: domain.CategorySavings, Method: domain.MethodDebit, Source: domain.SourceStatement},
				domain.Movement{Date: day(2026, 10, 4), Merchant: "DiDi", Amount: 7000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceStatement},
			)

			got, err := db.Recategorise(t.Context(), tt.merchant, tt.category)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			movements, err := db.ListByMonth(t.Context(), "2026-10")
			require.NoError(t, err)
			categories := map[string]string{}
			for _, m := range movements {
				categories[m.Merchant] = m.Category
			}
			assert.Equal(t, tt.wantGroups, categories)
		})
	}
}
