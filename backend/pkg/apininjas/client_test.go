package apininjas

import (
	"encoding/json"
	"testing"
)

func TestCarUnmarshal(t *testing.T) {
	// Sample JSON with mixed types (integers as strings, floats, etc.)
	// derived from real API behavior we've seen.
	jsonData := `[
		{
			"make": "TestMake",
			"model": "Model StringMPG",
			"year": 2022,
			"city_mpg": "25",
			"highway_mpg": "32",
			"cylinders": 4,
			"displacement": 2.0
		},
		{
			"make": "TestMake",
			"model": "Model IntMPG",
			"year": 2023,
			"city_mpg": 20,
			"highway_mpg": 28,
			"cylinders": "6",
			"displacement": "3.5"
		},
		{
			"make": "TestMake",
			"model": "Model Mixed",
			"year": 2024,
			"city_mpg": 22,
			"highway_mpg": "30",
			"cylinders": 8,
			"displacement": 5.0
		}
	]`

	var cars []Car
	if err := json.Unmarshal([]byte(jsonData), &cars); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(cars) != 3 {
		t.Errorf("Expected 3 cars, got %d", len(cars))
	}

	// Verify Check 1: String MPG
	if cars[0].Model != "Model StringMPG" {
		t.Errorf("Unexpected model name: %s", cars[0].Model)
	}
	// We are just checking if it Unmarshals without error.
	// The processing logic happens in the importer, but the Client struct
	// must be able to hold the data.

	// Check generic fields content
	t.Logf("Car 0 CityMPG type: %T value: %v", cars[0].CityMPG, cars[0].CityMPG)
	t.Logf("Car 1 Displacement type: %T value: %v", cars[1].Displacement, cars[1].Displacement)

	// Since we defined fields as interface{}, json.Unmarshal should succeed.
	// The real logic validation is ensuring our Importer can HANDLE these interfaces.
}
