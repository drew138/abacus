package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Top(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category Category
		want     string
	}{
		{name: "a subcategory belongs to its parent", category: Category{Name: "Groceries", Parent: "Food and eating out"}, want: "Food and eating out"},
		{name: "a top-level category is its own top", category: Category{Name: "Gaming"}, want: "Gaming"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.category.Top())
		})
	}
}

func TestMovement_Month(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		movement Movement
		want     string
	}{
		{name: "formats year and month", movement: Movement{Date: time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)}, want: "2026-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.movement.Month())
		})
	}
}

func TestMovement_IsSavings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		movement Movement
		want     bool
	}{
		{name: "savings category", movement: Movement{Category: CategorySavings}, want: true},
		{name: "spending category", movement: Movement{Category: "Groceries"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.movement.IsSavings())
		})
	}
}

func TestMethod_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method Method
		want   bool
	}{
		{name: "credit", method: MethodCredit, want: true},
		{name: "debit", method: MethodDebit, want: true},
		{name: "cash is not tracked", method: "cash", want: false},
		{name: "unknown", method: "card", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.method.Valid())
		})
	}
}

func TestSource_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source Source
		want   bool
	}{
		{name: "statement", source: SourceStatement, want: true},
		{name: "manual", source: SourceManual, want: true},
		{name: "unknown", source: "import", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.source.Valid())
		})
	}
}
