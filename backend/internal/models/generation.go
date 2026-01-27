package models

import "time"

// Generation represents a specific model generation/body type (Level 3)
// Sits between Model and Trim in the hierarchy
type Generation struct {
	ID      int64 `db:"id" json:"id"`
	ModelID int64 `db:"model_id" json:"model_id"`

	// Generation Identification
	Code      string  `db:"code" json:"code"`           // e.g., "F30", "G20", "8V", "Mk7"
	Name      *string `db:"name" json:"name,omitempty"` // e.g., "F30 (2012-2018)"
	StartYear int     `db:"start_year" json:"start_year"`
	EndYear   *int    `db:"end_year" json:"end_year,omitempty"` // NULL if current

	// Additional Info
	ImageURL    *string `db:"image_url" json:"image_url,omitempty"`
	Description *string `db:"description" json:"description,omitempty"`
	IsCurrent   bool    `db:"is_current" json:"is_current"`
	Platform    *string `db:"platform" json:"platform,omitempty"` // e.g., "MQB", "CLAR"

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relationships (populated via joins)
	Model *Model `db:"-" json:"model,omitempty"`
	Trims []Trim `db:"-" json:"trims,omitempty"`
}
