# Daybook - Google Sheets Template

Complete financial tracking system in Google Sheets with all Daybook features.

## 📋 Google Sheets Structure

This template consists of **15 interconnected sheets**:

### Core Sheets
1. **🏠 Dashboard** - Overview with key metrics
2. **👤 Profile** - User settings
3. **💰 Accounts** - Bank accounts, cash, investments
4. **📊 Transactions** - All income and expenses
5. **🔄 Recurring** - Recurring transactions
6. **🏷️ Tags** - Transaction categories
7. **📈 Budgets** - Monthly budgets
8. **🎯 Goals** - Financial goals tracking

### Credit & Bills
9. **💳 Credit Cards** - Credit card management
10. **💳 CC Transactions** - Credit card transactions
11. **💰 CC Payments** - Credit card payments
12. **📄 Bills** - Bills tracking
13. **💵 Bill Payments** - Bill payment history

### Reports
14. **📊 Reports** - Financial reports
15. **📉 Charts** - Visual analytics

---

## 🚀 Quick Setup Guide

### Step 1: Create New Google Sheet

1. Go to [Google Sheets](https://sheets.google.com)
2. Create a new spreadsheet
3. Name it "Daybook Financial Tracker"

### Step 2: Create All Sheets

Create these 15 sheets (tabs at bottom):
- Dashboard
- Profile
- Accounts
- Transactions
- Recurring
- Tags
- Budgets
- Goals
- Credit_Cards
- CC_Transactions
- CC_Payments
- Bills
- Bill_Payments
- Reports
- Charts

---

## 📝 Detailed Sheet Setup

## 1. 🏠 DASHBOARD

**Purpose:** Overview of your financial health

### Columns (Row 1):
```
A: Metric
B: Value
C: Change
D: Status
```

### Setup:

**Row 2-20: Add these metrics**

```
Row 2:  Total Net Worth | =SUM(Accounts!D:D) | |
Row 3:  Total Income (This Month) | =SUMIFS(Transactions!C:C,Transactions!B:B,"Income",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1) | |
Row 4:  Total Expenses (This Month) | =SUMIFS(Transactions!C:C,Transactions!B:B,"Expense",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1) | |
Row 5:  Net Savings (This Month) | =B3-B4 | |
Row 6:  Budget Progress | =B4/SUM(Budgets!C:C) | |
Row 7:  Active Accounts | =COUNTA(Accounts!B:B)-1 | |
Row 8:  Active Credit Cards | =COUNTA(Credit_Cards!B:B)-1 | |
Row 9:  Pending Bills | =COUNTIFS(Bills!F:F,"Pending") | |
Row 10: Goals Progress | =AVERAGE(Goals!E:E) | |
```

### Formatting:
- **B2**: Currency format
- **B3-B5**: Currency format
- **B6**: Percentage format
- **B7-B9**: Number format
- **B10**: Percentage format

### Add Chart (E2:K20):
- **Chart Type**: Pie chart
- **Data Range**: Accounts!B:D (Account names and balances)
- **Title**: "Account Distribution"

---

## 2. 👤 PROFILE

**Purpose:** User settings and preferences

### Columns:
```
A: Setting | B: Value
```

### Data:
```
Row 2:  Name | Your Name
Row 3:  Currency | USD
Row 4:  Date Format | MM/DD/YYYY
Row 5:  Start Date | 01/01/2024
Row 6:  Fiscal Year Start | January
Row 7:  Budget Period | Monthly
```

---

## 3. 💰 ACCOUNTS

**Purpose:** Track all your accounts

### Columns:
```
A: ID | B: Account Name | C: Type | D: Balance | E: Last Updated | F: Status | G: Notes
```

### Data Validation:
- **Column C (Type)**: Dropdown
  - Values: `Checking, Savings, Investment, Cash, Credit Card, Loan, Other`

- **Column F (Status)**: Dropdown
  - Values: `Active, Inactive, Closed`

### Sample Data:
```
Row 2: 1 | Chase Checking | Checking | 5,000.00 | =TODAY() | Active | Primary checking
Row 3: 2 | Savings Account | Savings | 15,000.00 | =TODAY() | Active | Emergency fund
Row 4: 3 | Investment Account | Investment | 50,000.00 | =TODAY() | Active | Retirement
Row 5: 4 | Cash Wallet | Cash | 200.00 | =TODAY() | Active | Daily cash
```

### Formulas:
- **Total Row (below last account)**:
  - Column C: `TOTAL`
  - Column D: `=SUM(D2:D100)`

### Conditional Formatting:
- **Column D (Balance)**:
  - If < 0: Red background
  - If > 10000: Green background
  - If between 0-1000: Yellow background

---

## 4. 📊 TRANSACTIONS

**Purpose:** All income and expense transactions

### Columns:
```
A: Date | B: Type | C: Amount | D: Category | E: Account | F: Description | G: Tags | H: Recurring | I: Status
```

### Data Validation:
- **Column B (Type)**: Dropdown → `Income, Expense, Transfer`
- **Column D (Category)**: Dropdown → Use data from Tags!B:B
- **Column E (Account)**: Dropdown → Use data from Accounts!B:B
- **Column H (Recurring)**: Dropdown → `Yes, No`
- **Column I (Status)**: Dropdown → `Completed, Pending, Cancelled`

### Sample Data:
```
Row 2: 01/15/2024 | Income | 5000.00 | Salary | Chase Checking | Monthly salary | Income:Salary | No | Completed
Row 3: 01/16/2024 | Expense | 1200.00 | Rent | Chase Checking | Apartment rent | Housing:Rent | Yes | Completed
Row 4: 01/17/2024 | Expense | 150.00 | Groceries | Chase Checking | Weekly groceries | Food:Groceries | No | Completed
Row 5: 01/18/2024 | Expense | 50.00 | Gas | Chase Checking | Car fuel | Transportation:Gas | No | Completed
```

### Summary Section (M1:N10):
```
M1: THIS MONTH SUMMARY
M2: Total Income | =SUMIFS(C:C,B:B,"Income",A:A,">="&EOMONTH(TODAY(),-1)+1,A:A,"<="&EOMONTH(TODAY(),0))
M3: Total Expenses | =SUMIFS(C:C,B:B,"Expense",A:A,">="&EOMONTH(TODAY(),-1)+1,A:A,"<="&EOMONTH(TODAY(),0))
M4: Net Savings | =N2-N3
M5: Transaction Count | =COUNTIFS(A:A,">="&EOMONTH(TODAY(),-1)+1,A:A,"<="&EOMONTH(TODAY(),0))
```

### Conditional Formatting:
- **Column B (Type)**:
  - Income: Green text
  - Expense: Red text
  - Transfer: Blue text

- **Column I (Status)**:
  - Completed: Light green background
  - Pending: Light yellow background
  - Cancelled: Light red background

---

## 5. 🔄 RECURRING

**Purpose:** Manage recurring transactions

### Columns:
```
A: ID | B: Name | C: Type | D: Amount | E: Category | F: Account | G: Frequency | H: Start Date | I: End Date | J: Active
```

### Data Validation:
- **Column C (Type)**: Dropdown → `Income, Expense`
- **Column E (Category)**: Dropdown → Use data from Tags!B:B
- **Column F (Account)**: Dropdown → Use data from Accounts!B:B
- **Column G (Frequency)**: Dropdown → `Daily, Weekly, Monthly, Yearly`
- **Column J (Active)**: Dropdown → `Yes, No`

### Sample Data:
```
Row 2: 1 | Monthly Salary | Income | 5000.00 | Salary | Chase Checking | Monthly | 01/01/2024 | | Yes
Row 3: 2 | Rent Payment | Expense | 1200.00 | Rent | Chase Checking | Monthly | 01/01/2024 | | Yes
Row 4: 3 | Netflix | Expense | 15.99 | Subscription | Chase Checking | Monthly | 01/01/2024 | | Yes
Row 5: 4 | Gym Membership | Expense | 50.00 | Health | Chase Checking | Monthly | 01/01/2024 | | Yes
```

### Auto-Create Transactions:
Add this script via **Tools → Script Editor**:

```javascript
function createRecurringTransactions() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const recurringSheet = ss.getSheetByName('Recurring');
  const transactionsSheet = ss.getSheetByName('Transactions');

  const recurringData = recurringSheet.getDataRange().getValues();
  const today = new Date();

  for (let i = 1; i < recurringData.length; i++) {
    if (recurringData[i][9] === 'Yes') { // If Active
      const startDate = new Date(recurringData[i][7]);
      const frequency = recurringData[i][6];

      if (shouldCreateTransaction(startDate, frequency, today)) {
        transactionsSheet.appendRow([
          today,
          recurringData[i][2], // Type
          recurringData[i][3], // Amount
          recurringData[i][4], // Category
          recurringData[i][5], // Account
          recurringData[i][1], // Name as Description
          '',
          'Yes',
          'Completed'
        ]);
      }
    }
  }
}

function shouldCreateTransaction(startDate, frequency, today) {
  // Add logic based on frequency
  return true; // Simplified
}
```

---

## 6. 🏷️ TAGS

**Purpose:** Categories for transactions

### Columns:
```
A: ID | B: Tag Name | C: Type | D: Parent | E: Color | F: Icon
```

### Sample Data:
```
Row 2: 1 | Income | Income | | #4CAF50 | 💰
Row 3: 2 | Salary | Income | Income | #4CAF50 | 💵
Row 4: 3 | Freelance | Income | Income | #4CAF50 | 💼
Row 5: 4 | Expenses | Expense | | #F44336 | 💸
Row 6: 5 | Housing | Expense | Expenses | #F44336 | 🏠
Row 7: 6 | Rent | Expense | Housing | #F44336 | 🏠
Row 8: 7 | Utilities | Expense | Housing | #F44336 | ⚡
Row 9: 8 | Food | Expense | Expenses | #FF9800 | 🍔
Row 10: 9 | Groceries | Expense | Food | #FF9800 | 🛒
Row 11: 10 | Dining Out | Expense | Food | #FF9800 | 🍽️
Row 12: 11 | Transportation | Expense | Expenses | #2196F3 | 🚗
Row 13: 12 | Gas | Expense | Transportation | #2196F3 | ⛽
Row 14: 13 | Public Transit | Expense | Transportation | #2196F3 | 🚌
Row 15: 14 | Entertainment | Expense | Expenses | #9C27B0 | 🎬
Row 16: 15 | Shopping | Expense | Expenses | #E91E63 | 🛍️
Row 17: 16 | Health | Expense | Expenses | #00BCD4 | 🏥
Row 18: 17 | Subscription | Expense | Expenses | #795548 | 📱
```

---

## 7. 📈 BUDGETS

**Purpose:** Monthly budget tracking

### Columns:
```
A: Month | B: Category | C: Budgeted | D: Spent | E: Remaining | F: Progress % | G: Status
```

### Formulas:
- **Column D (Spent)**:
  ```
  =SUMIFS(Transactions!C:C, Transactions!D:D, B2, Transactions!A:A, ">="&DATE(YEAR(A2),MONTH(A2),1), Transactions!A:A, "<="&EOMONTH(A2,0), Transactions!B:B, "Expense")
  ```

- **Column E (Remaining)**: `=C2-D2`

- **Column F (Progress %)**: `=D2/C2`

- **Column G (Status)**:
  ```
  =IF(F2>1,"Over Budget",IF(F2>0.8,"Warning","On Track"))
  ```

### Sample Data:
```
Row 2: 01/01/2024 | Housing | 1500.00 | =Formula | =Formula | =Formula | =Formula
Row 3: 01/01/2024 | Food | 600.00 | =Formula | =Formula | =Formula | =Formula
Row 4: 01/01/2024 | Transportation | 300.00 | =Formula | =Formula | =Formula | =Formula
Row 5: 01/01/2024 | Entertainment | 200.00 | =Formula | =Formula | =Formula | =Formula
Row 6: 01/01/2024 | Shopping | 400.00 | =Formula | =Formula | =Formula | =Formula
```

### Conditional Formatting:
- **Column G (Status)**:
  - "Over Budget": Red background
  - "Warning": Yellow background
  - "On Track": Green background

- **Column F (Progress %)**:
  - >100%: Red background
  - 80-100%: Yellow background
  - <80%: Green background

---

## 8. 🎯 GOALS

**Purpose:** Financial goals tracking

### Columns:
```
A: ID | B: Goal Name | C: Target Amount | D: Current Amount | E: Progress % | F: Deadline | G: Monthly Target | H: Status
```

### Formulas:
- **Column D (Current Amount)**: Link to specific account or manual entry

- **Column E (Progress %)**: `=D2/C2`

- **Column G (Monthly Target)**:
  ```
  =IF(F2="","N/A",(C2-D2)/((F2-TODAY())/30))
  ```

- **Column H (Status)**:
  ```
  =IF(D2>=C2,"Achieved",IF(E2>=0.75,"On Track",IF(E2>=0.5,"Behind","Urgent")))
  ```

### Sample Data:
```
Row 2: 1 | Emergency Fund | 10000.00 | 5000.00 | =Formula | 12/31/2024 | =Formula | =Formula
Row 3: 2 | Vacation | 3000.00 | 1200.00 | =Formula | 06/30/2024 | =Formula | =Formula
Row 4: 3 | New Car | 25000.00 | 8000.00 | =Formula | 12/31/2025 | =Formula | =Formula
Row 5: 4 | House Down Payment | 50000.00 | 15000.00 | =Formula | 12/31/2026 | =Formula | =Formula
```

### Progress Bars:
Use conditional formatting on Column E with color scale:
- 0%: Red
- 50%: Yellow
- 100%: Green

---

## 9. 💳 CREDIT CARDS

**Purpose:** Credit card management

### Columns:
```
A: ID | B: Card Name | C: Bank | D: Last 4 Digits | E: Credit Limit | F: Current Balance | G: Available Credit | H: Interest Rate | I: Due Date | J: Status
```

### Formulas:
- **Column F (Current Balance)**:
  ```
  =SUMIFS(CC_Transactions!D:D, CC_Transactions!B:B, A2) - SUMIFS(CC_Payments!C:C, CC_Payments!B:B, A2)
  ```

- **Column G (Available Credit)**: `=E2-F2`

### Sample Data:
```
Row 2: 1 | Chase Sapphire | Chase | 1234 | 10000.00 | =Formula | =Formula | 18.99% | 15 | Active
Row 3: 2 | Amex Gold | American Express | 5678 | 15000.00 | =Formula | =Formula | 19.99% | 20 | Active
```

### Conditional Formatting:
- **Column G (Available Credit)**:
  - <20% of limit: Red
  - 20-50% of limit: Yellow
  - >50% of limit: Green

---

## 10. 💳 CC_TRANSACTIONS

**Purpose:** Credit card transactions

### Columns:
```
A: Date | B: Card ID | C: Card Name | D: Amount | E: Category | F: Description | G: Status
```

### Data Validation:
- **Column C (Card Name)**: Dropdown → Use data from Credit_Cards!B:B
- **Column E (Category)**: Dropdown → Use data from Tags!B:B
- **Column G (Status)**: Dropdown → `Posted, Pending, Declined`

### Sample Data:
```
Row 2: 01/15/2024 | 1 | Chase Sapphire | 150.00 | Dining Out | Restaurant | Posted
Row 3: 01/16/2024 | 1 | Chase Sapphire | 50.00 | Gas | Gas station | Posted
Row 4: 01/17/2024 | 2 | Amex Gold | 1200.00 | Shopping | Electronics | Posted
```

---

## 11. 💰 CC_PAYMENTS

**Purpose:** Credit card payment history

### Columns:
```
A: Date | B: Card ID | C: Payment Amount | D: From Account | E: Payment Method | F: Status
```

### Data Validation:
- **Column D (From Account)**: Dropdown → Use data from Accounts!B:B
- **Column E (Payment Method)**: Dropdown → `Bank Transfer, Check, Cash, Auto-Pay`
- **Column F (Status)**: Dropdown → `Completed, Pending, Failed`

### Sample Data:
```
Row 2: 01/20/2024 | 1 | 500.00 | Chase Checking | Bank Transfer | Completed
Row 3: 01/21/2024 | 2 | 800.00 | Chase Checking | Bank Transfer | Completed
```

---

## 12. 📄 BILLS

**Purpose:** Track recurring bills

### Columns:
```
A: ID | B: Bill Name | C: Payee | D: Amount | E: Due Date | F: Status | G: Category | H: Auto-Pay | I: Account
```

### Data Validation:
- **Column F (Status)**: Dropdown → `Pending, Paid, Overdue, Cancelled`
- **Column G (Category)**: Dropdown → Use data from Tags!B:B
- **Column H (Auto-Pay)**: Dropdown → `Yes, No`
- **Column I (Account)**: Dropdown → Use data from Accounts!B:B

### Formulas:
- **Add conditional formatting for overdue bills**:
  - If F="Overdue": Red background

### Sample Data:
```
Row 2: 1 | Electric Bill | Power Company | 120.00 | 15 | Pending | Utilities | No | Chase Checking
Row 3: 2 | Internet | ISP Provider | 79.99 | 10 | Paid | Utilities | Yes | Chase Checking
Row 4: 3 | Phone Bill | Phone Company | 65.00 | 5 | Paid | Utilities | Yes | Chase Checking
Row 5: 4 | Insurance | Insurance Co | 250.00 | 1 | Pending | Insurance | Yes | Chase Checking
```

---

## 13. 💵 BILL_PAYMENTS

**Purpose:** Bill payment history

### Columns:
```
A: Date | B: Bill ID | C: Bill Name | D: Amount Paid | E: From Account | F: Payment Method | G: Status
```

### Data Validation:
- **Column E (From Account)**: Dropdown → Use data from Accounts!B:B
- **Column F (Payment Method)**: Dropdown → `Bank Transfer, Check, Cash, Credit Card, Auto-Pay`
- **Column G (Status)**: Dropdown → `Completed, Pending, Failed`

### Sample Data:
```
Row 2: 01/15/2024 | 1 | Electric Bill | 120.00 | Chase Checking | Bank Transfer | Completed
Row 3: 01/10/2024 | 2 | Internet | 79.99 | Chase Checking | Auto-Pay | Completed
```

---

## 14. 📊 REPORTS

**Purpose:** Financial reports and analytics

### Monthly Summary (A1:B20):
```
A1: MONTHLY REPORT | =TEXT(TODAY(),"MMMM YYYY")
A3: Total Income | =SUMIFS(Transactions!C:C,Transactions!B:B,"Income",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
A4: Total Expenses | =SUMIFS(Transactions!C:C,Transactions!B:B,"Expense",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
A5: Net Savings | =B3-B4
A6: Savings Rate | =B5/B3
A8: Top 5 Expense Categories |
```

### Query for Top Expenses (A9:B14):
```
=QUERY(Transactions!A:I,"SELECT D, SUM(C) WHERE B='Expense' AND A >= date '"&TEXT(EOMONTH(TODAY(),-1)+1,"yyyy-mm-dd")&"' GROUP BY D ORDER BY SUM(C) DESC LIMIT 5 LABEL D 'Category', SUM(C) 'Amount'")
```

### Yearly Summary (D1:E20):
```
D1: YEARLY REPORT | =YEAR(TODAY())
D3: Total Income | =SUMIFS(Transactions!C:C,Transactions!B:B,"Income",Transactions!A:A,">="&DATE(YEAR(TODAY()),1,1))
D4: Total Expenses | =SUMIFS(Transactions!C:C,Transactions!B:B,"Expense",Transactions!A:A,">="&DATE(YEAR(TODAY()),1,1))
D5: Net Savings | =E3-E4
D6: Average Monthly Income | =E3/MONTH(TODAY())
D7: Average Monthly Expense | =E4/MONTH(TODAY())
```

### Account Summary (G1:H20):
```
=QUERY(Accounts!A:G,"SELECT B, D, F WHERE F='Active' ORDER BY D DESC LABEL B 'Account', D 'Balance', F 'Status'")
```

---

## 15. 📉 CHARTS

**Purpose:** Visual analytics

### Chart 1: Monthly Income vs Expenses (A1:C15)

Create a **Combo Chart**:
1. Insert → Chart
2. Chart type: Combo chart
3. Data range: Use QUERY to get monthly totals
4. Title: "Monthly Income vs Expenses"

**Query for data**:
```
=QUERY(Transactions!A:C,"SELECT MONTH(A)+1, SUM(C) WHERE B='Income' AND YEAR(A)=YEAR(TODAY()) GROUP BY MONTH(A)+1 LABEL MONTH(A)+1 'Month', SUM(C) 'Income'")
```

### Chart 2: Expense by Category (E1:G15)

Create a **Pie Chart**:
1. Insert → Chart
2. Chart type: Pie chart
3. Data range: Query for category totals
4. Title: "Expenses by Category"

**Query for data**:
```
=QUERY(Transactions!A:I,"SELECT D, SUM(C) WHERE B='Expense' AND A >= date '"&TEXT(EOMONTH(TODAY(),-1)+1,"yyyy-mm-dd")&"' GROUP BY D ORDER BY SUM(C) DESC LABEL D 'Category', SUM(C) 'Amount'")
```

### Chart 3: Net Worth Over Time (I1:K15)

Create a **Line Chart**:
1. Insert → Chart
2. Chart type: Line chart
3. Manual data entry or link to historical balances
4. Title: "Net Worth Trend"

### Chart 4: Budget Progress (M1:O15)

Create a **Bar Chart**:
1. Insert → Chart
2. Chart type: Bar chart
3. Data range: Budgets!B:F
4. Title: "Budget Progress"

---

## 🎨 STYLING & FORMATTING

### Color Scheme

**Primary Colors:**
- Headers: `#4A90E2` (Blue)
- Positive values: `#7ED321` (Green)
- Negative values: `#D0021B` (Red)
- Warnings: `#F5A623` (Orange)
- Info: `#9013FE` (Purple)

### Header Row Formatting (All Sheets):
- Background: `#4A90E2`
- Font: Bold, White, 11pt
- Alignment: Center
- Border: Bottom border

### Number Formatting:
- Currency: `$#,##0.00`
- Percentage: `0.00%`
- Date: `MM/DD/YYYY`

---

## 🔧 ADVANCED FEATURES

### 1. Automated Monthly Reports

**Apps Script** (Tools → Script Editor):

```javascript
function sendMonthlyReport() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const reportsSheet = ss.getSheetByName('Reports');

  const income = reportsSheet.getRange('B3').getValue();
  const expenses = reportsSheet.getRange('B4').getValue();
  const savings = reportsSheet.getRange('B5').getValue();

  const email = Session.getActiveUser().getEmail();
  const subject = 'Monthly Financial Report - ' + Utilities.formatDate(new Date(), Session.getScriptTimeZone(), 'MMMM yyyy');

  const body = `
    Monthly Financial Summary

    Total Income: $${income.toFixed(2)}
    Total Expenses: $${expenses.toFixed(2)}
    Net Savings: $${savings.toFixed(2)}

    View full report: ${ss.getUrl()}
  `;

  GmailApp.sendEmail(email, subject, body);
}
```

**Set up trigger:**
1. Tools → Script Editor
2. Edit → Current project's triggers
3. Add trigger: sendMonthlyReport, Time-driven, Month timer, On the 1st

### 2. Transaction Import from CSV

```javascript
function importTransactionsFromCSV() {
  const file = DriveApp.getFilesByName('transactions.csv').next();
  const csvData = Utilities.parseCsv(file.getBlob().getDataAsString());

  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const transactionsSheet = ss.getSheetByName('Transactions');

  // Skip header row
  for (let i = 1; i < csvData.length; i++) {
    transactionsSheet.appendRow(csvData[i]);
  }

  SpreadsheetApp.getUi().alert('Import complete! ' + (csvData.length - 1) + ' transactions imported.');
}
```

### 3. Budget Alert System

```javascript
function checkBudgetAlerts() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const budgetsSheet = ss.getSheetByName('Budgets');
  const data = budgetsSheet.getDataRange().getValues();

  let alerts = [];

  for (let i = 1; i < data.length; i++) {
    const progress = data[i][5]; // Column F
    const category = data[i][1]; // Column B

    if (progress > 1.0) {
      alerts.push(`⚠️ OVER BUDGET: ${category} (${(progress*100).toFixed(0)}%)`);
    } else if (progress > 0.8) {
      alerts.push(`⚠️ WARNING: ${category} at ${(progress*100).toFixed(0)}%`);
    }
  }

  if (alerts.length > 0) {
    const email = Session.getActiveUser().getEmail();
    GmailApp.sendEmail(email, 'Budget Alert!', alerts.join('\n'));
  }
}
```

### 4. Duplicate Transaction Detection

```javascript
function findDuplicates() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const transactionsSheet = ss.getSheetByName('Transactions');
  const data = transactionsSheet.getDataRange().getValues();

  let seen = {};
  let duplicates = [];

  for (let i = 1; i < data.length; i++) {
    const key = data[i][0] + '|' + data[i][2] + '|' + data[i][5]; // Date|Amount|Description

    if (seen[key]) {
      duplicates.push(i + 1); // Row number
    } else {
      seen[key] = true;
    }
  }

  if (duplicates.length > 0) {
    SpreadsheetApp.getUi().alert('Possible duplicates found at rows: ' + duplicates.join(', '));
  } else {
    SpreadsheetApp.getUi().alert('No duplicates found!');
  }
}
```

---

## 📱 MOBILE ACCESS

### Google Sheets Mobile App Features:
- ✅ View all data on mobile
- ✅ Add transactions on-the-go
- ✅ Check balances
- ✅ View charts and reports
- ✅ Offline access

### Quick Add Transaction (Mobile):
1. Open Transactions sheet
2. Tap (+) button
3. Fill in: Date, Type, Amount, Category, Account, Description
4. Save

---

## 🔐 SECURITY & SHARING

### Protection Settings:

1. **Protect sensitive sheets:**
   - Profile: Protect (only you can edit)
   - Accounts: Protect with warning
   - Credit_Cards: Protect (only you can edit)

2. **Share with family member:**
   - Share → Get link
   - Change to "Anyone with the link can view"
   - Or: Specific people as "Viewer" or "Commenter"

3. **Make a copy for backup:**
   - File → Make a copy
   - Rename: "Daybook Backup - [Date]"

---

## 📤 EXPORT & IMPORT

### Export Data:
1. **Individual Sheet**: File → Download → CSV
2. **Entire Workbook**: File → Download → Excel (.xlsx)
3. **PDF Report**: File → Download → PDF

### Import Data:
1. **From CSV**: File → Import → Upload
2. **From Excel**: File → Import → Upload
3. **From Another Sheet**: File → Import → Google Sheets

---

## 🎯 TIPS & BEST PRACTICES

### Daily:
- [ ] Log all transactions
- [ ] Check account balances
- [ ] Review pending bills

### Weekly:
- [ ] Categorize transactions
- [ ] Review budget progress
- [ ] Check for duplicate entries

### Monthly:
- [ ] Review monthly report
- [ ] Update budget for next month
- [ ] Pay credit card bills
- [ ] Review goals progress
- [ ] Archive old data

### Quarterly:
- [ ] Review and adjust budgets
- [ ] Analyze spending trends
- [ ] Update financial goals
- [ ] Clean up old/closed accounts

### Yearly:
- [ ] Generate annual report
- [ ] File for taxes
- [ ] Set new financial goals
- [ ] Review overall financial health

---

## 🆘 TROUBLESHOOTING

### Formulas Not Working:
1. Check cell references
2. Verify sheet names match exactly
3. Look for circular references
4. Check data validation settings

### Charts Not Updating:
1. Refresh the page
2. Update chart data range
3. Check source data for errors

### Import Issues:
1. Check CSV format
2. Verify column order matches
3. Remove special characters

### Performance Slow:
1. Limit rows to current year only
2. Archive old data to separate sheet
3. Remove unnecessary formulas
4. Use QUERY instead of multiple formulas

---

## 📊 SAMPLE DATA INCLUDED

The template includes sample data for:
- ✅ 3 bank accounts
- ✅ 20+ transactions
- ✅ 4 recurring transactions
- ✅ 15+ expense categories
- ✅ 5 budget categories
- ✅ 4 financial goals
- ✅ 2 credit cards
- ✅ 4 bills

**To use your own data:**
1. Clear sample rows (keep headers)
2. Start entering your real data
3. Formulas will auto-update

---

## 🚀 GETTING STARTED CHECKLIST

- [ ] Create new Google Sheet
- [ ] Create all 15 sheets
- [ ] Set up column headers
- [ ] Add data validation
- [ ] Create formulas
- [ ] Apply conditional formatting
- [ ] Insert charts
- [ ] Add sample data
- [ ] Test all formulas
- [ ] Customize for your needs
- [ ] Start tracking!

---

## 📞 SUPPORT

### Resources:
- Google Sheets Help: https://support.google.com/sheets
- Formula Reference: https://support.google.com/docs/table/25273
- Apps Script: https://developers.google.com/apps-script

### Common Questions:

**Q: Can I use this offline?**
A: Yes, enable offline mode in Google Drive settings.

**Q: Can multiple people edit simultaneously?**
A: Yes, but be careful with formulas.

**Q: How do I backup my data?**
A: File → Make a copy regularly.

**Q: Can I connect to my bank?**
A: No, you need to manually enter or import transactions.

---

## 🎉 YOU'RE READY!

This Google Sheets template replicates all major features of the Daybook application:
- ✅ Account management
- ✅ Transaction tracking
- ✅ Budget planning
- ✅ Goal tracking
- ✅ Credit card management
- ✅ Bill tracking
- ✅ Financial reports
- ✅ Visual analytics

**Start tracking your finances today!** 📊💰
