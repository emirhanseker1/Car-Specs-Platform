# 4-Level Hierarchical Database Schema Design

## Overview
This document outlines a strict 4-level hierarchy for the Car Specification Platform:
1. **Brands** (Manufacturers)
2. **Models** (Product lines)
3. **Generations** (Body types/Year ranges)
4. **Trims** (Specific engine configurations)

---

## SQL Schema

### 1. Brands Table

```sql
-- ============================================
-- LEVEL 1: BRANDS (Manufacturers)
-- ============================================
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    logo_url TEXT,
    country VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_brands_slug ON brands(slug);
CREATE INDEX idx_brands_name ON brands(name);
```

**Example Data:**
- Audi, BMW, Volkswagen, Mercedes-Benz, Toyota

---

### 2. Models Table

```sql
-- ============================================
-- LEVEL 2: MODELS (Product Lines)
-- ============================================
CREATE TABLE IF NOT EXISTS models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    body_type VARCHAR(50), -- SUV, Sedan, Hatchback, Coupe, etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    UNIQUE(brand_id, slug)
);

CREATE INDEX idx_models_brand ON models(brand_id);
CREATE INDEX idx_models_slug ON models(slug);
CREATE INDEX idx_models_name ON models(name);
```

**Example Data:**
- Audi → A3, A4, Q5, Q7
- VW → Golf, Passat, Tiguan

---

### 3. Generations Table

```sql
-- ============================================
-- LEVEL 3: GENERATIONS (Body Types/Years)
-- ============================================
CREATE TABLE IF NOT EXISTS generations (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL,
    code VARCHAR(50) NOT NULL, -- e.g., "8V", "B8", "Mk7"
    name VARCHAR(150), -- Optional friendly name
    start_year INTEGER NOT NULL,
    end_year INTEGER, -- NULL if current generation
    image_url TEXT,
    description TEXT,
    is_current BOOLEAN DEFAULT FALSE,
    platform VARCHAR(100), -- e.g., "MQB", "MLB Evo"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE,
    CHECK (end_year IS NULL OR end_year >= start_year)
);

CREATE INDEX idx_generations_model ON generations(model_id);
CREATE INDEX idx_generations_code ON generations(code);
CREATE INDEX idx_generations_years ON generations(start_year, end_year);
```

**Example Data:**
- A3 → 8V (2012-2020), 8Y (2020-Present)
- Passat → B6 (2005-2010), B7 (2010-2014), B8 (2014-2019)
- Golf → Mk7 (2012-2019), Mk8 (2019-Present)

---

### 4. Trims Table

```sql
-- ============================================
-- LEVEL 4: TRIMS (Engine & Specifications)
-- ============================================
CREATE TABLE IF NOT EXISTS trims (
    id SERIAL PRIMARY KEY,
    generation_id INTEGER NOT NULL,
    
    -- Identification
    trim_name VARCHAR(150) NOT NULL, -- e.g., "1.6 TDI Ambition", "2.0 TFSI Quattro Sport"
    trim_level VARCHAR(50), -- Ambition, Sport, S-Line, etc.
    market VARCHAR(10) DEFAULT 'TR', -- TR, EU, US, etc.
    
    -- Engine Specifications
    fuel_type VARCHAR(30) NOT NULL, -- Diesel, Petrol, Hybrid, Plugin Hybrid, Electric
    engine_displacement_cc INTEGER, -- NULL for EVs
    engine_code VARCHAR(50),
    cylinders INTEGER,
    cylinder_layout VARCHAR(20), -- Inline, V, Boxer, etc.
    
    -- Power & Performance
    horsepower_hp INTEGER NOT NULL,
    horsepower_kw INTEGER, -- Can be calculated: hp * 0.7457
    torque_nm INTEGER,
    
    -- Transmission & Drivetrain
    transmission_type VARCHAR(50), -- Manual, Automatic, DSG, CVT
    gears INTEGER,
    drivetrain VARCHAR(30), -- FWD, RWD, AWD, 4WD
    
    -- Performance Metrics
    acceleration_0_to_100_kmh DECIMAL(4,2), -- e.g., 8.50 seconds
    top_speed_kmh INTEGER,
    
    -- Fuel Efficiency
    fuel_consumption_city_l_100km DECIMAL(4,2),
    fuel_consumption_highway_l_100km DECIMAL(4,2),
    fuel_consumption_combined_l_100km DECIMAL(4,2),
    
    -- Emissions
    co2_emissions_g_km INTEGER,
    emission_standard VARCHAR(30), -- Euro 6d, Euro 5, etc.
    
    -- Dimensions & Weight
    curb_weight_kg INTEGER,
    gross_weight_kg INTEGER,
    
    -- Pricing
    msrp_price DECIMAL(12,2),
    currency VARCHAR(10) DEFAULT 'TRY',
    
    -- Additional Info
    battery_capacity_kwh DECIMAL(5,2), -- For EVs and Hybrids
    electric_range_km INTEGER, -- For EVs and Hybrids
    
    -- Metadata
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (generation_id) REFERENCES generations(id) ON DELETE CASCADE
);

CREATE INDEX idx_trims_generation ON trims(generation_id);
CREATE INDEX idx_trims_fuel_type ON trims(fuel_type);
CREATE INDEX idx_trims_horsepower ON trims(horsepower_hp);
CREATE INDEX idx_trims_featured ON trims(is_featured);
CREATE INDEX idx_trims_market ON trims(market);
```

**Example Data:**
- A3 8V (2012-2020) → "1.6 TDI Ambition 110hp", "2.0 TDI S-Line 150hp", "1.4 TFSI Sport 125hp"
- Golf Mk7 → "1.6 TDI Bluemotion 110hp", "2.0 GTI 230hp"

---

## Go Structs with GORM Tags

### 1. Brand Struct

```go
package models

import "time"

// Brand represents a vehicle manufacturer (Level 1)
type Brand struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
    Slug        string    `gorm:"type:varchar(100);not null;unique;index" json:"slug"`
    LogoURL     *string   `gorm:"type:text" json:"logo_url,omitempty"`
    Country     *string   `gorm:"type:varchar(100)" json:"country,omitempty"`
    Description *string   `gorm:"type:text" json:"description,omitempty"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // Relationships
    Models []Model `gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE" json:"models,omitempty"`
}

// TableName overrides the default table name
func (Brand) TableName() string {
    return "brands"
}
```

---

### 2. Model Struct

```go
package models

import "time"

// Model represents a product line under a brand (Level 2)
type Model struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    BrandID     uint      `gorm:"not null;index" json:"brand_id"`
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`
    Slug        string    `gorm:"type:varchar(100);not null;index" json:"slug"`
    Description *string   `gorm:"type:text" json:"description,omitempty"`
    BodyType    *string   `gorm:"type:varchar(50)" json:"body_type,omitempty"` // SUV, Sedan, etc.
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // Relationships
    Brand       Brand        `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
    Generations []Generation `gorm:"foreignKey:ModelID;constraint:OnDelete:CASCADE" json:"generations,omitempty"`
}

// TableName overrides the default table name
func (Model) TableName() string {
    return "models"
}
```

---

### 3. Generation Struct

```go
package models

import "time"

// Generation represents a specific body type/year range (Level 3)
type Generation struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    ModelID     uint      `gorm:"not null;index" json:"model_id"`
    Code        string    `gorm:"type:varchar(50);not null;index" json:"code"` // 8V, B8, Mk7
    Name        *string   `gorm:"type:varchar(150)" json:"name,omitempty"`
    StartYear   int       `gorm:"not null;index:idx_gen_years" json:"start_year"`
    EndYear     *int      `gorm:"index:idx_gen_years" json:"end_year,omitempty"` // NULL if current
    ImageURL    *string   `gorm:"type:text" json:"image_url,omitempty"`
    Description *string   `gorm:"type:text" json:"description,omitempty"`
    IsCurrent   bool      `gorm:"default:false" json:"is_current"`
    Platform    *string   `gorm:"type:varchar(100)" json:"platform,omitempty"` // MQB, MLB, etc.
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // Relationships
    Model Model `gorm:"foreignKey:ModelID" json:"model,omitempty"`
    Trims []Trim `gorm:"foreignKey:GenerationID;constraint:OnDelete:CASCADE" json:"trims,omitempty"`
}

// TableName overrides the default table name
func (Generation) TableName() string {
    return "generations"
}
```

---

### 4. Trim Struct

```go
package models

import "time"

// Trim represents a specific engine configuration (Level 4)
type Trim struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    GenerationID uint      `gorm:"not null;index" json:"generation_id"`
    
    // Identification
    TrimName  string  `gorm:"type:varchar(150);not null" json:"trim_name"`
    TrimLevel *string `gorm:"type:varchar(50)" json:"trim_level,omitempty"`
    Market    string  `gorm:"type:varchar(10);default:'TR';index" json:"market"`
    
    // Engine Specifications
    FuelType             string  `gorm:"type:varchar(30);not null;index" json:"fuel_type"`
    EngineDisplacementCC *int    `gorm:"" json:"engine_displacement_cc,omitempty"` // NULL for EVs
    EngineCode           *string `gorm:"type:varchar(50)" json:"engine_code,omitempty"`
    Cylinders            *int    `gorm:"" json:"cylinders,omitempty"`
    CylinderLayout       *string `gorm:"type:varchar(20)" json:"cylinder_layout,omitempty"`
    
    // Power & Performance
    HorsepowerHP int  `gorm:"not null;index" json:"horsepower_hp"`
    HorsepowerKW *int `gorm:"" json:"horsepower_kw,omitempty"`
    TorqueNM     *int `gorm:"" json:"torque_nm,omitempty"`
    
    // Transmission & Drivetrain
    TransmissionType *string `gorm:"type:varchar(50)" json:"transmission_type,omitempty"`
    Gears            *int    `gorm:"" json:"gears,omitempty"`
    Drivetrain       *string `gorm:"type:varchar(30)" json:"drivetrain,omitempty"`
    
    // Performance Metrics
    Acceleration0To100Kmh *float64 `gorm:"type:decimal(4,2)" json:"acceleration_0_to_100_kmh,omitempty"`
    TopSpeedKmh           *int     `gorm:"" json:"top_speed_kmh,omitempty"`
    
    // Fuel Efficiency
    FuelConsumptionCityL100km     *float64 `gorm:"type:decimal(4,2)" json:"fuel_consumption_city_l_100km,omitempty"`
    FuelConsumptionHighwayL100km  *float64 `gorm:"type:decimal(4,2)" json:"fuel_consumption_highway_l_100km,omitempty"`
    FuelConsumptionCombinedL100km *float64 `gorm:"type:decimal(4,2)" json:"fuel_consumption_combined_l_100km,omitempty"`
    
    // Emissions
    CO2EmissionsGKm   *int    `gorm:"" json:"co2_emissions_g_km,omitempty"`
    EmissionStandard  *string `gorm:"type:varchar(30)" json:"emission_standard,omitempty"`
    
    // Dimensions & Weight
    CurbWeightKg  *int `gorm:"" json:"curb_weight_kg,omitempty"`
    GrossWeightKg *int `gorm:"" json:"gross_weight_kg,omitempty"`
    
    // Pricing
    MSRPPrice *float64 `gorm:"type:decimal(12,2)" json:"msrp_price,omitempty"`
    Currency  string   `gorm:"type:varchar(10);default:'TRY'" json:"currency"`
    
    // Additional Info for EVs/Hybrids
    BatteryCapacityKwh *float64 `gorm:"type:decimal(5,2)" json:"battery_capacity_kwh,omitempty"`
    ElectricRangeKm    *int     `gorm:"" json:"electric_range_km,omitempty"`
    
    // Metadata
    IsFeatured bool      `gorm:"default:false;index" json:"is_featured"`
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    
    // Relationships
    Generation Generation `gorm:"foreignKey:GenerationID" json:"generation,omitempty"`
}

// TableName overrides the default table name
func (Trim) TableName() string {
    return "trims"
}
```

---

## Complete Migration File

Save as: `backend/migrations/001_create_4_level_schema.sql`

```sql
-- ============================================
-- 4-LEVEL HIERARCHICAL SCHEMA
-- Car Specification Platform
-- ============================================

-- Drop existing tables if needed (use with caution in production)
-- DROP TABLE IF EXISTS trims CASCADE;
-- DROP TABLE IF EXISTS generations CASCADE;
-- DROP TABLE IF EXISTS models CASCADE;
-- DROP TABLE IF EXISTS brands CASCADE;

-- LEVEL 1: BRANDS
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    logo_url TEXT,
    country VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_brands_slug ON brands(slug);
CREATE INDEX idx_brands_name ON brands(name);

-- LEVEL 2: MODELS
CREATE TABLE IF NOT EXISTS models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    body_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    UNIQUE(brand_id, slug)
);

CREATE INDEX idx_models_brand ON models(brand_id);
CREATE INDEX idx_models_slug ON models(slug);
CREATE INDEX idx_models_name ON models(name);

-- LEVEL 3: GENERATIONS
CREATE TABLE IF NOT EXISTS generations (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150),
    start_year INTEGER NOT NULL,
    end_year INTEGER,
    image_url TEXT,
    description TEXT,
    is_current BOOLEAN DEFAULT FALSE,
    platform VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE,
    CHECK (end_year IS NULL OR end_year >= start_year)
);

CREATE INDEX idx_generations_model ON generations(model_id);
CREATE INDEX idx_generations_code ON generations(code);
CREATE INDEX idx_generations_years ON generations(start_year, end_year);

-- LEVEL 4: TRIMS
CREATE TABLE IF NOT EXISTS trims (
    id SERIAL PRIMARY KEY,
    generation_id INTEGER NOT NULL,
    
    -- Identification
    trim_name VARCHAR(150) NOT NULL,
    trim_level VARCHAR(50),
    market VARCHAR(10) DEFAULT 'TR',
    
    -- Engine Specifications
    fuel_type VARCHAR(30) NOT NULL,
    engine_displacement_cc INTEGER,
    engine_code VARCHAR(50),
    cylinders INTEGER,
    cylinder_layout VARCHAR(20),
    
    -- Power & Performance
    horsepower_hp INTEGER NOT NULL,
    horsepower_kw INTEGER,
    torque_nm INTEGER,
    
    -- Transmission & Drivetrain
    transmission_type VARCHAR(50),
    gears INTEGER,
    drivetrain VARCHAR(30),
    
    -- Performance Metrics
    acceleration_0_to_100_kmh DECIMAL(4,2),
    top_speed_kmh INTEGER,
    
    -- Fuel Efficiency
    fuel_consumption_city_l_100km DECIMAL(4,2),
    fuel_consumption_highway_l_100km DECIMAL(4,2),
    fuel_consumption_combined_l_100km DECIMAL(4,2),
    
    -- Emissions
    co2_emissions_g_km INTEGER,
    emission_standard VARCHAR(30),
    
    -- Dimensions & Weight
    curb_weight_kg INTEGER,
    gross_weight_kg INTEGER,
    
    -- Pricing
    msrp_price DECIMAL(12,2),
    currency VARCHAR(10) DEFAULT 'TRY',
    
    -- Additional Info
    battery_capacity_kwh DECIMAL(5,2),
    electric_range_km INTEGER,
    
    -- Metadata
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (generation_id) REFERENCES generations(id) ON DELETE CASCADE
);

CREATE INDEX idx_trims_generation ON trims(generation_id);
CREATE INDEX idx_trims_fuel_type ON trims(fuel_type);
CREATE INDEX idx_trims_horsepower ON trims(horsepower_hp);
CREATE INDEX idx_trims_featured ON trims(is_featured);
CREATE INDEX idx_trims_market ON trims(market);
```

---

## Example Data Hierarchy

```
Audi (Brand)
├── A3 (Model)
│   ├── 8V (2012-2020) (Generation)
│   │   ├── 1.6 TDI Ambition 110hp (Trim)
│   │   ├── 2.0 TDI S-Line 150hp (Trim)
│   │   └── 1.4 TFSI Sport 125hp (Trim)
│   └── 8Y (2020-Present) (Generation)
│       ├── 35 TFSI 150hp (Trim)
│       └── 40 TFSI Quattro 190hp (Trim)
├── A4 (Model)
│   └── B9 (2015-2019) (Generation)
│       ├── 2.0 TDI Ultra 150hp (Trim)
│       └── 3.0 TDI Quattro 272hp (Trim)

Volkswagen (Brand)
└── Golf (Model)
    ├── Mk7 (2012-2019) (Generation)
    │   ├── 1.6 TDI Bluemotion 110hp (Trim)
    │   └── 2.0 GTI 230hp (Trim)
    └── Mk8 (2019-Present) (Generation)
        └── 2.0 TDI 150hp (Trim)
```

---

## Naming Conventions Summary

| Element | SQL (snake_case) | Go (PascalCase) | JSON (snake_case) |
|---------|------------------|-----------------|-------------------|
| Primary Key | `id` | `ID` | `id` |
| Foreign Key | `brand_id` | `BrandID` | `brand_id` |
| Column | `fuel_type` | `FuelType` | `fuel_type` |
| Table | `generations` | `Generation` | N/A |

---

## Key Design Decisions

1. **PostgreSQL vs SQLite**: Schema uses `SERIAL` for auto-increment (PostgreSQL). For SQLite, change to `INTEGER PRIMARY KEY AUTOINCREMENT`.

2. **Nullable End Year**: `end_year` can be NULL for current/ongoing generations.

3. **Cascading Deletes**: All foreign keys use `ON DELETE CASCADE` to maintain referential integrity.

4. **Indexing**: Indexes added on frequently queried columns (slugs, foreign keys, years, fuel type, horsepower).

5. **Decimal Precision**: Performance metrics use `DECIMAL(4,2)` for values like `8.50` seconds.

6. **EV Support**: `battery_capacity_kwh` and `electric_range_km` fields support electric vehicles.

7. **Pointer Fields in Go**: Optional fields use pointers (`*string`, `*int`) to distinguish between zero values and NULL.

---

## Migration from Current 3-Level to 4-Level

If you have existing data in a 3-level structure (Brand → Model → Trim), you'll need a data migration:

1. Create `generations` table
2. For each existing `trim`, create a corresponding `generation` record
3. Update `trim.generation_id` to reference the new generation
4. Remove `trim.model_id` column

Would you like me to create the migration script as well?
