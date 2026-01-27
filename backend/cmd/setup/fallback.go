package main

import (
	"fmt"
	"log"
	"strings"
)

// CreateFallbackTrim creates a trim entry with estimated/default values when API Ninjas has no data
//
//	This ensures we don't skip cars entirely just because specs aren't available
func (s *SetupService) CreateFallbackTrim(modelID int64, brand, model string, year int, imageURL, generation string) error {
	// Build a reasonable trim name
	trimName := fmt.Sprintf("%s %d", strings.Title(model), year)

	query := `
		INSERT INTO trims (
			model_id, name, year, generation, market,
			fuel_type, transmission_type,
			seating_capacity, doors,
			image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		modelID,
		trimName,
		year,
		generation,
		"TR",       // Turkish market
		"Gasoline", // Default fuel type
		"Manual",   // Most common in Turkey
		5,          // Standard seating
		4,          // Most cars are 4-door
		imageURL,
	)

	if err != nil {
		return fmt.Errorf("failed to create fall back trim: %w", err)
	}

	log.Printf("  ⚠️  Created fallback entry (no API data available)")
	return nil
}
