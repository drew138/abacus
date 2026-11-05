package application

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func TestApplication_LogOut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		raw        string
		setupMocks func(mocks)
		wantErr    error
	}{
		{
			name: "deletes the session",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("DeleteSession", anyCtx, token.Hash("raw")).Return(nil)
			},
		},
		{
			name:       "does nothing without a token",
			raw:        "",
			setupMocks: func(mocks) {},
		},
		{
			name: "returns the error when the store fails",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("DeleteSession", anyCtx, token.Hash("raw")).Return(errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			err := m.app(0).LogOut(t.Context(), tt.raw)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
