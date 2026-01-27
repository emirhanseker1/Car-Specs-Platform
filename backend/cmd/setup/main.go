package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// TargetCar represents a specific car model to scrape for the Turkish market
type TargetCar struct {
	Brand      string // Brand name for API query
	Model      string // Model name for API query (Global name)
	Year       int    // Target year
	Alias      string // Turkish market name (if different from Model)
	Generation string // Optional: e.g., "F30", "G20" for BMW 3 Series
}

// API Ninjas response structure
type NinjasCarResponse struct {
	Make         string      `json:"make"`
	Model        string      `json:"model"`
	Year         int         `json:"year"`
	Class        string      `json:"class"`
	Cylinders    interface{} `json:"cylinders"`
	Displacement interface{} `json:"displacement"`
	Drive        string      `json:"drive"`
	FuelType     string      `json:"fuel_type"`
	Highway_MPG  interface{} `json:"highway_mpg"`
	City_MPG     interface{} `json:"city_mpg"`
	Transmission string      `json:"transmission"`
}

// SerpApi Google Images response structure
type SerpApiResponse struct {
	ImagesResults []struct {
		Original  string `json:"original"`
		Thumbnail string `json:"thumbnail"`
		Title     string `json:"title"`
	} `json:"images_results"`
}

type SetupService struct {
	db           *sql.DB
	ninjasAPIKey string
	serpApiKey   string
	httpClient   *http.Client
}

func NewSetupService(db *sql.DB) *SetupService {
	return &SetupService{
		db:           db,
		ninjasAPIKey: os.Getenv("NINJAS_API_KEY"),
		serpApiKey:   os.Getenv("SERPAPI_KEY"),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// GetTurkishMarketTargets returns the curated list of ~50 popular Turkish market cars
func GetTurkishMarketTargets() []TargetCar {
	return []TargetCar{
		// Sedans - Popular in Turkey
		{Brand: "fiat", Model: "egea", Year: 2023, Alias: "Tipo", Generation: ""}, // Try Egea (Turkish name) first
		{Brand: "fiat", Model: "egea", Year: 2024, Alias: "Tipo", Generation: ""},
		{Brand: "renault", Model: "megane", Year: 2022, Alias: "", Generation: ""},
		{Brand: "renault", Model: "megane", Year: 2023, Alias: "", Generation: ""},
		{Brand: "renault", Model: "taliant", Year: 2023, Alias: "Symbol", Generation: ""}, // Symbol = Taliant
		{Brand: "toyota", Model: "corolla", Year: 2023, Alias: "", Generation: ""},
		{Brand: "toyota", Model: "corolla", Year: 2024, Alias: "", Generation: ""},
		{Brand: "honda", Model: "civic", Year: 2022, Alias: "", Generation: ""},
		{Brand: "honda", Model: "civic", Year: 2023, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "passat", Year: 2021, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "passat", Year: 2022, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "jetta", Year: 2018, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "jetta", Year: 2022, Alias: "", Generation: ""},
		{Brand: "ford", Model: "focus", Year: 2022, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "elantra", Year: 2023, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "i30", Year: 2023, Alias: "", Generation: ""},

		// Hatchbacks - Very popular in Turkey
		{Brand: "renault", Model: "clio", Year: 2023, Alias: "", Generation: ""},
		{Brand: "renault", Model: "clio", Year: 2024, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "golf", Year: 2022, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "golf", Year: 2023, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "polo", Year: 2023, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "i20", Year: 2023, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "i20", Year: 2024, Alias: "", Generation: ""},
		{Brand: "opel", Model: "corsa", Year: 2023, Alias: "", Generation: ""},
		{Brand: "ford", Model: "fiesta", Year: 2021, Alias: "", Generation: ""},
		{Brand: "peugeot", Model: "208", Year: 2023, Alias: "", Generation: ""},

		// Premium Sedans - BMW, Mercedes, Audi
		{Brand: "bmw", Model: "3 series", Year: 2016, Alias: "", Generation: "F30"},
		{Brand: "bmw", Model: "3 series", Year: 2021, Alias: "", Generation: "G20"},
		{Brand: "bmw", Model: "5 series", Year: 2022, Alias: "", Generation: "G30"},
		{Brand: "mercedes-benz", Model: "c-class", Year: 2022, Alias: "C-Class", Generation: ""},
		{Brand: "mercedes-benz", Model: "e-class", Year: 2022, Alias: "E-Class", Generation: ""},
		{Brand: "audi", Model: "a3", Year: 2022, Alias: "", Generation: ""},
		{Brand: "audi", Model: "a4", Year: 2022, Alias: "", Generation: ""},
		{Brand: "audi", Model: "a6", Year: 2022, Alias: "", Generation: ""},

		// SUV/Crossover - Growing Turkish market
		{Brand: "nissan", Model: "qashqai", Year: 2022, Alias: "", Generation: ""},
		{Brand: "nissan", Model: "qashqai", Year: 2023, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "tiguan", Year: 2022, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "tiguan", Year: 2023, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "t-roc", Year: 2023, Alias: "", Generation: ""},
		{Brand: "peugeot", Model: "3008", Year: 2023, Alias: "", Generation: ""},
		{Brand: "dacia", Model: "duster", Year: 2023, Alias: "", Generation: ""},
		{Brand: "renault", Model: "captur", Year: 2023, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "tucson", Year: 2023, Alias: "", Generation: ""},
		{Brand: "toyota", Model: "c-hr", Year: 2023, Alias: "", Generation: ""},
		{Brand: "kia", Model: "sportage", Year: 2023, Alias: "", Generation: ""},
		{Brand: "skoda", Model: "karoq", Year: 2023, Alias: "", Generation: ""},

		// MPV/Family Cars
		{Brand: "renault", Model: "kangoo", Year: 2023, Alias: "", Generation: ""},
		{Brand: "citroen", Model: "berlingo", Year: 2023, Alias: "", Generation: ""},
		{Brand: "peugeot", Model: "rifter", Year: 2023, Alias: "", Generation: ""},
		// Note: Fiat Doblo may not be well-represented in API Ninjas (US-focused)

		// Electric/Hybrid - Emerging in Turkey
		{Brand: "renault", Model: "zoe", Year: 2023, Alias: "", Generation: ""},
		{Brand: "volkswagen", Model: "id.3", Year: 2023, Alias: "", Generation: ""},
		{Brand: "hyundai", Model: "kona", Year: 2023, Alias: "", Generation: ""},
		{Brand: "toyota", Model: "prius", Year: 2023, Alias: "", Generation: ""},
	}
}

// DropAllTables drops all existing tables
func (s *SetupService) DropAllTables() error {
	log.Println("🗑️  Dropping all tables...")

	tables := []string{"trim_features", "features", "trims", "models", "brands"}
	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		log.Printf("  ✓ Dropped table: %s", table)
	}

	return nil
}

// CreateSchema creates all tables from scratch
func (s *SetupService) CreateSchema() error {
	log.Println("🏗️  Creating fresh schema...")

	migrations := []string{
		// Brands table
		`CREATE TABLE IF NOT EXISTS brands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			country TEXT,
			logo_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Models table
		`CREATE TABLE IF NOT EXISTS models (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			brand_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			body_style TEXT,
			segment TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
			UNIQUE(brand_id, name)
		)`,

		// Trims table with image_url
		`CREATE TABLE IF NOT EXISTS trims (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			model_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			year INTEGER NOT NULL,
			generation TEXT,
			is_facelift BOOLEAN DEFAULT 0,
			market TEXT DEFAULT 'TR',
			engine_type TEXT,
			fuel_type TEXT,
			displacement_cc INTEGER,
			cylinders INTEGER,
			cylinder_layout TEXT,
			power_hp INTEGER,
			power_kw INTEGER,
			torque_nm INTEGER,
			engine_code TEXT,
			acceleration_0_100 REAL,
			top_speed_kmh INTEGER,
			fuel_consumption_city REAL,
			fuel_consumption_highway REAL,
			fuel_consumption_combined REAL,
			co2_emissions INTEGER,
			emission_standard TEXT,
			transmission_type TEXT,
			gears INTEGER,
			drivetrain TEXT,
			length_mm INTEGER,
			width_mm INTEGER,
			height_mm INTEGER,
			wheelbase_mm INTEGER,
			ground_clearance_mm INTEGER,
			curb_weight_kg INTEGER,
			gross_weight_kg INTEGER,
			luggage_capacity_l INTEGER,
			luggage_capacity_max_l INTEGER,
			fuel_tank_capacity_l INTEGER,
			tire_size_front TEXT,
			tire_size_rear TEXT,
			wheel_size_inches INTEGER,
			seating_capacity INTEGER DEFAULT 5,
			doors INTEGER,
			image_url TEXT,
			msrp_price REAL,
			currency TEXT DEFAULT 'TRY',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE
		)`,
	}

	for i, migration := range migrations {
		if _, err := s.db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i, err)
		}
	}

	log.Println("  ✓ Schema created successfully")
	return nil
}

// FetchCarFromNinjas queries API Ninjas with improved retry logic
func (s *SetupService) FetchCarFromNinjas(brand, model string, year int) (*NinjasCarResponse, error) {
	// Generate multiple query variations to improve match rate
	modelVariations := []string{
		model,                              // Original: "3 series"
		strings.ToLower(model),             // Lowercase: "3 series"
		strings.ReplaceAll(model, " ", ""), // No spaces: "3series"
		strings.ReplaceAll(strings.ToLower(model), " ", ""), // Lowercase no spaces: "3series"
	}

	// Add common abbreviations
	if strings.Contains(model, "series") {
		modelVariations = append(modelVariations, strings.ReplaceAll(model, "series", ""))  // "3 "
		modelVariations = append(modelVariations, strings.ReplaceAll(model, " series", "")) // "3"
	}

	var lastErr error
	attemptedQueries := []string{}

	for _, modelVar := range modelVariations {
		modelVar = strings.TrimSpace(modelVar)
		if modelVar == "" {
			continue
		}

		// Build query
		queryURL := fmt.Sprintf("https://api.api-ninjas.com/v1/cars?make=%s&model=%s&year=%d",
			url.QueryEscape(brand),
			url.QueryEscape(modelVar),
			year,
		)

		attemptedQueries = append(attemptedQueries, fmt.Sprintf("%s %s %d", brand, modelVar, year))

		req, err := http.NewRequest("GET", queryURL, nil)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("X-Api-Key", s.ninjasAPIKey)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		var cars []NinjasCarResponse
		if err := json.Unmarshal(body, &cars); err != nil {
			lastErr = err
			continue
		}

		if len(cars) > 0 {
			// Success! Found a match
			log.Printf("  ✓ Found %d variant(s) using query: %s %s %d", len(cars), brand, modelVar, year)
			return &cars[0], nil
		}

		// No results with this variation, try next
		lastErr = fmt.Errorf("no results found")
	}

	// All variations failed
	log.Printf("  ⚠️  Failed after trying %d query variations:", len(attemptedQueries))
	for _, q := range attemptedQueries {
		log.Printf("      - %s", q)
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all query variations failed: %w", lastErr)
	}
	return nil, fmt.Errorf("no results found for any query variation")
}

// FindCarImageSerpApi searches for car image using SerpApi (Google Images)
// Uses Turkish alias if available for better local market results
func (s *SetupService) FindCarImageSerpApi(brand, model, alias string, year int) (string, error) {
	if s.serpApiKey == "" {
		return "", nil // Skip if not configured
	}

	// Use Turkish alias for image search if available, otherwise use global model name
	searchModel := model
	if alias != "" {
		searchModel = alias
	}

	// Context-aware query optimized for Turkish market
	query := fmt.Sprintf("%s %s %d exterior white studio",
		brand, searchModel, year)

	// SerpApi endpoint for Google Images
	baseURL := "https://serpapi.com/search.json"
	params := url.Values{}
	params.Add("q", query)
	params.Add("tbm", "isch") // Target: Images
	params.Add("api_key", s.serpApiKey)
	params.Add("num", "3") // Get 3 results

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("SerpApi returned %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SerpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", err
	}

	// Return first valid image URL
	if len(searchResp.ImagesResults) > 0 {
		imageURL := searchResp.ImagesResults[0].Original
		if isValidImageURL(imageURL) {
			return imageURL, nil
		}
	}

	return "", fmt.Errorf("no valid image found")
}

// isValidImageURL checks if URL is valid (lenient for SerpApi)
func isValidImageURL(urlStr string) bool {
	// SerpApi returns valid image URLs that don't always end in .jpg/.png
	// Accept any non-empty URL that starts with http
	return len(urlStr) > 0 && (strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://"))
}

// Helper functions for type conversion
func toInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		var result int
		fmt.Sscanf(v, "%d", &result)
		return result
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
		var result float64
		fmt.Sscanf(v, "%f", &result)
		return result
	}
	return 0.0
}

// GetOrCreateBrand gets existing brand or creates new one
func (s *SetupService) GetOrCreateBrand(name string) (int64, error) {
	// Check if exists
	var id int64
	err := s.db.QueryRow("SELECT id FROM brands WHERE LOWER(name) = LOWER(?)", name).Scan(&id)
	if err == nil {
		return id, nil
	}

	// Create new
	result, err := s.db.Exec("INSERT INTO brands (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetOrCreateModel gets existing model or creates new one
// Uses Turkish alias for model name if provided
func (s *SetupService) GetOrCreateModel(brandID int64, name, alias, bodyStyle string) (int64, error) {
	// Use alias if available (Turkish market name)
	modelName := name
	if alias != "" {
		modelName = alias
	}

	// Check if exists
	var id int64
	err := s.db.QueryRow(
		"SELECT id FROM models WHERE brand_id = ? AND LOWER(name) = LOWER(?)",
		brandID, modelName,
	).Scan(&id)
	if err == nil {
		return id, nil
	}

	// Create new
	result, err := s.db.Exec(
		"INSERT INTO models (brand_id, name, body_style) VALUES (?, ?, ?)",
		brandID, modelName, bodyStyle,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// CreateTrim creates a new trim with Turkish market settings
func (s *SetupService) CreateTrim(modelID int64, car NinjasCarResponse, imageURL, generation string) error {
	cylinders := toInt(car.Cylinders)
	displacement := toFloat(car.Displacement)
	cityMPG := toInt(car.City_MPG)
	highwayMPG := toInt(car.Highway_MPG)

	displacementCC := int(displacement * 1000)

	var fuelConsumption float64
	if cityMPG > 0 && highwayMPG > 0 {
		avgMPG := float64(cityMPG+highwayMPG) / 2.0
		fuelConsumption = 235.214 / avgMPG
	}

	// Trim name includes generation if specified
	trimName := fmt.Sprintf("%s %d", car.Model, car.Year)
	if generation != "" {
		trimName = fmt.Sprintf("%s %s %d", car.Model, generation, car.Year)
	}

	_, err := s.db.Exec(`
		INSERT INTO trims (
			model_id, name, year, generation, fuel_type, displacement_cc, cylinders,
			transmission_type, drivetrain, fuel_consumption_combined,
			image_url, market, currency, seating_capacity
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, modelID, trimName, car.Year, generation, car.FuelType, displacementCC, cylinders,
		car.Transmission, car.Drive, fuelConsumption, imageURL, "TR", "TRY", 5)

	return err
}

// PopulateDatabase fetches and populates data using targeted Turkish market list
func (s *SetupService) PopulateDatabase() error {
	targets := GetTurkishMarketTargets()

	log.Printf("\n🎯 Targeting %d specific Turkish market vehicles\n", len(targets))

	stats := struct {
		brands    int
		models    int
		trims     int
		images    int
		errors    int
		skipped   int
		processed int
	}{}

	for _, target := range targets {
		stats.processed++
		log.Printf("\n[%d/%d] 📥 Fetching: %s %s %d (Alias: %s, Gen: %s)",
			stats.processed, len(targets),
			strings.ToUpper(target.Brand), strings.ToUpper(target.Model),
			target.Year, target.Alias, target.Generation)

		// Step 1: Fetch from API Ninjas
		car, err := s.FetchCarFromNinjas(target.Brand, target.Model, target.Year)
		if err != nil {
			log.Printf("  ⚠️  API search failed: %v", err)
			log.Printf("  💡 Will use fallback with basic info")
			car = nil // Set to nil so fallback is used below
		}

		// Step 2: Get or create brand (do this regardless of API success)
		brandName := target.Brand
		if car != nil {
			brandName = car.Make // Use API's brand name if available
		}

		brandID, err := s.GetOrCreateBrand(brandName)
		if err != nil {
			log.Printf("  ❌ Failed to create brand %s: %v", brandName, err)
			stats.errors++
			continue
		}
		if brandID > int64(stats.brands) {
			stats.brands++
			log.Printf("  ✓ Created brand: %s", brandName)
		}

		// Step 3: Get or create model
		modelName := target.Model
		if car != nil {
			modelName = car.Model // Use API's model name if available
		}

		modelID, err := s.GetOrCreateModel(brandID, modelName, target.Alias, "")
		if err != nil {
			log.Printf("  ❌ Failed to create model: %v", err)
			stats.errors++
			continue
		}
		if modelID > int64(stats.models) {
			stats.models++
			log.Printf("  ✓ Created model: %s", modelName)
		}

		// Step 4: Find image (use alias if available for Turkish market)
		imageURL := ""
		if s.serpApiKey != "" {
			log.Printf("  🖼️  Searching for image...")
			var imgErr error
			imageURL, imgErr = s.FindCarImageSerpApi(target.Brand, target.Model, target.Alias, target.Year)
			if imgErr != nil {
				log.Printf("  ⚠️  Image search failed: %v", imgErr)
			} else if imageURL != "" {
				log.Printf("  ✓ Found image: %.60s...", imageURL)
				stats.images++
			}
			// Rate limiting (SerpApi)
			time.Sleep(1 * time.Second)
		}

		// Step 5: Create trim - use API data if available, otherwise use fallback
		if car != nil {
			// We have API data - use it!
			log.Printf("  ✓ Got specs: %s %s", car.Make, car.Model)
			if err := s.CreateTrim(modelID, *car, imageURL, target.Generation); err != nil {
				log.Printf("  ❌ Failed to create trim: %v", err)
				stats.errors++
				continue
			}
		} else {
			// No API data - create fallback entry
			log.Printf("  ⚠️  API returned no data, creating fallback entry...")
			if err := s.CreateFallbackTrim(modelID, target.Brand, target.Model, target.Year, imageURL, target.Generation); err != nil {
				log.Printf("  ❌ Failed to create fallback trim: %v", err)
				stats.errors++
				continue
			}
		}

		stats.trims++
		log.Printf("  ✅ Successfully created trim")

		// Rate limiting between API Ninjas calls
		time.Sleep(2 * time.Second)
	}

	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("📊 Turkish Market Setup Summary:")
	log.Println(strings.Repeat("=", 50))
	log.Printf("  Targets processed:  %d / %d", stats.processed, len(targets))
	log.Printf("  Brands created:     %d", stats.brands)
	log.Printf("  Models created:     %d", stats.models)
	log.Printf("  Trims created:      %d", stats.trims)
	log.Printf("  Images found:       %d", stats.images)
	log.Printf("  Skipped (no data):  %d", stats.skipped)
	log.Printf("  Errors:             %d", stats.errors)
	log.Println(strings.Repeat("=", 50))

	return nil
}

func main() {
	log.Println("=== Turkish Market Database - Add New Vehicles ===\n")
	log.Println("ℹ️  This script will ADD new vehicles to the database.")
	log.Println("ℹ️  Existing data will NOT be deleted.\n")

	// Load environment variables - try multiple paths
	envPaths := []string{
		"../../.env.setup", // From cmd/setup -> backend/.env.setup
		"../../.env",       // From cmd/setup -> backend/.env
		".env.setup",       // If running from backend/
		".env",             // If running from backend/
	}

	loaded := false
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("✓ Loaded environment from: %s\n", envPath)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Validate required keys
	if os.Getenv("NINJAS_API_KEY") == "" {
		log.Fatal("❌ NINJAS_API_KEY is required")
	}

	// Connect to database - use absolute path to ensure we use the same DB as the API
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Get absolute path to backend/vehicles.db
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		// We're in backend/cmd/setup, go up 2 levels to backend/
		dbPath = filepath.Join(cwd, "..", "..", "vehicles.db")
		dbPath, err = filepath.Abs(dbPath)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
	}

	log.Printf("📁 Using database: %s\n", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize setup service
	service := NewSetupService(db)

	// Ensure schema exists (will create if not exists, skip if already exists)
	log.Println("📝 Ensuring database schema exists...")
	if err := service.CreateSchema(); err != nil {
		// This is OK - tables might already exist
		log.Printf("  ℹ️  Tables already exist (this is normal)\n")
	} else {
		log.Println("  ✓ Schema is ready")
	}

	// Populate with targeted Turkish market data (append mode)
	log.Println("\n🚗 Adding Turkish market vehicles...")
	if err := service.PopulateDatabase(); err != nil {
		log.Fatalf("Failed to populate database: %v", err)
	}

	log.Println("\n✅ Complete! New vehicles have been added to the database.")
	log.Println("💡 Tip: You can run this script again anytime to add more vehicles.")
}
