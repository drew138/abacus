//go:build integration

package rqlite

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/rqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rqlite/gorqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/drew138/abacus/internal/domain"
)

type testRqlite struct {
	*Repository
	conn *gorqlite.Connection
}

func newRqlite(t *testing.T) *testRqlite {
	t.Helper()

	container, err := testcontainers.Run(t.Context(), "rqlite/rqlite:10.6.6",
		testcontainers.WithEnv(map[string]string{"ENABLE_FK": "1"}),
		testcontainers.WithExposedPorts("4001/tcp"),
		testcontainers.WithWaitStrategy(wait.ForHTTP("/readyz").WithPort("4001/tcp")),
	)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, container.Terminate(context.Background())) })

	endpoint, err := container.PortEndpoint(t.Context(), "4001/tcp", "")
	require.NoError(t, err)

	m, err := migrate.New("file://../../../../migration", "rqlite://"+endpoint+"?x-connect-insecure=true&disableClusterDiscovery=true")
	require.NoError(t, err)
	require.NoError(t, m.Up())

	conn, err := gorqlite.Open("http://" + endpoint + "?disableClusterDiscovery=true")
	require.NoError(t, err)
	require.NoError(t, conn.SetExecutionWithTransaction(true))

	db := &testRqlite{Repository: New(conn), conn: conn}
	t.Cleanup(func() { db.reset(t) })
	return db
}

func (d *testRqlite) reset(t *testing.T) {
	t.Helper()
	_, err := d.conn.Write([]string{
		"DELETE FROM sessions",
		"DELETE FROM users",
		"DELETE FROM merchant_rules",
		"DELETE FROM movements",
	})
	require.NoError(t, err)
}

func (d *testRqlite) seed(t *testing.T, movements ...domain.Movement) []domain.Movement {
	t.Helper()
	saved := make([]domain.Movement, 0, len(movements))
	for _, m := range movements {
		created, err := d.Create(t.Context(), m)
		require.NoError(t, err)
		if m.Deleted {
			require.NoError(t, d.SetDeleted(t.Context(), created.ID, true))
			created.Deleted = true
		}
		saved = append(saved, created)
	}
	return saved
}

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
