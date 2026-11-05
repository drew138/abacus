package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/drew138/abacus/internal/domain"
)

var anyCtx any = mock.Anything

type MockRepository struct {
	mock.Mock
}

func newMockRepository(t *testing.T) *MockRepository {
	t.Helper()
	m := &MockRepository{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockRepository) Import(ctx context.Context, movements []domain.Movement) (int, int, int, error) {
	args := m.Called(ctx, movements)
	return args.Int(0), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockRepository) Months(ctx context.Context) ([]domain.MonthSummary, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.MonthSummary), args.Error(1)
}

func (m *MockRepository) ListByMonth(ctx context.Context, month string) ([]domain.Movement, error) {
	args := m.Called(ctx, month)
	return args.Get(0).([]domain.Movement), args.Error(1)
}

func (m *MockRepository) Totals(ctx context.Context) (map[string]map[string]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]map[string]int64), args.Error(1)
}

func (m *MockRepository) Find(ctx context.Context, id int64) (domain.Movement, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	args := m.Called(ctx, movement)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	args := m.Called(ctx, movement)
	return args.Get(0).(domain.Movement), args.Error(1)
}

func (m *MockRepository) SetDeleted(ctx context.Context, id int64, deleted bool) error {
	args := m.Called(ctx, id, deleted)
	return args.Error(0)
}

func (m *MockRepository) Recategorise(ctx context.Context, merchant, category string) (int, error) {
	args := m.Called(ctx, merchant, category)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context) ([]domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockRepository) Ensure(ctx context.Context, categories []domain.Category) error {
	args := m.Called(ctx, categories)
	return args.Error(0)
}

func (m *MockRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRepository) FindUser(ctx context.Context, id int64) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRepository) HasUsers(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateFirstUser(ctx context.Context, username, password string) (domain.User, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRepository) CreateSession(ctx context.Context, userID int64, lifetime time.Duration) (string, time.Time, error) {
	args := m.Called(ctx, userID, lifetime)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockRepository) FindSession(ctx context.Context, token string) (int64, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) DeleteSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

type mocks struct {
	repo *MockRepository
}

func newMocks(t *testing.T) mocks {
	t.Helper()
	return mocks{repo: newMockRepository(t)}
}

func (m mocks) app() *Application {
	return New(m.repo)
}

func day(y int, mo time.Month, d int) time.Time {
	return time.Date(y, mo, d, 0, 0, 0, 0, time.UTC)
}

var errStore = fmt.Errorf("store down: %w", domain.ErrUnavailable)
