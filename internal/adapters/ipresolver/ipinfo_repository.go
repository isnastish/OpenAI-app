package ipresolver

import "github.com/ipinfo/go/v2/ipinfo"

type IpinfoRepository struct {
	client *ipinfo.Client
}

func NewIpinfoRepository() *IpinfoRepository {
	return &IpinfoRepository{}
}
