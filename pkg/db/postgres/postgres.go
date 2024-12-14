package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	hashutils "github.com/isnastish/openai/pkg/hash_utils"
	"github.com/isnastish/openai/pkg/ipresolver"
	"github.com/isnastish/openai/pkg/log"
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

	query := `CREATE TABLE IF NOT EXISTS "users" (
		"id" SERIAL, 
		"first_name" VARCHAR(64) NOT NULL, 
		"last_name" VARCHAR(64) NOT NULL,
		"email" VARCHAR(320) NOT NULL UNIQUE,
		"password" CHARACTER(64) NOT NULL,
		"country" VARCHAR(64) NOT NULL, 
		"city" VARCHAR(64) NOT NULL, 
		"country_code" VARCHAR(32) NOT NULL,
		PRIMARY KEY("id")
	);`

	if _, err := conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("postgres: failed to create a table, error: %v", err)
	}

	log.Logger.Info("Successfully initialized postgres database controller")

	return nil
}

func (pc *PostgresController) AddUser(ctx context.Context, firstName, lastName, email, password string, geolocationData *ipresolver.GeolocationData) error {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres: failed to acquire connection from the pool, error: %v", err)
	}

	defer conn.Release()

	query := `INSERT INTO "users" (
		"first_name", "last_name", "email", "password", 
		"country", "city", "country_code"
	) values ($1, $2, $3, $4, $5, $6, $7);`

	hashedPassword := hashutils.NewSha256([]byte(password))

	if _, err := conn.Exec(ctx, query, firstName, lastName, email, hashedPassword, geolocationData.Country, geolocationData.City, geolocationData.CountryCode); err != nil {
		return fmt.Errorf("postgres: failed to add user, error: %v", err)
	}

	return nil
}

func (pc *PostgresController) HasUser(ctx context.Context, email string) (bool, error) {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("postgres: failed to acquire database connection, error: %v", err)
	}

	defer conn.Release()

	query := `SELECT "email" FROM "users" WHERE "email" = ($1);`
	row := conn.QueryRow(ctx, query, email)

	var result string
	if err := row.Scan(&result); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("postgres: failed to select user, error: %v", err)
	}

	return true, nil
}

// TODO: Maybe we can have a function which will return all the users
// in a database, and we can render them from react

func (db *PostgresController) Close() error {
	db.connPool.Close()
	return nil
}
