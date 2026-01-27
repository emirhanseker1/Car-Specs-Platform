-- ============================================
-- MIGRATION: 3-Level to 4-Level Structure
-- Goal: Insert Generations table between Models and Trims
-- ============================================

-- This migration script converts the existing structure:
--   brands → models → trims
-- To the new 4-level structure:
--   brands → models → generations → trims

-- ============================================
-- STEP 1: Create Generations Table
-- ============================================

CREATE TABLE IF NOT EXISTS generations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    
    -- Generation Identification
    code TEXT NOT NULL,              -- e.g., "F30", "G20", "8V", "Mk7"
    name TEXT,                        -- Optional friendly name
    start_year INTEGER NOT NULL,
    end_year INTEGER,                 -- NULL if current generation
    
    -- Additional Info
    image_url TEXT,
    description TEXT,
    is_current BOOLEAN DEFAULT 0,
    platform TEXT,                    -- e.g., "MQB", "MLB Evo", "CLAR"
    
    -- Metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE,
    CHECK (end_year IS NULL OR end_year >= start_year)
);

CREATE INDEX IF NOT EXISTS idx_generations_model ON generations(model_id);
CREATE INDEX IF NOT EXISTS idx_generations_code ON generations(code);
CREATE INDEX IF NOT EXISTS idx_generations_years ON generations(start_year, end_year);

-- ============================================
-- STEP 2: Populate Generations from Existing Trims
-- ============================================

-- Strategy: Extract unique (model_id, generation, year) combinations
-- and create generation records

-- Insert generations based on existing trim data
-- Group by model_id, generation code, and year ranges

INSERT INTO generations (model_id, code, name, start_year, end_year, is_current, platform)
SELECT DISTINCT
    t.model_id,
    -- Use generation if exists, otherwise use year as fallback code
    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT)) as code,
    -- Create a friendly name
    CASE 
        WHEN t.generation IS NOT NULL AND t.generation != '' 
        THEN t.generation || ' (' || CAST(MIN(t.year) AS TEXT) || 
             CASE 
                WHEN MAX(t.year) != MIN(t.year) 
                THEN '-' || CAST(MAX(t.year) AS TEXT) 
                ELSE '' 
             END || ')'
        ELSE 'Generation ' || CAST(t.year AS TEXT)
    END as name,
    -- Start year is the earliest year for this generation
    MIN(t.year) as start_year,
    -- End year is the latest year (or NULL if still current)
    CASE 
        WHEN MAX(t.year) >= 2024 THEN NULL  -- Current generation
        ELSE MAX(t.year)
    END as end_year,
    -- Mark as current if max year is 2024 or later
    CASE WHEN MAX(t.year) >= 2024 THEN 1 ELSE 0 END as is_current,
    -- Platform is NULL for now (can be populated manually later)
    NULL as platform
FROM trims t
WHERE t.model_id IS NOT NULL
GROUP BY 
    t.model_id,
    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT))
ORDER BY t.model_id, MIN(t.year);

-- ============================================
-- STEP 3: Add generation_id column to trims
-- ============================================

-- Add new column (initially nullable for migration)
ALTER TABLE trims ADD COLUMN generation_id INTEGER;

-- ============================================
-- STEP 4: Update trims to reference generations
-- ============================================

-- Update each trim to reference the appropriate generation
-- Match based on model_id, generation code, and year range

UPDATE trims
SET generation_id = (
    SELECT g.id
    FROM generations g
    WHERE g.model_id = trims.model_id
      AND g.code = COALESCE(NULLIF(trims.generation, ''), 'Y' || CAST(trims.year AS TEXT))
      AND trims.year >= g.start_year
      AND (g.end_year IS NULL OR trims.year <= g.end_year)
    LIMIT 1
);

-- ============================================
-- STEP 5: Verify data integrity
-- ============================================

-- Check if any trims don't have generation_id assigned
-- (This should return 0 rows if migration is successful)

SELECT COUNT(*) as orphaned_trims
FROM trims
WHERE generation_id IS NULL;

-- If orphaned_trims > 0, manual intervention needed

-- ============================================
-- STEP 6: Make generation_id NOT NULL and add foreign key
-- ============================================

-- Note: SQLite doesn't support adding NOT NULL constraint to existing column
-- We need to recreate the table or ensure all values are set

-- First verify all trims have generation_id
-- Then we can add the constraint

-- Create a backup of trims table
CREATE TABLE trims_backup AS SELECT * FROM trims;

-- Drop old indexes on trims
DROP INDEX IF EXISTS idx_trims_model;
DROP INDEX IF EXISTS idx_trims_year;
DROP INDEX IF EXISTS idx_trims_fuel;
DROP INDEX IF EXISTS idx_trims_transmission;
DROP INDEX IF EXISTS idx_trims_power;

-- Recreate trims table with new structure
DROP TABLE trims;

CREATE TABLE trims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    generation_id INTEGER NOT NULL,  -- Changed from model_id
    
    -- Identification
    name TEXT NOT NULL,
    year INTEGER NOT NULL,
    generation TEXT,  -- Keep for backwards compatibility, but redundant now
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
    
    FOREIGN KEY (generation_id) REFERENCES generations(id) ON DELETE CASCADE
);

-- Restore data from backup
INSERT INTO trims SELECT * FROM trims_backup;

-- Recreate indexes for trims (now referencing generation instead of model)
CREATE INDEX IF NOT EXISTS idx_trims_generation ON trims(generation_id);
CREATE INDEX IF NOT EXISTS idx_trims_year ON trims(year);
CREATE INDEX IF NOT EXISTS idx_trims_fuel ON trims(fuel_type);
CREATE INDEX IF NOT EXISTS idx_trims_transmission ON trims(transmission_type);
CREATE INDEX IF NOT EXISTS idx_trims_power ON trims(power_hp);

-- Drop backup table
DROP TABLE trims_backup;

-- ============================================
-- STEP 7: Verification Queries
-- ============================================

-- Count records at each level
SELECT 'Brands' as level, COUNT(*) as count FROM brands
UNION ALL
SELECT 'Models', COUNT(*) FROM models
UNION ALL
SELECT 'Generations', COUNT(*) FROM generations
UNION ALL
SELECT 'Trims', COUNT(*) FROM trims;

-- Show the new hierarchy for a sample
SELECT 
    b.name as brand,
    m.name as model,
    g.code as generation_code,
    g.name as generation_name,
    g.start_year || '-' || COALESCE(CAST(g.end_year AS TEXT), 'Present') as years,
    COUNT(t.id) as trim_count
FROM brands b
JOIN models m ON m.brand_id = b.id
JOIN generations g ON g.model_id = m.id
LEFT JOIN trims t ON t.generation_id = g.id
GROUP BY b.id, m.id, g.id
ORDER BY b.name, m.name, g.start_year;

-- ============================================
-- NOTES FOR MANUAL REVIEW
-- ============================================

-- 1. Some trims have empty generation field - these will get auto-generated
--    generation codes like "Y2024" (year-based)
--
-- 2. Current generations (2024+) are marked as ongoing (end_year = NULL)
--
-- 3. Platform information is not populated automatically - needs manual entry
--
-- 4. Generation images (image_url) should be added manually later
--
-- 5. After migration, you may want to:
--    - Update generation codes to proper format (F30, G20, Mk7, etc.)
--    - Add platform information
--    - Add generation descriptions
--    - Add generation images
--
-- 6. The old 'generation' column in trims is kept for backwards compatibility
--    but is now redundant (data is in generations table)

-- ============================================
-- ROLLBACK PROCEDURE (if needed)
-- ============================================

-- If you need to rollback this migration:
-- 1. Restore from backup
-- 2. Or manually:
--    ALTER TABLE trims ADD COLUMN model_id INTEGER;
--    UPDATE trims SET model_id = (SELECT model_id FROM generations WHERE id = trims.generation_id);
--    ALTER TABLE trims DROP COLUMN generation_id;
--    DROP TABLE generations;
