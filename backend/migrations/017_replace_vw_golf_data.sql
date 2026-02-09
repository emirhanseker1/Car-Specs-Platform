-- 017_replace_vw_golf_data.sql
BEGIN TRANSACTION;

-- 1. Clean existing Golf Data
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id IN (
        SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')
    )
);
DELETE FROM generations WHERE model_id IN (
    SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')
);
DELETE FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen');

-- 2. Ensure Brand Exists (Volkswagen)
INSERT OR IGNORE INTO brands (name, country, logo_url) 
VALUES ('Volkswagen', 'Germany', '/images/brands/volkswagen.png');

-- 3. Re-Create Model
INSERT INTO models (brand_id, name, image_url, body_style)
VALUES (
    (SELECT id FROM brands WHERE name = 'Volkswagen'),
    'Golf',
    '/images/models/vw-golf.png', 
    'Hatchback'
);

-- 4. Create Generations
-- MK8
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')),
    'Golf 8', 'MK8', 2020, 2024, '/images/generations/vw-golf-8.png', 0
);
-- MK7.5
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')),
    'Golf 7.5 (Makyajlı)', 'MK7.5', 2017, 2020, '/images/generations/vw-golf-7-5.png', 1
);
-- MK7
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')),
    'Golf 7', 'MK7', 2012, 2016, '/images/generations/vw-golf-7.png', 0
);
-- MK6
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')),
    'Golf 6', 'MK6', 2009, 2012, '/images/generations/vw-golf-6.png', 0
);
-- MK5
INSERT INTO generations (model_id, name, code, start_year, end_year, image_url, is_facelift)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf' AND brand_id = (SELECT id FROM brands WHERE name = 'Volkswagen')),
    'Golf 5', 'MK5', 2004, 2008, '/images/generations/vw-golf-5.png', 0
);

-- 5. Insert Trims

-- === GOLF 8 ===
-- 1.0 eTSI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 8'),
    '1.0 eTSI Life', 2020, 110, 200, 'Otomatik (DSG)', 'DQ200', '1.0 eTSI', 'Önden Çekiş', 'Mild Hybrid', 10.2, 202, '/images/trims/golf8-life.png'
);
-- 1.5 eTSI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 8'),
    '1.5 eTSI Style', 2020, 150, 250, 'Otomatik (DSG)', 'DQ200', '1.5 eTSI', 'Önden Çekiş', 'Mild Hybrid', 8.5, 224, '/images/trims/golf8-style.png'
);

-- === GOLF 7.5 ===
-- 1.0 TSI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7.5 (Makyajlı)'),
    '1.0 TSI Comfortline', 2017, 110, 200, 'Otomatik (DSG)', 'DQ200', '1.0 TSI', 'Önden Çekiş', 'Benzin', 9.9, 196, NULL
);
-- 1.4 TSI 125
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7.5 (Makyajlı)'),
    '1.4 TSI Comfortline', 2017, 125, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', 9.1, 204, NULL
);
-- 1.5 TSI Evo
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7.5 (Makyajlı)'),
    '1.5 TSI EVO Highline', 2018, 150, 250, 'Otomatik (DSG)', 'DQ200', '1.5 TSI Evo', 'Önden Çekiş', 'Benzin', 8.3, 216, NULL
);
-- 1.6 TDI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7.5 (Makyajlı)'),
    '1.6 TDI Comfortline', 2017, 115, 250, 'Otomatik (DSG)', 'DQ200', '1.6 TDI', 'Önden Çekiş', 'Dizel', 10.2, 198, NULL
);

-- === GOLF 7 ===
-- 1.2 TSI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7'),
    '1.2 TSI Midline Plus', 2013, 105, 175, 'Otomatik (DSG)', 'DQ200', '1.2 TSI', 'Önden Çekiş', 'Benzin', 10.2, 192, NULL
);
-- 1.4 TSI 122
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7'),
    '1.4 TSI Comfortline', 2013, 122, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', 9.3, 203, NULL
);
-- 1.4 TSI ACT
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7'),
    '1.4 TSI ACT Highline', 2014, 140, 250, 'Otomatik (DSG)', 'DQ200', '1.4 TSI ACT', 'Önden Çekiş', 'Benzin', 8.4, 212, NULL
);
-- 1.6 TDI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 7'),
    '1.6 TDI Comfortline', 2013, 105, 250, 'Otomatik (DSG)', 'DQ200', '1.6 TDI', 'Önden Çekiş', 'Dizel', 10.7, 190, NULL
);

-- === GOLF 6 ===
-- 1.6 Hitachi
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 6'),
    '1.6 Primeline', 2009, 102, 148, 'Otomatik (DSG)', 'DQ200', '1.6', 'Önden Çekiş', 'Benzin', 11.3, 188, NULL
);
-- 1.4 TSI 122
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 6'),
    '1.4 TSI Comfortline', 2009, 122, 200, 'Otomatik (DSG)', 'DQ200', '1.4 TSI', 'Önden Çekiş', 'Benzin', 9.5, 200, NULL
);
-- 1.4 TSI 160
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 6'),
    '1.4 TSI Highline (Twincharger)', 2009, 160, 240, 'Otomatik (DSG)', 'DQ200', '1.4 TSI Twincharger', 'Önden Çekiş', 'Benzin', 8.0, 220, NULL
);
-- 1.6 TDI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 6'),
    '1.6 TDI Comfortline', 2009, 105, 250, 'Otomatik (DSG)', 'DQ200', '1.6 TDI', 'Önden Çekiş', 'Dizel', 11.3, 189, NULL
);

-- === GOLF 5 ===
-- 1.6 Primeline
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 5'),
    '1.6 Primeline', 2004, 102, 148, 'Otomatik (Tiptronic)', '09G', '1.6', 'Önden Çekiş', 'Benzin', 12.5, 184, NULL
);
-- 1.6 FSI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 5'),
    '1.6 FSI Comfortline', 2004, 115, 155, 'Otomatik (Tiptronic)', '09G', '1.6 FSI', 'Önden Çekiş', 'Benzin', 11.5, 192, NULL
);
-- 1.4 TSI GT
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 5'),
    '1.4 TSI GT', 2006, 170, 240, 'Otomatik (DSG)', 'DQ250', '1.4 TSI', 'Önden Çekiş', 'Benzin', 7.9, 220, NULL
);
-- 1.9 TDI
INSERT INTO trims (model_id, generation_id, name, year, power_hp, torque_nm, transmission_type, transmission_code, engine_code, drivetrain, fuel_type, acceleration_0_100, top_speed_kmh, image_url)
VALUES (
    (SELECT id FROM models WHERE name = 'Golf'),
    (SELECT id FROM generations WHERE name = 'Golf 5'),
    '1.9 TDI Comfortline', 2004, 105, 250, 'Otomatik (DSG)', 'DQ250', '1.9 TDI', 'Önden Çekiş', 'Dizel', 11.3, 187, NULL
);

COMMIT;
