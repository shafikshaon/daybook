# Database Migrations - Daybook Backend

This directory contains database migrations for the Daybook application.

## Migration 001: Convert to TIMESTAMPTZ

This migration converts all date and timestamp columns from `TIMESTAMP` to `TIMESTAMPTZ` (timestamp with timezone) to ensure consistent UTC storage across the entire application.

### Why This Migration?

**Problem**: Timezone inconsistencies causing:
- Wrong dates when creating transactions
- Incorrect budget calculations
- Reports showing wrong data
- Transactions missing from date range queries

**Solution**: Use `TIMESTAMPTZ` for all date/time fields to store everything in UTC consistently.

---

## Running the Migration

### Prerequisites

1. **Backup your database first!**
   ```bash
   # Set your database URL
   export DATABASE_URL="postgresql://user:password@host:port/database"

   # Create backup
   pg_dump $DATABASE_URL > daybook_backup_$(date +%Y%m%d_%H%M%S).sql
   ```

2. **Ensure PostgreSQL is running and accessible**

3. **Stop your backend server** (to avoid conflicts during migration)

### Option 1: Using the Migration Script (Recommended)

```bash
cd backend/migrations

# Run the migration script
./run_migration.sh
```

The script will:
- Check for DATABASE_URL environment variable
- Prompt you to confirm you've backed up your database
- Run the migration
- Verify the changes
- Show you a summary of all timestamp columns

### Option 2: Manual Migration

```bash
# Set your database URL
export DATABASE_URL="postgresql://user:password@host:port/database"

# Run the migration SQL directly
psql $DATABASE_URL -f 001_convert_to_timestamptz.sql

# Verify the changes
psql $DATABASE_URL -c "
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = 'public'
  AND (column_name LIKE '%date%' OR column_name LIKE '%time%')
ORDER BY table_name, column_name;
"
```

---

## What This Migration Does

1. **Converts all date/time columns to TIMESTAMPTZ:**
   - Transactions table: `date` field
   - Credit cards: `due_date`, `statement_date`, `last_payment_date`
   - Budgets: `custom_start_date`, `custom_end_date`
   - Goals: `target_date`, `achieved_date`, `last_contribution_date`
   - And many more...

2. **Updates existing data:**
   - All dates stored with time components are preserved
   - Date-only values are converted to midnight UTC
   - All timestamps are normalized to UTC

3. **Ensures consistency:**
   - All `created_at`, `updated_at`, `deleted_at` columns are TIMESTAMPTZ
   - Future queries will automatically handle timezone conversions

---

## After Running the Migration

1. **Restart your backend server:**
   ```bash
   cd backend
   go run main.go
   # or if using a process manager
   systemctl restart daybook-backend
   ```

2. **Verify functionality:**
   - Create a new transaction and check the date is correct
   - View reports and verify totals match
   - Check budgets show correct spent amounts

3. **Monitor logs** for any timezone-related warnings or errors

---

## Rollback (If Needed)

⚠️ **WARNING**: Rolling back will lose timezone information!

If you need to rollback:

```bash
# Option 1: Restore from backup (RECOMMENDED)
psql $DATABASE_URL < your_backup_file.sql

# Option 2: Run rollback migration (loses timezone data)
psql $DATABASE_URL -f 000_rollback_timestamptz.sql
```

---

## Technical Details

### PostgreSQL TIMESTAMPTZ vs TIMESTAMP

- **TIMESTAMP**: Stores date and time without timezone info
  - `2024-11-23 14:30:00` (ambiguous - what timezone?)

- **TIMESTAMPTZ**: Stores date and time with timezone (always in UTC)
  - `2024-11-23 14:30:00+00` (unambiguous - UTC)

### How Go/GORM Handles It

```go
// Before (could be any timezone)
Date time.Time `gorm:"not null;index" json:"date"`

// After (explicitly UTC)
Date time.Time `gorm:"type:timestamptz;not null;index" json:"date"`
```

When you use `time.Time` in Go:
- Values are stored in UTC in the database
- When retrieved, Go automatically converts to your local timezone
- Comparisons and calculations are timezone-aware

### Model Updates Required

All models with date/time fields have been updated to include `type:timestamptz`:

**Updated models:**
- `Transaction`
- `CreditCard`, `Statement`, `Reward`, `CreditCardTransaction`, `CreditCardPayment`
- `Budget`
- `Reconciliation`
- `Goal`, `GoalHolding`, `GoalContribution`
- `DebtRecord`, `DebtPayment`
- `LendRecord`, `LendPayment`
- `Asset`, `ServiceRecord`
- And more...

---

## Verification

After migration, verify all timestamp columns:

```sql
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = 'public'
  AND data_type LIKE '%timestamp%'
ORDER BY table_name, column_name;
```

All should show `timestamp with time zone` or `timestamptz`.

---

## Troubleshooting

### Migration fails with "relation does not exist"

Some tables might not exist yet in your database. The migration uses `DO $$ ... END $$` blocks to conditionally alter tables only if they exist.

### Data looks different after migration

TIMESTAMPTZ automatically converts to your display timezone. The data is still correct in UTC, but your client might be showing it in a different timezone.

### Transactions still showing wrong dates

1. Clear your browser cache
2. Hard refresh (Ctrl+Shift+R or Cmd+Shift+R)
3. Restart your backend server
4. Check that the frontend is sending dates in UTC format

---

## Support

If you encounter issues:

1. Check the migration logs
2. Verify your database backup is valid
3. Restore from backup if needed
4. Contact the development team

---

## Future Migrations

When creating new migrations:

1. Name them sequentially: `002_description.sql`, `003_description.sql`, etc.
2. Always create a corresponding rollback file
3. Test on a development database first
4. Document what the migration does in this README
5. Update the `run_migration.sh` script to include the new migration

---

Last updated: 2024-11-23
