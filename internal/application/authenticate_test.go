package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func TestApplication_Authenticate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 4, 20, 0, 0, 0, time.UTC)
	expires := now.Add(time.Hour)
	tests := []struct {
		name       string
		raw        string
		setupMocks func(mocks)
		want       domain.User
		wantErr    error
	}{
		{
			name: "returns the user of a live session",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(1), expires, nil)
				m.clock.On("Now").Return(now)
				m.repo.On("FindUser", anyCtx, int64(1)).Return(domain.User{ID: 1, Username: "andres"}, nil)
			},
			want: domain.User{ID: 1, Username: "andres"},
		},
		{
			name:       "rejects an empty token without a lookup",
			raw:        "",
			setupMocks: func(mocks) {},
			wantErr:    domain.ErrNotAuthenticated,
		},
		{
			name: "rejects an unknown token",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(0), time.Time{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotAuthenticated,
		},
		{
			name: "rejects an expired session",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(1), expires, nil)
				m.clock.On("Now").Return(now.Add(2 * time.Hour))
			},
			wantErr: domain.ErrNotAuthenticated,
		},
		{
			name: "rejects a session whose user is gone",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(1), expires, nil)
				m.clock.On("Now").Return(now)
				m.repo.On("FindUser", anyCtx, int64(1)).Return(domain.User{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrNotAuthenticated,
		},
		{
			name: "returns the error when sessions cannot be read",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(0), time.Time{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name: "returns the error when the user cannot be read",
			raw:  "raw",
			setupMocks: func(m mocks) {
				m.repo.On("FindSession", anyCtx, token.Hash("raw")).Return(int64(1), expires, nil)
				m.clock.On("Now").Return(now)
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

			got, err := m.app(0).Authenticate(t.Context(), tt.raw)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
