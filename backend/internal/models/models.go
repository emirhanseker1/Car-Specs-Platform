package models

import "time"

// Brand represents a car manufacturer
type Brand struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Country   *string   `db:"country" json:"country,omitempty"`
	LogoURL   *string   `db:"logo_url" json:"logo_url,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Model represents a car model line
type Model struct {
	ID        int64     `db:"id" json:"id"`
	BrandID   int64     `db:"brand_id" json:"brand_id"`
	Name      string    `db:"name" json:"name"`
	BodyStyle *string   `db:"body_style" json:"body_style,omitempty"`
	Segment   *string   `db:"segment" json:"segment,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relationships (populated via joins)
	Brand *Brand `db:"-" json:"brand,omitempty"`
}

// Trim represents a specific car configuration with all technical specs
type Trim struct {
	ID           int64 `db:"id" json:"id"`
	ModelID      int64 `db:"model_id" json:"model_id"`           // Legacy column, kept for compatibility
	GenerationID int64 `db:"generation_id" json:"generation_id"` // New 4-level structure

	// Identification
	Name       string  `db:"name" json:"name"`
	Year       int     `db:"year" json:"year"`
	Generation *string `db:"generation" json:"generation,omitempty"` // Kept for backwards compatibility
	IsFacelift bool    `db:"is_facelift" json:"is_facelift"`
	Market     string  `db:"market" json:"market"`

	// Engine Specifications
	EngineType     *string `db:"engine_type" json:"engine_type,omitempty"`
	FuelType       *string `db:"fuel_type" json:"fuel_type,omitempty"`
	DisplacementCC *int    `db:"displacement_cc" json:"displacement_cc,omitempty"`
	Cylinders      *int    `db:"cylinders" json:"cylinders,omitempty"`
	CylinderLayout *string `db:"cylinder_layout" json:"cylinder_layout,omitempty"`
	PowerHP        *int    `db:"power_hp" json:"power_hp,omitempty"`
	PowerKW        *int    `db:"power_kw" json:"power_kw,omitempty"`
	TorqueNM       *int    `db:"torque_nm" json:"torque_nm,omitempty"`
	EngineCode     *string `db:"engine_code" json:"engine_code,omitempty"`

	// Performance
	Acceleration0To100  *float64 `db:"acceleration_0_100" json:"acceleration_0_100,omitempty"`
	TopSpeedKmh         *int     `db:"top_speed_kmh" json:"top_speed_kmh,omitempty"`
	FuelConsumptionCity *float64 `db:"fuel_consumption_city" json:"fuel_consumption_city,omitempty"`
	FuelConsumptionHwy  *float64 `db:"fuel_consumption_highway" json:"fuel_consumption_highway,omitempty"`
	FuelConsumptionComb *float64 `db:"fuel_consumption_combined" json:"fuel_consumption_combined,omitempty"`
	CO2Emissions        *int     `db:"co2_emissions" json:"co2_emissions,omitempty"`
	EmissionStandard    *string  `db:"emission_standard" json:"emission_standard,omitempty"`

	// Transmission & Drivetrain
	TransmissionType *string `db:"transmission_type" json:"transmission_type,omitempty"`
	Gears            *int    `db:"gears" json:"gears,omitempty"`
	Drivetrain       *string `db:"drivetrain" json:"drivetrain,omitempty"`

	// Dimensions & Weight
	LengthMM            *int `db:"length_mm" json:"length_mm,omitempty"`
	WidthMM             *int `db:"width_mm" json:"width_mm,omitempty"`
	HeightMM            *int `db:"height_mm" json:"height_mm,omitempty"`
	WheelbaseMM         *int `db:"wheelbase_mm" json:"wheelbase_mm,omitempty"`
	GroundClearanceMM   *int `db:"ground_clearance_mm" json:"ground_clearance_mm,omitempty"`
	CurbWeightKG        *int `db:"curb_weight_kg" json:"curb_weight_kg,omitempty"`
	GrossWeightKG       *int `db:"gross_weight_kg" json:"gross_weight_kg,omitempty"`
	LuggageCapacityL    *int `db:"luggage_capacity_l" json:"luggage_capacity_l,omitempty"`
	LuggageCapacityMaxL *int `db:"luggage_capacity_max_l" json:"luggage_capacity_max_l,omitempty"`
	FuelTankCapacityL   *int `db:"fuel_tank_capacity_l" json:"fuel_tank_capacity_l,omitempty"`

	// Wheels & Tires
	TireSizeFront   *string  `db:"tire_size_front" json:"tire_size_front,omitempty"`
	TireSizeRear    *string  `db:"tire_size_rear" json:"tire_size_rear,omitempty"`
	WheelSizeInches *float64 `db:"wheel_size_inches" json:"wheel_size_inches,omitempty"`

	// Additional
	SeatingCapacity int      `db:"seating_capacity" json:"seating_capacity"`
	Doors           *int     `db:"doors" json:"doors,omitempty"`
	ImageURL        *string  `db:"image_url" json:"image_url,omitempty"`
	MSRPPrice       *float64 `db:"msrp_price" json:"msrp_price,omitempty"`
	Currency        string   `db:"currency" json:"currency"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Relationships
	GenerationObj *Generation `db:"-" json:"generation_obj,omitempty"` // Changed from Model
	Model         *Model      `db:"-" json:"model,omitempty"`          // Kept via Generation
}

// Feature represents a car feature (safety, comfort, tech)
type Feature struct {
	ID       int64   `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Category *string `db:"category" json:"category,omitempty"`
}

// TrimFeature represents the many-to-many relationship
type TrimFeature struct {
	TrimID     int64 `db:"trim_id" json:"trim_id"`
	FeatureID  int64 `db:"feature_id" json:"feature_id"`
	IsStandard bool  `db:"is_standard" json:"is_standard"`
}
