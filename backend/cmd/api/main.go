package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/emirh/car-specs/backend/internal/database"
	"github.com/emirh/car-specs/backend/internal/handlers"
	"github.com/emirh/car-specs/backend/internal/repository"
	"github.com/emirh/car-specs/backend/internal/service"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize repositories
	brandRepo := repository.NewBrandRepository(database.DB)
	modelRepo := repository.NewModelRepository(database.DB)
	generationRepo := repository.NewGenerationRepository(database.DB)
	trimRepo := repository.NewTrimRepository(database.DB)

	// Initialize services
	brandService := service.NewBrandService(brandRepo)
	modelService := service.NewModelService(modelRepo, brandRepo)
	generationService := service.NewGenerationService(generationRepo, modelRepo)
	trimService := service.NewTrimService(trimRepo, modelRepo)

	// Initialize handlers
	brandHandler := handlers.NewBrandHandler(brandService)
	modelHandler := handlers.NewModelHandler(modelService, trimService, brandService)
	generationHandler := handlers.NewGenerationHandler(generationService)
	trimHandler := handlers.NewTrimHandler(trimService)

	// Setup routes
	mux := http.NewServeMux()

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Brand routes
	mux.HandleFunc("/api/brands", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			brandHandler.HandleCreateBrand(w, r)
		case http.MethodGet:
			brandHandler.HandleListBrands(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// IMPORTANT: More specific brand routes must come BEFORE generic {id} route
	mux.HandleFunc("/api/brands/{brandId}/models", modelHandler.HandleListModelsByBrand)
	mux.HandleFunc("/api/brands/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			brandHandler.HandleGetBrand(w, r)
		case http.MethodPut:
			brandHandler.HandleUpdateBrand(w, r)
		case http.MethodDelete:
			brandHandler.HandleDeleteBrand(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Model routes
	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			modelHandler.HandleCreateModel(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/models/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			modelHandler.HandleGetModel(w, r)
		case http.MethodPut:
			modelHandler.HandleUpdateModel(w, r)
		case http.MethodDelete:
			modelHandler.HandleDeleteModel(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Generation routes
	mux.HandleFunc("GET /api/models/{modelId}/generations", generationHandler.HandleListByModel)
	mux.HandleFunc("GET /api/generations/{generationId}", generationHandler.HandleGetGeneration)

	// Legacy/Frontend aggregate route
	mux.HandleFunc("/api/vehicles", modelHandler.HandleListVehicles)
	// Vehicle Details (Aggregation)
	mux.HandleFunc("GET /api/vehicles/{id}", modelHandler.HandleGetVehicleDetails)

	// Trim routes
	mux.HandleFunc("/api/trims", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			trimHandler.HandleCreateTrim(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("GET /api/trims/{id}", trimHandler.HandleGetTrim)
	mux.HandleFunc("DELETE /api/trims/{id}", trimHandler.HandleDeleteTrim)
	mux.HandleFunc("/api/models/{modelId}/trims", trimHandler.HandleListTrimsByModel)
	mux.HandleFunc("GET /api/generations/{generationId}/trims", trimHandler.HandleListTrimsByGeneration)

	// Search route
	mux.HandleFunc("/api/search", trimHandler.HandleSearchTrims)

	// Featured route for homepage
	mux.HandleFunc("/api/featured", trimHandler.HandleGetFeaturedTrims)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	absDB, _ := filepath.Abs(os.Getenv("DB_PATH"))
	if os.Getenv("DB_PATH") == "" {
		absDB, _ = filepath.Abs("./vehicles.db")
	}
	log.Printf("📊 Database: %s", absDB)
	log.Printf("🔗 API endpoints:")
	log.Printf("   - GET    /api/brands")
	log.Printf("   - POST   /api/brands")
	log.Printf("   - GET    /api/brands/{id}")
	log.Printf("   - GET    /api/brands/{brandId}/models")
	log.Printf("   - GET    /api/models/{modelId}/generations")
	log.Printf("   - GET    /api/generations/{generationId}")
	log.Printf("   - GET    /api/generations/{generationId}/trims")
	log.Printf("   - GET    /api/models/{modelId}/trims")
	log.Printf("   - GET    /api/search")
	log.Printf("   - GET    /health")

	if err := http.ListenAndServe(":"+port, corsHandler(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
