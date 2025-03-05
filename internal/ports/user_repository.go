package ports

import (
	"context"

	"github.com/isnastish/aiclient/internal/domain/ipresolver"
	"github.com/isnastish/aiclient/internal/domain/users"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*users.User, error)
	GetUserByID(ctx context.Context, id int) (*users.User, error)
	InsertUser(ctx context.Context, userData *users.User, geolocation *ipresolver.UserGeolocation) error
	HasUser(ctx context.Context, email string) (bool, error)
}
