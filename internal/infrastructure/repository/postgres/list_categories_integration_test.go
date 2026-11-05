//go:build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_ListCategories(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "lists the seeded categories", want: 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newPostgres(t)

			got, err := db.ListCategories(t.Context())

			require.NoError(t, err)
			require.Len(t, got, tt.want)
			assert.Contains(t, got, domain.Category{Name: "Restaurants and cafés", Parent: "Food and eating out"})
			assert.Contains(t, got, domain.Category{Name: domain.CategorySavings})
		})
	}
}
