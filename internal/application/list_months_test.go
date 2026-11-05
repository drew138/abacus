package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_ListMonths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       []domain.MonthSummary
		wantErr    error
	}{
		{
			name: "returns the months",
			setupMocks: func(m mocks) {
				m.repo.On("ListMonths", anyCtx).Return([]domain.MonthSummary{{Month: "2026-10", Spent: 1}}, nil)
			},
			want: []domain.MonthSummary{{Month: "2026-10", Spent: 1}},
		},
		{
			name: "returns the error",
			setupMocks: func(m mocks) {
				m.repo.On("ListMonths", anyCtx).Return([]domain.MonthSummary(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().ListMonths(t.Context())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
