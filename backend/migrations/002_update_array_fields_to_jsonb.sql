-- Migration: Update array fields from text[] to jsonb
-- This fixes PostgreSQL array handling issues by using jsonb instead

-- Update transactions table
DO $$
BEGIN
    -- Check if tags column is text[] and convert to jsonb
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'transactions'
        AND column_name = 'tags'
        AND data_type = 'ARRAY'
    ) THEN
        -- Drop the old column and create new one as jsonb
        ALTER TABLE transactions DROP COLUMN IF EXISTS tags;
        ALTER TABLE transactions ADD COLUMN tags jsonb DEFAULT '[]'::jsonb;
    END IF;

    -- Check if attachments column is text[] and convert to jsonb
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'transactions'
        AND column_name = 'attachments'
        AND data_type = 'ARRAY'
    ) THEN
        -- Drop the old column and create new one as jsonb
        ALTER TABLE transactions DROP COLUMN IF EXISTS attachments;
        ALTER TABLE transactions ADD COLUMN attachments jsonb DEFAULT '[]'::jsonb;
    END IF;
END $$;

-- Update credit_card_transactions table
DO $$
BEGIN
    -- Check if tags column exists and is text[]
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'credit_card_transactions'
        AND column_name = 'tags'
        AND data_type = 'ARRAY'
    ) THEN
        -- Drop the old column and create new one as jsonb
        ALTER TABLE credit_card_transactions DROP COLUMN IF EXISTS tags;
        ALTER TABLE credit_card_transactions ADD COLUMN tags jsonb DEFAULT '[]'::jsonb;
    END IF;

    -- Check if attachments column exists and is text[]
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'credit_card_transactions'
        AND column_name = 'attachments'
        AND data_type = 'ARRAY'
    ) THEN
        -- Drop the old column and create new one as jsonb
        ALTER TABLE credit_card_transactions DROP COLUMN IF EXISTS attachments;
        ALTER TABLE credit_card_transactions ADD COLUMN attachments jsonb DEFAULT '[]'::jsonb;
    END IF;
END $$;

-- Verify the changes
SELECT
    table_name,
    column_name,
    data_type
FROM information_schema.columns
WHERE table_name IN ('transactions', 'credit_card_transactions')
AND column_name IN ('tags', 'attachments')
ORDER BY table_name, column_name;
