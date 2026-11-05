//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_HasUsers(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		want   bool
	}{
		{name: "false while empty"},
		{name: "true once a user exists", exists: true, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			if tt.exists {
				_, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")
				require.NoError(t, err)
			}

			got, err := db.HasUsers(t.Context())

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
