package testutil

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/seeu359/go-project-278/links"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func StartPostgres() (*pgxpool.Pool, func(), error) {
	ctx := context.Background()
	container, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		testcontainers.TerminateContainer(container)
		return nil, nil, err
	}

	if err := runMigrations(dsn); err != nil {
		testcontainers.TerminateContainer(container)
		return nil, nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		testcontainers.TerminateContainer(container)
		return nil, nil, err
	}
	cleanup := func() {
		pool.Close()
		testcontainers.TerminateContainer(container)
	}
	return pool, cleanup, nil
}

func Tx(t *testing.T, pool *pgxpool.Pool) *links.Queries {
	t.Helper()

	tx, err := pool.Begin(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, tx.Rollback(context.Background()))
	})

	return links.New(tx)
}
func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	err = goose.Up(db, "../migrations/")
	if err != nil {
		return err
	}
	return nil
}
