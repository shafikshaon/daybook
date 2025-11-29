-- Migration: Convert transaction datetime to date only
-- Created: 2025-11-29

BEGIN;

-- Step 1: Add a new date column
ALTER TABLE transactions ADD COLUMN transaction_date_new DATE;

-- Step 2: Copy date portion from existing datetime column to new date column
-- Adjust the column name based on your actual datetime column name
-- Common names: transaction_datetime, created_at, transaction_time, date
UPDATE transactions 
SET transaction_date_new = transaction_datetime::DATE
WHERE transaction_datetime IS NOT NULL;

-- Step 3: Drop the old datetime column
ALTER TABLE transactions DROP COLUMN transaction_datetime;

-- Step 4: Rename the new column to the original name
ALTER TABLE transactions RENAME COLUMN transaction_date_new TO transaction_datetime;

-- Step 5: Add NOT NULL constraint if needed (uncomment if required)
-- ALTER TABLE transactions ALTER COLUMN transaction_datetime SET NOT NULL;

-- Step 6: Add index on the date column for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_datetime ON transactions(transaction_datetime);

COMMIT;
