package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/emirh/car-specs/backend/internal/models"
)

type ModelRepository struct {
	db *sql.DB
}

func NewModelRepository(db *sql.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

// Create inserts a new model
func (r *ModelRepository) Create(model *models.Model) error {
	query := `
		INSERT INTO models (brand_id, name, body_style, segment)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, model.BrandID, model.Name, model.BodyStyle, model.Segment)
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	model.ID = id
	return nil
}

// GetByID retrieves a model by ID with optional brand join
func (r *ModelRepository) GetByID(id int64, includeBrand bool) (*models.Model, error) {
	var query string
	if includeBrand {
		query = `
			SELECT 
				m.id, m.brand_id, m.name, m.body_style, m.segment, m.created_at, m.updated_at,
				b.id, b.name, b.country, b.logo_url, b.created_at, b.updated_at
			FROM models m
			LEFT JOIN brands b ON m.brand_id = b.id
			WHERE m.id = ?
		`
	} else {
		query = `
			SELECT id, brand_id, name, body_style, segment, created_at, updated_at
			FROM models
			WHERE id = ?
		`
	}

	model := &models.Model{}
	var err error

	if includeBrand {
		brand := &models.Brand{}
		err = r.db.QueryRow(query, id).Scan(
			&model.ID, &model.BrandID, &model.Name, &model.BodyStyle, &model.Segment,
			&model.CreatedAt, &model.UpdatedAt,
			&brand.ID, &brand.Name, &brand.Country, &brand.LogoURL,
			&brand.CreatedAt, &brand.UpdatedAt,
		)
		if err == nil {
			model.Brand = brand
		}
	} else {
		err = r.db.QueryRow(query, id).Scan(
			&model.ID, &model.BrandID, &model.Name, &model.BodyStyle, &model.Segment,
			&model.CreatedAt, &model.UpdatedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("model not found")
		}
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	return model, nil
}

// ListByBrand retrieves all models for a brand
func (r *ModelRepository) ListByBrand(brandID int64) ([]*models.Model, error) {
	query := `
		SELECT id, brand_id, name, body_style, segment, created_at, updated_at
		FROM models
		WHERE brand_id = ?
		ORDER BY name
	`
	rows, err := r.db.Query(query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer rows.Close()

	var modelsList []*models.Model
	for rows.Next() {
		model := &models.Model{}
		err := rows.Scan(
			&model.ID, &model.BrandID, &model.Name, &model.BodyStyle, &model.Segment,
			&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}
		modelsList = append(modelsList, model)
	}

	return modelsList, nil
}

// Update updates a model
func (r *ModelRepository) Update(model *models.Model) error {
	query := `
		UPDATE models
		SET brand_id = ?, name = ?, body_style = ?, segment = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, model.BrandID, model.Name, model.BodyStyle, model.Segment, model.ID)
	if err != nil {
		return fmt.Errorf("failed to update model: %w", err)
	}

	return nil
}

// Delete deletes a model
func (r *ModelRepository) Delete(id int64) error {
	query := `DELETE FROM models WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	return nil
}

// ListVehiclesByName retrieves aggregated vehicle list (Generation focused)
func (r *ModelRepository) ListVehiclesByName(brandName string) ([]*models.VehicleListItem, error) {
	query := `
		SELECT 
			COALESCE(g.id, 0), 
			m.id,
			MAX(b.name), 
			MAX(m.name), 
			COALESCE(g.code, ''), 
			g.start_year, g.end_year,
			GROUP_CONCAT(t.name)
		FROM models m
		JOIN brands b ON m.brand_id = b.id
		LEFT JOIN generations g ON g.model_id = m.id
		LEFT JOIN trims t ON t.generation_id = g.id
		WHERE b.name = ?
		GROUP BY m.id, g.id
		ORDER BY MAX(m.name), g.start_year DESC
	`

	rows, err := r.db.Query(query, brandName)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var list []*models.VehicleListItem
	for rows.Next() {
		v := &models.VehicleListItem{}
		var startYear, endYear sql.NullInt64
		// var isFacelift sql.NullBool // Removed
		var engines sql.NullString

		// Order: g.id, m.id, b.name, m.name, g.code, start, end, engines
		err := rows.Scan(
			&v.ID,
			&v.ModelID,
			&v.Brand,
			&v.Model,
			&v.Generation,
			&startYear, &endYear, &engines,
		)

		if err != nil {
			return nil, err
		}

		if startYear.Valid || endYear.Valid {
			v.GenerationMeta = &models.GenerationMeta{}
			if startYear.Valid {
				y := int(startYear.Int64)
				v.GenerationMeta.StartYear = &y
			}
			if endYear.Valid {
				y := int(endYear.Int64)
				v.GenerationMeta.EndYear = &y
			}
		}

		if engines.Valid {
			raw := strings.Split(engines.String, ",")
			unique := make(map[string]bool)
			for _, e := range raw {
				e = strings.TrimSpace(e)
				if e != "" && !unique[e] {
					unique[e] = true
					v.EngineOptions = append(v.EngineOptions, e)
				}
			}
		}
		list = append(list, v)
	}
	return list, nil
}

// GetGeneration retrieves a generation by ID with joined model/brand info
func (r *ModelRepository) GetGeneration(id int64) (*models.Generation, *models.Model, error) {
	query := `
		SELECT 
			g.id, g.model_id, g.code, g.name, g.start_year, g.end_year,
			m.id, m.brand_id, m.name, m.body_style,
			b.id, b.name, b.logo_url
		FROM models m
		JOIN generations g ON g.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		WHERE g.id = ?
	`

	g := &models.Generation{}
	m := &models.Model{}
	b := &models.Brand{}

	var startYear, endYear sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&g.ID, &g.ModelID, &g.Code, &g.Name, &startYear, &endYear,
		&m.ID, &m.BrandID, &m.Name, &m.BodyStyle,
		&b.ID, &b.Name, &b.LogoURL,
	)

	if err != nil {
		fmt.Printf("GetGeneration Error for ID %d: %v\n", id, err)
		return nil, nil, err
	}

	if startYear.Valid {
		g.StartYear = int(startYear.Int64)
	}
	if endYear.Valid {
		y := int(endYear.Int64)
		g.EndYear = &y
	}

	m.Brand = b
	g.ModelID = m.ID // Already set

	return g, m, nil
}
