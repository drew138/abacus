package application

import (
	"context"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type Repository interface {
	ImportMovements(ctx context.Context, movements []domain.Movement) (added, updated, kept int, err error)
	ListMonths(ctx context.Context) ([]domain.MonthSummary, error)
	ListMovements(ctx context.Context, month string) ([]domain.Movement, error)
	ListCategoryTotals(ctx context.Context) (map[string]map[string]int64, error)
	GetMovement(ctx context.Context, id int64) (domain.Movement, error)
	CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	UpdateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	DeleteMovement(ctx context.Context, id int64) error
	RestoreMovement(ctx context.Context, id int64) error
	RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	EnsureCategories(ctx context.Context, categories []domain.Category) error
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	CountUsers(ctx context.Context) (int, error)
	CreateFirstUser(ctx context.Context, username, password string) (domain.User, error)
	CreateSession(ctx context.Context, userID int64, lifetime time.Duration) (token string, expiresAt time.Time, err error)
	GetSession(ctx context.Context, token string) (userID int64, err error)
	DeleteSession(ctx context.Context, token string) error
}

type Application struct {
	repo Repository
}

func New(repo Repository) *Application {
	return &Application{repo: repo}
}
