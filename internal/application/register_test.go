package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		password   string
		setupMocks func(mocks)
		want       domain.User
		wantErr    error
	}{
		{
			name:     "creates the first user and clears the hash",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, nil)
				m.repo.On("CreateFirstUser", anyCtx, "andres", "correct horse battery").Return(domain.User{ID: 1, Username: "andres", PasswordHash: "hash"}, nil)
			},
			want: domain.User{ID: 1, Username: "andres"},
		},
		{
			name:     "refuses once a user exists",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(true, nil)
			},
			wantErr: domain.ErrRegistrationClosed,
		},
		{
			name:     "rejects a password shorter than twelve characters",
			password: "short pass",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, nil)
			},
			wantErr: domain.ErrPasswordTooShort,
		},
		{
			name:     "returns closed when another registration won the race",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, nil)
				m.repo.On("CreateFirstUser", anyCtx, "andres", "correct horse battery").Return(domain.User{}, domain.ErrRegistrationClosed)
			},
			wantErr: domain.ErrRegistrationClosed,
		},
		{
			name:     "returns the error when users cannot be counted",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
		{
			name:     "returns the error when the user cannot be saved",
			password: "correct horse battery",
			setupMocks: func(m mocks) {
				m.repo.On("HasUsers", anyCtx).Return(false, nil)
				m.repo.On("CreateFirstUser", anyCtx, "andres", "correct horse battery").Return(domain.User{}, errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app().Register(t.Context(), "andres", tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
