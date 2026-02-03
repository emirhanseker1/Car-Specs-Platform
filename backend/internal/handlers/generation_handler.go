package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/service"
)

type GenerationHandler struct {
	generationService *service.GenerationService
}

func NewGenerationHandler(generationService *service.GenerationService) *GenerationHandler {
	return &GenerationHandler{
		generationService: generationService,
	}
}

// GenerationDTO represents the response format for a generation
type GenerationDTO struct {
	ID          int64   `json:"id"`
	ModelID     int64   `json:"model_id"`
	Code        string  `json:"code"`
	Name        *string `json:"name,omitempty"`
	StartYear   int     `json:"start_year"`
	EndYear     *int    `json:"end_year,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	Description *string `json:"description,omitempty"`
	IsCurrent   bool    `json:"is_current"`
	Platform    *string `json:"platform,omitempty"`
	TrimCount   int     `json:"trim_count,omitempty"`
}

// HandleListByModel handles GET /api/models/{modelId}/generations
func (h *GenerationHandler) HandleListByModel(w http.ResponseWriter, r *http.Request) {
	modelIDStr := r.PathValue("modelId")
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	generations, err := h.generationService.ListGenerationsByModel(modelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to DTOs and add trim counts
	var dtos []GenerationDTO
	for _, gen := range generations {
		dto := GenerationDTO{
			ID:          gen.ID,
			ModelID:     gen.ModelID,
			Code:        gen.Code,
			Name:        gen.Name,
			StartYear:   gen.StartYear,
			EndYear:     gen.EndYear,
			ImageURL:    gen.ImageURL,
			Description: gen.Description,
			IsCurrent:   gen.IsCurrent,
			Platform:    gen.Platform,
		}

		// Get trim count for this generation
		count, err := h.generationService.GetTrimCount(gen.ID)
		if err == nil {
			dto.TrimCount = count
		}

		dtos = append(dtos, dto)
	}

	// Wrap in response object to match frontend expectations
	response := map[string]interface{}{
		"value": dtos,
		"Count": len(dtos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetGeneration handles GET /api/generations/{generationId}
func (h *GenerationHandler) HandleGetGeneration(w http.ResponseWriter, r *http.Request) {
	generationIDStr := r.PathValue("generationId")
	generationID, err := strconv.ParseInt(generationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid generation ID", http.StatusBadRequest)
		return
	}

	generation, err := h.generationService.GetGeneration(generationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert to DTO
	dto := GenerationDTO{
		ID:          generation.ID,
		ModelID:     generation.ModelID,
		Code:        generation.Code,
		Name:        generation.Name,
		StartYear:   generation.StartYear,
		EndYear:     generation.EndYear,
		ImageURL:    generation.ImageURL,
		Description: generation.Description,
		IsCurrent:   generation.IsCurrent,
		Platform:    generation.Platform,
	}

	// Get trim count
	count, err := h.generationService.GetTrimCount(generationID)
	if err == nil {
		dto.TrimCount = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
