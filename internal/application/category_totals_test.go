package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_CategoryTotals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       map[string]map[string]int64
		wantErr    error
	}{
		{
			name: "returns the totals by month and category",
			setupMocks: func(m mocks) {
				m.repo.On("Totals", anyCtx).Return(map[string]map[string]int64{"2026-10": {"Groceries": 80000}}, nil)
			},
			want: map[string]map[string]int64{"2026-10": {"Groceries": 80000}},
		},
		{
			name: "returns the error",
			setupMocks: func(m mocks) {
				m.repo.On("Totals", anyCtx).Return(map[string]map[string]int64(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).CategoryTotals(t.Context())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
