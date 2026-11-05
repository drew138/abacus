package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestApplication_ListCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mocks)
		want       []domain.Category
		wantErr    error
	}{
		{
			name: "returns the categories",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return([]domain.Category{{Name: "Gaming"}}, nil)
			},
			want: []domain.Category{{Name: "Gaming"}},
		},
		{
			name: "returns the error",
			setupMocks: func(m mocks) {
				m.categories.On("List", anyCtx).Return([]domain.Category(nil), errStore)
			},
			wantErr: domain.ErrUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMocks(t)
			tt.setupMocks(m)

			got, err := m.app(0).ListCategories(t.Context())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
