package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		token      string
		setupMocks func(mocks)
		want       domain.User
		wantErr    error
	}{
		{
			name:  "returns the user of a live session",
			token: "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, "raw").Return(int64(1), nil)
				m.repo.On("FindUser", anyCtx, int64(1)).Return(domain.User{ID: 1, Username: "andres"}, nil)
			},
			want: domain.User{ID: 1, Username: "andres"},
		},
		{
			name:       "rejects an empty token without a lookup",
			token:      "",
			setupMocks: func(mocks) {},
			wantErr:    domain.ErrNotAuthenticated,
		},
		{
			name:  "rejects an unknown or expired session",
			token: "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, "raw").Return(int64(0), domain.ErrNotFound)
			},
			wantErr: domain.ErrNotAuthenticated,
		},
		{
			name:  "rejects a session whose user is gone",
			token: "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, "raw").Return(int64(1), nil)
				m.repo.On("FindUser", anyCtx, int64(1)).Return(domain.User{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotAuthenticated,
		},
		{
			name:  "returns the error when sessions cannot be read",
			token: "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, "raw").Return(int64(0), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:  "returns the error when the user cannot be read",
			token: "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, "raw").Return(int64(1), nil)
				m.repo.On("FindUser", anyCtx, int64(1)).Return(domain.User{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().Authenticate(t.Context(), tt.token)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
