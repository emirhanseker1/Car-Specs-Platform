package service

import (
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
)

type TrimService struct {
	trimRepo  *repository.TrimRepository
	modelRepo *repository.ModelRepository
}

func NewTrimService(trimRepo *repository.TrimRepository, modelRepo *repository.ModelRepository) *TrimService {
	return &TrimService{
		trimRepo:  trimRepo,
		modelRepo: modelRepo,
	}
}

// CreateTrim creates a new trim with validation
func (s *TrimService) CreateTrim(trim *models.Trim) error {
	// Validation
	if trim.Name == "" {
		return fmt.Errorf("trim name is required")
	}
	if trim.Year == 0 {
		return fmt.Errorf("year is required")
	}
	if trim.GenerationID == 0 {
		return fmt.Errorf("generation_id is required")
	}

	// Note: We could validate generation exists here, but skipping for simplicity
	// since the database foreign key will enforce it anyway

	// Business logic: Set defaults
	if trim.Market == "" {
		trim.Market = "TR"
	}
	if trim.Currency == "" {
		trim.Currency = "TRY"
	}
	if trim.SeatingCapacity == 0 {
		trim.SeatingCapacity = 5
	}

	if err := s.trimRepo.Create(trim); err != nil {
		return fmt.Errorf("failed to create trim: %w", err)
	}

	return nil
}

// GetTrim retrieves a trim by ID with optional relationships
func (s *TrimService) GetTrim(id int64, includeRelations bool) (*models.Trim, error) {
	trim, err := s.trimRepo.GetByID(id, includeRelations)
	if err != nil {
		return nil, fmt.Errorf("failed to get trim: %w", err)
	}
	return trim, nil
}

// SearchTrims searches for trims with filters
func (s *TrimService) SearchTrims(filters map[string]interface{}) ([]*models.Trim, error) {
	trims, err := s.trimRepo.Search(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search trims: %w", err)
	}
	return trims, nil
}

// ListTrimsByModel retrieves all trims for a model
func (s *TrimService) ListTrimsByModel(modelID int64) ([]*models.Trim, error) {
	// Verify model exists
	_, err := s.modelRepo.GetByID(modelID, false)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	trims, err := s.trimRepo.ListByModel(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to list trims: %w", err)
	}
	return trims, nil
}

// ListTrimsByGeneration retrieves all trims for a generation
func (s *TrimService) ListTrimsByGeneration(genID int64) ([]*models.Trim, error) {
	// We might validate generation exists, but skipping for now
	trims, err := s.trimRepo.ListByGeneration(genID)
	if err != nil {
		return nil, fmt.Errorf("failed to list trims: %w", err)
	}
	return trims, nil
}

// DeleteTrim deletes a trim
func (s *TrimService) DeleteTrim(id int64) error {
	// Check if trim exists
	_, err := s.trimRepo.GetByID(id, false)
	if err != nil {
		return fmt.Errorf("trim not found: %w", err)
	}

	if err := s.trimRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete trim: %w", err)
	}

	return nil
}

// GetFeatured retrieves random trims with images for the homepage
func (s *TrimService) GetFeatured(limit int) ([]*models.Trim, error) {
	trims, err := s.trimRepo.GetFeaturedTrims(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured trims: %w", err)
	}
	return trims, nil
}

// GetSearchFacets returns available filter options for search
func (s *TrimService) GetSearchFacets() (map[string]interface{}, error) {
	// This would query distinct values for filters
	// For now, returning a placeholder
	facets := map[string]interface{}{
		"fuel_types":    []string{},
		"transmissions": []string{},
		"min_year":      0,
		"max_year":      0,
		"min_hp":        0,
		"max_hp":        0,
	}
	return facets, nil
}
