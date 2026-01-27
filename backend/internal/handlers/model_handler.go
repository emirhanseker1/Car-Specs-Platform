package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/service"
)

type ModelHandler struct {
	service *service.ModelService
}

func NewModelHandler(service *service.ModelService) *ModelHandler {
	return &ModelHandler{service: service}
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
func (h *ModelHandler) HandleListModelsByBrand(w http.ResponseWriter, r *http.Request) {
	brandIDStr := r.PathValue("brandId")
	brandID, err := strconv.ParseInt(brandIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	models, err := h.service.ListModelsByBrand(brandID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
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
