//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_FindUser(t *testing.T) {
	tests := []struct {
		name    string
		exists  bool
		wantErr error
	}{
		{name: "finds the user without the hash", exists: true},
		{name: "returns not found", wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			id := int64(999)
			if tt.exists {
				created, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")
				require.NoError(t, err)
				id = created.ID
			}

			got, err := db.FindUser(t.Context(), id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, domain.User{ID: id, Username: "andres"}, got)
		})
	}
}
