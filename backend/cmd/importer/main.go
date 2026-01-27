package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/database"
	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
	"github.com/emirh/car-specs/backend/internal/service"
)

func main() {
	log.Println("=== CSV Vehicle Importer ===")

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize repositories
	brandRepo := repository.NewBrandRepository(database.DB)
	modelRepo := repository.NewModelRepository(database.DB)
	trimRepo := repository.NewTrimRepository(database.DB)

	// Initialize services
	brandService := service.NewBrandService(brandRepo)
	modelService := service.NewModelService(modelRepo, brandRepo)
	trimService := service.NewTrimService(trimRepo, modelRepo)

	// Open CSV file
	csvFile := "vehicles.csv"
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Skip header row
	if len(records) < 2 {
		log.Fatal("CSV file is empty or has no data rows")
	}

	log.Printf("Found %d vehicles to import\n", len(records)-1)

	// Statistics
	stats := struct {
		brandsCreated int
		modelsCreated int
		trimsCreated  int
		errors        int
	}{}

	// Process each row
	for i, record := range records[1:] {
		if len(record) < 15 {
			log.Printf("Row %d: Invalid format (expected 15 columns, got %d), skipping\n", i+2, len(record))
			stats.errors++
			continue
		}

		// Parse CSV columns (now with Image_URL as column 14)
		brandName := record[0]
		modelName := record[1]
		trimName := record[2]
		year := parseInt(record[3])
		hp := parseInt(record[4])
		torque := parseInt(record[5])
		accel := parseFloat(record[6])
		fuelType := record[7]
		transmission := record[8]
		drivetrain := record[9]
		displacement := parseInt(record[10])
		topSpeed := parseInt(record[11])
		fuelConsumption := parseFloat(record[12])
		bodyStyle := record[13]
		imageURL := record[14]

		log.Printf("\n[%d/%d] Processing: %s %s %s", i+1, len(records)-1, brandName, modelName, trimName)

		// Step 1: Get or Create Brand
		brand, err := brandService.GetBrandByName(brandName)
		if err != nil {
			// Brand doesn't exist, create it
			brand, err = brandService.CreateBrand(brandName, nil, nil)
			if err != nil {
				log.Printf("  ❌ Failed to create brand: %v", err)
				stats.errors++
				continue
			}
			log.Printf("  ✓ Created brand: %s (ID: %d)", brandName, brand.ID)
			stats.brandsCreated++
		} else {
			log.Printf("  ✓ Found existing brand: %s (ID: %d)", brandName, brand.ID)
		}

		// Step 2: Get or Create Model
		modelsList, err := modelService.ListModelsByBrand(brand.ID)
		if err != nil {
			log.Printf("  ❌ Failed to list models: %v", err)
			stats.errors++
			continue
		}

		var model *models.Model
		for _, m := range modelsList {
			if m.Name == modelName {
				model = m
				break
			}
		}

		if model == nil {
			// Model doesn't exist, create it
			model, err = modelService.CreateModel(brand.ID, modelName, &bodyStyle, nil)
			if err != nil {
				log.Printf("  ❌ Failed to create model: %v", err)
				stats.errors++
				continue
			}
			log.Printf("  ✓ Created model: %s (ID: %d)", modelName, model.ID)
			stats.modelsCreated++
		} else {
			log.Printf("  ✓ Found existing model: %s (ID: %d)", modelName, model.ID)
		}

		// Step 3: Create Trim
		trim := &models.Trim{
			ModelID:             model.ID,
			Name:                trimName,
			Year:                year,
			FuelType:            &fuelType,
			PowerHP:             &hp,
			TorqueNM:            &torque,
			Acceleration0To100:  &accel,
			TransmissionType:    &transmission,
			Drivetrain:          &drivetrain,
			DisplacementCC:      &displacement,
			TopSpeedKmh:         &topSpeed,
			FuelConsumptionComb: &fuelConsumption,
			ImageURL:            &imageURL,
			Market:              "TR",
			Currency:            "TRY",
			SeatingCapacity:     5,
		}

		if err := trimService.CreateTrim(trim); err != nil {
			log.Printf("  ❌ Failed to create trim: %v", err)
			stats.errors++
			continue
		}
		log.Printf("  ✓ Created trim: %s (ID: %d)", trimName, trim.ID)
		stats.trimsCreated++
	}

	// Print summary
	log.Println("\n=== Import Summary ===")
	log.Printf("Brands created: %d", stats.brandsCreated)
	log.Printf("Models created: %d", stats.modelsCreated)
	log.Printf("Trims created:  %d", stats.trimsCreated)
	log.Printf("Errors:         %d", stats.errors)
	log.Printf("Total vehicles: %d", stats.trimsCreated)
	log.Println("\n✓ Import complete!")
}

// Helper functions
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
