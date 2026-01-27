package service

import (
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
)

type ModelService struct {
	modelRepo *repository.ModelRepository
	brandRepo *repository.BrandRepository
}

func NewModelService(modelRepo *repository.ModelRepository, brandRepo *repository.BrandRepository) *ModelService {
	return &ModelService{
		modelRepo: modelRepo,
		brandRepo: brandRepo,
	}
}

// CreateModel creates a new model with validation
func (s *ModelService) CreateModel(brandID int64, name string, bodyStyle, segment *string) (*models.Model, error) {
	// Validation
	if name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	// Verify brand exists
	_, err := s.brandRepo.GetByID(brandID)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	model := &models.Model{
		BrandID:   brandID,
		Name:      name,
		BodyStyle: bodyStyle,
		Segment:   segment,
	}

	if err := s.modelRepo.Create(model); err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return model, nil
}

// GetModel retrieves a model by ID with optional brand information
func (s *ModelService) GetModel(id int64, includeBrand bool) (*models.Model, error) {
	model, err := s.modelRepo.GetByID(id, includeBrand)
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}
	return model, nil
}

// ListModelsByBrand retrieves all models for a brand
func (s *ModelService) ListModelsByBrand(brandID int64) ([]*models.Model, error) {
	// Verify brand exists
	_, err := s.brandRepo.GetByID(brandID)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	models, err := s.modelRepo.ListByBrand(brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	return models, nil
}

// UpdateModel updates a model
func (s *ModelService) UpdateModel(id int64, brandID int64, name string, bodyStyle, segment *string) (*models.Model, error) {
	// Validation
	if name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	// Check if model exists
	model, err := s.modelRepo.GetByID(id, false)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	// Verify brand exists if changing
	if brandID != model.BrandID {
		_, err := s.brandRepo.GetByID(brandID)
		if err != nil {
			return nil, fmt.Errorf("brand not found: %w", err)
		}
	}

	// Update fields
	model.BrandID = brandID
	model.Name = name
	model.BodyStyle = bodyStyle
	model.Segment = segment

	if err := s.modelRepo.Update(model); err != nil {
		return nil, fmt.Errorf("failed to update model: %w", err)
	}

	return model, nil
}

// DeleteModel deletes a model
func (s *ModelService) DeleteModel(id int64) error {
	// Check if model exists
	_, err := s.modelRepo.GetByID(id, false)
	if err != nil {
		return fmt.Errorf("model not found: %w", err)
	}

	if err := s.modelRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	return nil
}
