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

func (m *MockRepository) Import(ctx context.Context, movements []domain.Movement) (domain.ImportResult, error) {
	args := m.Called(ctx, movements)
	return args.Get(0).(domain.ImportResult), args.Error(1)
}

func (m *MockRepository) Months(ctx context.Context) ([]domain.MonthSummary, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.MonthSummary), args.Error(1)
}

func (m *MockRepository) ListByMonth(ctx context.Context, month string) ([]domain.Movement, error) {
	args := m.Called(ctx, month)
	return args.Get(0).([]domain.Movement), args.Error(1)
}

func (m *MockRepository) Totals(ctx context.Context) ([]domain.CategoryTotal, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.CategoryTotal), args.Error(1)
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

func (m *MockRepository) Recategorise(ctx context.Context, merchantKey, category string) (int, error) {
	args := m.Called(ctx, merchantKey, category)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) MerchantRules(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
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

func (m *MockRepository) CreateFirstUser(ctx context.Context, username, passwordHash string) (domain.User, error) {
	args := m.Called(ctx, username, passwordHash)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRepository) CreateSession(ctx context.Context, session domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRepository) FindSession(ctx context.Context, tokenHash string) (domain.Session, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(domain.Session), args.Error(1)
}

func (m *MockRepository) DeleteSession(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

type MockPasswords struct {
	mock.Mock
}

func newMockPasswords(t *testing.T) *MockPasswords {
	t.Helper()
	m := &MockPasswords{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockPasswords) Hash(plain string) (string, error) {
	args := m.Called(plain)
	return args.String(0), args.Error(1)
}

func (m *MockPasswords) Matches(hash, plain string) (bool, error) {
	args := m.Called(hash, plain)
	return args.Bool(0), args.Error(1)
}

type MockTokens struct {
	mock.Mock
}

func newMockTokens(t *testing.T) *MockTokens {
	t.Helper()
	m := &MockTokens{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockTokens) New() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockClock struct {
	mock.Mock
}

func newMockClock(t *testing.T) *MockClock {
	t.Helper()
	m := &MockClock{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mocks struct {
	repo      *MockRepository
	passwords *MockPasswords
	tokens    *MockTokens
	clock     *MockClock
}

func newMocks(t *testing.T) mocks {
	t.Helper()
	return mocks{
		repo:      newMockRepository(t),
		passwords: newMockPasswords(t),
		tokens:    newMockTokens(t),
		clock:     newMockClock(t),
	}
}

func (m mocks) app(sessionTTL time.Duration) *Application {
	a := New(m.repo, sessionTTL)
	a.now = m.clock.Now
	a.newToken = m.tokens.New
	a.hashPassword = m.passwords.Hash
	a.passwordMatches = m.passwords.Matches
	return a
}

func day(y int, mo time.Month, d int) time.Time {
	return time.Date(y, mo, d, 0, 0, 0, 0, time.UTC)
}

var errStore = fmt.Errorf("store down: %w", domain.ErrUnavailable)
