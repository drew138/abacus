package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_ImportMovements(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{{Name: "Board games"}}
	tests := []struct {
		name        string
		movements   []domain.Movement
		setupMocks  func(mocks)
		wantAdded   int
		wantUpdated int
		wantKept    int
		wantErr     error
	}{
		{
			name: "ensures the categories and imports every movement as from a statement",
			movements: []domain.Movement{
				{Ref: "a", Merchant: "Rappi", Category: "Food delivery", Amount: 42900},
				{Ref: "b", Merchant: "To Hapi", Category: domain.CategorySavings, Amount: 1000000},
			},
			setupMocks: func(m mocks) {
				m.repo.On("EnsureCategories", anyCtx, categories).Return(nil)
				m.repo.On("ImportMovements", anyCtx, []domain.Movement{
					{Ref: "a", Merchant: "Rappi", Category: "Food delivery", Amount: 42900, Source: domain.SourceStatement},
					{Ref: "b", Merchant: "To Hapi", Category: domain.CategorySavings, Amount: 1000000, Source: domain.SourceStatement},
				}).Return(1, 0, 1, nil)
			},
			wantAdded: 1,
			wantKept:  1,
		},
		{
			name:      "returns the error when categories cannot be saved",
			movements: []domain.Movement{{Ref: "a"}},
			setupMocks: func(m mocks) {
				m.repo.On("EnsureCategories", anyCtx, categories).Return(errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:      "returns the error when the import fails",
			movements: []domain.Movement{{Ref: "a"}},
			setupMocks: func(m mocks) {
				m.repo.On("EnsureCategories", anyCtx, categories).Return(nil)
				m.repo.On("ImportMovements", anyCtx, []domain.Movement{{Ref: "a", Source: domain.SourceStatement}}).Return(0, 0, 0, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			added, updated, kept, err := m.app().ImportMovements(t.Context(), tt.movements, categories)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAdded, added)
			assert.Equal(t, tt.wantUpdated, updated)
			assert.Equal(t, tt.wantKept, kept)
		})
	}
}
