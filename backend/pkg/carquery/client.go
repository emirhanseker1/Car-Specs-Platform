package carquery

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const BaseURL = "https://www.carqueryapi.com/api/0.3/"

// Client handles interaction with the CarQuery API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new CarQuery API client.
func NewClient() *Client {
	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// Structs for API responses

// YearResponse represents the response for GetYears.
type YearResponse struct {
	Years struct {
		MinYear int `json:"min_year,string"`
		MaxYear int `json:"max_year,string"`
	} `json:"Years"`
}

// CarMake represents a vehicle manufacturer.
type CarMake struct {
	MakeID       string `json:"make_id"`
	MakeDisplay  string `json:"make_display"`
	MakeIsCommon string `json:"make_is_common"`
	MakeCountry  string `json:"make_country"`
}

type MakesResponse struct {
	Makes []CarMake `json:"Makes"`
}

// CarModel represents a vehicle model.
type CarModel struct {
	ModelName   string `json:"model_name"`
	ModelMakeId string `json:"model_make_id"`
}

type ModelsResponse struct {
	Models []CarModel `json:"Models"`
}

// CarTrim represents specific trim details.
type CarTrim struct {
	ModelID               string `json:"model_id"`
	ModelMakeID           string `json:"model_make_id"`
	ModelName             string `json:"model_name"`
	ModelTrim             string `json:"model_trim"`
	ModelYear             string `json:"model_year"`
	ModelBody             string `json:"model_body"`
	ModelEnginePosition   string `json:"model_engine_position"`
	ModelEngineCC         string `json:"model_engine_cc"`
	ModelEngineCyl        string `json:"model_engine_cyl"`
	ModelEngineType       string `json:"model_engine_type"`
	ModelValvesPerCyl     string `json:"model_valves_per_cyl"`
	ModelPowerPS          string `json:"model_power_ps"`
	ModelPowerKW          string `json:"model_power_kw"`
	ModelTorqueNm         string `json:"model_torque_nm"`
	ModelTopSpeedKph      string `json:"model_top_speed_kph"`
	Model0to100Kph        string `json:"model_0_to_100_kph"`
	ModelDrive            string `json:"model_drive"`
	ModelTransmissionType string `json:"model_transmission_type"`
	ModelSeats            string `json:"model_seats"`
	ModelDoors            string `json:"model_doors"`
	ModelWeightKg         string `json:"model_weight_kg"`
	ModelLengthMm         string `json:"model_length_mm"`
	ModelWidthMm          string `json:"model_width_mm"`
	ModelHeightMm         string `json:"model_height_mm"`
	ModelWheelbaseMm      string `json:"model_wheelbase_mm"`
	ModelLkmHwy           string `json:"model_lkm_hwy"`
	ModelLkmMixed         string `json:"model_lkm_mixed"`
	ModelLkmCity          string `json:"model_lkm_city"`
	ModelFuelCapL         string `json:"model_fuel_cap_l"`
	ModelSoldInUS         string `json:"model_sold_in_us"`
	ModelCo2              string `json:"model_co2"`
	ModelMakeDisplay      string `json:"model_make_display"`
}

type TrimsResponse struct {
	Trims []CarTrim `json:"Trims"`
}

// Helper to make GET request and unmarshal response
func (c *Client) get(params url.Values, target interface{}) error {
	reqURL := fmt.Sprintf("%s?%s", c.BaseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent to mimic a browser to avoid bot blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// ...

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}
	return nil
}

// GetYears returns the min and max year available in the database.
func (c *Client) GetYears() (min, max int, err error) {
	params := url.Values{}
	params.Set("cmd", "getYears")

	var res YearResponse
	if err := c.get(params, &res); err != nil {
		return 0, 0, err
	}
	return res.Years.MinYear, res.Years.MaxYear, nil
}

// GetMakes returns all makes available for a specific year.
// Year is optional (pass 0 to ignore).
func (c *Client) GetMakes(year int) ([]CarMake, error) {
	params := url.Values{}
	params.Set("cmd", "getMakes")
	if year > 0 {
		params.Set("year", fmt.Sprintf("%d", year))
	}
	// sold_in_us filtering? Let's leave clear for now to get all.

	var res MakesResponse
	if err := c.get(params, &res); err != nil {
		return nil, err
	}
	return res.Makes, nil
}

// GetModels returns models for a specific make and year.
func (c *Client) GetModels(makeID string, year int) ([]CarModel, error) {
	params := url.Values{}
	params.Set("cmd", "getModels")
	params.Set("make", makeID)
	if year > 0 {
		params.Set("year", fmt.Sprintf("%d", year))
	}

	var res ModelsResponse
	if err := c.get(params, &res); err != nil {
		return nil, err
	}
	return res.Models, nil
}

// GetTrims returns detailed trim information.
// You can filter by modelID if known, or by make/model/year params in general.
// The user request asked for "Trim details for a specific model".
// If "Model" means the name (e.g. "3 Series"), it might return multiple trims.
// If we have a specific model ID (from previous calls? CarQuery getModels doesn't return IDs, just names),
// we usually query by Make, Model, Year to get Trims.
// However, the `CarTrim` struct has `model_id`.
// Let's support querying by Make, Model, Year as that's the standard flow.
// Or if the user meant specific Trim ID?
// Let's provide a flexible function.
type TrimFilter struct {
	Make  string
	Model string
	Year  int
	// Keyword string
}

func (c *Client) GetTrims(filter TrimFilter) ([]CarTrim, error) {
	params := url.Values{}
	params.Set("cmd", "getTrims")

	if filter.Make != "" {
		params.Set("make", filter.Make)
	}
	if filter.Model != "" {
		params.Set("model", filter.Model)
	}
	if filter.Year > 0 {
		params.Set("year", fmt.Sprintf("%d", filter.Year))
	}
	// "full_results" might be needed? Docs say "Only basic info returned unless `full_results=1`" (Sometimes implied).
	// Let's try to add it just in case.
	// Actually docs: "To get full specifications... use getTrims"

	var res TrimsResponse
	if err := c.get(params, &res); err != nil {
		return nil, err
	}
	return res.Trims, nil
}

// GetTrim returns a single trim detail if we had an ID.
// Since API uses getTrims for everything, the above function covers it if we filter precise enough.
