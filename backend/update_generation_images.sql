-- Update generation image URLs
UPDATE generations SET image_url = '/images/generations/audi-a3-8p-sportback.png' WHERE code = '8P' AND model_id = 1;
UPDATE generations SET image_url = '/images/generations/audi-a3-8v-sedan.png' WHERE code = '8V' AND model_id = 1;
UPDATE generations SET image_url = '/images/generations/audi-a3-8y-sportback.png' WHERE code = '8Y' AND model_id = 1;

-- Verify
SELECT id, code, name, image_url FROM generations WHERE model_id = 1 ORDER BY start_year DESC;
