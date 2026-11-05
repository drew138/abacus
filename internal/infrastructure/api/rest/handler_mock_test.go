package rest

import (
	"context"
	"testing"
	"time"

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

func (m *MockApplication) ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (int, int, int, error) {
	args := m.Called(ctx, movements, categories)
	return args.Int(0), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockApplication) ListMonths(ctx context.Context) ([]domain.MonthSummary, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.MonthSummary), args.Error(1)
}

func (m *MockApplication) ListMovements(ctx context.Context, month string) ([]domain.Movement, error) {
	args := m.Called(ctx, month)
	return args.Get(0).([]domain.Movement), args.Error(1)
}

func (m *MockApplication) CategoryTotals(ctx context.Context) (map[string]map[string]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]map[string]int64), args.Error(1)
}

func (m *MockApplication) ListCategories(ctx context.Context) ([]domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockApplication) GetMovement(ctx context.Context, id int64) (domain.Movement, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockApplication) CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	args := m.Called(ctx, movement)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockApplication) UpdateMovement(ctx context.Context, movement domain.Movement, applyToMerchant bool) (domain.Movement, error) {
	args := m.Called(ctx, movement, applyToMerchant)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockApplication) RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error) {
	args := m.Called(ctx, merchant, category)
	return args.Int(0), args.Error(1)
}

func (m *MockApplication) DeleteMovement(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApplication) RestoreMovement(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApplication) Register(ctx context.Context, username, password string) (domain.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockApplication) RegistrationOpen(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockApplication) LogIn(ctx context.Context, username, password string) (string, time.Time, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockApplication) Authenticate(ctx context.Context, token string) (domain.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockApplication) LogOut(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
