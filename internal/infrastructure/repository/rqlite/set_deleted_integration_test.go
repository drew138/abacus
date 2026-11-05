//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_SetDeleted(t *testing.T) {
	tests := []struct {
		name        string
		deleted     bool
		exists      bool
		wantDeleted bool
		wantErr     error
	}{
		{name: "deletes", deleted: true, exists: true, wantDeleted: true},
		{name: "restores", deleted: false, exists: true, wantDeleted: false},
		{name: "returns not found", deleted: true, wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			id := int64(999)
			if tt.exists {
				id = db.seed(t, domain.Movement{Date: day(2026, 10, 1), Merchant: "Uber", Amount: 9000, Category: "Other", Method: domain.MethodCredit, Source: domain.SourceManual, Deleted: !tt.deleted})[0].ID
			}

			err := db.SetDeleted(t.Context(), id, tt.deleted)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			got, err := db.Find(t.Context(), id)
			require.NoError(t, err)
			assert.Equal(t, tt.wantDeleted, got.Deleted)
		})
	}
}
