-- ============================================
-- MODELS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    body_style TEXT,
    segment TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    UNIQUE(brand_id, name)
);

CREATE INDEX IF NOT EXISTS idx_models_brand ON models(brand_id);
CREATE INDEX IF NOT EXISTS idx_models_name ON models(name);
