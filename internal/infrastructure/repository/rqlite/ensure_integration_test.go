//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_Ensure(t *testing.T) {
	tests := []struct {
		name       string
		categories []domain.Category
		want       []domain.Category
	}{
		{
			name:       "adds new categories and leaves existing ones as they are",
			categories: []domain.Category{{Name: "Board games", Parent: "Leisure"}, {Name: "Groceries", Parent: "Changed"}},
			want:       []domain.Category{{Name: "Board games", Parent: "Leisure"}, {Name: "Groceries", Parent: "Food and eating out"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)

			require.NoError(t, db.Ensure(t.Context(), tt.categories))

			all, err := db.List(t.Context())
			require.NoError(t, err)
			byName := map[string]domain.Category{}
			for _, c := range all {
				byName[c.Name] = c
			}
			for _, w := range tt.want {
				assert.Equal(t, w, byName[w.Name])
			}
		})
	}
}
