package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_UpdateMovement(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{{Name: "Groceries", Parent: "Food and eating out"}, {Name: "Food delivery", Parent: "Food and eating out"}}
	date := day(2026, 10, 3)
	imported := domain.Movement{ID: 2, Ref: "r2", Date: date, Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceStatement}
	manual := domain.Movement{ID: 6, Date: date, Merchant: "Empanadas", Amount: 12000, Category: "Groceries", Method: domain.MethodDebit, Source: domain.SourceManual}
	tests := []struct {
		name            string
		movement        domain.Movement
		applyToMerchant bool
		setupMocks      func(mocks)
		want            domain.Movement
		wantErr         error
	}{
		{
			name:     "marks an imported movement as edited and keeps its ref",
			movement: domain.Movement{ID: 2, Date: date, Merchant: "Rappi", Amount: 45000, Category: "Groceries", Method: domain.MethodCredit, Note: "tip"},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(2)).Return(imported, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
				saved := domain.Movement{ID: 2, Ref: "r2", Date: date, Merchant: "Rappi", Amount: 45000, Category: "Groceries", Method: domain.MethodCredit, Note: "tip", Source: domain.SourceStatement, Edited: true}
				m.repo.On("UpdateMovement", anyCtx, saved).Return(saved, nil)
			},
			want: domain.Movement{ID: 2, Ref: "r2", Date: date, Merchant: "Rappi", Amount: 45000, Category: "Groceries", Method: domain.MethodCredit, Note: "tip", Source: domain.SourceStatement, Edited: true},
		},
		{
			name:     "does not mark a manual movement as edited",
			movement: domain.Movement{ID: 6, Date: date, Merchant: "Empanadas", Amount: 15000, Category: "Groceries", Method: domain.MethodDebit},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(6)).Return(manual, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
				saved := domain.Movement{ID: 6, Date: date, Merchant: "Empanadas", Amount: 15000, Category: "Groceries", Method: domain.MethodDebit, Source: domain.SourceManual}
				m.repo.On("UpdateMovement", anyCtx, saved).Return(saved, nil)
			},
			want: domain.Movement{ID: 6, Date: date, Merchant: "Empanadas", Amount: 15000, Category: "Groceries", Method: domain.MethodDebit, Source: domain.SourceManual},
		},
		{
			name:            "applies the category to every movement of the merchant",
			movement:        domain.Movement{ID: 2, Date: date, Merchant: "Rappi", Amount: 42900, Category: "Groceries", Method: domain.MethodCredit},
			applyToMerchant: true,
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(2)).Return(imported, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
				saved := domain.Movement{ID: 2, Ref: "r2", Date: date, Merchant: "Rappi", Amount: 42900, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceStatement, Edited: true}
				m.repo.On("UpdateMovement", anyCtx, saved).Return(saved, nil)
				m.repo.On("RecategoriseMerchant", anyCtx, "Rappi", "Groceries").Return(4, nil)
			},
			want: domain.Movement{ID: 2, Ref: "r2", Date: date, Merchant: "Rappi", Amount: 42900, Category: "Groceries", Method: domain.MethodCredit, Source: domain.SourceStatement, Edited: true},
		},
		{
			name:     "returns not found for a missing movement",
			movement: domain.Movement{ID: 99, Category: "Groceries"},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(99)).Return(domain.Movement{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name:     "rejects an unknown category",
			movement: domain.Movement{ID: 2, Category: "Nope"},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(2)).Return(imported, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
			},
			wantErr: domain.ErrUnknownCategory,
		},
		{
			name:     "returns the error when categories cannot be listed",
			movement: domain.Movement{ID: 2, Category: "Groceries"},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(2)).Return(imported, nil)
				m.repo.On("ListCategories", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when saving fails",
			movement: domain.Movement{ID: 6, Category: "Groceries"},
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(6)).Return(manual, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
				m.repo.On("UpdateMovement", anyCtx, domain.Movement{ID: 6, Category: "Groceries", Source: domain.SourceManual}).Return(domain.Movement{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:            "returns the error when the merchant cannot be recategorised",
			movement:        domain.Movement{ID: 6, Merchant: "Empanadas", Category: "Groceries"},
			applyToMerchant: true,
			setupMocks: func(m mocks) {
				m.repo.On("GetMovement", anyCtx, int64(6)).Return(manual, nil)
				m.repo.On("ListCategories", anyCtx).Return(categories, nil)
				saved := domain.Movement{ID: 6, Merchant: "Empanadas", Category: "Groceries", Source: domain.SourceManual}
				m.repo.On("UpdateMovement", anyCtx, saved).Return(saved, nil)
				m.repo.On("RecategoriseMerchant", anyCtx, "Empanadas", "Groceries").Return(0, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().UpdateMovement(t.Context(), tt.movement, tt.applyToMerchant)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
