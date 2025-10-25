# Account Balance System Documentation

## Overview

The Daybook application uses a **dual-balance accounting system** that separates:
- **Initial Balance** (Opening Balance) - Fixed, never changes
- **Current Balance** - Updated automatically with every transaction

This matches standard accounting practices where you track:
```
Current Balance = Initial Balance + Income - Expenses + Transfers In - Transfers Out
```

## Balance Fields

### Initial Balance (`initialBalance`)
- **Purpose**: Records the starting balance when the account is created
- **Behavior**:
  - Set once when account is created
  - **Never modified** by transactions
  - Cannot be changed through the UI after creation
  - Represents the opening balance for accounting periods

### Current Balance (`balance`)
- **Purpose**: Shows the real-time balance of the account
- **Behavior**:
  - Starts equal to initial balance
  - **Automatically updated** with every transaction:
    - `+` Income transactions
    - `-` Expense transactions
    - `+` Transfers IN
    - `-` Transfers OUT
    - `-` Credit card payments
  - Cannot be manually edited
  - Calculated automatically by the system

## Examples

### Example 1: New Bank Account
```
Day 1: Create account with $1000
  Initial Balance: $1000
  Current Balance: $1000
  Net Change: $0

Day 2: Receive salary +$3000
  Initial Balance: $1000 (unchanged)
  Current Balance: $4000 (1000 + 3000)
  Net Change: +$3000

Day 3: Pay rent -$1500
  Initial Balance: $1000 (unchanged)
  Current Balance: $2500 (4000 - 1500)
  Net Change: +$1500

Day 4: Buy groceries -$200
  Initial Balance: $1000 (unchanged)
  Current Balance: $2300 (2500 - 200)
  Net Change: +$1300
```

### Example 2: Zero Initial Balance
```
Day 1: Create new account with $0
  Initial Balance: $0
  Current Balance: $0
  Net Change: $0

Day 2: Deposit $500
  Initial Balance: $0 (unchanged)
  Current Balance: $500 (0 + 500)
  Net Change: +$500

Day 3: Spend $100
  Initial Balance: $0 (unchanged)
  Current Balance: $400 (500 - 100)
  Net Change: +$400
```

### Example 3: Credit Card Payment
```
Checking Account Initial Balance: $5000

Pay credit card bill $800:
  Initial Balance: $5000 (unchanged)
  Current Balance: $4200 (5000 - 800)
  Net Change: -$800

Transaction Type: Expense
Category: Credit Card Payment
Linked to: Credit Card
```

## Transaction Impact on Balances

### Income Transaction
```go
account.Balance += transactionAmount
// initialBalance remains unchanged
```

### Expense Transaction
```go
account.Balance -= transactionAmount
// initialBalance remains unchanged
```

### Transfer OUT
```go
fromAccount.Balance -= transferAmount
// initialBalance remains unchanged
```

### Transfer IN
```go
toAccount.Balance += transferAmount
// initialBalance remains unchanged
```

### Credit Card Payment
```go
// From account (bank)
paymentAccount.Balance -= paymentAmount
// initialBalance remains unchanged

// Credit card
creditCard.CurrentBalance -= paymentAmount
// Credit card doesn't have initialBalance
```

## Database Schema

```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    initial_balance DOUBLE PRECISION DEFAULT 0,  -- Opening balance
    balance DOUBLE PRECISION DEFAULT 0,          -- Current balance
    currency VARCHAR DEFAULT 'USD',
    description TEXT,
    institution VARCHAR,
    account_number VARCHAR,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## Backend Implementation

### Account Creation
```go
func (a *Account) BeforeCreate(tx *gorm.DB) error {
    if a.ID == uuid.Nil {
        a.ID = uuid.New()
    }
    // Set initial balance to provided balance on creation
    if a.InitialBalance == 0 && a.Balance != 0 {
        a.InitialBalance = a.Balance
    }
    return nil
}
```

### Account Update Handler
```go
// Update only allowed fields - Balance and InitialBalance are NOT included
existingAccount.Name = updateData.Name
existingAccount.Type = updateData.Type
existingAccount.Currency = updateData.Currency
existingAccount.Description = updateData.Description
existingAccount.Institution = updateData.Institution
existingAccount.AccountNumber = updateData.AccountNumber
existingAccount.Active = updateData.Active
// Balance is NEVER directly updated through API
```

### Transaction Creation
```go
// When creating an income transaction
account.Balance += transaction.Amount
// initialBalance is not touched

// When creating an expense transaction
account.Balance -= transaction.Amount
// initialBalance is not touched
```

## Frontend Display

### Accounts List
Shows current balance with net change indicator:
```
Account Name    Type       Balance
Main Checking   Bank       $2,300
                           (+$1,300)
```

Hover tooltip shows:
```
Initial: $1,000
Current: $2,300
Change: +$1,300
```

### Account Creation Form
```
Initial Balance (Opening Balance) *
[1000.00]
Starting balance when creating the account. This will not change.
```

### Account Edit Form
```
Initial Balance (Opening Balance) *
[1000.00] (disabled)
Initial balance cannot be changed. Current balance updates with transactions.
```

## API Behavior

### Create Account
```json
POST /api/v1/accounts
{
  "name": "Main Checking",
  "type": "bank",
  "balance": 1000.00  // This becomes both initial and current balance
}

Response:
{
  "id": "...",
  "initialBalance": 1000.00,
  "balance": 1000.00
}
```

### Update Account
```json
PUT /api/v1/accounts/:id
{
  "name": "Updated Name",
  "description": "New description"
  // balance and initialBalance are ignored even if sent
}
```

### Get Account
```json
GET /api/v1/accounts/:id

Response:
{
  "id": "...",
  "name": "Main Checking",
  "type": "bank",
  "initialBalance": 1000.00,  // Opening balance
  "balance": 2300.00,          // Current balance
  "currency": "USD"
}
```

## Calculating Net Change

```javascript
// Frontend calculation
const netChange = account.balance - account.initialBalance

// Display
if (netChange > 0) {
  // Show as: +$1,300 (in green)
} else if (netChange < 0) {
  // Show as: -$500 (in red)
} else {
  // No change indicator
}
```

## Migration Guide

### Adding InitialBalance to Existing System

Run the migration script:
```sql
-- Add column
ALTER TABLE accounts ADD COLUMN initial_balance DOUBLE PRECISION DEFAULT 0;

-- Copy current balance to initial balance for existing accounts
UPDATE accounts
SET initial_balance = balance
WHERE initial_balance = 0 OR initial_balance IS NULL;
```

This preserves the current balance as the opening balance for existing accounts.

## Benefits

1. **Accurate Accounting**: Separates opening balance from current state
2. **Period Reporting**: Easy to calculate net change over time
3. **Transaction Audit**: Current balance is always derived from transactions
4. **Data Integrity**: Initial balance cannot be accidentally modified
5. **Performance**: No need to sum all transactions to get current balance
6. **Reconciliation**: Can compare initial + transactions = current

## Common Questions

### Q: Why can't I change the initial balance?
**A**: The initial balance represents your opening balance at a point in time. Changing it would invalidate all transaction history and calculations. If you need to adjust, create an adjustment transaction instead.

### Q: What if I entered the wrong initial balance?
**A**: If the account is newly created with no transactions, delete and recreate it. If transactions exist, create an adjustment transaction to correct the balance.

### Q: Can current balance be less than initial balance?
**A**: Yes! If you have more expenses than income, your current balance can go below (or even negative compared to) the initial balance. This shows you've spent more than you started with.

### Q: How do I reset my balance to initial?
**A**: You cannot reset balances. The current balance reflects all transactions. To adjust, you would need to delete transactions or create offsetting transactions.

### Q: What happens to balance when I delete a transaction?
**A**: The system automatically adjusts the current balance by reversing the transaction's effect. Initial balance remains unchanged.

## Best Practices

1. **Set Initial Balance Carefully**: Double-check the amount when creating accounts
2. **Use Transactions**: Always record income/expenses as transactions, never try to manually adjust balance
3. **Reconcile Regularly**: Compare your current balance with bank statements
4. **Document Adjustments**: If you need to make corrections, create transactions with clear descriptions
5. **Monitor Net Change**: Use the net change indicator to track how your accounts are performing

## Technical Notes

- Both balances are stored as `DOUBLE PRECISION` in PostgreSQL
- Balance updates use database transactions to ensure atomicity
- No race conditions: Transaction handlers use database-level locking
- Historical balance can be calculated: `InitialBalance + SUM(transactions up to date)`
