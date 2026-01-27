package repository

import (
	"fmt" // Added for fmt.Errorf

	"github.com/emirh/car-specs/backend/internal/models"
)

// GetFeaturedTrims returns random trims with images for the homepage
func (r *TrimRepository) GetFeaturedTrims(limit int) ([]*models.Trim, error) {
	query := `
		SELECT 
			t.id, t.model_id, t.generation_id, t.name, t.year, t.generation, t.is_facelift, t.market,
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
		WHERE t.image_url IS NOT NULL AND t.image_url != ''
		ORDER BY RANDOM()
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured trims: %w", err)
	}
	defer rows.Close()

	var trims []*models.Trim
	for rows.Next() {
		trim := &models.Trim{}
		model := &models.Model{}
		brand := &models.Brand{}

		err := rows.Scan(
			&trim.ID, &trim.ModelID, &trim.GenerationID, &trim.Name, &trim.Year, &trim.Generation, &trim.IsFacelift, &trim.Market,
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
			return nil, fmt.Errorf("failed to scan featured trim: %w", err)
		}

		// Populate relationships
		model.Brand = brand
		trim.Model = model

		trims = append(trims, trim)
	}

	return trims, rows.Err()
}
