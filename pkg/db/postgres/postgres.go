package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/log"
)

type PostgresController struct {
	// context timeout for database queries
	ctxTimeout time.Duration
	// connection pool
	connPool *pgxpool.Pool
}

func NewPostgresController(ctx context.Context) (*PostgresController, error) {
	postgresUrl, set := os.LookupEnv("POSTGRES_URL")
	if !set || postgresUrl == "" {
		return nil, fmt.Errorf("POSTGRES_URL is not set")
	}

	config, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to parse a connection config, error: %v", err)
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

	log.Logger.Info("Successfully initialized postgres database")

	return postgres, nil
}

func (pc *PostgresController) createTable(ctx context.Context) error {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres: failed to acquire database connection, error: %v", err)
	}

	defer conn.Release()

	// TOOD: Figure out how to use timestamps
	// "created_at" TIMESTAMP NOT NULL,
	// "updated_at" TIMESTAMP NOT NULL,

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

func (pc *PostgresController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres: failed to acquire connection from the pool, error: %v", err)
	}

	defer conn.Release()

	query := `INSERT INTO "users" (
		"first_name", "last_name", "email", "password", 
		"country", "city", "country_code"
	) values ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := conn.Exec(ctx, query, userData.FirstName, userData.LastName,
		userData.Email, userData.Password, geolocationData.Country, geolocationData.City, geolocationData.CountryCode); err != nil {
		return fmt.Errorf("postgres: failed to add user, error: %v", err)
	}

	return nil
}

func (pc *PostgresController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to acquire database connection, error: %v", err)
	}

	defer conn.Release()

	query := `SELECT 
	"first_name", "last_name", "email", "password", "country", "city"
	FROM "users" WHERE "email" = ($1);`

	rows, _ := conn.Query(ctx, query, email)
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserData])
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to select user, error: %v", err)
	}

	// User doesn't exist
	if len(users) == 0 {
		return nil, nil
	}

	log.Logger.Info("Matched user: %v", users[0])

	return &users[0], nil
}

// NOTE: This method will be used when a refresh token contains an issuer.
// TODO: This has to be tested properly.
func (pc *PostgresController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	conn, err := pc.connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to acquire database connection, error: %v", err)
	}

	defer conn.Release()

	query := `SELECT 
	"first_name", "last_name", "email", "password", 
	"country", "city" FROM "users" WHERE "id" = ($1);`

	rows, _ := conn.Query(ctx, query, id)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserData])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			// User doesn't exist
			return nil, nil
		default:
			return nil, fmt.Errorf("postgres: failed to collect rows, error: %v", user)
		}
	}

	return &user, nil
}

func (pc *PostgresController) Close() error {
	pc.connPool.Close()
	return nil
}
