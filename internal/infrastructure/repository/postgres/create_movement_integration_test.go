//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_CreateMovement(t *testing.T) {
	tests := []struct {
		name     string
		movement domain.Movement
		wantErr  error
	}{
		{
			name:     "stores a movement and returns its id",
			movement: domain.Movement{Date: day(2026, 10, 18), Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés", Method: domain.MethodDebit, Note: "street", Source: domain.SourceManual},
		},
		{
			name:     "rejects a category that does not exist",
			movement: domain.Movement{Date: day(2026, 10, 18), Merchant: "X", Amount: 1, Category: "Nope", Method: domain.MethodDebit, Source: domain.SourceManual},
			wantErr:  domain.ErrInvalid,
		},
		{
			name:     "rejects a non-positive amount",
			movement: domain.Movement{Date: day(2026, 10, 18), Merchant: "X", Amount: 0, Category: "Other", Method: domain.MethodDebit, Source: domain.SourceManual},
			wantErr:  domain.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)

			created, err := db.CreateMovement(t.Context(), tt.movement)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotZero(t, created.ID)
			found, err := db.GetMovement(t.Context(), created.ID)
			require.NoError(t, err)
			want := tt.movement
			want.ID = created.ID
			assert.Equal(t, want, found)
		})
	}
}
