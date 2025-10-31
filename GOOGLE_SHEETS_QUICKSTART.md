# Daybook Google Sheets - Quick Start (5 Minutes)

Get your financial tracker up and running in 5 minutes!

## 🚀 Step-by-Step Setup

### 1. Create the Spreadsheet (30 seconds)

1. Go to https://sheets.google.com
2. Click **"+ Blank"** to create new sheet
3. Name it: **"Daybook Financial Tracker"**

### 2. Create Sheets (1 minute)

Click the **"+"** at bottom left to add sheets. Create these 15 sheets:

```
1.  Dashboard
2.  Profile
3.  Accounts
4.  Transactions
5.  Recurring
6.  Tags
7.  Budgets
8.  Goals
9.  Credit_Cards
10. CC_Transactions
11. CC_Payments
12. Bills
13. Bill_Payments
14. Reports
15. Charts
```

**Tip:** Right-click on "Sheet1" → Rename → "Dashboard"

### 3. Set Up Core Sheets (3 minutes)

#### 📊 TRANSACTIONS (Most Important!)

**Row 1 Headers:**
```
A: Date | B: Type | C: Amount | D: Category | E: Account | F: Description | G: Tags | H: Recurring | I: Status
```

**Add Sample Data (Row 2-5):**
```
01/15/2024 | Income | 5000 | Salary | Checking | Monthly salary | | No | Completed
01/16/2024 | Expense | 1200 | Rent | Checking | Rent payment | | Yes | Completed
01/17/2024 | Expense | 150 | Groceries | Checking | Weekly groceries | | No | Completed
01/18/2024 | Expense | 50 | Gas | Checking | Gas station | | No | Completed
```

**Add Data Validation:**
- Column B: Data → Data validation → Dropdown: `Income, Expense, Transfer`
- Column I: Data → Data validation → Dropdown: `Completed, Pending, Cancelled`

#### 💰 ACCOUNTS

**Row 1 Headers:**
```
A: ID | B: Account Name | C: Type | D: Balance | E: Last Updated | F: Status
```

**Add Your Accounts (Row 2-4):**
```
1 | Checking Account | Checking | 5000 | =TODAY() | Active
2 | Savings Account | Savings | 15000 | =TODAY() | Active
3 | Cash | Cash | 200 | =TODAY() | Active
```

**Add Total (Row 5):**
```
| TOTAL | | =SUM(D2:D4) | |
```

#### 🏠 DASHBOARD

**Row 1-10:**
```
A: Metric | B: Value

Total Net Worth | =SUM(Accounts!D:D)
Income (This Month) | =SUMIFS(Transactions!C:C,Transactions!B:B,"Income",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
Expenses (This Month) | =SUMIFS(Transactions!C:C,Transactions!B:B,"Expense",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
Net Savings | =B3-B4
Transaction Count | =COUNTA(Transactions!A:A)-1
```

**Format B2-B4 as Currency:** Format → Number → Currency

### 4. Start Using! (1 minute)

Go to **Transactions** sheet and add your first real transaction!

---

## ✨ Essential Formulas

Copy-paste these into the indicated cells:

### Dashboard Summary (Cell B2-B6)

```
B2: =SUM(Accounts!D:D)
B3: =SUMIFS(Transactions!C:C,Transactions!B:B,"Income",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
B4: =SUMIFS(Transactions!C:C,Transactions!B:B,"Expense",Transactions!A:A,">="&EOMONTH(TODAY(),-1)+1)
B5: =B3-B4
B6: =COUNTA(Transactions!A:A)-1
```

### Account Total (Accounts sheet, D5)

```
=SUM(D2:D4)
```

---

## 📱 Using on Mobile

1. Install **Google Sheets** app
2. Open your "Daybook Financial Tracker"
3. **Add transaction:** Tap (+) on Transactions sheet
4. **Check balance:** View Dashboard or Accounts

---

## 🎯 Daily Workflow

### Morning:
1. Open **Dashboard** - Check balances
2. Review any pending transactions

### When You Spend:
1. Open **Transactions** sheet
2. Add new row:
   - Date: Today
   - Type: Income or Expense
   - Amount: Enter amount
   - Category: Select category
   - Account: Select account
   - Description: What it was for
   - Status: Completed

### Evening:
1. Review today's transactions
2. Check if all entries are correct

---

## 🔧 Next Steps (Optional)

Once comfortable with basics, add:

### Week 1:
- [ ] Set up **Budgets** sheet
- [ ] Add more **Tags/Categories**
- [ ] Create a simple chart

### Week 2:
- [ ] Set up **Goals** sheet
- [ ] Add **Recurring** transactions
- [ ] Try the Reports sheet

### Week 3:
- [ ] Add **Credit Cards** (if you have them)
- [ ] Set up **Bills** tracking
- [ ] Explore automation with Apps Script

---

## 💡 Pro Tips

1. **Use Shortcuts:**
   - `Ctrl/Cmd + ;` = Insert today's date
   - `Ctrl/Cmd + D` = Fill down
   - `Alt + Shift + 1` = Insert row

2. **Quick Entry:**
   - Keep Transactions sheet always open
   - Use mobile app for on-the-go entries
   - Set up keyboard shortcuts for common entries

3. **Stay Organized:**
   - Use consistent category names
   - Enter transactions daily
   - Review weekly

4. **Backup:**
   - File → Make a copy (monthly)
   - File → Version history (automatic)

---

## 🆘 Common Issues

### Formula shows #REF!
- Check sheet names are spelled correctly
- Make sure referenced sheets exist

### Numbers showing as text
- Format → Number → Number or Currency

### Date not working
- Format → Number → Date
- Enter as MM/DD/YYYY

### Sum not calculating
- Check if cells contain actual numbers, not text
- Remove any spaces or special characters

---

## 📊 What's Included

**Fully Set Up:**
- ✅ Dashboard with key metrics
- ✅ Account tracking
- ✅ Transaction logging
- ✅ Basic formulas

**Ready to Add:**
- ⏳ Budget tracking
- ⏳ Goal setting
- ⏳ Credit card management
- ⏳ Bill tracking
- ⏳ Advanced reports

**Optional:**
- 💡 Automated reports
- 💡 Email alerts
- 💡 Advanced charts
- 💡 Import from CSV

---

## 🎉 You're Ready!

You now have a working financial tracker!

**Start simple:**
1. Log your transactions daily
2. Check your dashboard weekly
3. Add features as you need them

**Full documentation:** See `GOOGLE_SHEETS_TEMPLATE.md`

---

## 📋 Copy-Paste Headers

For quick setup, copy these headers:

### Transactions:
```
Date	Type	Amount	Category	Account	Description	Tags	Recurring	Status
```

### Accounts:
```
ID	Account Name	Type	Balance	Last Updated	Status
```

### Dashboard:
```
Metric	Value
```

### Budgets:
```
Month	Category	Budgeted	Spent	Remaining	Progress %	Status
```

### Goals:
```
ID	Goal Name	Target Amount	Current Amount	Progress %	Deadline	Monthly Target	Status
```

**Just paste into Row 1 of each sheet!**

---

## ⏱️ Time Commitment

- **Setup:** 5 minutes (one time)
- **Daily use:** 2 minutes (log transactions)
- **Weekly review:** 5 minutes
- **Monthly review:** 15 minutes

**Total:** ~30 minutes per month for complete financial visibility!

---

## 🚀 Ready to Track!

Open your Google Sheet and start adding transactions. It's that simple!

Questions? See the full guide: `GOOGLE_SHEETS_TEMPLATE.md`
