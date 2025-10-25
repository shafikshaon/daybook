-- Migration: Fix array columns by converting text[] to jsonb
-- Run this BEFORE restarting the backend server

-- Fix transactions table
ALTER TABLE transactions DROP COLUMN IF EXISTS tags CASCADE;
ALTER TABLE transactions ADD COLUMN tags jsonb DEFAULT '[]'::jsonb;

ALTER TABLE transactions DROP COLUMN IF EXISTS attachments CASCADE;
ALTER TABLE transactions ADD COLUMN attachments jsonb DEFAULT '[]'::jsonb;

-- Note: If you have existing data you want to preserve, uncomment and use these instead:
-- ALTER TABLE transactions ALTER COLUMN tags TYPE jsonb USING array_to_json(tags)::jsonb;
-- ALTER TABLE transactions ALTER COLUMN attachments TYPE jsonb USING array_to_json(attachments)::jsonb;

-- Verify the changes
SELECT
    table_name,
    column_name,
    data_type
FROM information_schema.columns
WHERE table_name = 'transactions'
AND column_name IN ('tags', 'attachments');
