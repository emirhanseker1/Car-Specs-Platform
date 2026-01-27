package service

import (
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
	"github.com/emirh/car-specs/backend/internal/repository"
)

type BrandService struct {
	repo *repository.BrandRepository
}

func NewBrandService(repo *repository.BrandRepository) *BrandService {
	return &BrandService{repo: repo}
}

// CreateBrand creates a new brand with validation
func (s *BrandService) CreateBrand(name string, country, logoURL *string) (*models.Brand, error) {
	// Validation
	if name == "" {
		return nil, fmt.Errorf("brand name is required")
	}

	// Check if brand already exists
	existing, _ := s.repo.GetByName(name)
	if existing != nil {
		return nil, fmt.Errorf("brand '%s' already exists", name)
	}

	brand := &models.Brand{
		Name:    name,
		Country: country,
		LogoURL: logoURL,
	}

	if err := s.repo.Create(brand); err != nil {
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	return brand, nil
}

// GetBrand retrieves a brand by ID
func (s *BrandService) GetBrand(id int64) (*models.Brand, error) {
	brand, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}
	return brand, nil
}

// GetBrandByName retrieves a brand by name
func (s *BrandService) GetBrandByName(name string) (*models.Brand, error) {
	brand, err := s.repo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}
	return brand, nil
}

// ListBrands retrieves all brands
func (s *BrandService) ListBrands() ([]*models.Brand, error) {
	brands, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	return brands, nil
}

// UpdateBrand updates a brand
func (s *BrandService) UpdateBrand(id int64, name string, country, logoURL *string) (*models.Brand, error) {
	// Validation
	if name == "" {
		return nil, fmt.Errorf("brand name is required")
	}

	// Check if brand exists
	brand, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}

	// Update fields
	brand.Name = name
	brand.Country = country
	brand.LogoURL = logoURL

	if err := s.repo.Update(brand); err != nil {
		return nil, fmt.Errorf("failed to update brand: %w", err)
	}

	return brand, nil
}

// DeleteBrand deletes a brand
func (s *BrandService) DeleteBrand(id int64) error {
	// Check if brand exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("brand not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}

// GetOrCreateBrand gets an existing brand or creates it if it doesn't exist
func (s *BrandService) GetOrCreateBrand(name string, country *string) (*models.Brand, error) {
	// Try to get existing brand
	brand, err := s.repo.GetByName(name)
	if err == nil {
		return brand, nil
	}

	// Create new brand
	return s.CreateBrand(name, country, nil)
}
