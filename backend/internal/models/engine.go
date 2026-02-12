package models

import "time"

// Engine represents a car engine with detailed technical specifications,
// common problems, and maintenance information
type Engine struct {
	ID           int64   `db:"id" json:"id"`
	Code         string  `db:"code" json:"code"`                           // Unique engine code (e.g., "EA888", "N47")
	Name         *string `db:"name" json:"name,omitempty"`                 // Engine name (e.g., "2.0 TSI Gen 3")
	Manufacturer *string `db:"manufacturer" json:"manufacturer,omitempty"` // Manufacturer (e.g., "VW Group", "BMW")

	// Technical Specifications
	EngineType        *string `db:"engine_type" json:"engine_type,omitempty"`                 // e.g., "I4", "V6", "V8"
	FuelType          *string `db:"fuel_type" json:"fuel_type,omitempty"`                     // "Petrol", "Diesel", "Hybrid", "Electric"
	DisplacementCC    *int    `db:"displacement_cc" json:"displacement_cc,omitempty"`         // Engine displacement in cc
	Cylinders         *int    `db:"cylinders" json:"cylinders,omitempty"`                     // Number of cylinders
	CylinderLayout    *string `db:"cylinder_layout" json:"cylinder_layout,omitempty"`         // "Inline", "V", "Boxer", "W"
	ValvesPerCylinder *int    `db:"valves_per_cylinder" json:"valves_per_cylinder,omitempty"` // Valves per cylinder
	Aspiration        *string `db:"aspiration" json:"aspiration,omitempty"`                   // "Turbo", "Supercharged", "Natural Aspirated"

	// Additional Information
	Description         *string `db:"description" json:"description,omitempty"`                 // General description
	TechnologyFeatures  *string `db:"technology_features" json:"technology_features,omitempty"` // JSON array of features
	ProductionStartYear *int    `db:"production_start_year" json:"production_start_year,omitempty"`
	ProductionEndYear   *int    `db:"production_end_year" json:"production_end_year,omitempty"` // NULL if still in production

	// Problems and Solutions
	CommonProblems    *string `db:"common_problems" json:"common_problems,omitempty"`       // JSON array of problem objects
	Solutions         *string `db:"solutions" json:"solutions,omitempty"`                   // JSON array of solution objects
	MaintenanceNotes  *string `db:"maintenance_notes" json:"maintenance_notes,omitempty"`   // General maintenance advice
	ReliabilityRating *int    `db:"reliability_rating" json:"reliability_rating,omitempty"` // 1-10 rating

	// Metadata
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// EngineProblem represents a common problem with an engine
// This is used for parsing JSON from CommonProblems field
type EngineProblem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "low", "medium", "high"
}

// EngineSolution represents a solution for a common problem
// This is used for parsing JSON from Solutions field
type EngineSolution struct {
	ProblemTitle  string  `json:"problemTitle"`
	Solution      string  `json:"solution"`
	EstimatedCost *string `json:"estimatedCost,omitempty"`
}

// TechnologyFeature represents a technology feature
// This is used for parsing JSON from TechnologyFeatures field
type TechnologyFeature struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
