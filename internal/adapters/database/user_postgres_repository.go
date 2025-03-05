package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/isnastish/aiclient/internal/domain/ipresolver"
	"github.com/isnastish/aiclient/internal/domain/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	conn *pgxpool.Pool
}

func NewPostgresUserRepository(connPool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		conn: connPool,
	}
}

func (p PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	conn, err := p.conn.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection, error %v", err)
	}
	defer conn.Release()

	query := `SELECT 
	"first_name", "last_name", "email", "password", "country", "city"
	FROM "users" WHERE "email" = ($1);`

	rows, _ := conn.Query(ctx, query, email)
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[users.User])
	if err != nil {
		return nil, fmt.Errorf("failed to select user, error %v", err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}

// NOTE: id should probably be UUID instead of int?

func (p PostgresUserRepository) GetUserByID(ctx context.Context, id int) (*users.User, error) {
	conn, err := p.conn.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection, error %v", err)
	}
	defer conn.Release()

	query := `SELECT 
	"first_name", "last_name", "email", "password", 
	"country", "city" FROM "users" WHERE "id" = ($1);`

	rows, _ := conn.Query(ctx, query, id)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[users.User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil
		default:
			return nil, fmt.Errorf("failed to collect rows, error %v", user)
		}
	}
	return &user, nil
}

func (p PostgresUserRepository) InsertUser(ctx context.Context, userData *users.User, geolocation *ipresolver.UserGeolocation) error {
	conn, err := p.conn.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres: failed to acquire connection from the pool, error: %v", err)
	}
	defer conn.Release()

	query := `INSERT INTO "users" (
		"first_name", "last_name", "email", "password", 
		"country", "city", "country_code"
	) values ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := conn.Exec(ctx, query, userData.FirstName, userData.LastName,
		userData.Email, userData.Password, geolocation.Country, geolocation.City, geolocation.CountryCode); err != nil {
		return fmt.Errorf("postgres: failed to add user, error: %v", err)
	}
	return nil
}

func (p PostgresUserRepository) HasUser(ctx context.Context, email string) (bool, error) {
	user, err := p.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return (user == nil), nil
}
