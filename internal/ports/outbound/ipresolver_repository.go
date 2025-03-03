package ports

import (
	"context"

	"github.com/isnastish/aiclient/internal/domain/ipresolver"
)

type IpResolverRepository interface {
	GetUserGeolocationData(ctx context.Context, ipAddress string) (*ipresolver.UserGeolocation, error)
}
