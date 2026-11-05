//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		exists  bool
		wantErr error
	}{
		{name: "saves the changes and returns the stored movement", exists: true},
		{name: "returns not found", wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			change := domain.Movement{ID: 999, Ref: "r1", Date: day(2026, 10, 2), Merchant: "Rappi", Amount: 50000, Category: "Groceries", Method: domain.MethodCredit, Note: "tip", Source: domain.SourceStatement, Edited: true}
			if tt.exists {
				created := db.seed(t, domain.Movement{Ref: "r1", Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement})[0]
				change.ID = created.ID
			}

			got, err := db.Update(t.Context(), change)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, change, got)
		})
	}
}
