package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_GetMovement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       domain.Movement
		wantErr    error
	}{
		{
			name: "returns the movement",
			setupMocks: func(m mocks) {
				m.repo.On("Find", anyCtx, int64(3)).Return(domain.Movement{ID: 3, Merchant: "Éxito"}, nil)
			},
			want: domain.Movement{ID: 3, Merchant: "Éxito"},
		},
		{
			name: "returns not found",
			setupMocks: func(m mocks) {
				m.repo.On("Find", anyCtx, int64(3)).Return(domain.Movement{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).GetMovement(t.Context(), 3)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
