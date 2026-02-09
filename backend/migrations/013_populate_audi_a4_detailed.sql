-- Migration: Populate detailed Audi A4 Trims (B7, B8, B8.5, B9, B9.5)
-- Clears existing trims for affected generations and inserts precise user data

BEGIN TRANSACTION;

-- Step 1: Clear existing trims for the targeted generations to avoid duplicates
-- We will re-insert them with higher fidelity data
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id = (SELECT id FROM models WHERE name = 'A4')
    AND code IN ('B7', 'B8', 'B8.5', 'B9', 'B9.5')
);

-- Step 2: Insert Trims

-- ==========================================
-- 1. Audi A4 B7 (2005 - 2008)
-- ==========================================
-- 1.6 (Hitachi Motor)
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year
)
SELECT m.id, g.id, '1.6', 'Benzin', 1595, 4, 
       102, 148, 'Manuel', 'Önden Çekiş', 
       12.6, 2005
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B7' AND m.name = 'A4';

-- 1.8 T (20V Turbo)
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '1.8 T', 'Benzin', 1781, 4, 
       163, 225, 'Multitronic', 'Önden Çekiş', 
       8.6, 2005, 'VL300'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B7' AND m.name = 'A4';

-- 2.0 TDI (PD)
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TDI', 'Dizel', 1968, 4, 
       140, 320, 'Multitronic', 'Önden Çekiş', 
       9.8, 2005, 'VL300'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B7' AND m.name = 'A4';

-- 2.0 TFSI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TFSI quattro', 'Benzin', 1984, 4, 
       200, 280, 'Tiptronic', 'Quattro', 
       7.7, 2005, 'ZF 6HP'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B7' AND m.name = 'A4';


-- ==========================================
-- 2. Audi A4 B8 (2008 - 2011)
-- ==========================================
-- 1.8 TFSI
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '1.8 TFSI', 'Benzin', 1798, 4, 
       160, 250, 'Multitronic', 'Önden Çekiş', 
       8.6, 2008, 'VL381'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8' AND m.name = 'A4';

-- 2.0 TDI
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TDI', 'Dizel', 1968, 4, 
       143, 320, 'Multitronic', 'Önden Çekiş', 
       9.4, 2008, 'VL381'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8' AND m.name = 'A4';

-- 2.0 TFSI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TFSI quattro', 'Benzin', 1984, 4, 
       211, 350, 'S tronic', 'Quattro', 
       6.5, 2008, 'DL501'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8' AND m.name = 'A4';


-- ==========================================
-- 3. Audi A4 B8.5 (2012 - 2015)
-- ==========================================
-- 1.8 TFSI (Güncel Motor)
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code, engine_code
)
SELECT m.id, g.id, '1.8 TFSI', 'Benzin', 1798, 4, 
       170, 320, 'Multitronic', 'Önden Çekiş', 
       8.3, 2012, 'VL381', 'CJEB'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8.5' AND m.name = 'A4';

-- 2.0 TDI (Clean Diesel)
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TDI', 'Dizel', 1968, 4, 
       150, 320, 'Multitronic', 'Önden Çekiş', 
       9.1, 2012, 'VL381'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8.5' AND m.name = 'A4';

-- 2.0 TFSI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TFSI quattro', 'Benzin', 1984, 4, 
       225, 350, 'S tronic', 'Quattro', 
       6.4, 2012, 'DL501'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B8.5' AND m.name = 'A4';


-- ==========================================
-- 4. Audi A4 B9 (2016 - 2019)
-- ==========================================
-- 1.4 TFSI
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '1.4 TFSI', 'Benzin', 1395, 4, 
       150, 250, 'S tronic', 'Önden Çekiş', 
       8.5, 2016, 'DL382'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B9' AND m.name = 'A4';

-- 2.0 TDI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TDI quattro', 'Dizel', 1968, 4, 
       190, 400, 'S tronic', 'Quattro', 
       7.2, 2016, 'DL382'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B9' AND m.name = 'A4';

-- 2.0 TFSI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '2.0 TFSI quattro', 'Benzin', 1984, 4, 
       252, 370, 'S tronic', 'Quattro', 
       5.8, 2016, 'DL382'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B9' AND m.name = 'A4';


-- ==========================================
-- 5. Audi A4 B9.5 (2020 - Present)
-- ==========================================
-- 40 TDI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '40 TDI quattro', 'Dizel (Mild-Hybrid)', 1968, 4, 
       204, 400, 'S tronic', 'Quattro', 
       6.9, 2020, 'DL382+'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B9.5' AND m.name = 'A4';

-- 45 TFSI quattro
INSERT INTO trims (
    model_id, generation_id, name, fuel_type, displacement_cc, cylinders, 
    power_hp, torque_nm, transmission_type, drivetrain, 
    acceleration_0_100, year, transmission_code
)
SELECT m.id, g.id, '45 TFSI quattro', 'Benzin (Mild-Hybrid)', 1984, 4, 
       265, 370, 'S tronic', 'Quattro', 
       5.5, 2020, 'DL382+'
FROM generations g 
JOIN models m ON g.model_id = m.id
WHERE g.code = 'B9.5' AND m.name = 'A4';

COMMIT;
