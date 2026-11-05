package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_SetPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		password   string
		setupMocks func(mocks)
		want       domain.User
		wantErr    error
	}{
		{
			name:     "saves the hash, signs every session out and clears the hash",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.passwords.On("Hash", anyCtx, "correct horse battery").Return("hash", nil)
				m.users.On("SavePassword", anyCtx, "andres", "hash").Return(domain.User{ID: 1, Username: "andres", PasswordHash: "hash"}, nil)
				m.sessions.On("DeleteSessionsForUser", anyCtx, int64(1)).Return(nil)
			},
			want: domain.User{ID: 1, Username: "andres"},
		},
		{
			name:       "rejects a password shorter than twelve characters",
			password:   "short pass",
			setupMocks: func(mocks) {},
			wantErr:    domain.ErrPasswordTooShort,
		},
		{
			name:     "returns the error when hashing fails",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.passwords.On("Hash", anyCtx, "correct horse battery").Return("", domain.ErrInternal)
			},
			wantErr: domain.ErrInternal,
		},
		{
			name:     "returns the error when the user cannot be saved",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.passwords.On("Hash", anyCtx, "correct horse battery").Return("hash", nil)
				m.users.On("SavePassword", anyCtx, "andres", "hash").Return(domain.User{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when sessions cannot be cleared",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.passwords.On("Hash", anyCtx, "correct horse battery").Return("hash", nil)
				m.users.On("SavePassword", anyCtx, "andres", "hash").Return(domain.User{ID: 1, Username: "andres"}, nil)
				m.sessions.On("DeleteSessionsForUser", anyCtx, int64(1)).Return(errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).SetPassword(t.Context(), "andres", tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
