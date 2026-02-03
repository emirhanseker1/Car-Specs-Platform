package repository

import (
	"database/sql"
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
)

type BrandRepository struct {
	db *sql.DB
}

func NewBrandRepository(db *sql.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// Create inserts a new brand
func (r *BrandRepository) Create(brand *models.Brand) error {
	query := `
		INSERT INTO brands (name, country, logo_url)
		VALUES (?, ?, ?)
	`
	result, err := r.db.Exec(query, brand.Name, brand.Country, brand.LogoURL)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	brand.ID = id
	return nil
}

// GetByID retrieves a brand by ID
func (r *BrandRepository) GetByID(id int64) (*models.Brand, error) {
	query := `
		SELECT id, name, country, logo_url, created_at, updated_at
		FROM brands
		WHERE id = ?
	`
	brand := &models.Brand{}
	err := r.db.QueryRow(query, id).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Country,
		&brand.LogoURL,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("brand not found")
		}
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	return brand, nil
}

// GetByName retrieves a brand by name (case-insensitive)
func (r *BrandRepository) GetByName(name string) (*models.Brand, error) {
	query := `
		SELECT id, name, country, logo_url, created_at, updated_at
		FROM brands
		WHERE LOWER(name) = LOWER(?)
	`
	brand := &models.Brand{}
	err := r.db.QueryRow(query, name).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Country,
		&brand.LogoURL,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("brand not found")
		}
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	return brand, nil
}

// List retrieves all brands
func (r *BrandRepository) List() ([]*models.Brand, error) {
	query := `
		SELECT id, name, country, logo_url, created_at, updated_at
		FROM brands
		ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var brands []*models.Brand
	for rows.Next() {
		brand := &models.Brand{}
		err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.Country,
			&brand.LogoURL,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// Update updates a brand
func (r *BrandRepository) Update(brand *models.Brand) error {
	query := `
		UPDATE brands
		SET name = ?, country = ?, logo_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, brand.Name, brand.Country, brand.LogoURL, brand.ID)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}

	return nil
}

// Delete deletes a brand
func (r *BrandRepository) Delete(id int64) error {
	query := `DELETE FROM brands WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}
