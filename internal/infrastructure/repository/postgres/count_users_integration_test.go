//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CountUsers(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		want   int
	}{
		{name: "zero while empty"},
		{name: "one once a user exists", exists: true, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)
			if tt.exists {
				_, err := db.CreateFirstUser(t.Context(), "andres", "correct horse battery")
				require.NoError(t, err)
			}

			got, err := db.CountUsers(t.Context())

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
