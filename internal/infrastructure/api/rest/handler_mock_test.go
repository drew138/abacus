package rest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/drew138/abacus/internal/domain"
)

var anyCtx any = mock.Anything

const testToken = "test-token-0123456789abcdef0123456789"

type MockApplication struct {
	mock.Mock
}

func newMockApplication(t *testing.T) *MockApplication {
	t.Helper()
	m := &MockApplication{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockApplication) ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (domain.ImportResult, error) {
	args := m.Called(ctx, movements, categories)
	return args.Get(0).(domain.ImportResult), args.Error(1)
}
