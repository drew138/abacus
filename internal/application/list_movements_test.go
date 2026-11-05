package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_ListMovements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       []domain.Movement
		wantErr    error
	}{
		{
			name: "returns the month's movements",
			setupMocks: func(m mocks) {
				m.repo.On("ListByMonth", anyCtx, "2026-10").Return([]domain.Movement{{ID: 1}}, nil)
			},
			want: []domain.Movement{{ID: 1}},
		},
		{
			name: "returns the error",
			setupMocks: func(m mocks) {
				m.repo.On("ListByMonth", anyCtx, "2026-10").Return([]domain.Movement(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().ListMovements(t.Context(), "2026-10")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
