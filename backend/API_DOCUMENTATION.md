# Daybook API Documentation

Complete API reference for the Daybook Personal Finance Management Backend.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## Endpoints

### Authentication

#### Signup
Create a new user account.

**Endpoint:** `POST /auth/signup`

**Request Body:**
```json
{
  "username": "string (required)",
  "email": "string (required, email format)",
  "password": "string (required, min 6 characters)",
  "fullName": "string (optional)"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "token": "jwt_token_here",
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "fullName": "string",
      "role": "user",
      "createdAt": "timestamp",
      "updatedAt": "timestamp"
    }
  }
}
```

#### Login
Authenticate and receive JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "jwt_token_here",
    "user": { ... }
  }
}
```

#### Get Profile
Get current user profile.

**Endpoint:** `GET /auth/me`

**Headers:** Authorization required

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "fullName": "string",
    "role": "string",
    "lastLogin": "timestamp"
  }
}
```

#### Update Profile
Update user profile information.

**Endpoint:** `PUT /auth/profile`

**Headers:** Authorization required

**Request Body:**
```json
{
  "fullName": "string (optional)",
  "email": "string (optional)"
}
```

**Response:** `200 OK`

#### Change Password
Change user password.

**Endpoint:** `PUT /auth/change-password`

**Headers:** Authorization required

**Request Body:**
```json
{
  "currentPassword": "string (required)",
  "newPassword": "string (required, min 6 characters)"
}
```

**Response:** `200 OK`

---

### Accounts

#### List Accounts
Get all accounts for authenticated user.

**Endpoint:** `GET /accounts`

**Headers:** Authorization required

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "userId": "uuid",
      "name": "string",
      "type": "string",
      "balance": 0.00,
      "currency": "USD",
      "description": "string",
      "institution": "string",
      "active": true,
      "createdAt": "timestamp",
      "updatedAt": "timestamp"
    }
  ]
}
```

#### Get Account
Get specific account by ID.

**Endpoint:** `GET /accounts/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Account
Create new account.

**Endpoint:** `POST /accounts`

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "string (required)",
  "type": "string (required)", // cash, checking, savings, credit_card, etc
  "balance": 0.00,
  "currency": "USD",
  "description": "string",
  "institution": "string"
}
```

**Response:** `201 Created`

#### Update Account
Update account details.

**Endpoint:** `PUT /accounts/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Account

**Response:** `200 OK`

#### Delete Account
Delete account (soft delete).

**Endpoint:** `DELETE /accounts/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Update Account Balance
Directly update account balance.

**Endpoint:** `PATCH /accounts/:id/balance`

**Headers:** Authorization required

**Request Body:**
```json
{
  "balance": 0.00,
  "operation": "set" // set, add, or subtract
}
```

**Response:** `200 OK`

---

### Transactions

#### List Transactions
Get all transactions with optional filters.

**Endpoint:** `GET /transactions`

**Headers:** Authorization required

**Query Parameters:**
- `type` - income, expense, or transfer
- `categoryId` - Filter by category
- `accountId` - Filter by account
- `startDate` - Start date (YYYY-MM-DD)
- `endDate` - End date (YYYY-MM-DD)

**Response:** `200 OK`

#### Get Transaction
Get specific transaction.

**Endpoint:** `GET /transactions/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Transaction
Create new transaction.

**Endpoint:** `POST /transactions`

**Headers:** Authorization required

**Request Body:**
```json
{
  "accountId": "uuid (required)",
  "toAccountId": "uuid (optional, for transfers)",
  "type": "income|expense|transfer (required)",
  "amount": 0.00,
  "categoryId": "string (required)",
  "date": "timestamp (required)",
  "description": "string",
  "tags": ["string"],
  "savingsGoalId": "uuid (optional)",
  "creditCardId": "uuid (optional)"
}
```

**Response:** `201 Created`

#### Update Transaction
Update transaction.

**Endpoint:** `PUT /transactions/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Transaction

**Response:** `200 OK`

#### Delete Transaction
Delete transaction.

**Endpoint:** `DELETE /transactions/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Bulk Create Transactions
Import multiple transactions at once.

**Endpoint:** `POST /transactions/bulk`

**Headers:** Authorization required

**Request Body:**
```json
{
  "transactions": [
    { /* transaction object */ },
    { /* transaction object */ }
  ]
}
```

**Response:** `201 Created`

#### Get Transaction Statistics
Get transaction statistics.

**Endpoint:** `GET /transactions/stats`

**Headers:** Authorization required

**Query Parameters:**
- `startDate` - Start date (YYYY-MM-DD)
- `endDate` - End date (YYYY-MM-DD)

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "totalIncome": 0.00,
    "totalExpense": 0.00,
    "totalTransfers": 0.00,
    "netIncome": 0.00,
    "transactionCount": 0
  }
}
```

---

### Credit Cards

#### List Credit Cards
Get all credit cards.

**Endpoint:** `GET /credit-cards`

**Headers:** Authorization required

**Response:** `200 OK`

#### Get Credit Card
Get specific credit card.

**Endpoint:** `GET /credit-cards/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Credit Card
Create new credit card.

**Endpoint:** `POST /credit-cards`

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "string (required)",
  "lastFourDigits": "string",
  "cardNetwork": "string",
  "creditLimit": 0.00,
  "currentBalance": 0.00,
  "apr": 0.00,
  "dueDate": "timestamp",
  "statementDate": "timestamp",
  "rewardsProgram": "string",
  "notes": "string"
}
```

**Response:** `201 Created`

#### Update Credit Card
Update credit card details.

**Endpoint:** `PUT /credit-cards/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Credit Card

**Response:** `200 OK`

#### Delete Credit Card
Delete credit card.

**Endpoint:** `DELETE /credit-cards/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Record Credit Card Payment
Record a payment for credit card.

**Endpoint:** `POST /credit-cards/:id/payment`

**Headers:** Authorization required

**Request Body:**
```json
{
  "amount": 0.00,
  "paymentDate": "timestamp",
  "paymentAccountId": "uuid (optional)"
}
```

**Response:** `200 OK`

#### Get Card Statements
Get statements for specific credit card.

**Endpoint:** `GET /credit-cards/:id/statements`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Statement
Create new statement.

**Endpoint:** `POST /statements`

**Headers:** Authorization required

**Request Body:**
```json
{
  "cardId": "uuid (required)",
  "statementDate": "timestamp (required)",
  "dueDate": "timestamp (required)",
  "openingBalance": 0.00,
  "closingBalance": 0.00,
  "minimumPayment": 0.00,
  "totalCharges": 0.00,
  "totalPayments": 0.00,
  "interestCharged": 0.00
}
```

**Response:** `201 Created`

#### List Rewards
Get rewards with optional filters.

**Endpoint:** `GET /rewards`

**Headers:** Authorization required

**Query Parameters:**
- `cardId` - Filter by credit card
- `redeemed` - Filter by redeemed status (true/false)

**Response:** `200 OK`

#### Create Reward
Record new reward.

**Endpoint:** `POST /rewards`

**Headers:** Authorization required

**Request Body:**
```json
{
  "cardId": "uuid (required)",
  "type": "cashback|points|miles",
  "amount": 0.00,
  "description": "string",
  "earnedDate": "timestamp (required)"
}
```

**Response:** `201 Created`

---

### Investments

#### List Investments
Get all investments with optional filters.

**Endpoint:** `GET /investments`

**Headers:** Authorization required

**Query Parameters:**
- `portfolioId` - Filter by portfolio
- `assetType` - Filter by asset type

**Response:** `200 OK`

#### Get Investment
Get specific investment.

**Endpoint:** `GET /investments/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Investment
Create new investment.

**Endpoint:** `POST /investments`

**Headers:** Authorization required

**Request Body:**
```json
{
  "portfolioId": "uuid (optional)",
  "symbol": "string (required)",
  "name": "string (required)",
  "assetType": "stocks|bonds|mutual_funds|etf|crypto|real_estate|commodities|other",
  "quantity": 0.00,
  "costBasis": 0.00,
  "currentPrice": 0.00,
  "purchaseDate": "timestamp",
  "notes": "string"
}
```

**Response:** `201 Created`

#### Update Investment
Update investment details.

**Endpoint:** `PUT /investments/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Investment

**Response:** `200 OK`

#### Delete Investment
Delete investment.

**Endpoint:** `DELETE /investments/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Buy Shares
Buy additional shares (updates weighted average cost basis).

**Endpoint:** `POST /investments/:id/buy`

**Headers:** Authorization required

**Request Body:**
```json
{
  "quantity": 0.00,
  "price": 0.00
}
```

**Response:** `200 OK`

#### Sell Shares
Sell shares (calculates realized gain/loss).

**Endpoint:** `POST /investments/:id/sell`

**Headers:** Authorization required

**Request Body:**
```json
{
  "quantity": 0.00,
  "price": 0.00
}
```

**Response:** `200 OK`

#### List Portfolios
Get all portfolios.

**Endpoint:** `GET /portfolios`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Portfolio
Create new portfolio.

**Endpoint:** `POST /portfolios`

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "string (required)",
  "description": "string",
  "type": "retirement|taxable|education|other"
}
```

**Response:** `201 Created`

#### List Dividends
Get all dividends.

**Endpoint:** `GET /dividends`

**Headers:** Authorization required

**Query Parameters:**
- `investmentId` - Filter by investment

**Response:** `200 OK`

#### Create Dividend
Record new dividend payment.

**Endpoint:** `POST /dividends`

**Headers:** Authorization required

**Request Body:**
```json
{
  "investmentId": "uuid (required)",
  "amount": 0.00,
  "paymentDate": "timestamp (required)",
  "reinvested": false
}
```

**Response:** `201 Created`

---

### Bills

#### List Bills
Get all bills with optional filters.

**Endpoint:** `GET /bills`

**Headers:** Authorization required

**Query Parameters:**
- `active` - Filter by active status (true/false)
- `category` - Filter by category

**Response:** `200 OK`

#### Get Bill
Get specific bill.

**Endpoint:** `GET /bills/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Bill
Create new bill.

**Endpoint:** `POST /bills`

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "string (required)",
  "category": "string (required)",
  "amount": 0.00,
  "frequency": "weekly|biweekly|monthly|quarterly|semi-annually|annually",
  "startDate": "timestamp (required)",
  "dueDay": 1,
  "autoPay": false,
  "reminderDays": 3,
  "notes": "string"
}
```

**Response:** `201 Created`

#### Update Bill
Update bill details.

**Endpoint:** `PUT /bills/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Bill

**Response:** `200 OK`

#### Delete Bill
Delete bill.

**Endpoint:** `DELETE /bills/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Pay Bill
Mark bill as paid.

**Endpoint:** `POST /bills/:id/pay`

**Headers:** Authorization required

**Request Body:**
```json
{
  "amount": 0.00,
  "paymentDate": "timestamp (optional)",
  "accountId": "uuid (optional)",
  "notes": "string"
}
```

**Response:** `200 OK`

#### Get Bill Payments
Get bill payment history.

**Endpoint:** `GET /bill-payments`

**Headers:** Authorization required

**Query Parameters:**
- `billId` - Filter by bill
- `startDate` - Start date
- `endDate` - End date

**Response:** `200 OK`

---

### Budgets

#### List Budgets
Get all budgets with optional filters.

**Endpoint:** `GET /budgets`

**Headers:** Authorization required

**Query Parameters:**
- `enabled` - Filter by enabled status
- `period` - Filter by period
- `categoryId` - Filter by category

**Response:** `200 OK`

#### Get Budget
Get specific budget.

**Endpoint:** `GET /budgets/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Budget
Create new budget.

**Endpoint:** `POST /budgets`

**Headers:** Authorization required

**Request Body:**
```json
{
  "categoryId": "string (required)",
  "amount": 0.00,
  "period": "weekly|monthly|quarterly|yearly|custom",
  "customStartDate": "timestamp (for custom period)",
  "customEndDate": "timestamp (for custom period)",
  "rollover": false,
  "alertThreshold": 80,
  "notes": "string"
}
```

**Response:** `201 Created`

#### Update Budget
Update budget details.

**Endpoint:** `PUT /budgets/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Budget

**Response:** `200 OK`

#### Delete Budget
Delete budget.

**Endpoint:** `DELETE /budgets/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Get Budget Progress
Get budget progress and spending.

**Endpoint:** `GET /budgets/:id/progress`

**Headers:** Authorization required

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "spent": 0.00,
    "remaining": 0.00,
    "percentage": 0.00,
    "isOverBudget": false
  }
}
```

---

### Savings Goals

#### List Savings Goals
Get all savings goals with optional filters.

**Endpoint:** `GET /savings-goals`

**Headers:** Authorization required

**Query Parameters:**
- `achieved` - Filter by achieved status
- `archived` - Filter by archived status
- `category` - Filter by category
- `priority` - Filter by priority

**Response:** `200 OK`

#### Get Savings Goal
Get specific savings goal.

**Endpoint:** `GET /savings-goals/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Savings Goal
Create new savings goal.

**Endpoint:** `POST /savings-goals`

**Headers:** Authorization required

**Request Body:**
```json
{
  "name": "string (required)",
  "description": "string",
  "targetAmount": 0.00,
  "currentAmount": 0.00,
  "targetDate": "timestamp",
  "monthlyContribution": 0.00,
  "category": "emergency|vacation|purchase|other",
  "priority": "high|medium|low"
}
```

**Response:** `201 Created`

#### Update Savings Goal
Update savings goal details.

**Endpoint:** `PUT /savings-goals/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Savings Goal

**Response:** `200 OK`

#### Delete Savings Goal
Delete savings goal.

**Endpoint:** `DELETE /savings-goals/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Contribute to Goal
Add contribution to savings goal.

**Endpoint:** `POST /savings-goals/:id/contribute`

**Headers:** Authorization required

**Request Body:**
```json
{
  "amount": 0.00,
  "date": "timestamp (optional)",
  "notes": "string"
}
```

**Response:** `200 OK`

#### Withdraw from Goal
Withdraw from savings goal.

**Endpoint:** `POST /savings-goals/:id/withdraw`

**Headers:** Authorization required

**Request Body:**
```json
{
  "amount": 0.00,
  "notes": "string"
}
```

**Response:** `200 OK`

#### List Automated Rules
Get all automated savings rules.

**Endpoint:** `GET /automated-rules`

**Headers:** Authorization required

**Query Parameters:**
- `goalId` - Filter by goal
- `enabled` - Filter by enabled status

**Response:** `200 OK`

#### Create Automated Rule
Create new automated savings rule.

**Endpoint:** `POST /automated-rules`

**Headers:** Authorization required

**Request Body:**
```json
{
  "goalId": "uuid (required)",
  "ruleType": "percentage_of_income|fixed_amount|round_up",
  "amount": 0.00,
  "percentage": 0.00,
  "frequency": "daily|weekly|monthly"
}
```

**Response:** `201 Created`

---

### Fixed Deposits

#### List Fixed Deposits
Get all fixed deposits.

**Endpoint:** `GET /fixed-deposits`

**Headers:** Authorization required

**Query Parameters:**
- `withdrawn` - Filter by withdrawn status

**Response:** `200 OK`

#### Get Fixed Deposit
Get specific fixed deposit.

**Endpoint:** `GET /fixed-deposits/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Create Fixed Deposit
Create new fixed deposit.

**Endpoint:** `POST /fixed-deposits`

**Headers:** Authorization required

**Request Body:**
```json
{
  "institution": "string (required)",
  "accountNumber": "string",
  "principal": 0.00,
  "interestRate": 0.00,
  "tenureMonths": 12,
  "compounding": "simple|daily|monthly|quarterly|semi-annually|annually",
  "startDate": "timestamp (required)",
  "maturityDate": "timestamp (required)",
  "autoRenew": false,
  "notes": "string"
}
```

**Response:** `201 Created`

#### Update Fixed Deposit
Update fixed deposit details.

**Endpoint:** `PUT /fixed-deposits/:id`

**Headers:** Authorization required

**Request Body:** Same as Create Fixed Deposit

**Response:** `200 OK`

#### Delete Fixed Deposit
Delete fixed deposit.

**Endpoint:** `DELETE /fixed-deposits/:id`

**Headers:** Authorization required

**Response:** `200 OK`

#### Withdraw Fixed Deposit
Withdraw fixed deposit (calculate final amount).

**Endpoint:** `POST /fixed-deposits/:id/withdraw`

**Headers:** Authorization required

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "maturityAmount": 0.00,
    "interestEarned": 0.00
  }
}
```

---

### Settings

#### Get Settings
Get user settings.

**Endpoint:** `GET /settings`

**Headers:** Authorization required

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "userId": "uuid",
    "currency": "USD",
    "darkMode": false,
    "dateFormat": "MM/DD/YYYY",
    "firstDayOfWeek": 0,
    "language": "en",
    "notifications": {
      "push": true,
      "email": true,
      "budgetAlerts": true,
      "billReminders": true
    }
  }
}
```

#### Update Settings
Update user settings.

**Endpoint:** `PUT /settings`

**Headers:** Authorization required

**Request Body:**
```json
{
  "currency": "USD",
  "darkMode": false,
  "dateFormat": "MM/DD/YYYY",
  "firstDayOfWeek": 0,
  "language": "en",
  "notifications": {
    "push": true,
    "email": true,
    "budgetAlerts": true,
    "billReminders": true
  }
}
```

**Response:** `200 OK`

---

### Health Check

#### Health
Check API health status.

**Endpoint:** `GET /health`

**Response:** `200 OK`
```json
{
  "status": "ok",
  "message": "Daybook API is running"
}
```

---

## HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

## Rate Limiting

Currently no rate limiting is implemented. Consider implementing rate limiting for production use.

## Pagination

List endpoints currently return all results. Consider implementing pagination for large datasets:

```
GET /transactions?page=1&limit=50
```

## Filtering & Sorting

Many list endpoints support filtering via query parameters. Future enhancements could include:

```
GET /transactions?sort=date&order=desc
```

## Webhook Support

Future feature: Webhook notifications for events like:
- Budget exceeded
- Bill due soon
- Goal achieved
- Fixed deposit maturity approaching
