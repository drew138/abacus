package application

import (
	"context"
	"time"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
	"github.com/drew138/abacus/pkg/token"
)

type Repository interface {
	Import(ctx context.Context, movements []domain.Movement) (added, updated, kept int, err error)
	Months(ctx context.Context) ([]domain.MonthSummary, error)
	ListByMonth(ctx context.Context, month string) ([]domain.Movement, error)
	Totals(ctx context.Context) (map[string]map[string]int64, error)
	Find(ctx context.Context, id int64) (domain.Movement, error)
	Create(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	Update(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	SetDeleted(ctx context.Context, id int64, deleted bool) error
	Recategorise(ctx context.Context, merchant, category string) (int, error)
	List(ctx context.Context) ([]domain.Category, error)
	Ensure(ctx context.Context, categories []domain.Category) error
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	FindUser(ctx context.Context, id int64) (domain.User, error)
	HasUsers(ctx context.Context) (bool, error)
	CreateFirstUser(ctx context.Context, username, passwordHash string) (domain.User, error)
	CreateSession(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error
	FindSession(ctx context.Context, tokenHash string) (userID int64, expiresAt time.Time, err error)
	DeleteSession(ctx context.Context, tokenHash string) error
}

type Application struct {
	repo            Repository
	sessionTTL      time.Duration
	now             func() time.Time
	newToken        func() (string, error)
	hashPassword    func(plain string) (string, error)
	passwordMatches func(hash, plain string) (bool, error)
}

func New(repo Repository, sessionTTL time.Duration) *Application {
	return &Application{
		repo:            repo,
		sessionTTL:      sessionTTL,
		now:             time.Now,
		newToken:        token.New,
		hashPassword:    password.Hash,
		passwordMatches: password.Matches,
	}
}
