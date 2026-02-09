-- ============================================
-- MIGRATION: Add missing image_url column to generations table
-- ============================================

-- Check if column exists before adding (SQLite doesn't support IF NOT EXISTS for columns)
-- This will add the column if it doesn't exist

ALTER TABLE generations ADD COLUMN image_url TEXT;
