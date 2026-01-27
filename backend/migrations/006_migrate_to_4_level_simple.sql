-- ============================================
-- SIMPLIFIED MIGRATION: 3-Level to 4-Level
-- Compatibility: SQLite with existing schema
-- ============================================

-- STEP 1: Create Generations Table
CREATE TABLE IF NOT EXISTS generations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    code TEXT NOT NULL,
    name TEXT,
    start_year INTEGER NOT NULL,
    end_year INTEGER,
    image_url TEXT,
    description TEXT,
    is_current BOOLEAN DEFAULT 0,
    platform TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_generations_model ON generations(model_id);
CREATE INDEX IF NOT EXISTS idx_generations_code ON generations(code);

-- STEP 2: Populate Generations from Trims
-- Extract unique combinations of model_id, generation code, and years
INSERT INTO generations (model_id, code, name, start_year, end_year, is_current)
SELECT DISTINCT
    t.model_id,
    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(MIN(t.year) AS TEXT)) as code,
    CASE 
        WHEN t.generation IS NOT NULL AND t.generation != '' 
        THEN t.generation || ' (' || CAST(MIN(t.year) AS TEXT) || 
             CASE WHEN MAX(t.year) != MIN(t.year) THEN '-' || CAST(MAX(t.year) AS TEXT) ELSE '' END || ')'
        ELSE 'Generation ' || CAST(MIN(t.year) AS TEXT)
    END as name,
    MIN(t.year) as start_year,
    CASE WHEN MAX(t.year) >= 2024 THEN NULL ELSE MAX(t.year) END as end_year,
    CASE WHEN MAX(t.year) >= 2024 THEN 1 ELSE 0 END as is_current
FROM trims t
WHERE t.model_id IS NOT NULL
GROUP BY 
    t.model_id,
    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(MIN(t.year) AS TEXT))
ORDER BY t.model_id, MIN(t.year);

-- STEP 3: Add generation_id column to trims
ALTER TABLE trims ADD COLUMN generation_id INTEGER;

-- STEP 4: Update trims to reference generations
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

-- STEP 5: Add foreign key index
CREATE INDEX IF NOT EXISTS idx_trims_generation ON trims(generation_id);

-- VERIFICATION: Check if all trims have generation_id
SELECT 
    'Migration Status' as status,
    CASE 
        WHEN COUNT(CASE WHEN generation_id IS NULL THEN 1 END) = 0 
        THEN '✓ Success - All trims linked'
        ELSE '⚠ Warning - ' || COUNT(CASE WHEN generation_id IS NULL THEN 1 END) || ' orphaned trims'
    END as result
FROM trims;
