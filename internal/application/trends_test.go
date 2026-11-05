package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_Trends(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{
		{Name: domain.CategorySavings},
		{Name: "Restaurants and cafés", Parent: domain.GroupFood},
		{Name: "Groceries", Parent: domain.GroupFood},
		{Name: "Gym and fitness"},
		{Name: "Gaming"},
		{Name: domain.CategoryOther},
	}
	months := []domain.MonthSummary{{Month: "2026-09", Spent: 425000}, {Month: "2026-10", Spent: 900000}}
	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       domain.Trends
		wantErr    error
	}{
		{
			name: "builds one row per group, ordered by tier then by the last month, without savings",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return([]domain.CategoryTotal{
					{Month: "2026-09", Category: "Restaurants and cafés", Total: 100000},
					{Month: "2026-09", Category: "Gaming", Total: 30000},
					{Month: "2026-09", Category: "Gym and fitness", Total: 295000},
					{Month: "2026-09", Category: domain.CategorySavings, Total: 500000},
					{Month: "2026-10", Category: "Restaurants and cafés", Total: 57000},
					{Month: "2026-10", Category: "Groceries", Total: 80000},
					{Month: "2026-10", Category: "Gym and fitness", Total: 295000},
					{Month: "2026-10", Category: domain.CategoryOther, Total: 468000},
					{Month: "2026-11", Category: "Gaming", Total: 999},
				}, nil)
			},
			want: domain.Trends{
				Months: months,
				Rows: []domain.TrendRow{
					{Name: domain.GroupFood, Values: []int64{100000, 137000}, Subcategories: []domain.TrendRow{
						{Name: "Groceries", Values: []int64{0, 80000}},
						{Name: "Restaurants and cafés", Values: []int64{100000, 57000}},
					}},
					{Name: "Gym and fitness", Values: []int64{295000, 295000}},
					{Name: "Gaming", Values: []int64{30000, 0}},
					{Name: domain.CategoryOther, Values: []int64{0, 468000}},
				},
			},
		},
		{
			name: "returns the error when categories cannot be listed",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name: "returns the error when months cannot be listed",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return([]domain.MonthSummary(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name: "returns the error when totals cannot be read",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return([]domain.CategoryTotal(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).Trends(t.Context())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
