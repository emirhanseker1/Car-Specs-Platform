package service

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/emirh/car-specs-ai/backend/internal/models"
	"github.com/emirh/car-specs-ai/backend/internal/repository"
)

type EngineService struct {
	repo *repository.EngineRepository
}

func NewEngineService(repo *repository.EngineRepository) *EngineService {
	return &EngineService{repo: repo}
}

// GetAll retrieves all engines with optional filtering
func (s *EngineService) GetAll(filters map[string]interface{}) ([]models.Engine, error) {
	return s.repo.GetAll(filters)
}

// GetByID retrieves an engine by ID
func (s *EngineService) GetByID(id int64) (*models.Engine, error) {
	return s.repo.GetByID(id)
}

// GetByCode retrieves an engine by its code
func (s *EngineService) GetByCode(code string) (*models.Engine, error) {
	return s.repo.GetByCode(code)
}

// Create creates a new engine with validation
func (s *EngineService) Create(engine *models.Engine) error {
	// Validation
	if err := s.validateEngine(engine); err != nil {
		return err
	}

	// Normalize code to uppercase
	engine.Code = strings.ToUpper(strings.TrimSpace(engine.Code))

	// Check for duplicate code
	existing, err := s.repo.GetByCode(engine.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("engine code already exists")
	}

	// Validate JSON fields
	if err := s.validateJSONFields(engine); err != nil {
		return err
	}

	return s.repo.Create(engine)
}

// Update updates an existing engine with validation
func (s *EngineService) Update(engine *models.Engine) error {
	// Validation
	if err := s.validateEngine(engine); err != nil {
		return err
	}

	// Check if engine exists
	existing, err := s.repo.GetByID(engine.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("engine not found")
	}

	// Normalize code to uppercase
	engine.Code = strings.ToUpper(strings.TrimSpace(engine.Code))

	// Check for duplicate code (excluding current engine)
	duplicateCheck, err := s.repo.GetByCode(engine.Code)
	if err != nil {
		return err
	}
	if duplicateCheck != nil && duplicateCheck.ID != engine.ID {
		return errors.New("engine code already exists")
	}

	// Validate JSON fields
	if err := s.validateJSONFields(engine); err != nil {
		return err
	}

	return s.repo.Update(engine)
}

// Delete deletes an engine by ID
func (s *EngineService) Delete(id int64) error {
	// Check if engine exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("engine not found")
	}

	return s.repo.Delete(id)
}

// validateEngine validates basic engine fields
func (s *EngineService) validateEngine(engine *models.Engine) error {
	if engine.Code == "" {
		return errors.New("engine code is required")
	}

	if len(engine.Code) < 2 || len(engine.Code) > 20 {
		return errors.New("engine code must be between 2 and 20 characters")
	}

	// Validate reliability rating if provided
	if engine.ReliabilityRating != nil {
		rating := *engine.ReliabilityRating
		if rating < 1 || rating > 10 {
			return errors.New("reliability rating must be between 1 and 10")
		}
	}

	// Validate production years if provided
	if engine.ProductionStartYear != nil && *engine.ProductionStartYear < 1900 {
		return errors.New("production start year must be after 1900")
	}

	if engine.ProductionEndYear != nil && *engine.ProductionEndYear < 1900 {
		return errors.New("production end year must be after 1900")
	}

	if engine.ProductionStartYear != nil && engine.ProductionEndYear != nil {
		if *engine.ProductionEndYear < *engine.ProductionStartYear {
			return errors.New("production end year must be after start year")
		}
	}

	return nil
}

// validateJSONFields validates that JSON string fields are valid JSON
func (s *EngineService) validateJSONFields(engine *models.Engine) error {
	// Validate TechnologyFeatures JSON
	if engine.TechnologyFeatures != nil && *engine.TechnologyFeatures != "" {
		if !isValidJSON(*engine.TechnologyFeatures) {
			return errors.New("technology_features must be valid JSON")
		}
	}

	// Validate CommonProblems JSON
	if engine.CommonProblems != nil && *engine.CommonProblems != "" {
		if !isValidJSON(*engine.CommonProblems) {
			return errors.New("common_problems must be valid JSON")
		}
	}

	// Validate Solutions JSON
	if engine.Solutions != nil && *engine.Solutions != "" {
		if !isValidJSON(*engine.Solutions) {
			return errors.New("solutions must be valid JSON")
		}
	}

	return nil
}

// isValidJSON checks if a string is valid JSON
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// ParseProblems parses the CommonProblems JSON field
func (s *EngineService) ParseProblems(engine *models.Engine) ([]models.EngineProblem, error) {
	if engine.CommonProblems == nil || *engine.CommonProblems == "" {
		return []models.EngineProblem{}, nil
	}

	var problems []models.EngineProblem
	if err := json.Unmarshal([]byte(*engine.CommonProblems), &problems); err != nil {
		return nil, err
	}

	return problems, nil
}

// ParseSolutions parses the Solutions JSON field
func (s *EngineService) ParseSolutions(engine *models.Engine) ([]models.EngineSolution, error) {
	if engine.Solutions == nil || *engine.Solutions == "" {
		return []models.EngineSolution{}, nil
	}

	var solutions []models.EngineSolution
	if err := json.Unmarshal([]byte(*engine.Solutions), &solutions); err != nil {
		return nil, err
	}

	return solutions, nil
}

// ParseTechnologyFeatures parses the TechnologyFeatures JSON field
func (s *EngineService) ParseTechnologyFeatures(engine *models.Engine) ([]models.TechnologyFeature, error) {
	if engine.TechnologyFeatures == nil || *engine.TechnologyFeatures == "" {
		return []models.TechnologyFeature{}, nil
	}

	var features []models.TechnologyFeature
	if err := json.Unmarshal([]byte(*engine.TechnologyFeatures), &features); err != nil {
		return nil, err
	}

	return features, nil
}
