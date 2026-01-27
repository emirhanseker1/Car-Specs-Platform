package importer

import (
	"fmt"

	"github.com/emirh/car-specs/backend/internal/service"
	"gorm.io/gorm"
)

func SeedDemoData(db *gorm.DB) error {
	fmt.Println("Seeding demo data...")

	demoData := []service.CarJSON{
		{Make: "BMW", Model: "3 Series", Trim: "320i Sport Line", Year: 2024},
		{Make: "BMW", Model: "3 Series", Trim: "330i M Sport", Year: 2024},
		{Make: "BMW", Model: "5 Series", Trim: "520i Luxury", Year: 2024},
		{Make: "Audi", Model: "A4", Trim: "40 TFSI Advanced", Year: 2024},
		{Make: "Audi", Model: "A6", Trim: "45 TFSI Quattro", Year: 2024},
		{Make: "Mercedes-Benz", Model: "C-Class", Trim: "C 200 AMG", Year: 2024},
		{Make: "Tesla", Model: "Model Y", Trim: "Long Range AWD", Year: 2024},
		{Make: "Tesla", Model: "Model 3", Trim: "Performance", Year: 2024},
	}

	return service.ImportCarData(db, demoData)
}
