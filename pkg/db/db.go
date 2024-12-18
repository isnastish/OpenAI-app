package db

import (
	"context"

	"github.com/isnastish/openai/pkg/api/models"
)

// NOTE: We want to have an interface in DB with Postgres being one of them
// The interface should support main methods for writing and reading the data
// So we could swith between SQL and NoSQL data storages.
// But for now lets stick with a single implementation.
// Later, when we have a clear structure, we could replace Postgres
// with MySQL for example, or keep both.

type DatabaseController interface {
	AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error
	GetUserByEmail(ctx context.Context, email string) (*models.UserData, error)
	// The subject from the claims should have user ID
	GetUserByID(ctx context.Context, id int) (*models.UserData, error)
	Close() error
}
