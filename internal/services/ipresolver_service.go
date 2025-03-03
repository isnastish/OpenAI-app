package services

import (
	ports "github.com/isnastish/openai/internal/ports/outbound"
)

type IpResolverService struct {
	ipResolverRepo ports.IpResolverRepository
}

func NewIpResolverService(ipResolverRepo ports.IpResolverRepository) *IpResolverService {
	return &IpResolverService{
		ipResolverRepo: ipResolverRepo,
	}
}
