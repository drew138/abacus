package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_MonthReport(t *testing.T) {
	t.Parallel()

	categories := []domain.Category{
		{Name: domain.CategorySavings},
		{Name: "Restaurants and cafés", Parent: domain.GroupFood},
		{Name: "Groceries", Parent: domain.GroupFood},
		{Name: "Gym and fitness"},
		{Name: "Gaming"},
		{Name: domain.CategoryOther},
	}
	months := []domain.MonthSummary{
		{Month: "2026-09", Spent: 400000, Saved: 500000, Movements: 3},
		{Month: "2026-10", Spent: 517000, Saved: 1000000, Movements: 5},
	}
	totals := []domain.CategoryTotal{
		{Month: "2026-09", Category: "Restaurants and cafés", Total: 100000},
		{Month: "2026-09", Category: "Gym and fitness", Total: 295000},
		{Month: "2026-09", Category: domain.CategorySavings, Total: 500000},
		{Month: "2026-10", Category: "Restaurants and cafés", Total: 57000},
		{Month: "2026-10", Category: "Groceries", Total: 80000},
		{Month: "2026-10", Category: "Gym and fitness", Total: 295000},
		{Month: "2026-10", Category: domain.CategoryOther, Total: 450000},
		{Month: "2026-10", Category: domain.CategorySavings, Total: 1000000},
	}
	crossfit := domain.Movement{ID: 1, Date: day(2026, 10, 1), Merchant: "CrossFit Deluxe", Amount: 295000, Category: "Gym and fitness"}
	rappi := domain.Movement{ID: 2, Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 45000, Category: "Restaurants and cafés"}
	exito := domain.Movement{ID: 3, Date: day(2026, 10, 5), Merchant: "Éxito", Amount: 80000, Category: "Groceries"}
	hapi := domain.Movement{ID: 4, Date: day(2026, 10, 7), Merchant: "To Hapi", Amount: 1000000, Category: domain.CategorySavings}
	parking := domain.Movement{ID: 5, Date: day(2026, 10, 9), Merchant: "Parqueadero", Amount: 450000, Category: domain.CategoryOther}
	empanadas := domain.Movement{ID: 6, Date: day(2026, 10, 3), Merchant: "Empanadas", Amount: 12000, Category: "Restaurants and cafés"}
	movements := []domain.Movement{crossfit, rappi, exito, hapi, parking, empanadas}

	tests := []struct {
		name       string
		month      string
		filter     domain.ExpenseFilter
		setupMocks func(mocks)
		want       domain.MonthReport
		wantErr    error
	}{
		{
			name:  "orders savings, food, the rest by amount and other last, with last month's figures",
			month: "2026-10",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return(totals, nil)
				m.movements.On("ListByMonth", anyCtx, "2026-10").Return(movements, nil)
			},
			want: domain.MonthReport{
				Month:       "2026-10",
				Summary:     months[1],
				Previous:    months[0],
				HasPrevious: true,
				Groups: []domain.GroupTotal{
					{
						CategoryAmount: domain.CategoryAmount{Name: domain.GroupFood, Total: 137000, Previous: 100000, HasPrevious: true},
						Subcategories: []domain.CategoryAmount{
							{Name: "Groceries", Total: 80000, HasPrevious: true},
							{Name: "Restaurants and cafés", Total: 57000, Previous: 100000, HasPrevious: true},
						},
					},
					{CategoryAmount: domain.CategoryAmount{Name: "Gym and fitness", Total: 295000, Previous: 295000, HasPrevious: true}},
					{CategoryAmount: domain.CategoryAmount{Name: domain.CategoryOther, Total: 450000, HasPrevious: true}},
				},
				Savings: []domain.Movement{hapi},
				Expenses: []domain.ExpenseGroup{
					{Category: domain.CategorySavings, Total: 1000000, Movements: []domain.Movement{hapi}},
					{Category: "Groceries", Parent: domain.GroupFood, Total: 80000, Movements: []domain.Movement{exito}},
					{Category: "Restaurants and cafés", Parent: domain.GroupFood, Total: 57000, Movements: []domain.Movement{rappi, empanadas}},
					{Category: "Gym and fitness", Total: 295000, Movements: []domain.Movement{crossfit}},
					{Category: domain.CategoryOther, Total: 450000, Movements: []domain.Movement{parking}},
				},
			},
		},
		{
			name:   "filters expenses by group and search without changing the totals",
			month:  "2026-10",
			filter: domain.ExpenseFilter{Category: domain.GroupFood, Search: " RAP "},
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months[1:], nil)
				m.movements.On("Totals", anyCtx).Return(totals[3:], nil)
				m.movements.On("ListByMonth", anyCtx, "2026-10").Return(movements, nil)
			},
			want: domain.MonthReport{
				Month:   "2026-10",
				Summary: months[1],
				Groups: []domain.GroupTotal{
					{
						CategoryAmount: domain.CategoryAmount{Name: domain.GroupFood, Total: 137000},
						Subcategories:  []domain.CategoryAmount{{Name: "Groceries", Total: 80000}, {Name: "Restaurants and cafés", Total: 57000}},
					},
					{CategoryAmount: domain.CategoryAmount{Name: "Gym and fitness", Total: 295000}},
					{CategoryAmount: domain.CategoryAmount{Name: domain.CategoryOther, Total: 450000}},
				},
				Savings:  []domain.Movement{hapi},
				Expenses: []domain.ExpenseGroup{{Category: "Restaurants and cafés", Parent: domain.GroupFood, Total: 45000, Movements: []domain.Movement{rappi}}},
				Filter:   domain.ExpenseFilter{Category: domain.GroupFood, Search: " RAP "},
			},
		},
		{
			name:  "returns an empty report for a month with no data",
			month: "2026-12",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return(totals, nil)
				m.movements.On("ListByMonth", anyCtx, "2026-12").Return([]domain.Movement{}, nil)
			},
			want: domain.MonthReport{Month: "2026-12", Summary: domain.MonthSummary{Month: "2026-12"}, Groups: []domain.GroupTotal{}, Expenses: []domain.ExpenseGroup{}},
		},
		{
			name:  "returns the error when categories cannot be listed",
			month: "2026-10",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:  "returns the error when months cannot be listed",
			month: "2026-10",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return([]domain.MonthSummary(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:  "returns the error when totals cannot be read",
			month: "2026-10",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return([]domain.CategoryTotal(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:  "returns the error when movements cannot be listed",
			month: "2026-10",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return(categories, nil)
				m.movements.On("Months", anyCtx).Return(months, nil)
				m.movements.On("Totals", anyCtx).Return(totals, nil)
				m.movements.On("ListByMonth", anyCtx, "2026-10").Return([]domain.Movement(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).MonthReport(t.Context(), tt.month, tt.filter)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
