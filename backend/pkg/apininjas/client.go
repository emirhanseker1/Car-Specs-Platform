package apininjas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const BaseURL = "https://api.api-ninjas.com/v1/cars"

// Car represents the flat JSON object returned by API Ninjas.
type Car struct {
	Make         string      `json:"make"`
	Model        string      `json:"model"`
	Year         int         `json:"year"`
	FuelType     string      `json:"fuel_type"`
	Drive        string      `json:"drive"`
	Cylinders    interface{} `json:"cylinders"`
	Transmission string      `json:"transmission"`
	CityMPG      interface{} `json:"city_mpg"`
	HighwayMPG   interface{} `json:"highway_mpg"`
	Class        string      `json:"class"`
	Displacement interface{} `json:"displacement"` // e.g. 2.2 or "2.2"
}

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	apiKey := os.Getenv("API_NINJAS_KEY")
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchCars fetches cars with optional filters.
// If year is 0, it is ignored.
func (c *Client) FetchCars(makeName string, year int) ([]Car, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("API_NINJAS_KEY environment variable is not set")
	}

	reqURL := fmt.Sprintf("%s?make=%s", BaseURL, makeName)
	if year > 0 {
		reqURL = fmt.Sprintf("%s&year=%d", reqURL, year)
	}
	// Limit is still restricted on free tier, avoiding it.

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read body to see error message
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var cars []Car
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return cars, nil
}
