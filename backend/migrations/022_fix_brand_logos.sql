-- 022_fix_brand_logos.sql
BEGIN TRANSACTION;

UPDATE brands SET logo_url = '/images/brands/volkswagen.png' WHERE name = 'Volkswagen';
UPDATE brands SET logo_url = '/images/brands/bmw.png' WHERE name = 'BMW';

COMMIT;
