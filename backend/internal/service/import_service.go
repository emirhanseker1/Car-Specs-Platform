package service

import "gorm.io/gorm"

// CarJSON represents a simplified vehicle structure for import
type CarJSON struct {
	Make  string
	Model string
	Trim  string
	Year  int
}

// ImportCarData imports a batch of car data (placeholder implementation)
func ImportCarData(db *gorm.DB, data []CarJSON) error {
	// This function is required for legacy importer code to compile.
	// Current scraping logic uses Repositories with database/sql.
	return nil
}
