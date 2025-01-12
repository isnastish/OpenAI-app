package ipresolver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/log"
)

// A user facing client for interacting with an external
// service for resolving geolocation.
type Client struct {
	// http client
	httpClient *http.Client
	// ipflare API key retrieved from env variable
	ipflareApiKey string
}

func NewClient() (*Client, error) {
	ipflareApiKey, set := os.LookupEnv("IPFLARE_API_KEY")
	if !set || ipflareApiKey == "" {
		return nil, fmt.Errorf("IPFLARE_API_KEY is not set")
	}

	return &Client{
		httpClient:    &http.Client{},
		ipflareApiKey: ipflareApiKey,
	}, nil
}

// Make a request to an external service `ipflare`: https://www.ipflare.io/
// Parse its response and retreive geolocation data based on provided
// ip address.
// Return an error, if any, otherwise an instance of `IpInfo` containing the necessary information.
func (c *Client) GetGeolocationData(ipAddr string) (*models.Geolocation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.ipflare.io/%s", ipAddr), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %s", err.Error())
	}

	req.Header.Add("X-API-Key", c.ipflareApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %s", err.Error())
	}

	var geolocation models.Geolocation
	if err := json.Unmarshal(body, &geolocation); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response data: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", geolocation.ErrorCode, geolocation.ErrorMsg)
	}

	log.Logger.Info("Got geolocation for IP: %s, city: %s, country: %s", ipAddr, geolocation.City, geolocation.Country)

	return &geolocation, nil
}
