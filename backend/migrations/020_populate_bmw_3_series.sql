-- 020_populate_bmw_3_series.sql
BEGIN TRANSACTION;

-- 1. Ensure Brand Exists (BMW)
INSERT OR IGNORE INTO brands (name, country, logo_url) 
VALUES ('BMW', 'Germany', '/images/brands/bmw.png');

-- 2. Clean existing 3 Series Data (if any)
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id IN (
        SELECT id FROM models WHERE name = '3 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
    )
);
DELETE FROM generations WHERE model_id IN (
    SELECT id FROM models WHERE name = '3 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
);
DELETE FROM models WHERE name = '3 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW');

-- 3. Create Model
INSERT INTO models (brand_id, name, image_url, body_style)
VALUES (
    (SELECT id FROM brands WHERE name = 'BMW'),
    '3 Serisi',
    '/images/models/bmw-3-series.png', 
    'Sedan'
);

-- 4. Create Generations
-- G20 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi G20 (LCI)', 'G20 LCI', 2022, 2024, '/images/generations/bmw-3-g20-lci.png', 1
);
-- G20
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi G20', 'G20', 2019, 2022, '/images/generations/bmw-3-g20.png', 0
);
-- F30 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi F30 (LCI)', 'F30 LCI', 2015, 2019, '/images/generations/bmw-3-f30-lci.png', 1
);
-- F30
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi F30', 'F30', 2012, 2015, '/images/generations/bmw-3-f30.png', 0
);
-- E90 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi E90 (LCI)', 'E90 LCI', 2008, 2012, '/images/generations/bmw-3-e90-lci.png', 1
);
-- E90
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '3 Serisi'),
    '3 Serisi E90', 'E90', 2005, 2008, '/images/generations/bmw-3-e90.png', 0
);

-- 5. Insert Trims

-- === G20 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi G20 (LCI)'),
    '320i', 2022, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', '1.6 Turbo', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === G20 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi G20'),
    '320i (Türkiye Paketi)', 2019, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', '1.6 Turbo (4 Silindir)', 'Arkadan İtiş', 'Benzin', 7.7, NULL
);

-- === F30 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30 (LCI)'),
    '318i', 2015, 136, 220, 'Otomatik (ZF)', 'ZF 8HP', 'B38', 'Arkadan İtiş', 'Benzin', 9.1, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30 (LCI)'),
    '320i ED', 2015, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', 'N13', 'Arkadan İtiş', 'Benzin', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30 (LCI)'),
    '320d', 2015, 190, 400, 'Otomatik (ZF)', 'ZF 8HP', 'B47', 'Arkadan İtiş', 'Dizel', NULL, NULL
);

-- === F30 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30'),
    '316i', 2012, 136, 220, 'Otomatik (ZF)', 'ZF 8HP', 'N13', 'Arkadan İtiş', 'Benzin', 9.2, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30'),
    '320i ED', 2012, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', 'N13', 'Arkadan İtiş', 'Benzin', 7.6, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi F30'),
    '320d', 2012, 184, 380, 'Otomatik (ZF)', 'ZF 8HP', 'N47', 'Arkadan İtiş', 'Dizel', NULL, NULL
);

-- === E90 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi E90 (LCI)'),
    '320i', 2008, 156, 200, 'Otomatik', 'GM 6L45', '2.0 Atmosferik', 'Arkadan İtiş', 'Benzin', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi E90 (LCI)'),
    '320d', 2008, 177, 350, 'Otomatik (ZF)', 'ZF 6HP', 'N47', 'Arkadan İtiş', 'Dizel', 7.6, NULL
);

-- === E90 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi E90'),
    '320i', 2005, 150, 200, 'Otomatik', 'GM 6L45 / ZF 6HP', '2.0 Atmosferik', 'Arkadan İtiş', 'Benzin', 9.8, NULL
),
(
    (SELECT id FROM models WHERE name = '3 Serisi'),
    (SELECT id FROM generations WHERE name = '3 Serisi E90'),
    '320d', 2005, 163, 340, 'Otomatik (ZF)', 'ZF 6HP', 'M47', 'Arkadan İtiş', 'Dizel', NULL, NULL
);

COMMIT;
