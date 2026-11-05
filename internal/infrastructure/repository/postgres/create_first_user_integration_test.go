//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

func TestRepository_CreateFirstUser(t *testing.T) {
	tests := []struct {
		name     string
		existing bool
		wantErr  error
	}{
		{name: "creates the first user with a hashed password"},
		{name: "refuses once a user exists", existing: true, wantErr: domain.ErrRegistrationClosed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			if tt.existing {
				_, err := db.CreateFirstUser(t.Context(), "someone", "correct horse battery")
				require.NoError(t, err)
			}

			got, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "andres", got.Username)
			stored, err := db.GetUserByUsername(t.Context(), "andres")
			require.NoError(t, err)
			ok, err := password.Matches(stored.PasswordHash, "correct horse battery")
			require.NoError(t, err)
			assert.True(t, ok)
		})
	}
}
