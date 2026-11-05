package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

func TestApplication_LogInUser(t *testing.T) {
	t.Parallel()

	hash, err := password.Hash("secret password")
	require.NoError(t, err)
	user := domain.User{ID: 1, Username: "andres", PasswordHash: hash}
	expires := time.Date(2026, 12, 4, 20, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		password    string
		setupMocks  func(mocks)
		wantToken   string
		wantExpires time.Time
		wantErr     error
	}{
		{
			name:     "opens a thirty-day session",
			password: "secret password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(user, nil)
				m.repo.On("CreateSession", anyCtx, int64(1), 30*24*time.Hour).Return("raw-token", expires, nil)
			},
			wantToken:   "raw-token",
			wantExpires: expires,
		},
		{
			name:     "rejects an unknown user with the same error as a wrong password",
			password: "secret password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(domain.User{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name:     "rejects a wrong password",
			password: "wrong password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(user, nil)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name:     "returns an internal error for a corrupt hash",
			password: "secret password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(domain.User{ID: 1, PasswordHash: "not a hash"}, nil)
			},
			wantErr: domain.ErrInternal,
		},
		{
			name:     "returns the error when the user cannot be loaded",
			password: "secret password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(domain.User{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when the session cannot be created",
			password: "secret password",
			setupMocks: func(m mocks) {
				m.repo.On("GetUserByUsername", anyCtx, "andres").Return(user, nil)
				m.repo.On("CreateSession", anyCtx, int64(1), 30*24*time.Hour).Return("", time.Time{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			gotToken, gotExpires, err := m.app().LogInUser(t.Context(), "andres", tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantToken, gotToken)
			assert.Equal(t, tt.wantExpires, gotExpires)
		})
	}
}
