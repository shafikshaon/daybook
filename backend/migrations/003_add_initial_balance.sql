-- Migration: Add initial_balance column to accounts table
-- This separates opening balance (fixed) from current balance (changes with transactions)

-- Step 1: Add initial_balance column
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS initial_balance DOUBLE PRECISION DEFAULT 0;

-- Step 2: Copy current balance to initial_balance for existing accounts
-- This preserves the opening balance at the time of migration
UPDATE accounts
SET initial_balance = balance
WHERE initial_balance = 0 OR initial_balance IS NULL;

-- Step 3: Verify the changes
SELECT
    id,
    name,
    initial_balance as "Opening Balance",
    balance as "Current Balance",
    (balance - initial_balance) as "Net Change"
FROM accounts
ORDER BY created_at
LIMIT 10;

-- Summary
SELECT
    COUNT(*) as total_accounts,
    SUM(initial_balance) as total_initial_balance,
    SUM(balance) as total_current_balance,
    SUM(balance - initial_balance) as total_net_change
FROM accounts;
