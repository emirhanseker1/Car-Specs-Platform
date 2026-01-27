package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/emirh/car-specs/backend/internal/database"
	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
	"github.com/emirh/car-specs/backend/internal/service"
	"github.com/joho/godotenv"
)

// API Ninjas response structure
type NinjasCarResponse struct {
	Make         string      `json:"make"`
	Model        string      `json:"model"`
	Year         int         `json:"year"`
	Class        string      `json:"class"`
	Cylinders    interface{} `json:"cylinders"`    // Can be int or string
	Displacement interface{} `json:"displacement"` // Can be float or string
	Drive        string      `json:"drive"`
	FuelType     string      `json:"fuel_type"`
	Highway_MPG  interface{} `json:"highway_mpg"` // Can be int or string
	City_MPG     interface{} `json:"city_mpg"`    // Can be int or string
	Transmission string      `json:"transmission"`
}

// Google Custom Search response structure
type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

type IngestionService struct {
	brandService   *service.BrandService
	modelService   *service.ModelService
	trimService    *service.TrimService
	ninjasAPIKey   string
	googleAPIKey   string
	searchEngineID string
	httpClient     *http.Client
}

func NewIngestionService(brandSvc *service.BrandService, modelSvc *service.ModelService, trimSvc *service.TrimService) *IngestionService {
	return &IngestionService{
		brandService:   brandSvc,
		modelService:   modelSvc,
		trimService:    trimSvc,
		ninjasAPIKey:   os.Getenv("NINJAS_API_KEY"),
		googleAPIKey:   os.Getenv("GOOGLE_API_KEY"),
		searchEngineID: os.Getenv("SEARCH_ENGINE_ID"),
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchCarsFromNinjas fetches cars from API Ninjas
func (s *IngestionService) FetchCarsFromNinjas(make string, year int) ([]NinjasCarResponse, error) {
	baseURL := "https://api.api-ninjas.com/v1/cars"
	params := url.Values{}
	params.Add("make", make)
	params.Add("year", fmt.Sprintf("%d", year))
	// Note: 'limit' parameter is only available for premium users

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", s.ninjasAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from API Ninjas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Ninjas returned status %d: %s", resp.StatusCode, string(body))
	}

	var cars []NinjasCarResponse
	if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return cars, nil
}

// FindCarImage searches for a car image using Google Custom Search
func (s *IngestionService) FindCarImage(make, model string, year int) (string, error) {
	if s.googleAPIKey == "" || s.searchEngineID == "" {
		log.Println("⚠️  Google API credentials not set, skipping image search")
		return "", nil
	}

	query := fmt.Sprintf("%s %s %d exterior studio white background", make, model, year)
	baseURL := "https://www.googleapis.com/customsearch/v1"

	params := url.Values{}
	params.Add("key", s.googleAPIKey)
	params.Add("cx", s.searchEngineID)
	params.Add("q", query)
	params.Add("searchType", "image")
	params.Add("num", "1")
	params.Add("imgSize", "large")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to search Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Google API returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp GoogleSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", fmt.Errorf("failed to decode Google response: %w", err)
	}

	if len(searchResp.Items) == 0 {
		return "", fmt.Errorf("no images found")
	}

	return searchResp.Items[0].Link, nil
}

// ProcessCar processes a single car from API Ninjas
func (s *IngestionService) ProcessCar(car NinjasCarResponse) error {
	log.Printf("\n[Processing] %s %s %d", car.Make, car.Model, car.Year)

	// Step 1: Get or Create Brand
	brand, err := s.brandService.GetBrandByName(car.Make)
	if err != nil {
		brand, err = s.brandService.CreateBrand(car.Make, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to create brand: %w", err)
		}
		log.Printf("  ✓ Created brand: %s (ID: %d)", car.Make, brand.ID)
	} else {
		log.Printf("  ✓ Found existing brand: %s (ID: %d)", car.Make, brand.ID)
	}

	// Step 2: Get or Create Model
	modelsList, err := s.modelService.ListModelsByBrand(brand.ID)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	var model *models.Model
	for _, m := range modelsList {
		if strings.EqualFold(m.Name, car.Model) {
			model = m
			break
		}
	}

	if model == nil {
		bodyStyle := car.Class
		model, err = s.modelService.CreateModel(brand.ID, car.Model, &bodyStyle, nil)
		if err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}
		log.Printf("  ✓ Created model: %s (ID: %d)", car.Model, model.ID)
	} else {
		log.Printf("  ✓ Found existing model: %s (ID: %d)", car.Model, model.ID)
	}

	// Step 3: Check if trim already exists (idempotency)
	existingTrims, err := s.trimService.ListTrimsByModel(model.ID)
	if err != nil {
		return fmt.Errorf("failed to list trims: %w", err)
	}

	for _, trim := range existingTrims {
		if trim.Year == car.Year {
			log.Printf("  ⏭️  Trim already exists for year %d, skipping", car.Year)
			return nil
		}
	}

	// Step 4: Map API Ninjas data to Trim
	fuelType := mapFuelType(car.FuelType)
	transmission := mapTransmission(car.Transmission)
	drivetrain := mapDrivetrain(car.Drive)

	// Safely convert fields that might be strings or numbers
	cylinders := toInt(car.Cylinders)
	displacement := toFloat(car.Displacement)
	cityMPG := toInt(car.City_MPG)
	highwayMPG := toInt(car.Highway_MPG)

	// Convert displacement from liters to CC
	displacementCC := int(displacement * 1000)

	// Convert MPG to L/100km (approximate)
	var fuelConsumption float64
	if cityMPG > 0 && highwayMPG > 0 {
		avgMPG := float64(cityMPG+highwayMPG) / 2.0
		fuelConsumption = 235.214 / avgMPG // MPG to L/100km conversion
	}

	trim := &models.Trim{
		ModelID:             model.ID,
		Name:                fmt.Sprintf("%s %d", car.Model, car.Year),
		Year:                car.Year,
		FuelType:            &fuelType,
		DisplacementCC:      &displacementCC,
		Cylinders:           &cylinders,
		TransmissionType:    &transmission,
		Drivetrain:          &drivetrain,
		FuelConsumptionComb: &fuelConsumption,
		Market:              "US",
		Currency:            "USD",
		SeatingCapacity:     5,
	}

	// Step 5: Create Trim
	if err := s.trimService.CreateTrim(trim); err != nil {
		return fmt.Errorf("failed to create trim: %w", err)
	}
	log.Printf("  ✓ Created trim: %s (ID: %d)", trim.Name, trim.ID)

	// Step 6: Find and update image (with rate limiting)
	time.Sleep(2 * time.Second) // Rate limiting for Google API

	imageURL, err := s.FindCarImage(car.Make, car.Model, car.Year)
	if err != nil {
		log.Printf("  ⚠️  Failed to find image: %v", err)
	} else if imageURL != "" {
		trim.ImageURL = &imageURL
		// Update trim with image URL
		// Note: You'd need to add an Update method to TrimService
		log.Printf("  ✓ Found image: %s", imageURL)
	}

	return nil
}

// Helper functions to safely convert interface{} to numbers
func toInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse string to int
		if i, err := fmt.Sscanf(v, "%d", new(int)); err == nil && i == 1 {
			var result int
			fmt.Sscanf(v, "%d", &result)
			return result
		}
	}
	return 0
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		// Try to parse string to float
		if f, err := fmt.Sscanf(v, "%f", new(float64)); err == nil && f == 1 {
			var result float64
			fmt.Sscanf(v, "%f", &result)
			return result
		}
	}
	return 0.0
}

// Helper functions to map API Ninjas data to our schema
func mapFuelType(apiType string) string {
	switch strings.ToLower(apiType) {
	case "gas", "gasoline":
		return "Benzin"
	case "diesel":
		return "Dizel"
	case "electric":
		return "Elektrik"
	case "hybrid":
		return "Hybrid"
	default:
		return "Benzin"
	}
}

func mapTransmission(apiTrans string) string {
	lower := strings.ToLower(apiTrans)
	if strings.Contains(lower, "auto") || strings.Contains(lower, "a") {
		return "Automatic"
	}
	if strings.Contains(lower, "manual") || strings.Contains(lower, "m") {
		return "Manual"
	}
	if strings.Contains(lower, "cvt") {
		return "CVT"
	}
	return "Automatic"
}

func mapDrivetrain(apiDrive string) string {
	switch strings.ToLower(apiDrive) {
	case "fwd":
		return "FWD"
	case "rwd":
		return "RWD"
	case "awd", "4wd":
		return "AWD"
	default:
		return "FWD"
	}
}

func main() {
	log.Println("=== Vehicle Ingestion Service ===")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Validate required environment variables
	ninjasKey := os.Getenv("NINJAS_API_KEY")
	if ninjasKey == "" {
		log.Fatal("❌ NINJAS_API_KEY is required")
	}

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

	// Initialize ingestion service
	ingestionService := NewIngestionService(brandService, modelService, trimService)

	// Target brands and years
	targetBrands := []string{"bmw", "audi", "volkswagen", "mercedes-benz", "toyota", "ford"}
	targetYears := []int{2023, 2024}

	stats := struct {
		totalProcessed int
		totalCreated   int
		totalSkipped   int
		totalErrors    int
	}{}

	// Process each brand and year
	for _, brand := range targetBrands {
		for _, year := range targetYears {
			log.Printf("\n📥 Fetching %s vehicles for year %d...", strings.ToUpper(brand), year)

			cars, err := ingestionService.FetchCarsFromNinjas(brand, year)
			if err != nil {
				log.Printf("❌ Failed to fetch %s %d: %v", brand, year, err)
				stats.totalErrors++
				continue
			}

			log.Printf("Found %d vehicles", len(cars))

			for _, car := range cars {
				stats.totalProcessed++

				if err := ingestionService.ProcessCar(car); err != nil {
					log.Printf("  ❌ Error: %v", err)
					stats.totalErrors++
				} else {
					stats.totalCreated++
				}

				// Rate limiting between API calls
				time.Sleep(2 * time.Second)
			}
		}
	}

	// Print summary
	log.Println("\n=== Ingestion Summary ===")
	log.Printf("Total processed: %d", stats.totalProcessed)
	log.Printf("Total created:   %d", stats.totalCreated)
	log.Printf("Total skipped:   %d", stats.totalSkipped)
	log.Printf("Total errors:    %d", stats.totalErrors)
	log.Println("\n✓ Ingestion complete!")
}
