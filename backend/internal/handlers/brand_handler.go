package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/service"
)

type BrandHandler struct {
	service *service.BrandService
}

func NewBrandHandler(service *service.BrandService) *BrandHandler {
	return &BrandHandler{service: service}
}

// CreateBrandRequest represents the request body for creating a brand
type CreateBrandRequest struct {
	Name    string  `json:"name"`
	Country *string `json:"country,omitempty"`
	LogoURL *string `json:"logo_url,omitempty"`
}

// HandleCreateBrand handles POST /api/brands
func (h *BrandHandler) HandleCreateBrand(w http.ResponseWriter, r *http.Request) {
	var req CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	brand, err := h.service.CreateBrand(req.Name, req.Country, req.LogoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(brand)
}

// HandleGetBrand handles GET /api/brands/:id
func (h *BrandHandler) HandleGetBrand(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	log.Printf("🔍 DEBUG HandleGetBrand: received id='%s'", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	brand, err := h.service.GetBrand(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brand)
}

// HandleListBrands handles GET /api/brands
func (h *BrandHandler) HandleListBrands(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 DEBUG HandleListBrands: called")
	brands, err := h.service.ListBrands()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brands)
}

// HandleUpdateBrand handles PUT /api/brands/:id
func (h *BrandHandler) HandleUpdateBrand(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	var req CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	brand, err := h.service.UpdateBrand(id, req.Name, req.Country, req.LogoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brand)
}

// HandleDeleteBrand handles DELETE /api/brands/:id
func (h *BrandHandler) HandleDeleteBrand(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBrand(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
