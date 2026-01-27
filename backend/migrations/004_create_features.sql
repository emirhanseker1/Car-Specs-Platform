-- ============================================
-- FEATURES TABLE (Optional Equipment/Features)
-- ============================================
CREATE TABLE IF NOT EXISTS features (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    category TEXT
);

CREATE TABLE IF NOT EXISTS trim_features (
    trim_id INTEGER NOT NULL,
    feature_id INTEGER NOT NULL,
    is_standard BOOLEAN DEFAULT 1,
    
    PRIMARY KEY (trim_id, feature_id),
    FOREIGN KEY (trim_id) REFERENCES trims(id) ON DELETE CASCADE,
    FOREIGN KEY (feature_id) REFERENCES features(id) ON DELETE CASCADE
);
