package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCategoryTier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		cat    string
		parent string
		want   int
	}{
		{name: "savings comes first", cat: CategorySavings, want: 0},
		{name: "the food group comes second", cat: GroupFood, want: 1},
		{name: "a food subcategory comes second", cat: "Groceries", parent: GroupFood, want: 1},
		{name: "other comes last", cat: CategoryOther, want: 3},
		{name: "everything else is in between", cat: "Gaming", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, CategoryTier(tt.cat, tt.parent))
		})
	}
}

func TestCategory_Top(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category Category
		want     string
	}{
		{name: "a subcategory belongs to its parent", category: Category{Name: "Groceries", Parent: GroupFood}, want: GroupFood},
		{name: "a top-level category is its own top", category: Category{Name: "Gaming"}, want: "Gaming"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.category.Top())
		})
	}
}

func TestMerchantKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		merchant string
		want     string
	}{
		{name: "lowercases accented letters", merchant: "ÉXITO Poblado", want: "éxito poblado"},
		{name: "collapses and trims spaces", merchant: "  Crepes   Y  Waffles ", want: "crepes y waffles"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, MerchantKey(tt.merchant))
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

func TestSession_Expired(t *testing.T) {
	t.Parallel()

	expires := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{name: "before the expiry", now: expires.Add(-time.Second), want: false},
		{name: "exactly at the expiry", now: expires, want: true},
		{name: "after the expiry", now: expires.Add(time.Second), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, Session{ExpiresAt: expires}.Expired(tt.now))
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
