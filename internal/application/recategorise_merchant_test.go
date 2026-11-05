package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_RecategoriseMerchant(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{{Name: "Other"}, {Name: domain.CategorySavings}}
	tests := []struct {
		name       string
		merchant   string
		category   string
		setupMocks func(mocks)
		want       int
		wantErr    error
	}{
		{
			name:     "recategorises the merchant",
			merchant: "  UBER   Rides",
			category: "Other",
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
				m.repo.On("Recategorise", anyCtx, "  UBER   Rides", "Other").Return(16, nil)
			},
			want: 16,
		},
		{
			name:     "refuses to turn a merchant into savings",
			merchant: "Uber",
			category: domain.CategorySavings,
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
			},
			wantErr: domain.ErrUnknownCategory,
		},
		{
			name:     "rejects an unknown category",
			merchant: "Uber",
			category: "Nope",
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
			},
			wantErr: domain.ErrUnknownCategory,
		},
		{
			name:     "returns the error when categories cannot be listed",
			merchant: "Uber",
			category: "Other",
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when the store fails",
			merchant: "Uber",
			category: "Other",
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
				m.repo.On("Recategorise", anyCtx, "Uber", "Other").Return(0, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().RecategoriseMerchant(t.Context(), tt.merchant, tt.category)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
