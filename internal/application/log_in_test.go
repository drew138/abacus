package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/token"
)

func TestApplication_LogIn(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 4, 20, 0, 0, 0, time.UTC)
	ttl := 30 * 24 * time.Hour
	user := domain.User{ID: 1, Username: "andres", PasswordHash: "hash"}
	tests := []struct {
		name        string
		setupMocks  func(mocks)
		wantToken   string
		wantSession domain.Session
		wantErr     error
	}{
		{
			name: "creates a session that expires after the ttl",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(user, nil)
				m.passwords.On("Matches", anyCtx, "hash", "secret password").Return(true, nil)
				m.tokens.On("New", anyCtx).Return("raw-token", nil)
				m.clock.On("Now").Return(now)
				m.sessions.On("CreateSession", anyCtx, domain.Session{TokenHash: token.Hash("raw-token"), UserID: 1, ExpiresAt: now.Add(ttl)}).Return(nil)
			},
			wantToken:   "raw-token",
			wantSession: domain.Session{TokenHash: token.Hash("raw-token"), UserID: 1, ExpiresAt: now.Add(ttl)},
		},
		{
			name: "rejects an unknown user with the same error as a wrong password",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(domain.User{}, domain.ErrNotFound)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "rejects a wrong password",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(user, nil)
				m.passwords.On("Matches", anyCtx, "hash", "secret password").Return(false, nil)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "returns the error when the user cannot be loaded",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(domain.User{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name: "returns the error when the password cannot be checked",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(user, nil)
				m.passwords.On("Matches", anyCtx, "hash", "secret password").Return(false, domain.ErrInternal)
			},
			wantErr: domain.ErrInternal,
		},
		{
			name: "returns the error when no token can be made",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(user, nil)
				m.passwords.On("Matches", anyCtx, "hash", "secret password").Return(true, nil)
				m.tokens.On("New", anyCtx).Return("", domain.ErrInternal)
			},
			wantErr: domain.ErrInternal,
		},
		{
			name: "returns the error when the session cannot be saved",
			setupMocks: func(m mocks) {
				m.users.On("FindByUsername", anyCtx, "andres").Return(user, nil)
				m.passwords.On("Matches", anyCtx, "hash", "secret password").Return(true, nil)
				m.tokens.On("New", anyCtx).Return("raw-token", nil)
				m.clock.On("Now").Return(now)
				m.sessions.On("CreateSession", anyCtx, domain.Session{TokenHash: token.Hash("raw-token"), UserID: 1, ExpiresAt: now.Add(ttl)}).Return(errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			gotToken, gotSession, err := m.app(ttl).LogIn(t.Context(), "andres", "secret password")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantToken, gotToken)
			assert.Equal(t, tt.wantSession, gotSession)
		})
	}
}
