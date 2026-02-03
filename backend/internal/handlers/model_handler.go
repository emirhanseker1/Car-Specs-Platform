package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/service"
)

type ModelHandler struct {
	service      *service.ModelService
	trimService  *service.TrimService
	brandService *service.BrandService
}

func NewModelHandler(service *service.ModelService, trimService *service.TrimService, brandService *service.BrandService) *ModelHandler {
	return &ModelHandler{
		service:      service,
		trimService:  trimService,
		brandService: brandService,
	}
}

// CreateModelRequest represents the request body for creating a model
type CreateModelRequest struct {
	BrandID   int64   `json:"brand_id"`
	Name      string  `json:"name"`
	BodyStyle *string `json:"body_style,omitempty"`
	Segment   *string `json:"segment,omitempty"`
}

// HandleCreateModel handles POST /api/models
func (h *ModelHandler) HandleCreateModel(w http.ResponseWriter, r *http.Request) {
	var req CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	model, err := h.service.CreateModel(req.BrandID, req.Name, req.BodyStyle, req.Segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model)
}

// HandleGetModel handles GET /api/models/:id
func (h *ModelHandler) HandleGetModel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	includeBrand := r.URL.Query().Get("include_brand") == "true"

	model, err := h.service.GetModel(id, includeBrand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

// HandleListModelsByBrand handles GET /api/brands/:brandId/models
// Supports both numeric ID and brand name (e.g., "audi", "bmw")
func (h *ModelHandler) HandleListModelsByBrand(w http.ResponseWriter, r *http.Request) {
	brandIDStr := r.PathValue("brandId")

	var brandID int64
	var err error

	// Try to parse as numeric ID first
	brandID, err = strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		// Not a number, treat as brand name - lookup brand by name
		brand, err := h.brandService.GetBrandByName(brandIDStr)
		if err != nil {
			http.Error(w, "Brand not found", http.StatusNotFound)
			return
		}
		brandID = brand.ID
	}

	models, err := h.service.ListModelsByBrand(brandID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wrap in response object to match frontend expectations
	response := map[string]interface{}{
		"value": models,
		"Count": len(models),
	}

	log.Printf("🔍 DEBUG: Sending wrapped response with %d models", len(models))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleUpdateModel handles PUT /api/models/:id
func (h *ModelHandler) HandleUpdateModel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	var req CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	model, err := h.service.UpdateModel(id, req.BrandID, req.Name, req.BodyStyle, req.Segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

// HandleDeleteModel handles DELETE /api/models/:id
func (h *ModelHandler) HandleDeleteModel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteModel(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleListVehicles handles GET /api/vehicles?brand=...
func (h *ModelHandler) HandleListVehicles(w http.ResponseWriter, r *http.Request) {
	brandName := r.URL.Query().Get("brand")
	if brandName == "" {
		http.Error(w, "brand query parameter is required", http.StatusBadRequest)
		return
	}

	vehicles, err := h.service.ListVehiclesByName(brandName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

// HandleGetVehicleDetails handles GET /api/vehicles/:id (Aggregation for Frontend)
func (h *ModelHandler) HandleGetVehicleDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 1. Try to fetch as Generation (since frontend uses Generation ID)
	gen, model, err := h.service.GetGeneration(id)
	if err != nil {
		// If not found, maybe fallback to Model?
		// For now, assume Generation ID as per new frontend flow.
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	// 2. Fetch Trims for this Generation
	trims, err := h.trimService.ListTrimsByGeneration(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Construct Response
	// Format Trims to include powertrain_meta and map to frontend key names if needed
	// Frontend PowertrainSelect.tsx expects:
	// vehicle: { id, brand, model, generation, image_url, generation_meta }
	// trims: [ { id, name, powertrain_meta: { ... } } ]

	type PowertrainMeta struct {
		EngineCode       *string `json:"engine_code,omitempty"`
		FuelType         *string `json:"fuel_type,omitempty"`
		DisplacementCC   *int    `json:"displacement_cc,omitempty"`
		PowerHP          *int    `json:"power_hp,omitempty"`
		TorqueNM         *int    `json:"torque_nm,omitempty"`
		TransmissionType *string `json:"transmission_type,omitempty"`
		Gears            *int    `json:"gears,omitempty"`
		Drive            *string `json:"drive,omitempty"`
	}

	type TrimResponse struct {
		ID             int64           `json:"id"`
		Name           string          `json:"name"`
		PowertrainMeta *PowertrainMeta `json:"powertrain_meta,omitempty"`
	}

	var trimList []TrimResponse
	for _, t := range trims {
		meta := &PowertrainMeta{
			EngineCode:       t.EngineCode,
			FuelType:         t.FuelType,
			DisplacementCC:   t.DisplacementCC,
			PowerHP:          t.PowerHP,
			TorqueNM:         t.TorqueNM,
			TransmissionType: t.TransmissionType,
			Gears:            t.Gears,
			Drive:            t.Drivetrain,
		}

		trimList = append(trimList, TrimResponse{
			ID:             t.ID,
			Name:           t.Name,
			PowertrainMeta: meta,
		})
	}

	// Construct Vehicle object
	vehicleObj := map[string]interface{}{
		"id":         gen.ID, // Use GenID as ID
		"brand":      model.Brand.Name,
		"model":      model.Name,
		"generation": gen.Code,
		"image_url":  "",
		"generation_meta": map[string]interface{}{
			"start_year":  gen.StartYear,
			"end_year":    gen.EndYear,
			"is_facelift": false,
		},
	}

	response := map[string]interface{}{
		"vehicle": vehicleObj,
		"trims":   trimList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
