//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{name: "finds the user with the hash", username: "andres"},
		{name: "returns not found", username: "nobody", wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			created, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")
			require.NoError(t, err)

			got, err := db.GetUserByUsername(t.Context(), tt.username)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, created.ID, got.ID)
			assert.NotEmpty(t, got.PasswordHash)
		})
	}
}
