//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_ImportMovements(t *testing.T) {
	rappi := domain.Movement{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}
	tests := []struct {
		name        string
		setup       func(t *testing.T, db *testPostgres)
		movements   []domain.Movement
		wantAdded   int
		wantUpdated int
		wantKept    int
		wantRows    []domain.Movement
	}{
		{
			name:      "adds new movements",
			setup:     func(*testing.T, *testPostgres) {},
			movements: []domain.Movement{rappi},
			wantAdded: 1,
			wantRows:  []domain.Movement{{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}},
		},
		{
			name: "keeps a movement that did not change",
			setup: func(t *testing.T, db *testPostgres) {
				_, _, _, err := db.ImportMovements(t.Context(), []domain.Movement{rappi})
				require.NoError(t, err)
			},
			movements: []domain.Movement{rappi},
			wantKept:  1,
			wantRows:  []domain.Movement{{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}},
		},
		{
			name: "updates a movement the statement changed",
			setup: func(t *testing.T, db *testPostgres) {
				_, _, _, err := db.ImportMovements(t.Context(), []domain.Movement{rappi})
				require.NoError(t, err)
			},
			movements:   []domain.Movement{{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 45000, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}},
			wantUpdated: 1,
			wantRows:    []domain.Movement{{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 45000, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}},
		},
		{
			name: "keeps the user's edit",
			setup: func(t *testing.T, db *testPostgres) {
				db.seed(t, domain.Movement{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 50000, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceStatement, Edited: true})
			},
			movements: []domain.Movement{rappi},
			wantKept:  1,
			wantRows:  []domain.Movement{{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 50000, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceStatement, Edited: true}},
		},
		{
			name: "does not bring back a deleted movement",
			setup: func(t *testing.T, db *testPostgres) {
				db.seed(t, domain.Movement{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement, Deleted: true})
			},
			movements: []domain.Movement{rappi},
			wantKept:  1,
			wantRows:  []domain.Movement{},
		},
		{
			name: "applies the merchant's rule but never to savings",
			setup: func(t *testing.T, db *testPostgres) {
				_, err := db.RecategoriseMerchant(t.Context(), "  RAPPI ", "Groceries")
				require.NoError(t, err)
			},
			movements: []domain.Movement{rappi, {Ref: "s1", Date: day(2026, 10, 4), Merchant: "Rappi", Amount: 100000, Category: domain.CategorySavings, Method: domain.MethodDebit, Source: domain.SourceStatement}},
			wantAdded: 2,
			wantRows: []domain.Movement{
				{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceStatement},
				{Ref: "s1", Date: day(2026, 10, 4), Merchant: "Rappi", Amount: 100000, Category: domain.CategorySavings, Method: domain.MethodDebit, Source: domain.SourceStatement},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			tt.setup(t, db)

			added, updated, kept, err := db.ImportMovements(t.Context(), tt.movements)
			require.NoError(t, err)
			assert.Equal(t, tt.wantAdded, added)
			assert.Equal(t, tt.wantUpdated, updated)
			assert.Equal(t, tt.wantKept, kept)

			got, err := db.ListMovements(t.Context(), "2026-10")
			require.NoError(t, err)
			for i := range got {
				got[i].ID = 0
			}
			assert.Equal(t, tt.wantRows, got)
		})
	}
}
