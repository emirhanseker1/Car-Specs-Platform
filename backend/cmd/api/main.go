package main

import (
	"log"
	"net/http"
	"os"

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
	trimRepo := repository.NewTrimRepository(database.DB)

	// Initialize services
	brandService := service.NewBrandService(brandRepo)
	modelService := service.NewModelService(modelRepo, brandRepo)
	trimService := service.NewTrimService(trimRepo, modelRepo)

	// Initialize handlers
	brandHandler := handlers.NewBrandHandler(brandService)
	modelHandler := handlers.NewModelHandler(modelService)
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
	mux.HandleFunc("/api/brands/{brandId}/models", modelHandler.HandleListModelsByBrand)

	// Trim routes
	mux.HandleFunc("/api/trims", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			trimHandler.HandleCreateTrim(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/trims/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			trimHandler.HandleGetTrim(w, r)
		case http.MethodDelete:
			trimHandler.HandleDeleteTrim(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/models/{modelId}/trims", trimHandler.HandleListTrimsByModel)

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
	log.Printf("📊 Database: %s", os.Getenv("DB_PATH"))
	log.Printf("🔗 API endpoints:")
	log.Printf("   - GET    /api/brands")
	log.Printf("   - POST   /api/brands")
	log.Printf("   - GET    /api/brands/{id}")
	log.Printf("   - GET    /api/brands/{brandId}/models")
	log.Printf("   - GET    /api/models/{modelId}/trims")
	log.Printf("   - GET    /api/search")
	log.Printf("   - GET    /health")

	if err := http.ListenAndServe(":"+port, corsHandler(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
