-- ROLLBACK Migration: Revert TIMESTAMPTZ back to TIMESTAMP (if needed)
-- WARNING: This will lose timezone information!

-- ============================================================================
-- TRANSACTIONS TABLE
-- ============================================================================
ALTER TABLE transactions
    ALTER COLUMN date TYPE TIMESTAMP USING date::timestamp;

-- ============================================================================
-- RECURRING TRANSACTIONS TABLE
-- ============================================================================
ALTER TABLE recurring_transactions
    ALTER COLUMN template_date TYPE TIMESTAMP USING template_date::timestamp,
    ALTER COLUMN last_processed TYPE TIMESTAMP USING last_processed::timestamp;

-- ============================================================================
-- CREDIT CARDS TABLE
-- ============================================================================
ALTER TABLE credit_cards
    ALTER COLUMN due_date TYPE TIMESTAMP USING due_date::timestamp,
    ALTER COLUMN statement_date TYPE TIMESTAMP USING statement_date::timestamp,
    ALTER COLUMN last_payment_date TYPE TIMESTAMP USING last_payment_date::timestamp;

-- ============================================================================
-- Continue for all other tables...
-- ============================================================================
