-- 018_populate_vw_passat.sql
BEGIN TRANSACTION;

-- 1. Clean existing Passat Data (if any)
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id IN (
        SELECT id FROM models WHERE name = 'Passat' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')
    )
);
DELETE FROM generations WHERE model_id IN (
    SELECT id FROM models WHERE name = 'Passat' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')
);
DELETE FROM models WHERE name = 'Passat' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen');

-- 2. Ensure Brand Exists (Volkswagen) - Should already exist from Golf migration
INSERT OR IGNORE INTO brands (name, country, logo_url) 
VALUES ('Volkswagen', 'Germany', '/images/brands/volkswagen.png');

-- 3. Create Model
INSERT INTO models (brand_id, name, image_url, body_style)
VALUES (
    (SELECT id FROM brands WHERE name = 'Volkswagen'),
    'Passat',
    '/images/models/vw-passat.png', 
    'Sedan'
);

-- 4. Create Generations
-- B8.5
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Passat'),
    'Passat B8.5 (Makyajlı)', 'B8.5', 2019, 2024, '/images/generations/vw-passat-b8.5.png', 1
);
-- B8
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Passat'),
    'Passat B8', 'B8', 2014, 2019, '/images/generations/vw-passat-b8.png', 0
);
-- B7
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Passat'),
    'Passat B7', 'B7', 2010, 2014, '/images/generations/vw-passat-b7.png', 0
);
-- B6
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Passat'),
    'Passat B6', 'B6', 2005, 2010, '/images/generations/vw-passat-b6.png', 0
);
-- B5.5
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Passat'),
    'Passat B5.5', 'B5.5', 2000, 2005, '/images/generations/vw-passat-b5.5.png', 1
);

-- 5. Insert Trims

-- === Passat B8.5 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8.5 (Makyajlı)'),
    '1.5 TSI ACT Evo', 2019, 150, 250, 'Otomatik (DSG)', 'DQ200', '1.5 TSI Evo', 'Önden Çekiş', 'Benzin', 8.7, '/images/trims/passat-b8.5-tsi.png'
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8.5 (Makyajlı)'),
    '2.0 TDI Evo', 2019, 150, 360, 'Otomatik (DSG)', 'DQ381', '2.0 TDI Evo', 'Önden Çekiş', 'Dizel', 8.9, '/images/trims/passat-b8.5-tdi.png'
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8.5 (Makyajlı)'),
    '2.0 TDI 4Motion', 2019, 200, 400, 'Otomatik (DSG)', 'DQ381', '2.0 TDI', 'Quattro (4Motion)', 'Dizel', NULL, NULL
);

-- === Passat B8 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8'),
    '1.4 TSI', 2014, 125, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', 9.7, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8'),
    '1.4 TSI ACT', 2014, 150, 250, 'Otomatik (DSG)', 'DQ200', '1.4 TSI ACT', 'Önden Çekiş', 'Benzin', 8.4, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8'),
    '1.6 TDI', 2014, 120, 250, 'Otomatik (DSG)', 'DQ200', '1.6 TDI', 'Önden Çekiş', 'Dizel', 10.8, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8'),
    '2.0 TDI', 2014, 150, 340, 'Otomatik (DSG)', 'DQ250', '2.0 TDI', 'Önden Çekiş', 'Dizel', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B8'),
    '2.0 Bi-TDI 4Motion', 2014, 240, 500, 'Otomatik (DSG)', 'DQ500', '2.0 Bi-TDI', 'Quattro (4Motion)', 'Dizel', 6.1, NULL
);

-- === Passat B7 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B7'),
    '1.4 TSI BlueMotion', 2010, 122, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', 10.3, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B7'),
    '1.4 TSI Twincharger', 2010, 160, 240, 'Otomatik (DSG)', 'DQ200', '1.4 TSI Twincharger', 'Önden Çekiş', 'Benzin', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B7'),
    '1.6 TDI BlueMotion', 2010, 105, 250, 'Otomatik (DSG)', 'DQ200', '1.6 TDI', 'Önden Çekiş', 'Dizel', 12.2, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B7'),
    '2.0 TDI', 2010, 140, 320, 'Otomatik (DSG)', 'DQ250', '2.0 TDI', 'Önden Çekiş', 'Dizel', NULL, NULL
);

-- === Passat B6 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B6'),
    '1.6 FSI', 2005, 115, 155, 'Otomatik (Tiptronic)', '09G / AQ250', '1.6 FSI', 'Önden Çekiş', 'Benzin', 11.5, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B6'),
    '1.4 TSI (2008+)', 2008, 122, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B6'),
    '2.0 TDI', 2005, 140, 320, 'Otomatik (DSG)', 'DQ250', '2.0 TDI', 'Önden Çekiş', 'Dizel', NULL, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B6'),
    '2.0 FSI', 2005, 150, 200, 'Otomatik (Tiptronic)', '09G', '2.0 FSI', 'Önden Çekiş', 'Benzin', NULL, NULL
);

-- === Passat B5.5 ===
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, image_url)
VALUES 
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B5.5'),
    '1.8 T', 2000, 150, 210, 'Otomatik (Tiptronic)', 'ZF 5HP19', '1.8 T', 'Önden Çekiş', 'Benzin', 9.2, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B5.5'),
    '1.9 TDI', 2000, 130, 310, 'Otomatik (Tiptronic)', 'ZF 5HP19', '1.9 TDI', 'Önden Çekiş', 'Dizel', 9.9, NULL
),
(
    (SELECT id FROM models WHERE name = 'Passat'),
    (SELECT id FROM generations WHERE name = 'Passat B5.5'),
    '1.6', 2000, 102, 148, 'Otomatik', '4-İleri', '1.6', 'Önden Çekiş', 'Benzin', 12.7, NULL
);

COMMIT;
