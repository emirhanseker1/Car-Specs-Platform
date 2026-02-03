package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// GenerationOverride represents manual generation data
type GenerationOverride struct {
	StartYear int    `json:"start_year"`
	EndYear   *int   `json:"end_year"` // nil means current/ongoing
	Notes     string `json:"notes"`
}

// ModelOverrides contains generation overrides for a model
type ModelOverrides struct {
	Generations map[string]GenerationOverride `json:"generations"`
}

// BrandOverrides contains model overrides for a brand
type BrandOverrides struct {
	Models map[string]ModelOverrides
}

// ManualOverrides is the root structure
type ManualOverrides struct {
	Brands map[string]map[string]ModelOverrides
}

var overrides *ManualOverrides

// LoadOverrides loads manual overrides from JSON file
func LoadOverrides(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		// If file doesn't exist, that's okay - no overrides
		if os.IsNotExist(err) {
			overrides = &ManualOverrides{
				Brands: make(map[string]map[string]ModelOverrides),
			}
			return nil
		}
		return fmt.Errorf("failed to read overrides file: %w", err)
	}

	overrides = &ManualOverrides{}
	err = json.Unmarshal(data, &overrides.Brands)
	if err != nil {
		return fmt.Errorf("failed to parse overrides JSON: %w", err)
	}

	fmt.Printf("✅ Loaded manual overrides from %s\n", filename)
	return nil
}

// GetGenerationYears returns the year range for a generation from overrides
// Returns (startYear, endYear, found)
func GetGenerationYears(brand, model, generationCode string) (int, *int, bool) {
	if overrides == nil {
		return 0, nil, false
	}

	brandData, ok := overrides.Brands[brand]
	if !ok {
		return 0, nil, false
	}

	modelData, ok := brandData[model]
	if !ok {
		return 0, nil, false
	}

	genData, ok := modelData.Generations[generationCode]
	if !ok {
		return 0, nil, false
	}

	return genData.StartYear, genData.EndYear, true
}
