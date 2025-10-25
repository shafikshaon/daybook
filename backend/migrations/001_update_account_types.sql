-- Migration to update account_types table structure
-- This removes the category and parent_id columns

-- Drop the existing table (this will delete all existing account types)
-- Users will get new account types when they next login or you can trigger a re-seed
DROP TABLE IF EXISTS account_types CASCADE;

-- The table will be recreated automatically by GORM AutoMigrate with the new structure
-- New structure:
-- - id (uuid)
-- - user_id (uuid, not null)
-- - name (string, not null)
-- - icon (string)
-- - description (string)
-- - active (boolean, default true)
-- - sort_order (integer, default 0)
-- - created_at (timestamp)
-- - updated_at (timestamp)
-- - deleted_at (timestamp)
