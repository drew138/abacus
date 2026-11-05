//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/drew138/abacus/internal/domain"
)

type testPostgres struct {
	*Repository
	db *sql.DB
}

func newPostgres(t *testing.T) *testPostgres {
	t.Helper()

	container, err := tcpostgres.Run(t.Context(), "postgres:17-alpine",
		tcpostgres.WithDatabase("abacus"),
		tcpostgres.WithUsername("abacus"),
		tcpostgres.WithPassword("abacus"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, container.Terminate(context.Background())) })

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	require.NoError(t, err)

	m, err := migrate.New("file://../../../../migration", "pgx5://"+strings.TrimPrefix(dsn, "postgres://"))
	require.NoError(t, err)
	require.NoError(t, m.Up())

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	pg := &testPostgres{Repository: New(db), db: db}
	t.Cleanup(func() { pg.reset(t) })
	return pg
}

func (p *testPostgres) reset(t *testing.T) {
	t.Helper()
	_, err := p.db.Exec("TRUNCATE sessions, users, merchant_rules, movements RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func (p *testPostgres) seed(t *testing.T, movements ...domain.Movement) []domain.Movement {
	t.Helper()
	saved := make([]domain.Movement, 0, len(movements))
	for _, m := range movements {
		created, err := p.CreateMovement(t.Context(), m)
		require.NoError(t, err)
		if m.Deleted {
			require.NoError(t, p.DeleteMovement(t.Context(), created.ID))
			created.Deleted = true
		}
		saved = append(saved, created)
	}
	return saved
}

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
