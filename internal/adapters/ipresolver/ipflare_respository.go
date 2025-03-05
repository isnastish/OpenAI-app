package ipresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/isnastish/aiclient/internal/domain/ipresolver"
)

type IpflareRepository struct {
	httpClient *http.Client
	apiKey     string
}

func NewIpflareRespository(ipflareApiKey string) IpflareRepository {
	return IpflareRepository{
		httpClient: &http.Client{},
		apiKey:     ipflareApiKey,
	}
}

type errorResponse struct {
	ErrorCode string `json:"code,omitempty"`
	ErrorMsg  string `json:"error,omitempty"`
}

type ipflareResponse struct {
	Ip          string `json:"ip,omitempty"`
	City        string `json:"city,omitempty"`
	Region      string `json:"region,omitempty"`
	RegionCode  string `json:"region_code,omitempty"`
	Country     string `json:"country_name,omitempty"`
	CountryCode string `json:"country_code,omitempty"`

	errorResponse
}

const ipflareBaseURL = `https://api.ipflare.io/`

func (i IpflareRepository) GetUserGeolocationData(ctx context.Context, ipAddress string) (*ipresolver.UserGeolocation, error) {
	req, err := http.NewRequest("GET", ipflareBaseURL+ipAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %s", err.Error())
	}
	req.Header.Add("X-API-Key", i.apiKey)

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, error %v", err)
	}

	var res ipflareResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data, error %v", err)
	}

	// TODO: Work on this...
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", res.ErrorCode, res.ErrorMsg)
	}

	return &ipresolver.UserGeolocation{
		Ip:          res.Ip,
		City:        res.City,
		Country:     res.Country,
		CountryCode: res.CountryCode,
	}, nil
}
