-- Migration Rollback: Revert transaction date back to datetime
-- Created: 2025-11-29

BEGIN;

-- Step 1: Add a new timestamp column
ALTER TABLE transactions ADD COLUMN transaction_datetime_new TIMESTAMP;

-- Step 2: Copy date values and set time to midnight
UPDATE transactions 
SET transaction_datetime_new = transaction_datetime::TIMESTAMP
WHERE transaction_datetime IS NOT NULL;

-- Step 3: Drop the index
DROP INDEX IF EXISTS idx_transactions_datetime;

-- Step 4: Drop the date column
ALTER TABLE transactions DROP COLUMN transaction_datetime;

-- Step 5: Rename the new column back to original name
ALTER TABLE transactions RENAME COLUMN transaction_datetime_new TO transaction_datetime;

-- Step 6: Recreate index
CREATE INDEX IF NOT EXISTS idx_transactions_datetime ON transactions(transaction_datetime);

COMMIT;
