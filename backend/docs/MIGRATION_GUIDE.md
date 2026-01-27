# Migration Guide: 3-Level to 4-Level Database Structure

## Overview
This guide explains how to migrate your existing database from the current 3-level structure to a new 4-level hierarchy by introducing the **Generations** table.

### Current Structure (3-Level)
```
Brands → Models → Trims
```

### New Structure (4-Level)
```
Brands → Models → Generations → Trims
```

---

## What Changed?

### Before
- **Brands**: Manufacturers (Audi, BMW, VW)
- **Models**: Product lines (A3, 3 Series, Golf)
- **Trims**: Specific configurations with year and optional generation field

### After
- **Brands**: Manufacturers (unchanged)
- **Models**: Product lines (unchanged)
- **Generations**: 🆕 **NEW** - Specific body types/year ranges (F30, G20, 8V, Mk7)
- **Trims**: Engine configurations (now linked to Generation instead of Model)

---

## Why This Change?

1. **Better Organization**: Groups trims by generation/body type rather than just model
2. **Accurate Year Ranges**: Each generation has clear start/end years
3. **Platform Tracking**: Can track which platform each generation uses (MQB, CLAR, etc.)
4. **Scalability**: Easier to add new generations without cluttering the trims table

---

## Migration Steps

### 📋 Prerequisites

1. **Backup your database**:
   ```bash
   cp backend/vehicles.db backend/vehicles.db.backup
   ```

2. **Check current data**:
   ```bash
   python -c "import sqlite3; conn = sqlite3.connect('backend/vehicles.db'); cursor = conn.cursor(); cursor.execute('SELECT COUNT(*) FROM trims'); print(f'Total trims: {cursor.fetchone()[0]}')"
   ```

### 🚀 Step 1: Run Migration Script

```bash
# Navigate to backend directory
cd backend

# Run the migration SQL script
sqlite3 vehicles.db < migrations/006_migrate_to_4_level.sql
```

### ✅ Step 2: Verify Migration

After running the migration, verify the results:

```sql
-- Check record counts at each level
SELECT 'Brands' as level, COUNT(*) as count FROM brands
UNION ALL
SELECT 'Models', COUNT(*) FROM models
UNION ALL
SELECT 'Generations', COUNT(*) FROM generations
UNION ALL
SELECT 'Trims', COUNT(*) FROM trims;
```

Expected output:
```
level        | count
-------------+-------
Brands       | 3-5
Models       | 10-20
Generations  | 15-30
Trims        | 60
```

### 🔍 Step 3: Inspect Generated Data

View the new hierarchy:

```sql
SELECT 
    b.name as brand,
    m.name as model,
    g.code as gen_code,
    g.name as generation,
    g.start_year || '-' || COALESCE(CAST(g.end_year AS TEXT), 'Present') as years,
    COUNT(t.id) as trims
FROM brands b
JOIN models m ON m.brand_id = b.id
JOIN generations g ON g.model_id = m.id
LEFT JOIN trims t ON t.generation_id = g.id
GROUP BY b.id, m.id, g.id
ORDER BY b.name, m.name, g.start_year;
```

---

## How the Migration Works

### 1. Generation Extraction Logic

The migration script automatically creates generation records based on existing trim data:

```sql
-- Groups trims by model_id and generation code
-- If generation is empty, uses year as fallback (e.g., "Y2024")
INSERT INTO generations (model_id, code, name, start_year, end_year)
SELECT
    t.model_id,
    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT)) as code,
    -- Creates name like "F30 (2012-2018)" or "Generation 2024"
    CASE 
        WHEN t.generation IS NOT NULL THEN t.generation || ' (' || ... ')'
        ELSE 'Generation ' || CAST(t.year AS TEXT)
    END as name,
    MIN(t.year) as start_year,
    CASE WHEN MAX(t.year) >= 2024 THEN NULL ELSE MAX(t.year) END as end_year
FROM trims t
GROUP BY t.model_id, COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT));
```

### 2. Relationship Update

Each trim is updated to reference its appropriate generation:

```sql
UPDATE trims
SET generation_id = (
    SELECT g.id FROM generations g
    WHERE g.model_id = trims.model_id
      AND g.code = COALESCE(NULLIF(trims.generation, ''), 'Y' || CAST(trims.year AS TEXT))
      AND trims.year BETWEEN g.start_year AND COALESCE(g.end_year, 9999)
);
```

### 3. Table Recreation

The trims table is recreated with `generation_id` as a foreign key instead of `model_id`:

```sql
-- Old structure
model_id → models(id)

-- New structure  
generation_id → generations(id) → model_id → models(id)
```

---

## Data Validation

### Check for Orphaned Trims

Run this query to ensure all trims are properly linked:

```sql
SELECT COUNT(*) as orphaned_trims
FROM trims
WHERE generation_id IS NULL;
```

**Expected**: 0 orphaned trims

### Check Generation Distribution

```sql
SELECT 
    g.code,
    g.name,
    COUNT(t.id) as trim_count
FROM generations g
LEFT JOIN trims t ON t.generation_id = g.id
GROUP BY g.id
ORDER BY trim_count DESC;
```

---

## Post-Migration Tasks

### 1. Update Generation Data (Optional but Recommended)

After migration, you can manually improve generation data:

```sql
-- Example: Update BMW F30 generation with proper details
UPDATE generations
SET 
    name = 'F30 (2012-2019)',
    platform = 'F Platform',
    description = 'Sixth generation BMW 3 Series',
    image_url = 'https://example.com/bmw-f30.jpg'
WHERE code = 'F30';

-- Example: Update Golf Mk7 generation
UPDATE generations
SET 
    name = 'Golf Mk7 (2012-2019)',
    platform = 'MQB',
    description = 'Seventh generation Volkswagen Golf',
    is_current = 0
WHERE code = 'Mk7';
```

### 2. Add Images for Generations

```sql
-- Add generation images (one image per generation)
UPDATE generations
SET image_url = 'url_to_generation_image'
WHERE id = generation_id;
```

### 3. Mark Current Generations

```sql
-- Update current generation flags
UPDATE generations
SET is_current = 1
WHERE end_year IS NULL OR end_year >= 2024;
```

---

## Update Go Code

### 1. Add Generation Model

Create `backend/internal/models/generation.go`:

```go
package models

import "time"

type Generation struct {
    ID       int64  `db:"id" json:"id"`
    ModelID  int64  `db:"model_id" json:"model_id"`
    Code      string   `db:"code" json:"code"`
    Name      *string  `db:"name" json:"name,omitempty"`
    StartYear int      `db:"start_year" json:"start_year"`
    EndYear   *int     `db:"end_year" json:"end_year,omitempty"`
    ImageURL    *string `db:"image_url" json:"image_url,omitempty"`
    Description *string `db:"description" json:"description,omitempty"`
    IsCurrent   bool    `db:"is_current" json:"is_current"`
    Platform    *string `db:"platform" json:"platform,omitempty"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    
    // Relationships
    Model *Model `db:"-" json:"model,omitempty"`
    Trims []Trim `db:"-" json:"trims,omitempty"`
}
```

### 2. Update Trim Model

Modify `backend/internal/models/models.go`:

```go
type Trim struct {
    ID           int64 `db:"id" json:"id"`
    GenerationID int64 `db:"generation_id" json:"generation_id"` // Changed from ModelID
    
    // ... rest of fields remain the same ...
    
    // Relationships
    Generation *Generation `db:"-" json:"generation,omitempty"` // New relationship
}
```

### 3. Create Repository Methods

Create `backend/internal/repository/generation_repository.go`:

```go
package repository

import (
    "database/sql"
    "github.com/emirh/car-specs/backend/internal/models"
)

type GenerationRepository struct {
    db *sql.DB
}

func NewGenerationRepository(db *sql.DB) *GenerationRepository {
    return &GenerationRepository{db: db}
}

func (r *GenerationRepository) GetByID(id int64) (*models.Generation, error) {
    query := `SELECT * FROM generations WHERE id = ?`
    // Implementation...
}

func (r *GenerationRepository) GetByModel(modelID int64) ([]models.Generation, error) {
    query := `SELECT * FROM generations WHERE model_id = ? ORDER BY start_year DESC`
    // Implementation...
}

func (r *GenerationRepository) GetCurrent() ([]models.Generation, error) {
    query := `SELECT * FROM generations WHERE is_current = 1`
    // Implementation...
}
```

### 4. Update API Endpoints

Add new endpoints:

```go
// GET /api/models/{modelId}/generations
func (h *ModelHandler) HandleGetGenerations(w http.ResponseWriter, r *http.Request) {
    // Return all generations for a model
}

// GET /api/generations/{generationId}
func (h *GenerationHandler) HandleGetGeneration(w http.ResponseWriter, r *http.Request) {
    // Return specific generation with trims
}

// GET /api/generations/{generationId}/trims
func (h *GenerationHandler) HandleGetTrims(w http.ResponseWriter, r *http.Request) {
    // Return all trims for a generation
}
```

---

## API Response Changes

### Before Migration

```json
{
  "id": 1,
  "model_id": 5,
  "name": "320d",
  "year": 2018,
  "generation": "F30",
  "fuel_type": "Diesel",
  "power_hp": 190
}
```

### After Migration

```json
{
  "id": 1,
  "generation_id": 12,
  "name": "320d",
  "year": 2018,
  "fuel_type": "Diesel",
  "power_hp": 190,
  "generation": {
    "id": 12,
    "model_id": 5,
    "code": "F30",
    "name": "F30 (2012-2019)",
    "start_year": 2012,
    "end_year": 2019,
    "platform": "F Platform",
    "model": {
      "id": 5,
      "name": "3 Series",
      "brand": {
        "name": "BMW"
      }
    }
  }
}
```

---

## Rollback Instructions

If you need to rollback the migration:

```sql
-- Add model_id back to trims
ALTER TABLE trims ADD COLUMN model_id INTEGER;

-- Populate model_id from generations
UPDATE trims
SET model_id = (
    SELECT g.model_id
    FROM generations g
    WHERE g.id = trims.generation_id
);

-- Drop generation_id column
ALTER TABLE trims DROP COLUMN generation_id;

-- Drop generations table
DROP TABLE generations;

-- Restore from backup if needed
-- cp backend/vehicles.db.backup backend/vehicles.db
```

---

## Summary

✅ **Created**: `generations` table  
✅ **Updated**: `trims` table to reference `generations` instead of `models`  
✅ **Maintained**: All existing data and relationships  
✅ **Added**: Generation grouping, year ranges, and platform tracking

### Next Steps

1. Run migration on development database
2. Verify data integrity
3. Update Go models and repositories
4. Update API endpoints to use 4-level hierarchy
5. Test thoroughly
6. Update frontend to display generation information
7. Run migration on production (after testing!)

---

## Support

If you encounter issues:

1. Check the verification queries above
2. Review the orphaned trims query
3. Inspect the generated generations table
4. Compare with the migration script logic
5. Restore from backup if needed

The migration is designed to be safe and preserve all existing data. However, always test on a backup first!
