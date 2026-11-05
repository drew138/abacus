//go:build integration

package rqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestRepository_ListByMonth(t *testing.T) {
	tests := []struct {
		name    string
		month   string
		want    []string
		wantErr error
	}{
		{name: "lists only that month's live movements, by date", month: "2026-10", want: []string{"Uber", "Rappi"}},
		{name: "rejects a malformed month", month: "2026-13", wantErr: domain.ErrInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newRqlite(t)
			db.seed(t,
				domain.Movement{Date: day(2026, 10, 3), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit, Source: domain.SourceManual},
				domain.Movement{Date: day(2026, 10, 1), Merchant: "Uber", Amount: 9000, Category: "Other", Method: domain.MethodCredit, Source: domain.SourceManual},
				domain.Movement{Date: day(2026, 11, 1), Merchant: "November", Amount: 1, Category: "Other", Method: domain.MethodCredit, Source: domain.SourceManual},
				domain.Movement{Date: day(2026, 10, 2), Merchant: "Gone", Amount: 1, Category: "Other", Method: domain.MethodCredit, Source: domain.SourceManual, Deleted: true},
			)

			got, err := db.ListByMonth(t.Context(), tt.month)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			merchants := make([]string, 0, len(got))
			for _, m := range got {
				merchants = append(merchants, m.Merchant)
			}
			assert.Equal(t, tt.want, merchants)
		})
	}
}
