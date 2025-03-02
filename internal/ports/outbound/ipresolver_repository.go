package ports

import (
	"context"

	"github.com/isnastish/openai/internal/domain/ipresolver"
)

type IpResolverRepository interface {
	GetUserGeolocationData(ctx context.Context, ipAddress string) (*ipresolver.UserGeolocation, error)
}
