package application

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_RestoreMovement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         int64
		setupMocks func(mocks)
		wantErr    error
	}{
		{
			name: "clears the deletion",
			id:   7,
			setupMocks: func(m mocks) {
				m.repo.On("SetDeleted", anyCtx, int64(7), false).Return(nil)
			},
		},
		{
			name: "returns not found",
			id:   8,
			setupMocks: func(m mocks) {
				m.repo.On("SetDeleted", anyCtx, int64(8), false).Return(domain.ErrNotFound)
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			err := m.app().RestoreMovement(t.Context(), tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
