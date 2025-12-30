# Export Feature Integration Guide

The data export feature is now fully integrated and ready to use!

## 🚀 Services Running

- **Backend API**: http://localhost:8080 ✅
- **Frontend App**: http://localhost:3000 ✅

## 📋 How to Test the Export Feature

### 1. Access the Settings Page

1. Open your browser and go to: **http://localhost:3000**
2. Login with your credentials
3. Navigate to **Settings** from the sidebar menu
4. Scroll down to the **Data Management** section

### 2. Available Export Options

You'll see **6 export cards**:

#### **Transactions**
- Export all your income, expenses, and transfers
- Formats: CSV or JSON
- Date Range: Last 12 months (automatically set)

#### **Accounts**
- Export all accounts and their current balances
- Formats: CSV or JSON

#### **Budgets**
- Export all budget configurations
- Formats: CSV or JSON

#### **Goals**
- Export all savings and investment goals
- Formats: CSV or JSON

#### **Categories**
- Export all income and expense categories
- Formats: CSV or JSON

#### **All Data** (Complete Backup)
- Exports everything in one file:
  - All accounts
  - All categories
  - Last 2 years of transactions
  - All budgets
  - All goals
  - All credit cards
  - All debts/lends
- Format: JSON only
- Perfect for complete data backup

### 3. How to Export

1. **Choose a data type** (e.g., Transactions)
2. **Click CSV or JSON button**
3. The file will automatically download
4. Check your Downloads folder

### 4. Example Filenames

Exported files are automatically named:

```
transactions_2024-12-30_2025-12-30.csv
transactions_2024-12-30_2025-12-30.json
accounts.csv
accounts.json
budgets.csv
budgets.json
goals.csv
goals.json
categories.csv
categories.json
daybook_export_2025-12-30.json (for "All Data")
```

## 🔍 Testing the API Directly

You can also test the export endpoints directly using curl:

### Get an Auth Token First

```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"yourpassword"}' \
  | jq -r '.token'
```

Copy the token and use it in the examples below (replace `YOUR_TOKEN`).

### Export Transactions as CSV

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/export?type=transactions&format=csv&start_date=2024-01-01&end_date=2025-12-30" \
  -o transactions.csv
```

### Export Accounts as JSON

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/export?type=accounts&format=json" \
  -o accounts.json
```

### Export All Data (Complete Backup)

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/export?type=all&format=json" \
  -o complete_backup.json
```

## 📊 Export Data Structures

### CSV Format (Example: Transactions)

```csv
ID,Date,Type,Amount,Category ID,Account ID,To Account ID,Description,Tags,Created At
1,2025-01-15,expense,50.00,5,2,,Groceries,"[""food"",""monthly""]",2025-01-15 10:30:00
2,2025-01-16,income,1000.00,1,2,,Salary,"[""income""]",2025-01-16 09:00:00
```

### JSON Format (Example: Accounts)

```json
[
  {
    "id": 1,
    "userId": 123,
    "name": "Main Checking",
    "type": "checking",
    "balance": 5000.50,
    "currency": "BDT",
    "institution": "ABC Bank",
    "accountNumber": "****1234",
    "description": "Primary account",
    "createdAt": "2025-01-01T00:00:00Z"
  }
]
```

### Complete Backup (All Data JSON)

```json
{
  "export_date": "2025-12-30T21:38:00Z",
  "export_version": "1.0",
  "user_id": 123,
  "data": {
    "accounts": [...],
    "categories": [...],
    "transactions": [...],
    "budgets": [...],
    "goals": [...],
    "credit_cards": [...],
    "debts": [...],
    "lends": [...]
  },
  "metadata": {
    "accounts_count": 5,
    "categories_count": 25,
    "transactions_count": 1234,
    "budgets_count": 10,
    "goals_count": 3,
    "credit_cards_count": 2,
    "debts_count": 1,
    "lends_count": 0,
    "date_range": {
      "start": "2023-12-30",
      "end": "2025-12-30"
    }
  }
}
```

## 🎨 UI Preview

The Settings page now has a beautiful Data Management section with:

```
┌─────────────────────────────────────────────────────────┐
│                    Data Management                       │
├─────────────────────────────────────────────────────────┤
│  Export Your Data                                        │
│  Download your financial data in CSV or JSON format     │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ Transactions │  │   Accounts   │  │   Budgets    │ │
│  │ CSV  │ JSON  │  │ CSV  │ JSON  │  │ CSV  │ JSON  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │    Goals     │  │  Categories  │  │  🌟 All Data │ │
│  │ CSV  │ JSON  │  │ CSV  │ JSON  │  │     JSON     │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## ✅ Features Implemented

- ✅ 6 different export types
- ✅ CSV and JSON formats
- ✅ Automatic file downloads
- ✅ Proper authentication
- ✅ Activity logging (all exports are logged)
- ✅ User-specific data isolation
- ✅ Descriptive filenames
- ✅ Success/error notifications
- ✅ Complete data backup option

## 🔒 Security

- All exports require valid JWT authentication
- Users can only export their own data
- Activity logging tracks all export operations
- No sensitive data in URLs or filenames

## 🐛 Troubleshooting

### Export Button Not Working

1. Make sure you're logged in
2. Check browser console for errors (F12 → Console)
3. Verify backend is running: `curl http://localhost:8080/health`

### Download Not Starting

1. Check if browser is blocking downloads
2. Allow downloads from localhost:3000
3. Check browser's download settings

### Getting 401 Unauthorized

1. Your session may have expired - login again
2. Check if token exists: `localStorage.getItem('token')` in browser console

### Export Returns Empty Data

1. Make sure you have data in the selected category
2. For transactions, check the date range
3. Verify data exists in the database

## 📝 Next Steps

You can now:

1. **Test the feature** in the browser at http://localhost:3000
2. **Export your data** for backup or analysis
3. **Use CSV files** in Excel, Google Sheets, etc.
4. **Use JSON files** for programmatic processing
5. **Schedule regular backups** using the "All Data" export

## 🎯 Use Cases

### Backup Your Data
Use the "All Data" export to create complete backups of your financial information.

### Tax Preparation
Export transactions for a specific year to CSV and open in Excel for tax filing.

### Data Migration
Export all data as JSON to migrate to another system or restore later.

### Financial Analysis
Export transactions to CSV and analyze in spreadsheet software with pivot tables and charts.

### Audit Trail
Export account history to verify balances and reconcile accounts.

---

**Need Help?**
If you encounter any issues, check the browser console (F12) for error messages or review the backend logs at `/tmp/daybook-backend.log`.
