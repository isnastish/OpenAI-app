package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NOTE: We want to have an interface in DB with Postgres being one of them
// The interface should support main methods for writing and reading the data
// So we could swith between SQL and NoSQL data storages.
// But for now lets stick with a single implementation.
// Later, when we have a clear structure, we could replace Postgres
// with MySQL for example, or keep both.

type PostgresDB struct {
	ConnPool *pgxpool.Pool
}

func NewPostgresDB() (*PostgresDB, error) {
	const DATABASE_URL string = "postgres://postgres:12345678@localhost:5432/postgres?"

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return &PostgresDB{
		ConnPool: connPool,
	}, nil
}

func (db *PostgresDB) Close() {
	db.ConnPool.Close()
}
