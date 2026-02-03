-- Add Audi A3 Generation Data
-- Model: Audi A3 (ID: 1)
-- Generations: 8L, 8P, 8V, 8Y

INSERT INTO generations (model_id, code, name, start_year, end_year, is_current) VALUES
-- 1. Nesil: Typ 8L (1996-2003)
(1, '8L', 'Typ 8L', 1996, 2003, 0),

-- 2. Nesil: Typ 8P (2003-2012)
(1, '8P', 'Typ 8P', 2003, 2012, 0),

-- 3. Nesil: Typ 8V (2012-2020)
(1, '8V', 'Typ 8V', 2012, 2020, 0),

-- 4. Nesil: Typ 8Y (2020-Günümüz)
(1, '8Y', 'Typ 8Y', 2020, NULL, 1);
