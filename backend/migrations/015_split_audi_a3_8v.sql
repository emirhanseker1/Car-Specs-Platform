-- Migration: Split Audi A3 8V into Pre-Facelift (8V1) and Facelift (8V2)
-- Based on user details for 2013-2016 and 2016-2019 periods.

BEGIN TRANSACTION;

-- 0. Ensure is_facelift column exists in generations
-- SQLite doesn't support IF NOT EXISTS in ADD COLUMN well in older versions, 
-- but we can ignore error or use a separate script.
-- However, since we got "no such column", we know it's missing.
-- We'll try to add it. If it fails (already exists), the whole transaction might abort depending on driver.
-- But since we are fairly sure it's missing, this is the fix.
ALTER TABLE generations ADD COLUMN is_facelift BOOLEAN DEFAULT 0;

-- Variable to store model_id for Audi A3
-- SQLite doesn't support variables in pure SQL scripts easily without tools, 
-- so we'll use subqueries e.g., (SELECT id FROM models WHERE name = 'A3' AND brand_id = ...)

-- 1. Identify Model ID
-- Assuming Audi A3 is model_id 1 based on previous files, but using subquery for safety
-- (SELECT id FROM models WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi'))

-- 2. Clean up existing 8V trims to avoid duplication
DELETE FROM trims 
WHERE generation_id IN (
    SELECT id FROM generations 
    WHERE code LIKE '8V%' 
    AND model_id = (SELECT id FROM models WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi'))
);

-- 3. Modify existing 8V generation to be 8V1 (Pre-Facelift)
UPDATE generations 
SET 
    code = '8V1',
    name = 'Tip 8V - 3. Nesil (2013 - 2016) | Makyajsız',
    start_year = 2013,
    end_year = 2016,
    is_facelift = 0,
    description = 'Makyajsız kasa. Cylinder on Demand (COD) teknolojisi belirli motorlarda sunulmaya başlanmıştır.'
WHERE code = '8V' 
AND model_id = (SELECT id FROM models WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi'));

-- 4. Create new 8V2 generation (Facelift)
INSERT INTO generations (model_id, code, name, start_year, end_year, is_facelift, description, image_url)
SELECT 
    id, 
    '8V2', 
    'Tip 8V - 3. Nesil (2016 - 2019) | Makyajlı', 
    2016, 
    2019, 
    1, 
    'Makyajlı kasa. "Şimşek" formlu LED farlar ve Sanal Kokpit (Hayalet Ekran) opsiyonu.',
    image_url -- Copy image from parent model or previous generation if needed, or leave null
FROM models 
WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi');


-- ==========================================
-- 5. Insert Trims for 8V1 (2013 - 2016)
-- ==========================================

-- 1.2 TFSI (105 PS) - 2013-2014
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.2 TFSI', 2013, 2013, 2014,
    'Turbo Benzinli', 'Benzin', 1197, 105, 175,
    'S tronic', 'DQ200', 10.3, 5.0,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.2 TFSI (110 PS) - 2014-2016
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.2 TFSI', 2014, 2014, 2016,
    'Turbo Benzinli', 'Benzin', 1197, 110, 175,
    'S tronic', 'DQ200', 10.3, 5.0,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.4 TFSI (122 PS) - 2013-2014
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.4 TFSI', 2013, 2013, 2014,
    'Turbo Benzinli', 'Benzin', 1395, 122, 200,
    'S tronic', 'DQ200', 9.3, 5.1,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.4 TFSI (125 PS) - 2014-2016
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.4 TFSI', 2014, 2014, 2016,
    'Turbo Benzinli', 'Benzin', 1395, 125, 200,
    'S tronic', 'DQ200', 9.3, 5.1,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.4 TFSI COD (140 PS) - 2013-2016
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.4 TFSI COD', 2013, 2013, 2016,
    'Turbo Benzinli (Silindir Kapama)', 'Benzin', 1395, 140, 250,
    'S tronic', 'DQ200', 8.4, 4.8,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.6 TDI (105 PS) - 2013-2015
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.6 TDI', 2013, 2013, 2015,
    'Turbo Dizel', 'Dizel', 1598, 105, 250,
    'S tronic', 'DQ200', 10.7, 3.9,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';

-- 1.6 TDI (110 PS) - 2015-2016
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.6 TDI', 2015, 2015, 2016,
    'Turbo Dizel (Euro 6)', 'Dizel', 1598, 110, 250,
    'S tronic', 'DQ200', 10.7, 3.9,
    4, 'Önden Çekiş', 0
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V1' AND m.name = 'A3';


-- ==========================================
-- 6. Insert Trims for 8V2 (2016 - 2019)
-- ==========================================

-- 1.0 TFSI (30 TFSI) - 2016-2019
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.0 TFSI (30 TFSI)', 2016, 2016, 2019,
    'Turbo Benzinli', 'Benzin', 999, 116, 200,
    'S tronic', 'DQ200', 9.9, 4.6,
    3, 'Önden Çekiş', 1
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V2' AND m.name = 'A3';

-- 1.4 TFSI COD - 2016-2017
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.4 TFSI COD', 2016, 2016, 2017,
    'Turbo Benzinli (Silindir Kapama)', 'Benzin', 1395, 150, 250,
    'S tronic', 'DQ200', 8.2, 4.8, -- Tüketim yakındır diye tahmin ediyoruz veya null
    4, 'Önden Çekiş', 1
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V2' AND m.name = 'A3';

-- 1.5 TFSI (35 TFSI) - 2017-2019
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.5 TFSI (35 TFSI)', 2017, 2017, 2019,
    'Turbo Benzinli (Evo)', 'Benzin', 1498, 150, 250,
    'S tronic', 'DQ200', 8.2, 5.0,
    4, 'Önden Çekiş', 1
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V2' AND m.name = 'A3';

-- 1.6 TDI (30 TDI) - 2016-2019
INSERT INTO trims (
    generation_id, model_id, name, year, start_year, end_year, 
    engine_type, fuel_type, displacement_cc, power_hp, torque_nm, 
    transmission_type, transmission_code, acceleration_0_100, fuel_consumption_combined,
    cylinders, drivetrain, is_facelift
)
SELECT 
    g.id, m.id, '1.6 TDI (30 TDI)', 2016, 2016, 2019,
    'Turbo Dizel', 'Dizel', 1598, 116, 250,
    'S tronic', 'DQ200', 10.4, 4.0,
    4, 'Önden Çekiş', 1
FROM generations g
JOIN models m ON g.model_id = m.id
WHERE g.code = '8V2' AND m.name = 'A3';


COMMIT;
