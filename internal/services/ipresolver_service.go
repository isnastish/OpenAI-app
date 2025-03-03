package services

import (
	"context"

	"github.com/isnastish/aiclient/internal/domain/ipresolver"
	ports "github.com/isnastish/aiclient/internal/ports/outbound"
)

type IpResolverService struct {
	ipResolverRepo ports.IpResolverRepository
}

func NewIpResolverService(ipResolverRepo ports.IpResolverRepository) *IpResolverService {
	return &IpResolverService{
		ipResolverRepo: ipResolverRepo,
	}
}

func (i IpResolverService) GetUserGeolocation(ipAddress string) (*ipresolver.UserGeolocation, error) {
	location, err := i.ipResolverRepo.GetUserGeolocationData(context.Background(), ipAddress)
	if err != nil {
		return nil, err
	}
	return location, nil
}
