package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_RegistrationOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       bool
		wantErr    error
	}{
		{
			name: "open while there are no users",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, nil)
			},
			want: true,
		},
		{
			name: "closed once a user exists",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(true, nil)
			},
			want: false,
		},
		{
			name: "returns the error",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).RegistrationOpen(t.Context())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
