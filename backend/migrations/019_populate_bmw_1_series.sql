-- 019_populate_bmw_1_series.sql
BEGIN TRANSACTION;

-- 1. Ensure Brand Exists (BMW)
INSERT OR IGNORE INTO brands (name, country, logo_url) 
VALUES ('BMW', 'Germany', '/images/brands/bmw.png');

-- 2. Clean existing 1 Series Data (if any)
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id IN (
        SELECT id FROM models WHERE name = '1 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
    )
);
DELETE FROM generations WHERE model_id IN (
    SELECT id FROM models WHERE name = '1 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
);
DELETE FROM models WHERE name = '1 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW');

-- 3. Create Model
INSERT INTO models (brand_id, name, image_url, body_style)
VALUES (
    (SELECT id FROM brands WHERE name = 'BMW'),
    '1 Serisi',
    '/images/models/bmw-1-series.png', 
    'Hatchback'
);

-- 4. Create Generations
-- F40 (MK3)
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '1 Serisi'),
    '1 Serisi F40', 'F40', 2019, 2024, '/images/generations/bmw-1-f40.png', 0
);
-- F20 LCI (MK2 Makyajlı)
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '1 Serisi'),
    '1 Serisi F20 (LCI)', 'F20 LCI', 2015, 2019, '/images/generations/bmw-1-f20-lci.png', 1
);
-- F20 (MK2)
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '1 Serisi'),
    '1 Serisi F20', 'F20', 2011, 2015, '/images/generations/bmw-1-f20.png', 0
);
-- E87 LCI (MK1 Makyajlı)
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '1 Serisi'),
    '1 Serisi E87 (LCI)', 'E87 LCI', 2007, 2011, '/images/generations/bmw-1-e87-lci.png', 1
);
-- E87 (MK1)
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '1 Serisi'),
    '1 Serisi E87', 'E87', 2004, 2007, '/images/generations/bmw-1-e87.png', 0
);

-- 5. Insert Trims

-- === F40 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F40'),
    '118i M Sport', 2019, 140, 220, 'Otomatik (Çift Kavrama)', 'Getrag 7DCT300', '1.5 Turbo', 'Önden Çekiş', 'Benzin', 8.5, NULL
),
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F40'),
    '116d', 2019, 116, 270, 'Otomatik (Çift Kavrama)', 'Getrag 7DCT300', '1.5 Turbo Dizel', 'Önden Çekiş', 'Dizel', NULL, NULL
);

-- === F20 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F20 (LCI)'),
    '118i', 2015, 136, 220, 'Otomatik (ZF)', 'ZF 8HP', 'B38', 'Arkadan İtiş', 'Benzin', 8.7, NULL
),
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F20 (LCI)'),
    '116d', 2015, 116, 270, 'Otomatik (ZF)', 'ZF 8HP', 'B37', 'Arkadan İtiş', 'Dizel', NULL, NULL
);

-- === F20 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F20'),
    '116i', 2011, 136, 220, 'Otomatik (ZF)', 'ZF 8HP', 'N13', 'Arkadan İtiş', 'Benzin', 8.7, NULL
),
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi F20'),
    '118i', 2011, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', 'N13', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === E87 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi E87 (LCI)'),
    '116i', 2007, 115, 150, 'Otomatik', 'GM 6L45', '1.6 Atmosferik', 'Arkadan İtiş', 'Benzin', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi E87 (LCI)'),
    '118i', 2007, 143, 190, 'Otomatik', 'ZF 6HP / GM', '2.0 Atmosferik', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === E87 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '1 Serisi'),
    (SELECT id FROM generations WHERE name = '1 Serisi E87'),
    '116i', 2004, 115, 150, 'Otomatik', 'GM 6L45', '1.6 Atmosferik', 'Arkadan İtiş', 'Benzin', 10.8, NULL
);

COMMIT;
