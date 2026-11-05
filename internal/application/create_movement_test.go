package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_CreateMovement(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{{Name: "Restaurants and cafés", Parent: domain.GroupFood}, {Name: domain.CategorySavings}}
	date := day(2026, 10, 18)
	tests := []struct {
		name       string
		movement   domain.Movement
		setupMocks func(mocks)
		want       domain.Movement
		wantErr    error
	}{
		{
			name:     "saves a manual movement and ignores any ref or flags sent",
			movement: domain.Movement{Ref: "forged", Date: date, Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés", Method: domain.MethodDebit, Edited: true, Deleted: true},
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
				m.repo.On("Create", anyCtx, domain.Movement{Date: date, Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés", Method: domain.MethodDebit, Source: domain.SourceManual}).
					Return(domain.Movement{ID: 9, Date: date, Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés", Method: domain.MethodDebit, Source: domain.SourceManual}, nil)
			},
			want: domain.Movement{ID: 9, Date: date, Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés", Method: domain.MethodDebit, Source: domain.SourceManual},
		},
		{
			name:     "rejects a category that does not exist",
			movement: domain.Movement{Date: date, Merchant: "X", Amount: 1, Category: "Nope", Method: domain.MethodDebit},
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
			},
			wantErr: domain.ErrUnknownCategory,
		},
		{
			name:     "returns the error when categories cannot be listed",
			movement: domain.Movement{Category: "Restaurants and cafés"},
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when saving fails",
			movement: domain.Movement{Category: "Restaurants and cafés"},
			setupMocks: func(m mocks) {
				m.repo.On("List", anyCtx).Return(categories, nil)
				m.repo.On("Create", anyCtx, domain.Movement{Category: "Restaurants and cafés", Source: domain.SourceManual}).Return(domain.Movement{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).CreateMovement(t.Context(), tt.movement)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
