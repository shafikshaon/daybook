-- Migration: Add transaction_id column to credit_card_transactions table
-- This links credit card transactions to the main transactions table
-- so credit card expenses appear in the main transaction list

-- Step 1: Add transaction_id column
ALTER TABLE credit_card_transactions ADD COLUMN IF NOT EXISTS transaction_id UUID;

-- Step 2: Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_credit_card_transactions_transaction_id
ON credit_card_transactions(transaction_id);

-- Step 3: Verify the changes
SELECT
    table_name,
    column_name,
    data_type
FROM information_schema.columns
WHERE table_name = 'credit_card_transactions'
AND column_name = 'transaction_id';

-- Note: Existing credit card transactions will have NULL transaction_id
-- Only new transactions created after this migration will have the link
