package application

import (
	"context"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type Movements interface {
	Import(ctx context.Context, movements []domain.Movement) (domain.ImportResult, error)
	Months(ctx context.Context) ([]domain.MonthSummary, error)
	ListByMonth(ctx context.Context, month string) ([]domain.Movement, error)
	Totals(ctx context.Context) ([]domain.CategoryTotal, error)
	Find(ctx context.Context, id int64) (domain.Movement, error)
	Create(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	Update(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	SetDeleted(ctx context.Context, id int64, deleted bool) error
	Recategorise(ctx context.Context, merchantKey, category string) (int, error)
	MerchantRules(ctx context.Context) (map[string]string, error)
}

type Categories interface {
	List(ctx context.Context) ([]domain.Category, error)
	Ensure(ctx context.Context, categories []domain.Category) error
}

type Users interface {
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindUser(ctx context.Context, id int64) (domain.User, error)
	SavePassword(ctx context.Context, username, passwordHash string) (domain.User, error)
}

type Sessions interface {
	CreateSession(ctx context.Context, session domain.Session) error
	FindSession(ctx context.Context, tokenHash string) (domain.Session, error)
	DeleteSession(ctx context.Context, tokenHash string) error
	DeleteSessionsForUser(ctx context.Context, userID int64) error
}

type Passwords interface {
	Hash(ctx context.Context, plain string) (string, error)
	Matches(ctx context.Context, hash, plain string) (bool, error)
}

type Tokens interface {
	New(ctx context.Context) (string, error)
}

type Clock interface {
	Now() time.Time
}

type Application struct {
	movements  Movements
	categories Categories
	users      Users
	sessions   Sessions
	passwords  Passwords
	tokens     Tokens
	clock      Clock
	sessionTTL time.Duration
}

func New(movements Movements, categories Categories, users Users, sessions Sessions, passwords Passwords, tokens Tokens, clock Clock, sessionTTL time.Duration) *Application {
	return &Application{movements: movements, categories: categories, users: users, sessions: sessions, passwords: passwords, tokens: tokens, clock: clock, sessionTTL: sessionTTL}
}
