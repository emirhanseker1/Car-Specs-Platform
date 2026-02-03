package repository

import (
	"database/sql"
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
)

type GenerationRepository struct {
	db *sql.DB
}

func NewGenerationRepository(db *sql.DB) *GenerationRepository {
	return &GenerationRepository{db: db}
}

// ListByModel retrieves all generations for a specific model
func (r *GenerationRepository) ListByModel(modelID int64) ([]*models.Generation, error) {
	query := `
		SELECT 
			g.id, 
			g.model_id, 
			g.code, 
			g.name, 
			g.start_year, 
			g.end_year,
			g.image_url,
			g.created_at,
			g.updated_at,
			COUNT(t.id) as trim_count
		FROM generations g
		LEFT JOIN trims t ON t.generation_id = g.id
		WHERE g.model_id = ?
		GROUP BY g.id
		ORDER BY g.start_year DESC
	`

	rows, err := r.db.Query(query, modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to list generations: %w", err)
	}
	defer rows.Close()

	var generations []*models.Generation
	for rows.Next() {
		g := &models.Generation{}
		var trimCount int
		var endYear sql.NullInt64
		var name sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&g.ID,
			&g.ModelID,
			&g.Code,
			&name,
			&g.StartYear,
			&endYear,
			&imageURL,
			&g.CreatedAt,
			&g.UpdatedAt,
			&trimCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan generation: %w", err)
		}

		// Handle nullable fields
		if name.Valid {
			g.Name = &name.String
		}
		if endYear.Valid {
			y := int(endYear.Int64)
			g.EndYear = &y
		}
		if imageURL.Valid {
			g.ImageURL = &imageURL.String
		}

		generations = append(generations, g)
	}

	return generations, nil
}

func (r *GenerationRepository) GetByID(id int64) (*models.Generation, error) {
	query := `
		SELECT 
			id, 
			model_id, 
			code, 
			name, 
			start_year, 
			end_year,
			created_at,
			updated_at
		FROM generations
		WHERE id = ?
	`

	g := &models.Generation{}
	var endYear sql.NullInt64
	var name sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&g.ID,
		&g.ModelID,
		&g.Code,
		&name,
		&g.StartYear,
		&endYear,
		&g.CreatedAt,
		&g.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("generation not found")
		}
		return nil, fmt.Errorf("failed to get generation: %w", err)
	}

	// Handle nullable fields
	if name.Valid {
		g.Name = &name.String
	}
	if endYear.Valid {
		y := int(endYear.Int64)
		g.EndYear = &y
	}

	return g, nil
}

// GetTrimCount returns the number of trims for a generation
func (r *GenerationRepository) GetTrimCount(generationID int64) (int, error) {
	query := `SELECT COUNT(*) FROM trims WHERE generation_id = ?`
	var count int
	err := r.db.QueryRow(query, generationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get trim count: %w", err)
	}
	return count, nil
}
