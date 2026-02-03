CREATE TABLE IF NOT EXISTS brands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    country TEXT,
    logo_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    body_style TEXT,
    segment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(brand_id) REFERENCES brands(id),
    UNIQUE(brand_id, name)
);

CREATE TABLE IF NOT EXISTS generations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    code TEXT,
    name TEXT,
    start_year INTEGER,
    end_year INTEGER,
    is_facelift BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(model_id) REFERENCES models(id)
);

CREATE TABLE IF NOT EXISTS trims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    generation_id INTEGER NOT NULL,
    model_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    year INTEGER,
    start_year INTEGER,
    end_year INTEGER,
    generation TEXT,
    engine_type TEXT,
    fuel_type TEXT,
    is_facelift BOOLEAN DEFAULT 0,
    market TEXT DEFAULT 'Global',
    
    -- Engine & Performance
    power_hp INTEGER,
    power_kw INTEGER,
    torque_nm INTEGER,
    displacement_cc INTEGER,
    cylinders INTEGER,
    cylinder_layout TEXT,
    engine_code TEXT,
    acceleration_0_100 REAL,
    top_speed_kmh INTEGER,
    
    -- Consumption & Emissions
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
    
    -- Additional
    seating_capacity INTEGER DEFAULT 5,
    doors INTEGER DEFAULT 5,
    image_url TEXT,
    msrp_price REAL,
    currency TEXT DEFAULT 'USD',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(generation_id) REFERENCES generations(id),
    FOREIGN KEY(model_id) REFERENCES models(id)
);

CREATE TABLE IF NOT EXISTS specs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trim_id INTEGER NOT NULL,
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY(trim_id) REFERENCES trims(id)
);
