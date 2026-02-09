-- 021_populate_bmw_5_series.sql
BEGIN TRANSACTION;

-- 1. Ensure Brand Exists (BMW)
INSERT OR IGNORE INTO brands (name, country, logo_url) 
VALUES ('BMW', 'Germany', '/images/brands/bmw.png');

-- 2. Clean existing 5 Series Data (if any)
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id IN (
        SELECT id FROM models WHERE name = '5 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
    )
);
DELETE FROM generations WHERE model_id IN (
    SELECT id FROM models WHERE name = '5 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW')
);
DELETE FROM models WHERE name = '5 Serisi' AND brand_id = (SELECT id FROM brands WHERE name = 'BMW');

-- 3. Create Model
INSERT INTO models (brand_id, name, image_url, body_style)
VALUES (
    (SELECT id FROM brands WHERE name = 'BMW'),
    '5 Serisi',
    '/images/models/bmw-5-series.png', 
    'Sedan'
);

-- 4. Create Generations
-- G60
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi G60', 'G60', 2024, 2024, '/images/generations/bmw-5-g60.png', 0
);
-- G30 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi G30 (LCI)', 'G30 LCI', 2020, 2023, '/images/generations/bmw-5-g30-lci.png', 1
);
-- G30
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi G30', 'G30', 2017, 2020, '/images/generations/bmw-5-g30.png', 0
);
-- F10 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi F10 (LCI)', 'F10 LCI', 2013, 2017, '/images/generations/bmw-5-f10-lci.png', 1
);
-- F10
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi F10', 'F10', 2010, 2013, '/images/generations/bmw-5-f10.png', 0
);
-- E60 LCI
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi E60 (LCI)', 'E60 LCI', 2007, 2010, '/images/generations/bmw-5-e60-lci.png', 1
);
-- E60
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = '5 Serisi'),
    '5 Serisi E60', 'E60', 2003, 2007, '/images/generations/bmw-5-e60.png', 0
);

-- 5. Insert Trims

-- === G60 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi G60'),
    '520d xDrive', 2024, 197, 400, 'Otomatik (ZF)', 'ZF 8HP', '2.0 TDI Mild Hybrid', 'xDrive (4x4)', 'Dizel', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi G60'),
    'i5 eDrive40', 2024, 340, 430, 'Otomatik', 'Tek Oranlı', 'Elektrik Motoru', 'Arkadan İtiş', 'Elektrik', NULL, NULL
);

-- === G30 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi G30 (LCI)'),
    '520i Special Edition', 2020, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', '1.6 Turbo Mild Hybrid', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === G30 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi G30'),
    '520i', 2017, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', 'B48B16', 'Arkadan İtiş', 'Benzin', 8.3, NULL
),
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi G30'),
    '530i', 2017, 252, 350, 'Otomatik (ZF)', 'ZF 8HP', '2.0 Turbo', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === F10 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi F10 (LCI)'),
    '520i', 2013, 170, 250, 'Otomatik (ZF)', 'ZF 8HP', 'N20B16', 'Arkadan İtiş', 'Benzin', 8.7, NULL
),
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi F10 (LCI)'),
    '525d xDrive', 2013, 218, 450, 'Otomatik (ZF)', 'ZF 8HP', '2.0 Bi-Turbo Dizel', 'xDrive (4x4)', 'Dizel', NULL, NULL
);

-- === F10 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi F10'),
    '520d', 2010, 184, 380, 'Otomatik (ZF)', 'ZF 8HP', '2.0 TDI', 'Arkadan İtiş', 'Dizel', 8.1, NULL
);

-- === E60 LCI ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi E60 (LCI)'),
    '520d', 2007, 177, 350, 'Otomatik (ZF)', 'ZF 6HP', 'N47', 'Arkadan İtiş', 'Dizel', 8.4, NULL
),
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi E60 (LCI)'),
    '520i', 2007, 170, 210, 'Otomatik (ZF)', 'ZF 6HP', '2.0 Atmosferik', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

-- === E60 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi E60'),
    '520d', 2003, 163, 340, 'Otomatik (ZF)', 'ZF 6HP', 'M47', 'Arkadan İtiş', 'Dizel', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = '5 Serisi'),
    (SELECT id FROM generations WHERE name = '5 Serisi E60'),
    '520i', 2003, 170, 210, 'Otomatik (ZF)', 'ZF 6HP', '2.2 6 Silindir', 'Arkadan İtiş', 'Benzin', NULL, NULL
);

COMMIT;
