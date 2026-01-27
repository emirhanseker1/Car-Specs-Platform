-- Add image_url column to trims table if it doesn't exist
ALTER TABLE trims ADD COLUMN image_url TEXT;
