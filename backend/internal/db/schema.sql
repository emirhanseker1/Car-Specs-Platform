CREATE TABLE IF NOT EXISTS vehicles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand TEXT DEFAULT 'Fiat',
    model TEXT NOT NULL,
    generation TEXT NOT NULL,
    image_url TEXT,
    link TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicle_generation_meta (
    vehicle_id INTEGER PRIMARY KEY,
    start_year INTEGER,
    end_year INTEGER,
    is_facelift INTEGER DEFAULT 0,
    market TEXT,
    FOREIGN KEY(vehicle_id) REFERENCES vehicles(id)
);

CREATE TABLE IF NOT EXISTS trims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vehicle_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    FOREIGN KEY(vehicle_id) REFERENCES vehicles(id)
);

CREATE TABLE IF NOT EXISTS trim_powertrain_meta (
    trim_id INTEGER PRIMARY KEY,
    engine_code TEXT,
    fuel_type TEXT,
    displacement_cc INTEGER,
    power_hp INTEGER,
    torque_nm INTEGER,
    transmission_type TEXT,
    gears INTEGER,
    drive TEXT,
    market_scope TEXT,
    FOREIGN KEY(trim_id) REFERENCES trims(id)
);

CREATE TABLE IF NOT EXISTS specs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trim_id INTEGER NOT NULL,
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY(trim_id) REFERENCES trims(id)
);

CREATE TABLE IF NOT EXISTS source_documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    source_type TEXT NOT NULL,
    market_scope TEXT,
    retrieved_at TEXT
);

CREATE TABLE IF NOT EXISTS spec_sources (
    spec_id INTEGER NOT NULL,
    source_document_id INTEGER NOT NULL,
    page TEXT,
    note TEXT,
    PRIMARY KEY(spec_id, source_document_id, page),
    FOREIGN KEY(spec_id) REFERENCES specs(id),
    FOREIGN KEY(source_document_id) REFERENCES source_documents(id)
);
