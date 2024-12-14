package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/isnastish/openai/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresController struct {
	connPool *pgxpool.Pool
}

func NewPostgresController(ctx context.Context) (*PostgresController, error) {
	postgresUrl, set := os.LookupEnv("POSTGRES_URL")
	if !set || postgresUrl == "" {
		return nil, fmt.Errorf("POSTGRES_URL is not set")
	}

	config, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to create connection pool, error: %s", err.Error())
	}

	postgres := &PostgresController{
		connPool: connPool,
	}

	if err := postgres.createTable(ctx); err != nil {
		return nil, err
	}

	return postgres, nil
}

func (pc *PostgresController) createTable(ctx context.Context) error {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres: failed to acquire database connection, error: %v", err)
	}

	defer conn.Release()

	query := `
	`

	if _, err := conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("postgres: failed to create a table, error: %v", err)
	}

	log.Logger.Info("Successfully initialized postgres database controller")

	return nil
}

func (pc *PostgresController) AddUser(ctx context.Context) error {
	return nil
}

func (pc *PostgresController) HasUser(ctx context.Context) (bool, error) {
	return false, nil
}

func (db *PostgresController) Close() error {
	db.connPool.Close()
	return nil
}
