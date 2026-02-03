package repository

import (
	"database/sql"
	"fmt"

	"github.com/emirh/car-specs/backend/internal/models"
)

type TrimRepository struct {
	db *sql.DB
}

func NewTrimRepository(db *sql.DB) *TrimRepository {
	return &TrimRepository{db: db}
}

// Create inserts a new trim
func (r *TrimRepository) Create(trim *models.Trim) error {
	query := `
		INSERT INTO trims (
			generation_id, name, year, generation, is_facelift, market,
			engine_type, fuel_type, displacement_cc, cylinders, cylinder_layout,
			power_hp, power_kw, torque_nm, engine_code,
			acceleration_0_100, top_speed_kmh,
			fuel_consumption_city, fuel_consumption_highway, fuel_consumption_combined,
			co2_emissions, emission_standard,
			transmission_type, gears, drivetrain,
			length_mm, width_mm, height_mm, wheelbase_mm, ground_clearance_mm,
			curb_weight_kg, gross_weight_kg,
			luggage_capacity_l, luggage_capacity_max_l, fuel_tank_capacity_l,
			tire_size_front, tire_size_rear, wheel_size_inches,
			seating_capacity, doors, image_url, msrp_price, currency
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?, ?, ?,
			?, ?,
			?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?, ?
		)
	`
	result, err := r.db.Exec(query,
		trim.GenerationID, trim.Name, trim.Year, trim.Generation, trim.IsFacelift, trim.Market,
		trim.EngineType, trim.FuelType, trim.DisplacementCC, trim.Cylinders, trim.CylinderLayout,
		trim.PowerHP, trim.PowerKW, trim.TorqueNM, trim.EngineCode,
		trim.Acceleration0To100, trim.TopSpeedKmh,
		trim.FuelConsumptionCity, trim.FuelConsumptionHwy, trim.FuelConsumptionComb,
		trim.CO2Emissions, trim.EmissionStandard,
		trim.TransmissionType, trim.Gears, trim.Drivetrain,
		trim.LengthMM, trim.WidthMM, trim.HeightMM, trim.WheelbaseMM, trim.GroundClearanceMM,
		trim.CurbWeightKG, trim.GrossWeightKG,
		trim.LuggageCapacityL, trim.LuggageCapacityMaxL, trim.FuelTankCapacityL,
		trim.TireSizeFront, trim.TireSizeRear, trim.WheelSizeInches,
		trim.SeatingCapacity, trim.Doors, trim.ImageURL, trim.MSRPPrice, trim.Currency,
	)
	if err != nil {
		return fmt.Errorf("failed to create trim: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	trim.ID = id
	return nil
}

// GetByID retrieves a trim by ID with optional model/brand join
func (r *TrimRepository) GetByID(id int64, includeRelations bool) (*models.Trim, error) {
	var query string
	if includeRelations {
		query = `
			SELECT 
				t.id, t.generation_id, t.model_id, t.name, t.year, t.start_year, t.end_year, t.generation, t.is_facelift, t.market,
				t.engine_type, t.fuel_type, t.displacement_cc, t.cylinders, t.cylinder_layout,
				t.power_hp, t.power_kw, t.torque_nm, t.engine_code,
				t.acceleration_0_100, t.top_speed_kmh,
				t.fuel_consumption_city, t.fuel_consumption_highway, t.fuel_consumption_combined,
				t.co2_emissions, t.emission_standard,
				t.transmission_type, t.transmission_code, t.gears, t.drivetrain,
				t.length_mm, t.width_mm, t.height_mm, t.wheelbase_mm, t.ground_clearance_mm,
				t.curb_weight_kg, t.gross_weight_kg,
				t.luggage_capacity_l, t.luggage_capacity_max_l, t.fuel_tank_capacity_l,
				t.tire_size_front, t.tire_size_rear, t.wheel_size_inches,
				t.seating_capacity, t.doors, t.image_url, t.msrp_price, t.currency,
				t.created_at, t.updated_at,
				m.id, m.brand_id, m.name, m.body_style, m.segment, m.created_at, m.updated_at,
				b.id, b.name, b.country, b.logo_url, b.created_at, b.updated_at
			FROM trims t
			LEFT JOIN generations g ON t.generation_id = g.id
			LEFT JOIN models m ON g.model_id = m.id
			LEFT JOIN brands b ON m.brand_id = b.id
			WHERE t.id = ?
		`
	} else {
		query = `
			SELECT 
				id, generation_id, model_id, name, year, start_year, end_year, generation, is_facelift, market,
				engine_type, fuel_type, displacement_cc, cylinders, cylinder_layout,
				power_hp, power_kw, torque_nm, engine_code,
				acceleration_0_100, top_speed_kmh,
				fuel_consumption_city, fuel_consumption_highway, fuel_consumption_combined,
				co2_emissions, emission_standard,
				transmission_type, transmission_code, gears, drivetrain,
				length_mm, width_mm, height_mm, wheelbase_mm, ground_clearance_mm,
				curb_weight_kg, gross_weight_kg,
				luggage_capacity_l, luggage_capacity_max_l, fuel_tank_capacity_l,
				tire_size_front, tire_size_rear, wheel_size_inches,
				seating_capacity, doors, image_url, msrp_price, currency,
				created_at, updated_at
			FROM trims WHERE id = ?
		`
	}

	trim := &models.Trim{}
	var err error

	if includeRelations {
		model := &models.Model{}
		brand := &models.Brand{}
		err = r.db.QueryRow(query, id).Scan(
			&trim.ID, &trim.GenerationID, &trim.ModelID, &trim.Name, &trim.Year, &trim.StartYear, &trim.EndYear, &trim.Generation, &trim.IsFacelift, &trim.Market,
			&trim.EngineType, &trim.FuelType, &trim.DisplacementCC, &trim.Cylinders, &trim.CylinderLayout,
			&trim.PowerHP, &trim.PowerKW, &trim.TorqueNM, &trim.EngineCode,
			&trim.Acceleration0To100, &trim.TopSpeedKmh,
			&trim.FuelConsumptionCity, &trim.FuelConsumptionHwy, &trim.FuelConsumptionComb,
			&trim.CO2Emissions, &trim.EmissionStandard,
			&trim.TransmissionType, &trim.TransmissionCode, &trim.Gears, &trim.Drivetrain,
			&trim.LengthMM, &trim.WidthMM, &trim.HeightMM, &trim.WheelbaseMM, &trim.GroundClearanceMM,
			&trim.CurbWeightKG, &trim.GrossWeightKG,
			&trim.LuggageCapacityL, &trim.LuggageCapacityMaxL, &trim.FuelTankCapacityL,
			&trim.TireSizeFront, &trim.TireSizeRear, &trim.WheelSizeInches,
			&trim.SeatingCapacity, &trim.Doors, &trim.ImageURL, &trim.MSRPPrice, &trim.Currency,
			&trim.CreatedAt, &trim.UpdatedAt,
			&model.ID, &model.BrandID, &model.Name, &model.BodyStyle, &model.Segment, &model.CreatedAt, &model.UpdatedAt,
			&brand.ID, &brand.Name, &brand.Country, &brand.LogoURL, &brand.CreatedAt, &brand.UpdatedAt,
		)
		if err == nil {
			model.Brand = brand
			trim.Model = model
		}
	} else {
		// Simplified scan for trim-only query matching ListByGeneration
		err = r.db.QueryRow(query, id).Scan(
			&trim.ID, &trim.GenerationID, &trim.ModelID, &trim.Name, &trim.Year, &trim.StartYear, &trim.EndYear, &trim.Generation, &trim.IsFacelift, &trim.Market,
			&trim.EngineType, &trim.FuelType, &trim.DisplacementCC, &trim.Cylinders, &trim.CylinderLayout,
			&trim.PowerHP, &trim.PowerKW, &trim.TorqueNM, &trim.EngineCode,
			&trim.Acceleration0To100, &trim.TopSpeedKmh,
			&trim.FuelConsumptionCity, &trim.FuelConsumptionHwy, &trim.FuelConsumptionComb,
			&trim.CO2Emissions, &trim.EmissionStandard,
			&trim.TransmissionType, &trim.TransmissionCode, &trim.Gears, &trim.Drivetrain,
			&trim.LengthMM, &trim.WidthMM, &trim.HeightMM, &trim.WheelbaseMM, &trim.GroundClearanceMM,
			&trim.CurbWeightKG, &trim.GrossWeightKG,
			&trim.LuggageCapacityL, &trim.LuggageCapacityMaxL, &trim.FuelTankCapacityL,
			&trim.TireSizeFront, &trim.TireSizeRear, &trim.WheelSizeInches,
			&trim.SeatingCapacity, &trim.Doors, &trim.ImageURL, &trim.MSRPPrice, &trim.Currency,
			&trim.CreatedAt, &trim.UpdatedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trim not found")
		}
		return nil, fmt.Errorf("failed to get trim: %w", err)
	}

	return trim, nil
}

// Search searches trims with filters and includes brand/model data
func (r *TrimRepository) Search(filters map[string]interface{}) ([]*models.Trim, error) {
	query := `
		SELECT 
			t.id, t.generation_id, t.name, t.year, t.start_year, t.end_year, t.generation, t.is_facelift, t.market,
			t.engine_type, t.fuel_type, t.displacement_cc, t.cylinders, t.cylinder_layout,
			t.power_hp, t.power_kw, t.torque_nm, t.engine_code,
			t.acceleration_0_100, t.top_speed_kmh,
			t.fuel_consumption_city, t.fuel_consumption_highway, t.fuel_consumption_combined,
			t.co2_emissions, t.emission_standard,
			t.transmission_type, t.gears, t.drivetrain,
			t.length_mm, t.width_mm, t.height_mm, t.wheelbase_mm, t.ground_clearance_mm,
			t.curb_weight_kg, t.gross_weight_kg,
			t.luggage_capacity_l, t.luggage_capacity_max_l, t.fuel_tank_capacity_l,
			t.tire_size_front, t.tire_size_rear, t.wheel_size_inches,
			t.seating_capacity, t.doors, t.image_url, t.msrp_price, t.currency,
			t.created_at, t.updated_at,
			m.id as model_id, m.brand_id, m.name as model_name, m.body_style, m.segment,
			m.created_at as model_created_at, m.updated_at as model_updated_at,
			b.id as brand_id, b.name as brand_name, b.country, b.logo_url,
			b.created_at as brand_created_at, b.updated_at as brand_updated_at
		FROM trims t
		LEFT JOIN generations g ON t.generation_id = g.id
		LEFT JOIN models m ON g.model_id = m.id
		LEFT JOIN brands b ON m.brand_id = b.id
		WHERE 1=1
	`
	args := []interface{}{}

	// Build dynamic WHERE clause
	if brandName, ok := filters["brand"]; ok {
		query += " AND LOWER(b.name) = LOWER(?)"
		args = append(args, brandName)
	}
	if modelName, ok := filters["model"]; ok {
		query += " AND LOWER(m.name) LIKE LOWER(?)"
		args = append(args, "%"+modelName.(string)+"%")
	}
	if fuelType, ok := filters["fuel_type"]; ok {
		query += " AND LOWER(t.fuel_type) = LOWER(?)"
		args = append(args, fuelType)
	}
	if transmission, ok := filters["transmission"]; ok {
		query += " AND LOWER(t.transmission_type) = LOWER(?)"
		args = append(args, transmission)
	}
	if year, ok := filters["year"]; ok {
		query += " AND t.year = ?"
		args = append(args, year)
	}

	query += " ORDER BY b.name, m.name, t.year DESC, t.name"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search trims: %w", err)
	}
	defer rows.Close()

	var trims []*models.Trim
	for rows.Next() {
		trim := &models.Trim{}
		model := &models.Model{}
		brand := &models.Brand{}

		err := rows.Scan(
			&trim.ID, &trim.GenerationID, &trim.Name, &trim.Year, &trim.StartYear, &trim.EndYear, &trim.Generation, &trim.IsFacelift, &trim.Market,
			&trim.EngineType, &trim.FuelType, &trim.DisplacementCC, &trim.Cylinders, &trim.CylinderLayout,
			&trim.PowerHP, &trim.PowerKW, &trim.TorqueNM, &trim.EngineCode,
			&trim.Acceleration0To100, &trim.TopSpeedKmh,
			&trim.FuelConsumptionCity, &trim.FuelConsumptionHwy, &trim.FuelConsumptionComb,
			&trim.CO2Emissions, &trim.EmissionStandard,
			&trim.TransmissionType, &trim.Gears, &trim.Drivetrain,
			&trim.LengthMM, &trim.WidthMM, &trim.HeightMM, &trim.WheelbaseMM, &trim.GroundClearanceMM,
			&trim.CurbWeightKG, &trim.GrossWeightKG,
			&trim.LuggageCapacityL, &trim.LuggageCapacityMaxL, &trim.FuelTankCapacityL,
			&trim.TireSizeFront, &trim.TireSizeRear, &trim.WheelSizeInches,
			&trim.SeatingCapacity, &trim.Doors, &trim.ImageURL, &trim.MSRPPrice, &trim.Currency,
			&trim.CreatedAt, &trim.UpdatedAt,
			&model.ID, &model.BrandID, &model.Name, &model.BodyStyle, &model.Segment,
			&model.CreatedAt, &model.UpdatedAt,
			&brand.ID, &brand.Name, &brand.Country, &brand.LogoURL,
			&brand.CreatedAt, &brand.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trim: %w", err)
		}

		// Populate relationships
		model.Brand = brand
		trim.Model = model

		trims = append(trims, trim)
	}

	return trims, nil
}

// ListByModel retrieves all trims for a model (via generations)
func (r *TrimRepository) ListByModel(modelID int64) ([]*models.Trim, error) {
	query := `
		SELECT 
			t.id, t.generation_id, t.model_id, t.name, t.year, t.start_year, t.end_year, t.generation, t.is_facelift, t.market,
			t.engine_type, t.fuel_type, t.displacement_cc, t.cylinders, t.cylinder_layout,
			t.power_hp, t.power_kw, t.torque_nm, t.engine_code,
			t.acceleration_0_100, t.top_speed_kmh,
			t.fuel_consumption_city, t.fuel_consumption_highway, t.fuel_consumption_combined,
			t.co2_emissions, t.emission_standard,
			t.transmission_type, t.gears, t.drivetrain,
			t.length_mm, t.width_mm, t.height_mm, t.wheelbase_mm, t.ground_clearance_mm,
			t.curb_weight_kg, t.gross_weight_kg,
			t.luggage_capacity_l, t.luggage_capacity_max_l, t.fuel_tank_capacity_l,
			t.tire_size_front, t.tire_size_rear, t.wheel_size_inches,
			t.seating_capacity, t.doors, t.image_url, t.msrp_price, t.currency,
			t.created_at, t.updated_at
		FROM trims t
		LEFT JOIN generations g ON t.generation_id = g.id
		WHERE g.model_id = ?
		ORDER BY t.year DESC, t.name
	`
	rows, err := r.db.Query(query, modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to list trims: %w", err)
	}
	defer rows.Close()

	var trims []*models.Trim
	for rows.Next() {
		trim := &models.Trim{}
		err := rows.Scan(
			&trim.ID, &trim.GenerationID, &trim.ModelID, &trim.Name, &trim.Year, &trim.StartYear, &trim.EndYear, &trim.Generation, &trim.IsFacelift, &trim.Market,
			&trim.EngineType, &trim.FuelType, &trim.DisplacementCC, &trim.Cylinders, &trim.CylinderLayout,
			&trim.PowerHP, &trim.PowerKW, &trim.TorqueNM, &trim.EngineCode,
			&trim.Acceleration0To100, &trim.TopSpeedKmh,
			&trim.FuelConsumptionCity, &trim.FuelConsumptionHwy, &trim.FuelConsumptionComb,
			&trim.CO2Emissions, &trim.EmissionStandard,
			&trim.TransmissionType, &trim.Gears, &trim.Drivetrain,
			&trim.LengthMM, &trim.WidthMM, &trim.HeightMM, &trim.WheelbaseMM, &trim.GroundClearanceMM,
			&trim.CurbWeightKG, &trim.GrossWeightKG,
			&trim.LuggageCapacityL, &trim.LuggageCapacityMaxL, &trim.FuelTankCapacityL,
			&trim.TireSizeFront, &trim.TireSizeRear, &trim.WheelSizeInches,
			&trim.SeatingCapacity, &trim.Doors, &trim.ImageURL, &trim.MSRPPrice, &trim.Currency,
			&trim.CreatedAt, &trim.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trim: %w", err)
		}
		trims = append(trims, trim)
	}

	return trims, nil
}

// ListByGeneration retrieves all trims for a generation
func (r *TrimRepository) ListByGeneration(genID int64) ([]*models.Trim, error) {
	query := `
		SELECT 
			id, generation_id, model_id, name, year, start_year, end_year, generation, is_facelift, market,
			engine_type, fuel_type, displacement_cc, cylinders, cylinder_layout,
			power_hp, power_kw, torque_nm, engine_code,
			acceleration_0_100, top_speed_kmh,
			fuel_consumption_city, fuel_consumption_highway, fuel_consumption_combined,
			co2_emissions, emission_standard,
			transmission_type, transmission_code, gears, drivetrain,
			length_mm, width_mm, height_mm, wheelbase_mm, ground_clearance_mm,
			curb_weight_kg, gross_weight_kg,
			luggage_capacity_l, luggage_capacity_max_l, fuel_tank_capacity_l,
			tire_size_front, tire_size_rear, wheel_size_inches,
			seating_capacity, doors, image_url, msrp_price, currency,
			created_at, updated_at
		FROM trims
		WHERE generation_id = ?
		ORDER BY year DESC, name
	`
	rows, err := r.db.Query(query, genID)
	if err != nil {
		return nil, fmt.Errorf("failed to list trims: %w", err)
	}
	defer rows.Close()

	var trims []*models.Trim
	for rows.Next() {
		trim := &models.Trim{}
		err := rows.Scan(
			&trim.ID, &trim.GenerationID, &trim.ModelID, &trim.Name, &trim.Year, &trim.StartYear, &trim.EndYear, &trim.Generation, &trim.IsFacelift, &trim.Market,
			&trim.EngineType, &trim.FuelType, &trim.DisplacementCC, &trim.Cylinders, &trim.CylinderLayout,
			&trim.PowerHP, &trim.PowerKW, &trim.TorqueNM, &trim.EngineCode,
			&trim.Acceleration0To100, &trim.TopSpeedKmh,
			&trim.FuelConsumptionCity, &trim.FuelConsumptionHwy, &trim.FuelConsumptionComb,
			&trim.CO2Emissions, &trim.EmissionStandard,
			&trim.TransmissionType, &trim.TransmissionCode, &trim.Gears, &trim.Drivetrain,
			&trim.LengthMM, &trim.WidthMM, &trim.HeightMM, &trim.WheelbaseMM, &trim.GroundClearanceMM,
			&trim.CurbWeightKG, &trim.GrossWeightKG,
			&trim.LuggageCapacityL, &trim.LuggageCapacityMaxL, &trim.FuelTankCapacityL,
			&trim.TireSizeFront, &trim.TireSizeRear, &trim.WheelSizeInches,
			&trim.SeatingCapacity, &trim.Doors, &trim.ImageURL, &trim.MSRPPrice, &trim.Currency,
			&trim.CreatedAt, &trim.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trim: %w", err)
		}
		trims = append(trims, trim)
	}

	return trims, nil
}

// Delete deletes a trim
func (r *TrimRepository) Delete(id int64) error {
	query := `DELETE FROM trims WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trim: %w", err)
	}

	return nil
}
