package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/emirh/car-specs/backend/internal/formatter"
	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/service"
)

type TrimHandler struct {
	service *service.TrimService
}

func NewTrimHandler(service *service.TrimService) *TrimHandler {
	return &TrimHandler{service: service}
}

// HandleCreateTrim handles POST /api/trims
func (h *TrimHandler) HandleCreateTrim(w http.ResponseWriter, r *http.Request) {
	var trim models.Trim
	if err := json.NewDecoder(r.Body).Decode(&trim); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTrim(&trim); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(trim)
}

// HandleGetTrim handles GET /api/trims/:id
func (h *TrimHandler) HandleGetTrim(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trim ID", http.StatusBadRequest)
		return
	}

	includeRelations := r.URL.Query().Get("include_relations") == "true"

	trim, err := h.service.GetTrim(id, includeRelations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format data for professional display
	formatter.FormatTrim(trim)

	// Frontend expects a wrapper object: { vehicle: {...}, trims: [...] }

	// 1. Fetch all trims for this model (for the dropdown)
	var siblingTrims []*models.Trim
	if trim.ModelID != 0 {
		siblingTrims, _ = h.service.ListTrimsByModel(trim.ModelID)
		formatter.FormatTrims(siblingTrims)
	} else {
		// Fallback if no ModelID
		siblingTrims = []*models.Trim{trim}
	}

	// 2. Construct Vehicle object
	brandName := ""
	modelName := ""
	if trim.Model != nil {
		modelName = trim.Model.Name
		if trim.Model.Brand != nil {
			brandName = trim.Model.Brand.Name
		}
	}

	vehicleObj := map[string]interface{}{
		"id":         trim.ModelID,
		"brand":      brandName,
		"model":      modelName,
		"generation": trim.Generation,
		"image_url":  trim.ImageURL,
	}

	// 3. Construct specific Trim objects for the list
	// We marshal the full trim object to map to include all fields (specs are flattened in Trim struct)
	trimList := make([]map[string]interface{}, 0)
	for _, t := range siblingTrims {
		// Marshal to JSON then Unmarshal to map to get all fields with correct json tags
		b, _ := json.Marshal(t)
		var m map[string]interface{}
		json.Unmarshal(b, &m)

		// Frontend expects "specs" field.
		// Since our fields are flattened, we can try to send them as is,
		// OR we rely on helper to group them.
		// For now send flattened, Frontend might not show specs but at least won't crash.
		trimList = append(trimList, m)
	}

	response := map[string]interface{}{
		"vehicle": vehicleObj,
		"trims":   trimList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSearchTrims handles GET /api/search
func (h *TrimHandler) HandleSearchTrims(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := make(map[string]interface{})
	if brand := query.Get("brand"); brand != "" {
		filters["brand"] = brand
	}
	if model := query.Get("model"); model != "" {
		filters["model"] = model
	}
	if fuelType := query.Get("fuel_type"); fuelType != "" {
		filters["fuel_type"] = fuelType
	}
	if transmission := query.Get("transmission"); transmission != "" {
		filters["transmission"] = transmission
	}
	if yearStr := query.Get("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filters["year"] = year
		}
	}

	trims, err := h.service.SearchTrims(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format all trims for professional display
	formatter.FormatTrims(trims)

	// Get facets for filters
	facets, _ := h.service.GetSearchFacets()

	response := map[string]interface{}{
		"results": trims,
		"facets":  facets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListTrimsByModel handles GET /api/models/:modelId/trims
func (h *TrimHandler) HandleListTrimsByModel(w http.ResponseWriter, r *http.Request) {
	modelIDStr := r.PathValue("modelId")
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	trims, err := h.service.ListTrimsByModel(modelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format all trims for professional display
	formatter.FormatTrims(trims)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trims)
}

// HandleDeleteTrim handles DELETE /api/trims/:id
func (h *TrimHandler) HandleDeleteTrim(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid trim ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTrim(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetFeaturedTrims handles GET /api/featured
func (h *TrimHandler) HandleGetFeaturedTrims(w http.ResponseWriter, r *http.Request) {
	// Get 4 random trims with images for homepage
	trims, err := h.service.GetFeatured(4)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format all trims for professional display
	formatter.FormatTrims(trims)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trims)
}

// HandleListTrimsByGeneration handles GET /api/generations/{generationId}/trims
func (h *TrimHandler) HandleListTrimsByGeneration(w http.ResponseWriter, r *http.Request) {
	generationIDStr := r.PathValue("generationId")
	generationID, err := strconv.ParseInt(generationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid generation ID", http.StatusBadRequest)
		return
	}

	trims, err := h.service.ListTrimsByGeneration(generationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format all trims for professional display
	formatter.FormatTrims(trims)

	// Create simplified response for trim selection
	type TrimListDTO struct {
		ID                 int64    `json:"id"`
		Name               string   `json:"name"`
		PowerHP            *int     `json:"power_hp,omitempty"`
		TorqueNM           *int     `json:"torque_nm,omitempty"`
		Acceleration0To100 *float64 `json:"acceleration_0_100,omitempty"`
		FuelType           *string  `json:"fuel_type,omitempty"`
		TransmissionType   *string  `json:"transmission_type,omitempty"`
		TransmissionCode   *string  `json:"transmission_code,omitempty"`
		Drivetrain         *string  `json:"drivetrain,omitempty"`
		Year               int      `json:"year"`
		StartYear          *int     `json:"start_year,omitempty"`
		EndYear            *int     `json:"end_year,omitempty"`
	}

	var trimDTOs []TrimListDTO
	for _, trim := range trims {
		dto := TrimListDTO{
			ID:                 trim.ID,
			Name:               trim.Name,
			PowerHP:            trim.PowerHP,
			TorqueNM:           trim.TorqueNM,
			Acceleration0To100: trim.Acceleration0To100,
			FuelType:           trim.FuelType,
			TransmissionType:   trim.TransmissionType,
			TransmissionCode:   trim.TransmissionCode,
			Drivetrain:         trim.Drivetrain,
			Year:               trim.Year,
			StartYear:          trim.StartYear,
			EndYear:            trim.EndYear,
		}
		trimDTOs = append(trimDTOs, dto)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trimDTOs)
}
