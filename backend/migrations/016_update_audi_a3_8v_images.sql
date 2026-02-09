-- Migration: Update Audi A3 8V Images
-- Based on user request: 8V1 gets new image, 8V2 gets old 8V1 image.

BEGIN TRANSACTION;

-- 1. Identify Model information for A3
-- We use subqueries to target the rows safely.

-- 2. Update 8V2 (Facelift) to use the OLD 8V image
-- The old 8V image path was '/images/generations/audi-a3-8v-sedan.png'
UPDATE generations 
SET image_url = '/images/generations/audi-a3-8v-sedan.png'
WHERE code = '8V2' 
AND model_id = (SELECT id FROM models WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi'));

-- 3. Update 8V1 (Pre-Facelift) to use the NEW uploaded image
-- The new image was saved as 'audi-a3-8v1.png'
UPDATE generations 
SET image_url = '/images/generations/audi-a3-8v1.png'
WHERE code = '8V1' 
AND model_id = (SELECT id FROM models WHERE name = 'A3' AND brand_id = (SELECT id FROM brands WHERE name = 'Audi'));

COMMIT;
