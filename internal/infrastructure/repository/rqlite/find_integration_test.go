//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_Find(t *testing.T) {
	tests := []struct {
		name    string
		exists  bool
		wantErr error
	}{
		{name: "finds a deleted movement too", exists: true},
		{name: "returns not found", wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			id := int64(999)
			want := domain.Movement{Date: day(2026, 10, 1), Merchant: "Uber", Amount: 9000, Category: "Rides and transport", Method: domain.MethodCredit, Source: domain.SourceManual, Deleted: true}
			if tt.exists {
				want = db.seed(t, want)[0]
				id = want.ID
			}

			got, err := db.Find(t.Context(), id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}
