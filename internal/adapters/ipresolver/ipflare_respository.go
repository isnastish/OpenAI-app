package ipresolver

import (
	"net/http"
	"os"
)

// NOTE: This could be either in a subdirectory /adapters/ipresolver/
// or the file itself could be called ipresolver_ipflare_repository.go

// NOTE: We should process all environment variables in one place,
// and pass a config, probably.

type IpflareRepository struct {
	httpClient *http.Client
	apiKey     string
}

func NewIpflareRespository() *IpflareRepository {
	apiKey, set := os.LookupEnv("IPFLARE_API_KEY")
	if !set || apiKey == "" {
		panic("IPFLARE_API_KEY is not set")
	}

	return &IpflareRepository{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}
