-- Migration: Fix Audi A3 Generation Data and Add Comprehensive Motor Trims
-- Based on user-provided historical data for all 4 generations

BEGIN TRANSACTION;

-- Step 1: Update Generation Information with Correct Years and Names
UPDATE generations SET 
    name = 'Tip 8L - İlk Nesil (1996 - 2003)',
    start_year = 1996,
    end_year = 2003
WHERE code = '8L' AND model_id = 1;

UPDATE generations SET 
    name = 'Tip 8P - İkinci Nesil (2003 - 2012)',
    start_year = 2003,
    end_year = 2012
WHERE code = '8P' AND model_id = 1;

UPDATE generations SET 
    name = 'Tip 8V - Üçüncü Nesil (2012 - 2020)',
    start_year = 2012,
    end_year = 2020
WHERE code = '8V' AND model_id = 1;

UPDATE generations SET 
    name = 'Tip 8Y - Dördüncü Nesil (2020 - Günümüz)',
    start_year = 2020,
    end_year = NULL
WHERE code = '8Y' AND model_id = 1;

-- Step 2: Clear existing incorrect trims for A3
DELETE FROM trims WHERE generation_id IN (
    SELECT id FROM generations WHERE model_id = 1
);

-- Step 3: Insert Trims for Generation 8L (1996-2003)
-- Benzinli motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6', 'Benzin', 1595, 101, 145, 'Manuel', 'Önden Çekiş', 1996, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.8', 'Benzin', 1781, 125, 170, 'Manuel', 'Önden Çekiş', 1996, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.8 Turbo', 'Benzin', 1781, 150, 210, 'Manuel', 'Önden Çekiş', 1996, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.8 Turbo', 'Benzin', 1781, 180, 235, 'Manuel', 'Önden Çekiş', 1999, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

-- Dizel motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.9 TDI', 'Dizel', 1896, 90, 210, 'Manuel', 'Önden Çekiş', 1996, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.9 TDI', 'Dizel', 1896, 110, 235, 'Manuel', 'Önden Çekiş', 1997, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.9 TDI', 'Dizel', 1896, 130, 310, 'Manuel', 'Önden Çekiş', 2000, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

-- Performans: S3
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'S3 1.8 Turbo', 'Benzin', 1781, 210, 270, 'Manuel', 'Quattro', 1999, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'S3 1.8 Turbo', 'Benzin', 1781, 225, 280, 'Manuel', 'Quattro', 2001, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8L' AND g.model_id = 1;


-- Step 4: Insert Trims for Generation 8P (2003-2012)
-- Benzinli motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6 MPI', 'Benzin', 1595, 102, 148, 'Manuel', 'Önden Çekiş', 2003, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6 FSI', 'Benzin', 1598, 115, 155, 'Manuel', 'Önden Çekiş', 2003, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.4 TFSI', 'Benzin', 1390, 125, 200, 'Manuel', 'Önden Çekiş', 2007, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '2.0 TFSI', 'Benzin', 1984, 200, 280, 'S tronic', 'Önden Çekiş', 2004, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '3.2 V6 Quattro', 'Benzin', 3189, 250, 330, 'S tronic', 'Quattro', 2003, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

-- Dizel motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6 TDI', 'Dizel', 1598, 105, 250, 'Manuel', 'Önden Çekiş', 2009, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '2.0 TDI', 'Dizel', 1968, 140, 320, 'Manuel', 'Önden Çekiş', 2003, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '2.0 TDI', 'Dizel', 1968, 170, 350, 'S tronic', 'Önden Çekiş', 2005, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;

-- Performans: RS3
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'RS3 2.5 TFSI', 'Benzin', 2480, 340, 450, 'S tronic', 'Quattro', 2011, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8P' AND g.model_id = 1;


-- Step 5: Insert Trims for Generation 8V (2012-2020)
-- Benzinli motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.0 TFSI', 'Benzin', 999, 116, 200, 'Manuel', 'Önden Çekiş', 2016, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.2 TFSI', 'Benzin', 1197, 105, 175, 'Manuel', 'Önden Çekiş', 2012, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.4 TFSI', 'Benzin', 1395, 125, 200, 'Manuel', 'Önden Çekiş', 2012, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.4 TFSI COD', 'Benzin', 1395, 150, 250, 'S tronic', 'Önden Çekiş', 2013, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.5 TFSI', 'Benzin', 1498, 150, 250, 'S tronic', 'Önden Çekiş', 2017, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '2.0 TFSI', 'Benzin', 1984, 190, 320, 'S tronic', 'Önden Çekiş', 2012, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

-- Dizel motorlar
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6 TDI', 'Dizel', 1598, 105, 250, 'Manuel', 'Önden Çekiş', 2012, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '1.6 TDI', 'Dizel', 1598, 110, 250, 'Manuel', 'Önden Çekiş', 2013, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '2.0 TDI', 'Dizel', 1968, 150, 340, 'S tronic', 'Önden Çekiş', 2013, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

-- Performans: S3 & RS3
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'S3 2.0 TFSI', 'Benzin', 1984, 300, 380, 'S tronic', 'Quattro', 2013, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'RS3 2.5 TFSI', 'Benzin', 2480, 367, 465, 'S tronic', 'Quattro', 2015, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'RS3 2.5 TFSI', 'Benzin', 2480, 400, 480, 'S tronic', 'Quattro', 2017, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8V' AND g.model_id = 1;


-- Step 6: Insert Trims for Generation 8Y (2020-Present)
-- Yeni isimlendirme sistemi ile
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '30 TFSI', 'Benzin', 999, 110, 200, 'Manuel', 'Önden Çekiş', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '35 TFSI', 'Benzin (Mild-Hybrid)', 1498, 150, 250, 'S tronic', 'Önden Çekiş', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

-- Dizel
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '30 TDI', 'Dizel', 1968, 116, 300, 'Manuel', 'Önden Çekiş', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '35 TDI', 'Dizel', 1968, 150, 360, 'S tronic', 'Önden Çekiş', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

-- Plug-in Hybrid
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '40 TFSI e', 'Plug-in Hybrid', 1395, 204, 350, 'S tronic', 'Önden Çekiş', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, '45 TFSI e', 'Plug-in Hybrid', 1395, 245, 400, 'S tronic', 'Quattro', 2021, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

-- Performans: S3 & RS3
INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'S3 2.0 TFSI', 'Benzin', 1984, 310, 400, 'S tronic', 'Quattro', 2020, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

INSERT INTO trims (model_id, generation_id, name, fuel_type, displacement_cc, power_hp, torque_nm, transmission_type, drivetrain, year, created_at, updated_at)
SELECT 1, g.id, 'RS3 2.5 TFSI', 'Benzin', 2480, 400, 500, 'S tronic', 'Quattro', 2021, datetime('now'), datetime('now')
FROM generations g WHERE g.code = '8Y' AND g.model_id = 1;

COMMIT;
