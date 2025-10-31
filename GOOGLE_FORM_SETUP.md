# Add Transaction via Google Form

Create a quick-entry form for adding transactions to your Daybook spreadsheet!

---

## 🚀 **Quick Setup (5 Minutes)**

### **Step 1: Create Google Form**

1. Go to: https://forms.google.com
2. Click: **"+ Blank form"**
3. Title: **"Add Transaction - Daybook"**
4. Description: **"Quick entry for income and expenses"**

---

### **Step 2: Add Form Fields**

Add these questions in order:

#### **1. Date** (Date)
- Question: **"Transaction Date"**
- Type: **Date**
- Required: ✅ Yes

#### **2. Type** (Multiple choice)
- Question: **"Type"**
- Type: **Multiple choice**
- Options:
  ```
  ○ Income
  ○ Expense
  ○ Transfer
  ```
- Required: ✅ Yes

#### **3. Amount** (Short answer)
- Question: **"Amount"**
- Type: **Short answer**
- Validation: **Number** (Greater than 0)
- Required: ✅ Yes

#### **4. Category** (Dropdown)
- Question: **"Category"**
- Type: **Dropdown**
- Options:
  ```
  Salary
  Freelance
  Rent
  Utilities
  Groceries
  Dining Out
  Gas
  Public Transit
  Shopping
  Entertainment
  Health
  Subscription
  Other
  ```
- Required: ✅ Yes

#### **5. Account** (Dropdown)
- Question: **"Account"**
- Type: **Dropdown**
- Options:
  ```
  Chase Checking
  Savings Account
  Investment Account
  Cash Wallet
  Credit Card
  ```
- Required: ✅ Yes

#### **6. Description** (Short answer)
- Question: **"Description"**
- Type: **Short answer**
- Required: ❌ No

#### **7. Tags** (Short answer)
- Question: **"Tags (optional)"**
- Type: **Short answer**
- Placeholder: **"e.g., Food:Groceries"**
- Required: ❌ No

#### **8. Recurring** (Multiple choice)
- Question: **"Is this recurring?"**
- Type: **Multiple choice**
- Options:
  ```
  ○ Yes
  ○ No
  ```
- Default: **No**
- Required: ✅ Yes

---

### **Step 3: Configure Settings**

1. Click **Settings** (gear icon)
2. **Presentation tab:**
   - ✅ Show progress bar
   - ✅ Confirmation message: **"Transaction added successfully!"**
   - ✅ Show link to submit another response

3. **Responses tab:**
   - Select: **"Create new spreadsheet"** OR **"Select existing spreadsheet"**
   - If selecting existing: Choose your **Daybook_Financial_Tracker**

---

### **Step 4: Link to Your Spreadsheet**

#### **Method A: New Responses Sheet (Recommended)**

1. Click **Responses** tab in form
2. Click green Sheets icon
3. Select: **"Create a new spreadsheet"**
4. Name: **"Transaction Form Responses"**
5. Click **Create**

#### **Method B: Link to Existing Sheet**

1. Click **Responses** tab
2. Click three dots → **Select response destination**
3. Choose: **"Select existing spreadsheet"**
4. Select your **Daybook_Financial_Tracker**
5. Responses go to new "Form Responses 1" sheet

---

### **Step 5: Auto-Transfer to Transactions Sheet**

Add this **Apps Script** to automatically transfer form data to Transactions sheet:

1. **Open your spreadsheet**
2. **Extensions → Apps Script**
3. **Paste this code:**

```javascript
function onFormSubmit(e) {
  try {
    const ss = SpreadsheetApp.getActiveSpreadsheet();
    const formSheet = ss.getSheetByName('Form Responses 1'); // Form responses
    const transactionsSheet = ss.getSheetByName('Transactions'); // Your Transactions sheet

    if (!formSheet || !transactionsSheet) {
      Logger.log('Required sheets not found');
      return;
    }

    // Get the last submitted row
    const lastRow = formSheet.getLastRow();
    const formData = formSheet.getRange(lastRow, 2, 1, 8).getValues()[0]; // Skip timestamp column

    // Prepare transaction data (matching Transactions sheet columns)
    const transactionRow = [
      formData[0],  // Date
      formData[1],  // Type
      formData[2],  // Amount
      formData[3],  // Category
      formData[4],  // Account
      formData[5],  // Description
      formData[6],  // Tags
      formData[7],  // Recurring
      'Completed'   // Status (auto-set to Completed)
    ];

    // Add to Transactions sheet
    transactionsSheet.appendRow(transactionRow);

    // Optional: Delete from form responses (to keep it clean)
    // formSheet.deleteRow(lastRow);

    Logger.log('Transaction added successfully');

  } catch (error) {
    Logger.log('Error: ' + error.toString());
  }
}

// Set up the trigger
function createFormSubmitTrigger() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();

  // Delete existing triggers (to avoid duplicates)
  const triggers = ScriptApp.getUserTriggers(ss);
  triggers.forEach(trigger => ScriptApp.deleteTrigger(trigger));

  // Create new trigger
  ScriptApp.newTrigger('onFormSubmit')
    .forSpreadsheet(ss)
    .onFormSubmit()
    .create();

  SpreadsheetApp.getUi().alert('Trigger created successfully!');
}
```

4. **Save**: Click disk icon or `Ctrl+S`
5. **Run `createFormSubmitTrigger`**:
   - Select: `createFormSubmitTrigger` from dropdown
   - Click: **Run** (play button)
   - **Authorize**: Click "Review Permissions" → Select your account → Allow

6. **Done!** Now form submissions auto-add to Transactions sheet

---

## 📱 **Mobile Quick-Add**

### **Get Direct Link:**

1. **In your form**, click **Send**
2. Click **link icon** (🔗)
3. Check: **"Shorten URL"**
4. **Copy link**

### **Add to Home Screen (iOS/Android):**

**iOS:**
1. Open link in Safari
2. Tap **Share** button
3. Select **"Add to Home Screen"**
4. Name: **"Add Transaction"**
5. Tap **Add**

**Android:**
1. Open link in Chrome
2. Tap **menu** (⋮)
3. Select **"Add to Home screen"**
4. Name: **"Add Transaction"**
5. Tap **Add**

Now you have a one-tap transaction entry app! 📱💰

---

## 🎯 **Usage**

### **Desktop:**
1. Open form link
2. Fill in transaction details
3. Click **Submit**
4. Done! Check your Transactions sheet

### **Mobile:**
1. Tap home screen icon
2. Fill form (takes 30 seconds)
3. Submit
4. Transaction automatically added to spreadsheet

---

## 🔧 **Advanced Features**

### **1. Add Email Notification**

Add this to your script:

```javascript
function onFormSubmit(e) {
  // ... existing code ...

  // Send confirmation email
  const email = Session.getActiveUser().getEmail();
  const subject = 'Transaction Added';
  const body = `
    Transaction added successfully!

    Date: ${formData[0]}
    Type: ${formData[1]}
    Amount: $${formData[2]}
    Category: ${formData[3]}

    View: ${ss.getUrl()}
  `;

  GmailApp.sendEmail(email, subject, body);
}
```

### **2. Add Data Validation**

In the form, add validation:
- **Amount**: Must be number > 0
- **Date**: Must be within last 30 days (or future for planned)

### **3. Add Conditional Logic**

Make fields show/hide based on Type:
- If "Income" → Show income categories only
- If "Expense" → Show expense categories only

Steps:
1. Click on **Type** question
2. Click **⋮** (three dots)
3. Select **"Go to section based on answer"**
4. Create sections for Income/Expense categories

### **4. Pre-fill Form**

Create quick links with pre-filled values:

Example: **Quick Grocery Entry**
```
https://docs.google.com/forms/d/YOUR_FORM_ID/viewform?entry.123456=Groceries&entry.789012=Expense
```

Get entry IDs:
1. Open form
2. Click **⋮** → **Get pre-filled link**
3. Fill sample values
4. Copy link

---

## 📊 **Form Response Sheet Columns**

Your form creates these columns:

```
A: Timestamp (auto)
B: Transaction Date
C: Type
D: Amount
E: Category
F: Account
G: Description
H: Tags
I: Recurring
```

The script maps these to Transactions sheet:

```
A: Date (from B)
B: Type (from C)
C: Amount (from D)
D: Category (from E)
E: Account (from F)
F: Description (from G)
G: Tags (from H)
H: Recurring (from I)
I: Status (auto: "Completed")
```

---

## 🎨 **Customize Form**

### **Theme:**
1. Click **Customize theme** (palette icon)
2. Choose colors
3. Add header image

### **Sections:**
1. Add sections for better organization
2. Group related questions

### **Logic:**
1. Skip sections based on answers
2. Show relevant questions only

---

## 🔐 **Privacy & Sharing**

### **Form Permissions:**

**Public Form** (Anyone can submit):
1. Click **Send**
2. Anyone with link can submit
3. No Google account needed

**Restricted** (Only your organization):
1. **Settings** → **Responses**
2. Check: **"Limit to 1 response"** (optional)
3. Check: **"Collect email addresses"** (optional)

### **Spreadsheet Permissions:**

Keep spreadsheet private:
- Only YOU have edit access
- Form submitters don't see spreadsheet
- They only fill the form

---

## 📱 **Widget for iPhone/Android**

### **iOS Shortcuts:**

Create a Siri Shortcut:
1. Open **Shortcuts** app
2. Create new shortcut
3. Add action: **"Open URL"**
4. URL: Your form link
5. Name: **"Add Transaction"**
6. Say: **"Hey Siri, add transaction"**

### **Android Widgets:**

1. Long-press home screen
2. Add widget
3. Select **Chrome** bookmarks widget
4. Select your form bookmark

---

## 🎯 **Best Practices**

### **Daily Workflow:**

**Morning:**
- Open form on mobile
- Already have it handy

**Throughout Day:**
- Spent money? → Open form → 30 seconds to log
- Multiple purchases? → Submit multiple times

**Evening:**
- Review Transactions sheet
- Check if all entries are correct

### **Tips:**

1. **Use Categories Consistently**
   - Always use "Groceries" not "Grocery" or "groceries"
   - This helps with reports

2. **Add Description**
   - Helps remember later
   - "Target" vs "Walmart" vs "Costco"

3. **Use Tags for Detail**
   - "Food:Groceries:Weekly"
   - "Transportation:Gas:Commute"

4. **Regular Entry**
   - Daily entry > Weekly catchup
   - Less to remember

---

## 🐛 **Troubleshooting**

### **Script Not Running:**

1. **Check trigger:**
   - Extensions → Apps Script
   - Triggers (clock icon on left)
   - Should see `onFormSubmit` trigger

2. **Check permissions:**
   - May need to re-authorize
   - Run `createFormSubmitTrigger` again

3. **Check logs:**
   - In Apps Script, click **Executions**
   - See any errors

### **Data Not Transferring:**

1. **Check sheet names:**
   - Form responses: "Form Responses 1"
   - Target: "Transactions"
   - Must match exactly!

2. **Check column order:**
   - Form questions order = column order
   - Script expects specific order

3. **Test manually:**
   - Submit test transaction
   - Check if it appears in Transactions

### **Form Not Found:**

1. **Check link:**
   - Must be form edit link for editing
   - Use response link for submitting

2. **Check permissions:**
   - Make sure form is published
   - Check who can respond

---

## 📚 **Summary**

### **What You Get:**

✅ **Quick-entry form** - Add transactions in 30 seconds
✅ **Mobile-friendly** - Works on phone
✅ **Auto-sync** - Automatically adds to spreadsheet
✅ **No duplicates** - Each submission = one transaction
✅ **Always available** - Cloud-based, access anywhere
✅ **Home screen icon** - One-tap access
✅ **Email notifications** - Optional confirmations

### **Setup Time:**
- Form creation: 5 minutes
- Script setup: 3 minutes
- **Total: 8 minutes** for lifetime of easy entry!

---

## 🎉 **You're Done!**

Now you can add transactions from anywhere:
- 📱 Phone (home screen icon)
- 💻 Desktop (bookmark)
- 🗣️ Voice (Siri/Google Assistant)

**No more excuses for not tracking expenses!** 😊💰📊

---

## 📎 **Quick Links**

- Create Form: https://forms.google.com
- Your Spreadsheet: (your Google Sheets link)
- Apps Script: Extensions → Apps Script
- Form Responses: Check "Form Responses 1" sheet

**Happy tracking!** 🎯
