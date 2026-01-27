-- ============================================
-- TRIMS TABLE (Source of Truth for All Specs)
-- ============================================
CREATE TABLE IF NOT EXISTS trims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    
    -- Identification
    name TEXT NOT NULL,
    year INTEGER NOT NULL,
    generation TEXT,
    is_facelift BOOLEAN DEFAULT 0,
    market TEXT DEFAULT 'TR',
    
    -- Engine Specifications
    engine_type TEXT,
    fuel_type TEXT,
    displacement_cc INTEGER,
    cylinders INTEGER,
    cylinder_layout TEXT,
    power_hp INTEGER,
    power_kw INTEGER,
    torque_nm INTEGER,
    engine_code TEXT,
    
    -- Performance
    acceleration_0_100 REAL,
    top_speed_kmh INTEGER,
    fuel_consumption_city REAL,
    fuel_consumption_highway REAL,
    fuel_consumption_combined REAL,
    co2_emissions INTEGER,
    emission_standard TEXT,
    
    -- Transmission & Drivetrain
    transmission_type TEXT,
    gears INTEGER,
    drivetrain TEXT,
    
    -- Dimensions & Weight
    length_mm INTEGER,
    width_mm INTEGER,
    height_mm INTEGER,
    wheelbase_mm INTEGER,
    ground_clearance_mm INTEGER,
    curb_weight_kg INTEGER,
    gross_weight_kg INTEGER,
    luggage_capacity_l INTEGER,
    luggage_capacity_max_l INTEGER,
    fuel_tank_capacity_l INTEGER,
    
    -- Wheels & Tires
    tire_size_front TEXT,
    tire_size_rear TEXT,
    wheel_size_inches REAL,
    
    -- Additional Info
    seating_capacity INTEGER DEFAULT 5,
    doors INTEGER,
    image_url TEXT,
    msrp_price REAL,
    currency TEXT DEFAULT 'TRY',
    
    -- Metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_trims_model ON trims(model_id);
CREATE INDEX IF NOT EXISTS idx_trims_year ON trims(year);
CREATE INDEX IF NOT EXISTS idx_trims_fuel ON trims(fuel_type);
CREATE INDEX IF NOT EXISTS idx_trims_transmission ON trims(transmission_type);
CREATE INDEX IF NOT EXISTS idx_trims_power ON trims(power_hp);
