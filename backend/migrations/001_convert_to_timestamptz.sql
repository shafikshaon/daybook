-- Migration: Convert all date/timestamp columns to TIMESTAMPTZ (timestamp with timezone)
-- This ensures all timestamps are stored in UTC consistently across the application

-- ============================================================================
-- TRANSACTIONS TABLE
-- ============================================================================
-- Convert date column to TIMESTAMPTZ
ALTER TABLE transactions
    ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'UTC';

-- Update any date-only values to have midnight UTC time component
UPDATE transactions
SET date = date::date::timestamp AT TIME ZONE 'UTC'
WHERE EXTRACT(HOUR FROM date) = 0
  AND EXTRACT(MINUTE FROM date) = 0
  AND EXTRACT(SECOND FROM date) = 0;

-- ============================================================================
-- RECURRING TRANSACTIONS TABLE
-- ============================================================================
ALTER TABLE recurring_transactions
    ALTER COLUMN template_date TYPE TIMESTAMPTZ USING template_date AT TIME ZONE 'UTC';

ALTER TABLE recurring_transactions
    ALTER COLUMN last_processed TYPE TIMESTAMPTZ USING last_processed AT TIME ZONE 'UTC';

-- ============================================================================
-- CREDIT CARDS TABLE
-- ============================================================================
ALTER TABLE credit_cards
    ALTER COLUMN due_date TYPE TIMESTAMPTZ USING due_date AT TIME ZONE 'UTC',
    ALTER COLUMN statement_date TYPE TIMESTAMPTZ USING statement_date AT TIME ZONE 'UTC',
    ALTER COLUMN last_payment_date TYPE TIMESTAMPTZ USING last_payment_date AT TIME ZONE 'UTC';

-- ============================================================================
-- CREDIT CARD TRANSACTIONS TABLE
-- ============================================================================
ALTER TABLE credit_card_transactions
    ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'UTC';

-- ============================================================================
-- CREDIT CARD PAYMENTS TABLE
-- ============================================================================
ALTER TABLE credit_card_payments
    ALTER COLUMN payment_date TYPE TIMESTAMPTZ USING payment_date AT TIME ZONE 'UTC';

-- ============================================================================
-- STATEMENTS TABLE
-- ============================================================================
ALTER TABLE statements
    ALTER COLUMN statement_date TYPE TIMESTAMPTZ USING statement_date AT TIME ZONE 'UTC',
    ALTER COLUMN due_date TYPE TIMESTAMPTZ USING due_date AT TIME ZONE 'UTC',
    ALTER COLUMN paid_date TYPE TIMESTAMPTZ USING paid_date AT TIME ZONE 'UTC';

-- ============================================================================
-- REWARDS TABLE
-- ============================================================================
ALTER TABLE rewards
    ALTER COLUMN earned_date TYPE TIMESTAMPTZ USING earned_date AT TIME ZONE 'UTC';

-- ============================================================================
-- BUDGETS TABLE
-- ============================================================================
ALTER TABLE budgets
    ALTER COLUMN custom_start_date TYPE TIMESTAMPTZ USING custom_start_date AT TIME ZONE 'UTC',
    ALTER COLUMN custom_end_date TYPE TIMESTAMPTZ USING custom_end_date AT TIME ZONE 'UTC';

-- ============================================================================
-- RECONCILIATIONS TABLE
-- ============================================================================
ALTER TABLE reconciliations
    ALTER COLUMN reconciliation_date TYPE TIMESTAMPTZ USING reconciliation_date AT TIME ZONE 'UTC';

-- ============================================================================
-- GOALS TABLE (formerly savings_goals)
-- ============================================================================
ALTER TABLE goals
    ALTER COLUMN target_date TYPE TIMESTAMPTZ USING target_date AT TIME ZONE 'UTC',
    ALTER COLUMN achieved_date TYPE TIMESTAMPTZ USING achieved_date AT TIME ZONE 'UTC',
    ALTER COLUMN last_contribution_date TYPE TIMESTAMPTZ USING last_contribution_date AT TIME ZONE 'UTC';

-- ============================================================================
-- GOAL HOLDINGS TABLE
-- ============================================================================
ALTER TABLE goal_holdings
    ALTER COLUMN purchase_date TYPE TIMESTAMPTZ USING purchase_date AT TIME ZONE 'UTC',
    ALTER COLUMN maturity_date TYPE TIMESTAMPTZ USING maturity_date AT TIME ZONE 'UTC';

-- ============================================================================
-- GOAL CONTRIBUTIONS TABLE
-- ============================================================================
ALTER TABLE goal_contributions
    ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'UTC';

-- ============================================================================
-- SAVINGS GOALS TABLE (if still exists)
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'savings_goals') THEN
        ALTER TABLE savings_goals
            ALTER COLUMN target_date TYPE TIMESTAMPTZ USING target_date AT TIME ZONE 'UTC',
            ALTER COLUMN achieved_date TYPE TIMESTAMPTZ USING achieved_date AT TIME ZONE 'UTC',
            ALTER COLUMN archived_date TYPE TIMESTAMPTZ USING archived_date AT TIME ZONE 'UTC',
            ALTER COLUMN last_contribution_date TYPE TIMESTAMPTZ USING last_contribution_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- SAVINGS CONTRIBUTIONS TABLE (if still exists)
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'savings_contributions') THEN
        ALTER TABLE savings_contributions
            ALTER COLUMN date TYPE TIMESTAMPTZ USING date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- FIXED DEPOSITS TABLE (if exists)
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'fixed_deposits') THEN
        ALTER TABLE fixed_deposits
            ALTER COLUMN start_date TYPE TIMESTAMPTZ USING start_date AT TIME ZONE 'UTC',
            ALTER COLUMN maturity_date TYPE TIMESTAMPTZ USING maturity_date AT TIME ZONE 'UTC',
            ALTER COLUMN withdrawn_date TYPE TIMESTAMPTZ USING withdrawn_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- INVESTMENTS TABLE (if exists)
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'investments') THEN
        ALTER TABLE investments
            ALTER COLUMN purchase_date TYPE TIMESTAMPTZ USING purchase_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- DEBT RECORDS TABLE
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'debt_records') THEN
        ALTER TABLE debt_records
            ALTER COLUMN due_date TYPE TIMESTAMPTZ USING due_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- DEBT PAYMENTS TABLE
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'debt_payments') THEN
        ALTER TABLE debt_payments
            ALTER COLUMN payment_date TYPE TIMESTAMPTZ USING payment_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- LEND RECORDS TABLE
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'lend_records') THEN
        ALTER TABLE lend_records
            ALTER COLUMN due_date TYPE TIMESTAMPTZ USING due_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- LEND PAYMENTS TABLE
-- ============================================================================
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'lend_payments') THEN
        ALTER TABLE lend_payments
            ALTER COLUMN payment_date TYPE TIMESTAMPTZ USING payment_date AT TIME ZONE 'UTC';
    END IF;
END $$;

-- ============================================================================
-- STANDARD TIMESTAMP COLUMNS (created_at, updated_at, deleted_at)
-- ============================================================================
-- These are usually already TIMESTAMPTZ from GORM, but ensure consistency

DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN
        SELECT table_name, column_name
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND column_name IN ('created_at', 'updated_at', 'deleted_at')
          AND data_type != 'timestamp with time zone'
    LOOP
        EXECUTE format('ALTER TABLE %I ALTER COLUMN %I TYPE TIMESTAMPTZ USING %I AT TIME ZONE ''UTC''',
                      r.table_name, r.column_name, r.column_name);
    END LOOP;
END $$;

-- ============================================================================
-- VERIFICATION QUERY
-- ============================================================================
-- Run this to verify all timestamp columns are now TIMESTAMPTZ:
--
-- SELECT table_name, column_name, data_type
-- FROM information_schema.columns
-- WHERE table_schema = 'public'
--   AND (column_name LIKE '%date%' OR column_name LIKE '%time%' OR data_type LIKE '%timestamp%')
-- ORDER BY table_name, column_name;
