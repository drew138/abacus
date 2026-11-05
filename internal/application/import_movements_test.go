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
		name       string
		movements  []domain.Movement
		setupMocks func(mocks)
		want       domain.ImportResult
		wantErr    error
	}{
		{
			name: "applies merchant rules, keeps savings and marks every movement as from a statement",
			movements: []domain.Movement{
				{Ref: "a", Merchant: "Éxito  Poblado", Category: "Groceries", Amount: 80000},
				{Ref: "b", Merchant: "Rappi", Category: "Food delivery", Amount: 42900},
				{Ref: "c", Merchant: "To Hapi", Category: domain.CategorySavings, Amount: 1000000},
			},
			setupMocks: func(m mocks) {
				m.categories.On("Ensure", anyCtx, categories).Return(nil)
				m.movements.On("MerchantRules", anyCtx).Return(map[string]string{"éxito poblado": "Snacks and candy", "to hapi": "Other"}, nil)
				m.movements.On("Import", anyCtx, []domain.Movement{
					{Ref: "a", Merchant: "Éxito  Poblado", Category: "Snacks and candy", Amount: 80000, Source: domain.SourceStatement},
					{Ref: "b", Merchant: "Rappi", Category: "Food delivery", Amount: 42900, Source: domain.SourceStatement},
					{Ref: "c", Merchant: "To Hapi", Category: domain.CategorySavings, Amount: 1000000, Source: domain.SourceStatement},
				}).Return(domain.ImportResult{Added: 3}, nil)
			},
			want: domain.ImportResult{Added: 3},
		},
		{
			name:      "returns the error when categories cannot be saved",
			movements: []domain.Movement{{Ref: "a"}},
			setupMocks: func(m mocks) {
				m.categories.On("Ensure", anyCtx, categories).Return(errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:      "returns the error when rules cannot be loaded",
			movements: []domain.Movement{{Ref: "a"}},
			setupMocks: func(m mocks) {
				m.categories.On("Ensure", anyCtx, categories).Return(nil)
				m.movements.On("MerchantRules", anyCtx).Return(map[string]string(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:      "returns the error when the import fails",
			movements: []domain.Movement{{Ref: "a", Merchant: "Uber", Category: "Rides and transport"}},
			setupMocks: func(m mocks) {
				m.categories.On("Ensure", anyCtx, categories).Return(nil)
				m.movements.On("MerchantRules", anyCtx).Return(map[string]string{}, nil)
				m.movements.On("Import", anyCtx, []domain.Movement{{Ref: "a", Merchant: "Uber", Category: "Rides and transport", Source: domain.SourceStatement}}).
					Return(domain.ImportResult{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).ImportMovements(t.Context(), tt.movements, categories)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
