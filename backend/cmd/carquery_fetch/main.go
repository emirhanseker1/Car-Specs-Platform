package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// CarQuery API response structs
type TrimsResponse struct {
	Trims []CarTrim `json:"Trims"`
}

type CarTrim struct {
	ModelID          string `json:"model_id"`
	ModelMakeID      string `json:"model_make_id"`
	ModelName        string `json:"model_name"`
	ModelTrim        string `json:"model_trim"`
	ModelYear        string `json:"model_year"`
	ModelEngineCC    string `json:"model_engine_cc"`
	ModelEngineCyl   string `json:"model_engine_cyl"`
	ModelPowerPS     string `json:"model_power_ps"`
	ModelMakeDisplay string `json:"model_make_display"`
}

func main() {
	client := &http.Client{Timeout: 30 * time.Second}

	// Fetch Audi A3 trims for years 2010-2024
	allTrims := []CarTrim{}

	for year := 2010; year <= 2024; year++ {
		url := fmt.Sprintf("https://www.carqueryapi.com/api/0.3/?cmd=getTrims&make=audi&model=A3&year=%d", year)
		fmt.Printf("Fetching year %d...\n", year)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		var result TrimsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("  Decode error: %v\n", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		fmt.Printf("  Found %d trims\n", len(result.Trims))
		allTrims = append(allTrims, result.Trims...)

		time.Sleep(500 * time.Millisecond) // Rate limiting
	}

	// Save to JSON
	output, _ := json.MarshalIndent(allTrims, "", "  ")
	os.WriteFile("carquery_audi_a3.json", output, 0644)

	fmt.Printf("\nTotal trims fetched: %d\n", len(allTrims))
	fmt.Println("Saved to carquery_audi_a3.json")

	// Print sample
	fmt.Println("\nSample trims:")
	for i, t := range allTrims {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s %s %s (%s) - %s cc, %s PS\n",
			t.ModelMakeDisplay, t.ModelName, t.ModelTrim, t.ModelYear,
			t.ModelEngineCC, t.ModelPowerPS)
	}
}
