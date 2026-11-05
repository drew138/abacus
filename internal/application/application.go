package application

import (
	"context"
	"time"

	"github.com/drew138/abacus/internal/domain"
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
	CreateFirstUser(ctx context.Context, username, password string) (domain.User, error)
	CreateSession(ctx context.Context, userID int64, lifetime time.Duration) (token string, expiresAt time.Time, err error)
	FindSession(ctx context.Context, token string) (userID int64, err error)
	DeleteSession(ctx context.Context, token string) error
}

type Application struct {
	repo Repository
}

func New(repo Repository) *Application {
	return &Application{repo: repo}
}
