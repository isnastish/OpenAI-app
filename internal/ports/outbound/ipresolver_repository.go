package ports

import "context"

type IpResolverRepository interface {
	GetGeolocationData(ctx context.Context, ipAddress string)
}
