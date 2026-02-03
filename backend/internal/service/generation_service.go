package service

import (
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
)

type GenerationService struct {
	generationRepo *repository.GenerationRepository
	modelRepo      *repository.ModelRepository
}

func NewGenerationService(generationRepo *repository.GenerationRepository, modelRepo *repository.ModelRepository) *GenerationService {
	return &GenerationService{
		generationRepo: generationRepo,
		modelRepo:      modelRepo,
	}
}

// ListGenerationsByModel retrieves all generations for a model with validation
func (s *GenerationService) ListGenerationsByModel(modelID int64) ([]*models.Generation, error) {
	// Verify model exists
	_, err := s.modelRepo.GetByID(modelID, false)
	if err != nil {
		return nil, fmt.Errorf("model not found: %w", err)
	}

	generations, err := s.generationRepo.ListByModel(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to list generations: %w", err)
	}

	return generations, nil
}

// GetGeneration retrieves a generation by ID
func (s *GenerationService) GetGeneration(id int64) (*models.Generation, error) {
	generation, err := s.generationRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get generation: %w", err)
	}
	return generation, nil
}

// GetTrimCount returns the number of trims for a generation
func (s *GenerationService) GetTrimCount(generationID int64) (int, error) {
	return s.generationRepo.GetTrimCount(generationID)
}
