# Credit Card Transaction & Payment Flow Documentation

## Overview
This document explains how credit card transactions (purchases, fees, interest) and payments work in the Daybook application and where they are reflected across the system.

## Important: Dual Recording System

**All credit card activity appears in BOTH places:**
1. **Main Transactions List** - For complete financial overview
2. **Credit Card History** - For card-specific tracking

This ensures credit card expenses are included in:
- Budget calculations
- Category-wise spending analysis
- Overall transaction history
- Monthly expense reports

## Credit Card Transaction Flow (Purchases, Fees, Interest)

When you record a credit card transaction (purchase, fee, interest, or refund):

### 1. **Backend Processing** (`RecordCreditCardTransaction` handler)

```
User records expense: $150 grocery purchase on "Chase Visa"
```

**Step 1: Validate Transaction**
- ✅ Check credit card exists and belongs to user
- ✅ Validate transaction type (purchase, fee, interest, refund)

**Step 2: Create Main Transaction Record**
```go
Transaction {
    UserID:       user123
    CreditCardID: chase_visa_id       // Links to credit card
    AccountID:    nil                 // No account for credit card purchases
    CategoryID:   "groceries"
    Amount:       150.00
    Type:         "expense"           // "income" for refunds
    Date:         "2025-10-25"
    Description:  "Whole Foods"
    Tags:         ["groceries", "food"]
}
```

**Step 3: Create Credit Card Transaction Record**
```go
CreditCardTransaction {
    CardID:        chase_visa_id
    TransactionID: main_transaction_id  // Links to main transaction
    Amount:        150.00
    Type:          "purchase"           // purchase, fee, interest, refund
    Date:          "2025-10-25"
    Description:   "Whole Foods"
    Merchant:      "Whole Foods"
    CategoryID:    "groceries"
}
```

**Step 4: Update Credit Card Balance**
```
Chase Visa: $700 → $850 (+$150)
```

### 2. **Where Transactions Are Reflected**

#### A. **Transactions View** (`/transactions`)
Shows in main transaction list:
```
Date         Description    Category    Account          Type      Amount
2025-10-25   Whole Foods    Groceries   💳 Chase Visa    Expense   -$150
```

**Features:**
- Shows ALL credit card purchases
- Displays credit card icon 💳 with card name
- Includes in category totals
- Filterable and searchable
- Included in budget tracking

#### B. **Credit Cards View** (`/credit-cards`)

**Card Overview:**
```
Chase Visa
Balance: $850 (was $700)
Available: $4150 / $5000
```

**Transaction History (View Details):**
```
Purchase   $150    2025-10-25    Whole Foods    Groceries
```

## Payment Flow

When you pay a credit card bill from your bank account/wallet, the following happens:

### 1. **Backend Processing** (`RecordPayment` handler)

```
User initiates payment: Pay $500 to "Chase Visa" from "Main Checking"
```

**Step 1: Validate Payment**
- ✅ Check credit card exists and belongs to user
- ✅ Check payment account exists and belongs to user
- ✅ Validate payment amount ≤ credit card balance
- ✅ Validate payment account has sufficient funds

**Step 2: Create Transaction Record**
```go
Transaction {
    UserID:       user123
    AccountID:    main_checking_id
    CreditCardID: chase_visa_id       // Links to credit card
    CategoryID:   "credit_card_payment"
    Amount:       500.00
    Type:         "expense"
    Date:         "2025-10-25"
    Description:  "Credit card payment: Chase Visa"
    Tags:         ["credit_card_payment"]
}
```

**Step 3: Update Account Balance**
```
Main Checking: $2000 → $1500 (-$500)
```

**Step 4: Update Credit Card Balance**
```
Chase Visa: $1200 → $700 (-$500)
Card Last Payment: $500
Card Last Payment Date: 2025-10-25
```

**Step 5: Create Payment Record**
```go
CreditCardPayment {
    CardID:        chase_visa_id
    AccountID:     main_checking_id
    Amount:        500.00
    PaymentDate:   "2025-10-25"
    TransactionID: transaction_id    // Links to main transaction
}
```

**Step 6: Create Credit Card Transaction**
```go
CreditCardTransaction {
    CardID:      chase_visa_id
    Amount:      500.00
    Type:        "payment"
    Date:        "2025-10-25"
    Description: "Credit card payment: Chase Visa"
}
```

### 2. **Where Payments Are Reflected**

#### A. **Transactions View** (`/transactions`)
Shows in main transaction list:
```
Date         Description                      Category              Account          Type      Amount
2025-10-25   Credit card payment: Chase Visa  Credit Card Payment   Main Checking    Expense   -$500
             💳 Chase Visa
```

**Features:**
- Shows expense from your bank account
- Displays credit card icon 💳 with card name
- Tagged as "credit_card_payment"
- Filterable by category

#### B. **Accounts View** (`/accounts`)
In "Main Checking" account details:
```
Recent Transactions:
- Credit card payment: Chase Visa    -$500    (2025-10-25)
  💳 Chase Visa
```

**Impact:**
- Account balance decreases
- Shows as expense transaction
- Linked to credit card via creditCardId field

#### C. **Credit Cards View** (`/credit-cards`)

**Card Overview:**
```
Chase Visa
Balance: $700 (was $1200)
Available: $4300 / $5000
Last Payment: $500 on 2025-10-25
```

**Transaction History (View Details):**
```
Payment    $500    2025-10-25    -$500
```

**Payment History:**
```
Payment from Main Checking    $500    2025-10-25
```

#### D. **Dashboard** (if applicable)
- Updates total account balances
- Updates credit card outstanding balances
- Shows in recent transactions

## Data Flow Diagram

```
┌─────────────────┐
│  User Action    │
│  Pay $500 to    │
│  Chase Visa     │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  Backend Transaction (Atomic)           │
├─────────────────────────────────────────┤
│  1. Create Transaction Record           │
│     - Links to both account and card    │
│     - Type: expense                     │
│     - Category: credit_card_payment     │
│                                         │
│  2. Update Account Balance              │
│     Main Checking: -$500                │
│                                         │
│  3. Update Credit Card Balance          │
│     Chase Visa: -$500                   │
│                                         │
│  4. Create CreditCardPayment Record     │
│     - Links transaction to card payment │
│                                         │
│  5. Create CreditCardTransaction        │
│     - Shows in card's transaction log   │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  Frontend Updates                       │
├─────────────────────────────────────────┤
│  ✓ Transactions list shows payment      │
│  ✓ Account balance updated              │
│  ✓ Credit card balance updated          │
│  ✓ Payment history updated              │
└─────────────────────────────────────────┘
```

## Example Scenarios

### Scenario 1: Full Payment
```
Credit Card Balance: $1000
Payment Amount: $1000
From Account: Checking ($5000)

Results:
- Checking: $5000 → $4000
- Credit Card: $1000 → $0
- Transaction created (expense, $1000)
```

### Scenario 2: Partial Payment
```
Credit Card Balance: $1000
Payment Amount: $300
From Account: Savings ($2000)

Results:
- Savings: $2000 → $1700
- Credit Card: $1000 → $700
- Transaction created (expense, $300)
```

### Scenario 3: Multiple Payments
```
Payment 1: $500 from Checking
Payment 2: $300 from Savings

Results:
- Two separate transaction records
- Two payment records
- Credit card balance reduced by $800 total
```

## API Endpoints

### Make Payment
```
POST /api/v1/credit-cards/:cardId/payment

Request:
{
  "amount": 500,
  "accountId": "account-uuid",
  "paymentDate": "2025-10-25T10:00:00Z",
  "description": "Optional description"
}

Response:
{
  "success": true,
  "data": {
    "card": { /* updated card */ },
    "payment": { /* payment record */ }
  }
}
```

### Get Payment History
```
GET /api/v1/credit-cards/:cardId/payments

Response:
{
  "success": true,
  "data": [
    {
      "id": "payment-uuid",
      "cardId": "card-uuid",
      "accountId": "account-uuid",
      "amount": 500,
      "paymentDate": "2025-10-25",
      "transactionId": "transaction-uuid"
    }
  ]
}
```

## Frontend Components

### TransactionsView.vue
- Shows all transactions including credit card payments
- Displays credit card icon and name for linked payments
- Allows filtering by category "credit_card_payment"

### CreditCardsView.vue
- "Pay Bill" button opens payment modal
- Select source account from dropdown
- Validates payment amount
- Shows payment confirmation
- Displays transaction history

### AccountsView.vue
- Shows credit card payments as expenses
- Links to credit card in transaction details

## Database Schema

### transactions table
```sql
- id (uuid)
- user_id (uuid)
- account_id (uuid)          -- Source account
- credit_card_id (uuid)      -- Link to credit card (nullable)
- type (text)                -- 'expense' for payments
- amount (decimal)
- category_id (text)         -- 'credit_card_payment'
- date (timestamp)
- description (text)
- tags (text[])              -- ['credit_card_payment']
```

### credit_card_payments table
```sql
- id (uuid)
- user_id (uuid)
- card_id (uuid)
- account_id (uuid)          -- Payment source
- amount (decimal)
- payment_date (timestamp)
- description (text)
- transaction_id (uuid)      -- Links to transactions table
```

### credit_card_transactions table
```sql
- id (uuid)
- user_id (uuid)
- card_id (uuid)
- transaction_id (uuid)      -- Links to main transactions table
- type (text)                -- 'payment', 'purchase', 'fee', 'interest', 'refund'
- amount (decimal)
- date (timestamp)
- description (text)
- merchant (text)
- category_id (text)
- tags (jsonb)
- attachments (jsonb)
```

**Important:** The `transaction_id` field links credit card transactions to the main transactions table, ensuring all credit card activity appears in the main transaction list.

## Key Benefits

1. **Complete Financial Overview**: All credit card activity (purchases AND payments) appears in main transactions list
2. **Budget Tracking**: Credit card expenses automatically included in budget calculations
3. **Category Analysis**: Credit card purchases counted in category-wise spending reports
4. **Complete Audit Trail**: Every transaction and payment tracked in multiple places
5. **Account Integration**: Payments automatically update account balances
6. **Credit Card Tracking**: Card balances stay accurate with all activity
7. **Transaction History**: Full visibility in multiple views
8. **Data Integrity**: Atomic transactions ensure consistency across all tables
9. **Comprehensive Reporting**: Easy to generate complete financial reports including all credit card activity

## Usage Tips

### For Credit Card Transactions (Purchases)
1. All credit card purchases automatically appear in main transactions list
2. Use categories to organize credit card expenses (groceries, dining, etc.)
3. Add merchant names for better tracking
4. Use tags to add additional context
5. Refunds appear as income transactions

### For Credit Card Payments
1. Always ensure payment account has sufficient funds
2. Payment amount cannot exceed credit card balance
3. Payments are immediately reflected across all views
4. Use payment history to track when/where payments were made
5. Transaction tags make it easy to filter/report on payments

### General
1. Check main transactions view to see ALL financial activity in one place
2. Use credit card detail view for card-specific history
3. Deleting a credit card transaction removes it from both places (main transactions + card history)
4. All operations are atomic - either all updates succeed or all fail (no partial updates)
