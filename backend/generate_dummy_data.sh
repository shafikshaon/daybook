#!/bin/bash

# Script to generate dummy data for a specific user in Daybook application
# Usage:
#   Generate data: ./generate_dummy_data.sh <user_email>
#   Clean data: ./generate_dummy_data.sh --clean <user_email>

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

# Check for clean flag
CLEAN_MODE=false
if [ "$1" = "--clean" ]; then
    CLEAN_MODE=true
    USER_EMAIL="$2"
else
    USER_EMAIL="$1"
fi

# Check if email parameter is provided
if [ -z "$USER_EMAIL" ]; then
    echo -e "${RED}Error: Email parameter is required${NC}"
    echo "Usage:"
    echo "  Generate data: $0 <user_email>"
    echo "  Clean data: $0 --clean <user_email>"
    exit 1
fi

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-daybook}"

# Set PGPASSWORD for psql
export PGPASSWORD="$DB_PASSWORD"

echo -e "${YELLOW}Checking if user exists with email: $USER_EMAIL${NC}"

# Query user table for the email
USER_DATA=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -A -c "SELECT id, username, full_name FROM users WHERE email = '$USER_EMAIL' AND deleted_at IS NULL;")

if [ -z "$USER_DATA" ]; then
    echo -e "${RED}Error: User with email $USER_EMAIL not found${NC}"
    exit 1
fi

# Extract user details
USER_ID=$(echo "$USER_DATA" | cut -d'|' -f1)
USERNAME=$(echo "$USER_DATA" | cut -d'|' -f2)
FULL_NAME=$(echo "$USER_DATA" | cut -d'|' -f3)

echo -e "${GREEN}User found:${NC}"
echo "  ID: $USER_ID"
echo "  Username: $USERNAME"
echo "  Full Name: $FULL_NAME"
echo ""

# If clean mode, delete all data except user table
if [ "$CLEAN_MODE" = true ]; then
    echo -e "${YELLOW}Cleaning all data for user: $USER_EMAIL${NC}"
    echo -e "${RED}WARNING: This will delete all data except the user record!${NC}"
    read -p "Are you sure you want to continue? (yes/no): " CONFIRM

    if [ "$CONFIRM" != "yes" ]; then
        echo -e "${YELLOW}Operation cancelled${NC}"
        exit 0
    fi

    echo -e "${YELLOW}Deleting all user data...${NC}"

    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Delete in order to respect foreign key constraints
DELETE FROM credit_card_transactions WHERE user_id = '$USER_ID';
DELETE FROM transactions WHERE user_id = '$USER_ID';
DELETE FROM bills WHERE user_id = '$USER_ID';
DELETE FROM budgets WHERE user_id = '$USER_ID';
DELETE FROM savings_goals WHERE user_id = '$USER_ID';
DELETE FROM credit_cards WHERE user_id = '$USER_ID';
DELETE FROM accounts WHERE user_id = '$USER_ID';
EOF

    echo -e "${GREEN}✓ All data cleaned for user: $USER_EMAIL${NC}"
    echo -e "${GREEN}User record preserved${NC}"

    # Clean up
    unset PGPASSWORD
    exit 0
fi

# Function to generate random UUID
generate_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || uuidgen
}

# Function to get random date in past N days
random_past_date() {
    local days=$1
    local random_days=$((RANDOM % days))
    date -u -d "$random_days days ago" +"%Y-%m-%d %H:%M:%S" 2>/dev/null || date -u -v-${random_days}d +"%Y-%m-%d %H:%M:%S"
}

# Function to get future date N days ahead
future_date() {
    local days=$1
    date -u -d "$days days" +"%Y-%m-%d %H:%M:%S" 2>/dev/null || date -u -v+${days}d +"%Y-%m-%d %H:%M:%S"
}

echo -e "${YELLOW}Starting data generation...${NC}"
echo ""

# ==================== ACCOUNTS ====================
echo -e "${YELLOW}Creating accounts...${NC}"

CHECKING_ACCOUNT_ID=$(generate_uuid)
SAVINGS_ACCOUNT_ID=$(generate_uuid)
CASH_ACCOUNT_ID=$(generate_uuid)

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Checking Account
INSERT INTO accounts (id, user_id, name, type, initial_balance, balance, currency, description, institution, account_number, active, created_at, updated_at)
VALUES (
    '$CHECKING_ACCOUNT_ID',
    '$USER_ID',
    'Main Checking Account',
    'checking',
    5000.00,
    5000.00,
    'BDT',
    'Primary checking account',
    'Chase Bank',
    '****1234',
    true,
    NOW(),
    NOW()
);

-- Savings Account
INSERT INTO accounts (id, user_id, name, type, initial_balance, balance, currency, description, institution, account_number, active, created_at, updated_at)
VALUES (
    '$SAVINGS_ACCOUNT_ID',
    '$USER_ID',
    'Emergency Savings',
    'savings',
    10000.00,
    10000.00,
    'BDT',
    'Emergency fund savings',
    'Chase Bank',
    '****5678',
    true,
    NOW(),
    NOW()
);

-- Cash Account
INSERT INTO accounts (id, user_id, name, type, initial_balance, balance, currency, description, active, created_at, updated_at)
VALUES (
    '$CASH_ACCOUNT_ID',
    '$USER_ID',
    'Cash Wallet',
    'cash',
    500.00,
    500.00,
    'USD',
    'Physical cash',
    true,
    NOW(),
    NOW()
);
EOF

echo -e "${GREEN}✓ Created 3 accounts${NC}"

# ==================== CREDIT CARDS ====================
echo -e "${YELLOW}Creating credit cards...${NC}"

CREDIT_CARD_1_ID=$(generate_uuid)
CREDIT_CARD_2_ID=$(generate_uuid)

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Visa Credit Card
INSERT INTO credit_cards (id, user_id, name, last_four_digits, card_network, credit_limit, current_balance, apr, due_date, statement_date, minimum_payment, rewards_program, active, created_at, updated_at)
VALUES (
    '$CREDIT_CARD_1_ID',
    '$USER_ID',
    'Chase Freedom Unlimited',
    '4532',
    'Visa',
    15000.00,
    2500.00,
    16.99,
    '$(future_date 15)',
    '$(future_date 5)',
    75.00,
    'Cash Back Rewards',
    true,
    NOW(),
    NOW()
);

-- Mastercard Credit Card
INSERT INTO credit_cards (id, user_id, name, last_four_digits, card_network, credit_limit, current_balance, apr, due_date, statement_date, minimum_payment, rewards_program, active, created_at, updated_at)
VALUES (
    '$CREDIT_CARD_2_ID',
    '$USER_ID',
    'Capital One Quicksilver',
    '8765',
    'Mastercard',
    10000.00,
    1200.00,
    18.49,
    '$(future_date 20)',
    '$(future_date 10)',
    35.00,
    'Cashback',
    true,
    NOW(),
    NOW()
);
EOF

echo -e "${GREEN}✓ Created 2 credit cards${NC}"

# ==================== TRANSACTIONS ====================
echo -e "${YELLOW}Creating transactions...${NC}"

# Categories for transactions
CATEGORIES=("groceries" "dining" "transportation" "utilities" "entertainment" "shopping" "healthcare" "salary" "investment")

# Generate 50 transactions
for i in {1..500}; do
    TRANS_ID=$(generate_uuid)
    TRANS_DATE=$(random_past_date 90)
    CATEGORY=${CATEGORIES[$RANDOM % ${#CATEGORIES[@]}]}

    # Determine transaction type and amount based on category
    if [ "$CATEGORY" = "salary" ]; then
        TRANS_TYPE="income"
        AMOUNT=$((3000 + RANDOM % 2000))
        ACCOUNT_ID="$CHECKING_ACCOUNT_ID"
    else
        TRANS_TYPE="expense"
        AMOUNT=$((10 + RANDOM % 500))
        # Randomly select account
        ACCOUNTS=("$CHECKING_ACCOUNT_ID" "$CASH_ACCOUNT_ID")
        ACCOUNT_ID=${ACCOUNTS[$RANDOM % ${#ACCOUNTS[@]}]}
    fi

    DESCRIPTIONS=("Payment for $CATEGORY" "Monthly $CATEGORY" "$CATEGORY expense" "Weekly $CATEGORY" "$CATEGORY purchase")
    DESCRIPTION=${DESCRIPTIONS[$RANDOM % ${#DESCRIPTIONS[@]}]}

    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF > /dev/null 2>&1
INSERT INTO transactions (id, user_id, account_id, type, amount, category_id, date, description, tags, created_at, updated_at)
VALUES (
    '$TRANS_ID',
    '$USER_ID',
    '$ACCOUNT_ID',
    '$TRANS_TYPE',
    $AMOUNT,
    '$CATEGORY',
    '$TRANS_DATE',
    '$DESCRIPTION',
    '[]'::jsonb,
    NOW(),
    NOW()
);
EOF
done

echo -e "${GREEN}✓ Created 50 transactions${NC}"

# ==================== CREDIT CARD TRANSACTIONS ====================
echo -e "${YELLOW}Creating credit card transactions...${NC}"

CC_CATEGORIES=("groceries" "dining" "shopping" "entertainment" "gas" "utilities")

for i in {1..30}; do
    CC_TRANS_ID=$(generate_uuid)
    TRANS_ID=$(generate_uuid)
    CC_DATE=$(random_past_date 60)
    CC_CATEGORY=${CC_CATEGORIES[$RANDOM % ${#CC_CATEGORIES[@]}]}
    CC_AMOUNT=$((20 + RANDOM % 300))

    # Randomly select credit card
    CARDS=("$CREDIT_CARD_1_ID" "$CREDIT_CARD_2_ID")
    CARD_ID=${CARDS[$RANDOM % ${#CARDS[@]}]}

    MERCHANTS=("Amazon" "Walmart" "Target" "Costco" "Whole Foods" "Shell" "Chevron" "Netflix" "Spotify")
    MERCHANT=${MERCHANTS[$RANDOM % ${#MERCHANTS[@]}]}

    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF > /dev/null 2>&1
INSERT INTO credit_card_transactions (id, user_id, card_id, transaction_id, category_id, amount, description, merchant, date, type, tags, created_at, updated_at)
VALUES (
    '$CC_TRANS_ID',
    '$USER_ID',
    '$CARD_ID',
    '$TRANS_ID',
    '$CC_CATEGORY',
    $CC_AMOUNT,
    'Purchase at $MERCHANT',
    '$MERCHANT',
    '$CC_DATE',
    'purchase',
    '[]'::jsonb,
    NOW(),
    NOW()
);
EOF
done

echo -e "${GREEN}✓ Created 30 credit card transactions${NC}"

# ==================== SAVINGS GOALS ====================
echo -e "${YELLOW}Creating savings goals...${NC}"

GOAL_1_ID=$(generate_uuid)
GOAL_2_ID=$(generate_uuid)
GOAL_3_ID=$(generate_uuid)

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Emergency Fund
INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount, target_date, monthly_contribution, category, priority, created_at, updated_at)
VALUES (
    '$GOAL_1_ID',
    '$USER_ID',
    'Emergency Fund',
    'Six months of expenses',
    20000.00,
    8500.00,
    '$(future_date 365)',
    500.00,
    'emergency',
    'high',
    NOW(),
    NOW()
);

-- Vacation Fund
INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount, target_date, monthly_contribution, category, priority, created_at, updated_at)
VALUES (
    '$GOAL_2_ID',
    '$USER_ID',
    'Europe Vacation',
    'Trip to Europe next summer',
    8000.00,
    2400.00,
    '$(future_date 270)',
    300.00,
    'vacation',
    'medium',
    NOW(),
    NOW()
);

-- New Car
INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount, target_date, monthly_contribution, category, priority, created_at, updated_at)
VALUES (
    '$GOAL_3_ID',
    '$USER_ID',
    'New Car Down Payment',
    'Save for car down payment',
    15000.00,
    4500.00,
    '$(future_date 540)',
    400.00,
    'purchase',
    'medium',
    NOW(),
    NOW()
);
EOF

echo -e "${GREEN}✓ Created 3 savings goals${NC}"

# ==================== BUDGETS ====================
echo -e "${YELLOW}Creating budgets...${NC}"

BUDGET_CATEGORIES=("groceries" "dining" "entertainment" "transportation" "shopping" "utilities")

for BUDGET_CAT in "${BUDGET_CATEGORIES[@]}"; do
    BUDGET_ID=$(generate_uuid)
    BUDGET_AMOUNT=$((200 + RANDOM % 800))

    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF > /dev/null 2>&1
INSERT INTO budgets (id, user_id, category_id, amount, period, alert_threshold, enabled, created_at, updated_at)
VALUES (
    '$BUDGET_ID',
    '$USER_ID',
    '$BUDGET_CAT',
    $BUDGET_AMOUNT,
    'monthly',
    80.0,
    true,
    NOW(),
    NOW()
);
EOF
done

echo -e "${GREEN}✓ Created 6 budgets${NC}"

# ==================== BILLS ====================
echo -e "${YELLOW}Creating bills...${NC}"

BILL_1_ID=$(generate_uuid)
BILL_2_ID=$(generate_uuid)
BILL_3_ID=$(generate_uuid)
BILL_4_ID=$(generate_uuid)
BILL_5_ID=$(generate_uuid)

psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Electricity Bill
INSERT INTO bills (id, user_id, name, category, amount, frequency, start_date, due_day, auto_pay, reminder_days, active, created_at, updated_at)
VALUES (
    '$BILL_1_ID',
    '$USER_ID',
    'Electric Company',
    'utilities',
    120.00,
    'monthly',
    '2024-01-01',
    15,
    false,
    3,
    true,
    NOW(),
    NOW()
);

-- Internet Bill
INSERT INTO bills (id, user_id, name, category, amount, frequency, start_date, due_day, auto_pay, reminder_days, active, created_at, updated_at)
VALUES (
    '$BILL_2_ID',
    '$USER_ID',
    'Internet Provider',
    'utilities',
    80.00,
    'monthly',
    '2024-01-01',
    5,
    true,
    3,
    true,
    NOW(),
    NOW()
);

-- Netflix Subscription
INSERT INTO bills (id, user_id, name, category, amount, frequency, start_date, due_day, auto_pay, reminder_days, active, created_at, updated_at)
VALUES (
    '$BILL_3_ID',
    '$USER_ID',
    'Netflix',
    'subscriptions',
    15.99,
    'monthly',
    '2024-01-01',
    10,
    true,
    0,
    true,
    NOW(),
    NOW()
);

-- Car Insurance
INSERT INTO bills (id, user_id, name, category, amount, frequency, start_date, due_day, auto_pay, reminder_days, active, created_at, updated_at)
VALUES (
    '$BILL_4_ID',
    '$USER_ID',
    'Auto Insurance',
    'insurance',
    250.00,
    'monthly',
    '2024-01-01',
    1,
    true,
    5,
    true,
    NOW(),
    NOW()
);

-- Gym Membership
INSERT INTO bills (id, user_id, name, category, amount, frequency, start_date, due_day, auto_pay, reminder_days, active, created_at, updated_at)
VALUES (
    '$BILL_5_ID',
    '$USER_ID',
    'Fitness Gym',
    'subscriptions',
    45.00,
    'monthly',
    '2024-01-01',
    20,
    true,
    3,
    true,
    NOW(),
    NOW()
);
EOF

echo -e "${GREEN}✓ Created 5 bills${NC}"

# ==================== UPDATE ACCOUNT BALANCES ====================
echo -e "${YELLOW}Updating account balances based on transactions...${NC}"

# Calculate balance changes
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF > /dev/null 2>&1
-- Update balances based on income and expenses
UPDATE accounts a
SET balance = a.initial_balance + COALESCE(
    (SELECT SUM(CASE
        WHEN t.type = 'income' THEN t.amount
        WHEN t.type = 'expense' THEN -t.amount
        ELSE 0
    END)
    FROM transactions t
    WHERE t.account_id = a.id AND t.user_id = '$USER_ID'),
    0
)
WHERE a.user_id = '$USER_ID';
EOF

echo -e "${GREEN}✓ Updated account balances${NC}"

# ==================== SUMMARY ====================
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Dummy Data Generation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Summary for user: $FULL_NAME ($USER_EMAIL)"
echo "  - 3 Accounts (Checking, Savings, Cash)"
echo "  - 2 Credit Cards"
echo "  - 500 Transactions"
echo "  - 30 Credit Card Transactions"
echo "  - 3 Savings Goals"
echo "  - 6 Budgets"
echo "  - 5 Bills"
echo ""
echo -e "${GREEN}You can now use the application with this data!${NC}"

# Clean up
unset PGPASSWORD
