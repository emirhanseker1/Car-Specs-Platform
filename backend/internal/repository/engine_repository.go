package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/emirh/car-specs-ai/backend/internal/models"
)

type EngineRepository struct {
	db *sql.DB
}

func NewEngineRepository(db *sql.DB) *EngineRepository {
	return &EngineRepository{db: db}
}

// GetAll retrieves all engines with optional filtering
func (r *EngineRepository) GetAll(filters map[string]interface{}) ([]models.Engine, error) {
	query := `SELECT * FROM engines WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	// Apply filters
	if fuelType, ok := filters["fuel_type"].(string); ok && fuelType != "" {
		query += fmt.Sprintf(" AND fuel_type = ?%d", argPos)
		args = append(args, fuelType)
		argPos++
	}

	if manufacturer, ok := filters["manufacturer"].(string); ok && manufacturer != "" {
		query += fmt.Sprintf(" AND manufacturer = ?%d", argPos)
		args = append(args, manufacturer)
		argPos++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (LOWER(code) LIKE ?%d OR LOWER(name) LIKE ?%d)", argPos, argPos+1)
		searchPattern := "%" + strings.ToLower(search) + "%"
		args = append(args, searchPattern, searchPattern)
		argPos += 2
	}

	query += " ORDER BY code ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engines []models.Engine
	for rows.Next() {
		var engine models.Engine
		err := rows.Scan(
			&engine.ID,
			&engine.Code,
			&engine.Name,
			&engine.Manufacturer,
			&engine.EngineType,
			&engine.FuelType,
			&engine.DisplacementCC,
			&engine.Cylinders,
			&engine.CylinderLayout,
			&engine.ValvesPerCylinder,
			&engine.Aspiration,
			&engine.Description,
			&engine.TechnologyFeatures,
			&engine.ProductionStartYear,
			&engine.ProductionEndYear,
			&engine.CommonProblems,
			&engine.Solutions,
			&engine.MaintenanceNotes,
			&engine.ReliabilityRating,
			&engine.CreatedAt,
			&engine.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		engines = append(engines, engine)
	}

	return engines, nil
}

// GetByID retrieves an engine by ID
func (r *EngineRepository) GetByID(id int64) (*models.Engine, error) {
	query := `SELECT * FROM engines WHERE id = ?`

	var engine models.Engine
	err := r.db.QueryRow(query, id).Scan(
		&engine.ID,
		&engine.Code,
		&engine.Name,
		&engine.Manufacturer,
		&engine.EngineType,
		&engine.FuelType,
		&engine.DisplacementCC,
		&engine.Cylinders,
		&engine.CylinderLayout,
		&engine.ValvesPerCylinder,
		&engine.Aspiration,
		&engine.Description,
		&engine.TechnologyFeatures,
		&engine.ProductionStartYear,
		&engine.ProductionEndYear,
		&engine.CommonProblems,
		&engine.Solutions,
		&engine.MaintenanceNotes,
		&engine.ReliabilityRating,
		&engine.CreatedAt,
		&engine.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

// GetByCode retrieves an engine by its code
func (r *EngineRepository) GetByCode(code string) (*models.Engine, error) {
	query := `SELECT * FROM engines WHERE UPPER(code) = UPPER(?)`

	var engine models.Engine
	err := r.db.QueryRow(query, code).Scan(
		&engine.ID,
		&engine.Code,
		&engine.Name,
		&engine.Manufacturer,
		&engine.EngineType,
		&engine.FuelType,
		&engine.DisplacementCC,
		&engine.Cylinders,
		&engine.CylinderLayout,
		&engine.ValvesPerCylinder,
		&engine.Aspiration,
		&engine.Description,
		&engine.TechnologyFeatures,
		&engine.ProductionStartYear,
		&engine.ProductionEndYear,
		&engine.CommonProblems,
		&engine.Solutions,
		&engine.MaintenanceNotes,
		&engine.ReliabilityRating,
		&engine.CreatedAt,
		&engine.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

// Create creates a new engine
func (r *EngineRepository) Create(engine *models.Engine) error {
	query := `
		INSERT INTO engines (
			code, name, manufacturer, engine_type, fuel_type,
			displacement_cc, cylinders, cylinder_layout, valves_per_cylinder,
			aspiration, description, technology_features,
			production_start_year, production_end_year,
			common_problems, solutions, maintenance_notes, reliability_rating
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		engine.Code,
		engine.Name,
		engine.Manufacturer,
		engine.EngineType,
		engine.FuelType,
		engine.DisplacementCC,
		engine.Cylinders,
		engine.CylinderLayout,
		engine.ValvesPerCylinder,
		engine.Aspiration,
		engine.Description,
		engine.TechnologyFeatures,
		engine.ProductionStartYear,
		engine.ProductionEndYear,
		engine.CommonProblems,
		engine.Solutions,
		engine.MaintenanceNotes,
		engine.ReliabilityRating,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	engine.ID = id
	return nil
}

// Update updates an existing engine
func (r *EngineRepository) Update(engine *models.Engine) error {
	query := `
		UPDATE engines SET
			code = ?, name = ?, manufacturer = ?, engine_type = ?, fuel_type = ?,
			displacement_cc = ?, cylinders = ?, cylinder_layout = ?, valves_per_cylinder = ?,
			aspiration = ?, description = ?, technology_features = ?,
			production_start_year = ?, production_end_year = ?,
			common_problems = ?, solutions = ?, maintenance_notes = ?, reliability_rating = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		engine.Code,
		engine.Name,
		engine.Manufacturer,
		engine.EngineType,
		engine.FuelType,
		engine.DisplacementCC,
		engine.Cylinders,
		engine.CylinderLayout,
		engine.ValvesPerCylinder,
		engine.Aspiration,
		engine.Description,
		engine.TechnologyFeatures,
		engine.ProductionStartYear,
		engine.ProductionEndYear,
		engine.CommonProblems,
		engine.Solutions,
		engine.MaintenanceNotes,
		engine.ReliabilityRating,
		engine.ID,
	)

	return err
}

// Delete deletes an engine by ID
func (r *EngineRepository) Delete(id int64) error {
	query := `DELETE FROM engines WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
