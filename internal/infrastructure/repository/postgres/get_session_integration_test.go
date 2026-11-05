//go:build integration

package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_GetSession(t *testing.T) {
	tests := []struct {
		name     string
		lifetime time.Duration
		delete   bool
		wantErr  error
	}{
		{name: "finds a live session", lifetime: time.Hour},
		{name: "does not find an expired session", lifetime: -time.Hour, wantErr: domain.ErrNotFound},
		{name: "does not find a deleted session", lifetime: time.Hour, delete: true, wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			user, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")
			require.NoError(t, err)
			token, expires, err := db.CreateSession(t.Context(), user.ID, tt.lifetime)
			require.NoError(t, err)
			require.NotEmpty(t, token)
			assert.WithinDuration(t, time.Now().Add(tt.lifetime), expires, 5*time.Second)
			if tt.delete {
				require.NoError(t, db.DeleteSession(t.Context(), token))
			}

			got, err := db.GetSession(t.Context(), token)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, user.ID, got)
		})
	}
}
