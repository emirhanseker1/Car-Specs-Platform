package importer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/pkg/apininjas"
	"gorm.io/gorm"
)

// SyncApiNinjasData fetches data for provided makes and saves them ensuring relational integrity.
func SyncApiNinjasData(db *gorm.DB, client *apininjas.Client, makes []string) error {
	for _, makeName := range makes {
		log.Printf("Fetching cars for make: %s...", makeName)
		// Use 0 for year to indicate "any year" (legacy behavior)
		cars, err := client.FetchCars(makeName, 0)
		if err != nil {
			log.Printf("Error fetching %s: %v", makeName, err)
			continue
		}

		if len(cars) == 0 {
			log.Printf("No cars found for %s", makeName)
			continue
		}

		// Transaction for each Make to ensure consistency
		err = db.Transaction(func(tx *gorm.DB) error {
			// 1. Find or Create Make
			var makeModel models.Make
			if err := tx.FirstOrCreate(&makeModel, models.Make{Name: makeName}).Error; err != nil {
				return err
			}

			for _, car := range cars {
				// 2. Find or Create Model
				var modelModel models.Model
				if err := tx.FirstOrCreate(&modelModel, models.Model{Name: car.Model, MakeID: makeModel.ID}).Error; err != nil {
					return err
				}

				// 3. Create Trim
				// Construct a descriptive Trim name
				trimName := fmt.Sprintf("%d %s %s", car.Year, car.Transmission, car.Drive)

				// Helper to get float from interface{}
				dispFloat := toFloat(car.Displacement)

				if dispFloat > 0 {
					trimName = fmt.Sprintf("%s %.1fL", trimName, dispFloat)
				}

				var trimModel models.Trim
				// Avoid duplicates based on Name + ModelID?
				// Since we don't have unique constraints perfectly set up, we check manually or use FirstOrCreate
				err := tx.Where(models.Trim{ModelID: modelModel.ID, Name: trimName, Year: car.Year}).
					FirstOrCreate(&trimModel, models.Trim{
						ModelID: modelModel.ID,
						Name:    trimName,
						Year:    car.Year,
					}).Error
				if err != nil {
					return err
				}

				// 4. Add Specs (Technical Details)
				specs := []models.Spec{
					{TrimID: int64(trimModel.ID), Category: "General", Name: "Class", Value: car.Class},
					{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Cylinders", Value: toString(car.Cylinders)},
					{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Displacement", Value: toString(car.Displacement)},
					{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Fuel Type", Value: car.FuelType},
					{TrimID: int64(trimModel.ID), Category: "Transmission", Name: "Transmission", Value: car.Transmission},
					{TrimID: int64(trimModel.ID), Category: "Drivetrain", Name: "Drive", Value: car.Drive},
					{TrimID: int64(trimModel.ID), Category: "Consumption", Name: "City MPG", Value: toString(car.CityMPG)},
					{TrimID: int64(trimModel.ID), Category: "Consumption", Name: "Highway MPG", Value: toString(car.HighwayMPG)},
				}

				for _, spec := range specs {
					// Upsert Spec? Or just Create (can lead to dupes if run multiple times)
					// Let's use FirstOrCreate
					if err := tx.Where(&models.Spec{TrimID: spec.TrimID, Category: spec.Category, Name: spec.Name}).
						Assign(models.Spec{Value: spec.Value}). // Update value if exists
						FirstOrCreate(&spec).Error; err != nil {
						return err
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Failed to save data for %s: %v", makeName, err)
		} else {
			log.Printf("Successfully imported %d cars for %s", len(cars), makeName)
		}
	}
	return nil
}

// SyncDetailedYears fetches data for specific years (e.g. 2024, 2025) to showcase detailed new models.
func SyncDetailedYears(db *gorm.DB, client *apininjas.Client, makes []string, years []int) error {
	for _, makeName := range makes {
		for _, year := range years {
			log.Printf(">>> Demo Mode: Fetching %s models for %d...", makeName, year)
			cars, err := client.FetchCars(makeName, year)
			if err != nil {
				log.Printf("Error fetching %s %d: %v", makeName, year, err)
				continue
			}
			if len(cars) == 0 {
				log.Printf("   [INFO] No models returned for %s in %d (likely not yet in API).", makeName, year)
				continue
			}
			log.Printf("   Found %d models! Importing...", len(cars))

			// Reuse the same transaction logic?
			// We can extract the saving logic to a helper "importCars(tx, cars)" to avoid code duplication.
			// For now, to keep it simple and robust, I'll inline the save logic again or call a helper.
			// Let's refactor the save logic into `saveCars` private helper.
			if err := saveCars(db, cars); err != nil {
				log.Printf("Failed to save batch: %v", err)
			}
		}
	}
	return nil
}

// saveCars handles the DB transaction for a list of cars
func saveCars(db *gorm.DB, cars []apininjas.Car) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, car := range cars {
			// 1. Find or Create Make
			var makeModel models.Make
			if err := tx.FirstOrCreate(&makeModel, models.Make{Name: car.Make}).Error; err != nil {
				return err
			}

			// 2. Find or Create Model
			var modelModel models.Model
			if err := tx.FirstOrCreate(&modelModel, models.Model{Name: car.Model, MakeID: makeModel.ID}).Error; err != nil {
				return err
			}

			// 3. Create Trim
			trimName := fmt.Sprintf("%d %s %s", car.Year, car.Transmission, car.Drive)
			dispFloat := toFloat(car.Displacement)
			if dispFloat > 0 {
				trimName = fmt.Sprintf("%s %.1fL", trimName, dispFloat)
			}

			var trimModel models.Trim
			err := tx.Where(models.Trim{ModelID: modelModel.ID, Name: trimName, Year: car.Year}).
				FirstOrCreate(&trimModel, models.Trim{
					ModelID: modelModel.ID,
					Name:    trimName,
					Year:    car.Year,
				}).Error
			if err != nil {
				return err
			}

			// 4. Specs
			specs := []models.Spec{
				{TrimID: int64(trimModel.ID), Category: "General", Name: "Class", Value: normalizeSpecValue("General", "Class", car.Class)},
				{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Cylinders", Value: toString(car.Cylinders)},
				{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Displacement", Value: toString(car.Displacement)},
				{TrimID: int64(trimModel.ID), Category: "Engine", Name: "Fuel Type", Value: normalizeSpecValue("Engine", "Fuel Type", car.FuelType)},
				{TrimID: int64(trimModel.ID), Category: "Transmission", Name: "Transmission", Value: normalizeSpecValue("Transmission", "Transmission", car.Transmission)},
				{TrimID: int64(trimModel.ID), Category: "Drivetrain", Name: "Drive", Value: normalizeSpecValue("Drivetrain", "Drive", car.Drive)},
				{TrimID: int64(trimModel.ID), Category: "Consumption", Name: "City MPG", Value: toString(car.CityMPG)},
				{TrimID: int64(trimModel.ID), Category: "Consumption", Name: "Highway MPG", Value: toString(car.HighwayMPG)},
			}

			for _, spec := range specs {
				if err := tx.Where(&models.Spec{TrimID: spec.TrimID, Category: spec.Category, Name: spec.Name}).
					Assign(models.Spec{Value: spec.Value}).
					FirstOrCreate(&spec).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Helper to safely get float64 from interface{} (which might be float64, int, string)
func toFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	// Fallback/Log?
	return 0
}

// Helper to safely get string from interface{}
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	s := fmt.Sprintf("%v", v)
	// Filter out premium messages
	if s == "The limit parameter is for premium users only." || s == "Premium Subscribers Only" {
		return "-"
	}
	return s
}

func normalizeSpecValue(category, name, value string) string {
	v := strings.ToLower(strings.TrimSpace(value))

	if strings.Contains(v, "premium") {
		return "-"
	}

	if name == "Transmission" {
		if v == "a" || v == "automatic" {
			return "Otomatik"
		}
		if v == "m" || v == "manual" {
			return "Manuel"
		}
		if v == "cvt" {
			return "Otomatik (CVT)"
		}
	}

	if name == "Drive" {
		if v == "fwd" {
			return "Önden Çekiş"
		}
		if v == "rwd" {
			return "Arkadan İtiş"
		}
		if v == "awd" || v == "4wd" {
			return "4 Çeker"
		}
	}

	if name == "Fuel Type" {
		switch v {
		case "gas", "gasoline":
			return "Benzin"
		case "diesel":
			return "Dizel"
		case "electricity", "electric":
			return "Elektrik"
		case "hybrid":
			return "Hibrit"
		}
	}

	if name == "Class" {
		// Translate common classes
		if strings.Contains(v, "sedan") {
			return "Sedan"
		}
		if strings.Contains(v, "coupe") {
			return "Coupe"
		}
		if strings.Contains(v, "wagon") {
			return "Station Wagon"
		}
		if strings.Contains(v, "suv") {
			return "SUV"
		}
		if strings.Contains(v, "pickup") {
			return "Pickup"
		}
		if strings.Contains(v, "convertible") {
			return "Cabrio"
		}
		if strings.Contains(v, "minivan") {
			return "Minivan"
		}
	}

	// Capitalize first letter for others
	if len(v) > 1 {
		return strings.ToUpper(v[:1]) + v[1:]
	}
	return value
}
