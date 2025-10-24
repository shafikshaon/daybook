# Navigation Structure

## Menu Grouping (Standard Personal Finance App)

The sidebar navigation is now organized like professional finance trackers (Mint, YNAB, Personal Capital, etc.)

### 📊 Overview
Quick access to the main dashboard.

**Menu Items:**
- Dashboard - Overview of all finances

---

### 🏦 Accounts & Transactions
Day-to-day money management.

**Menu Items:**
- Accounts - All bank accounts, cash, etc.
- Transactions - Income and expense tracking
- Credit Cards - Credit card management

**Use Case:**
- Daily transaction entry
- Account balance checks
- Credit card payments

---

### 📈 Planning & Budgets
Financial planning and budget management.

**Menu Items:**
- Budgets - Category-based budgets
- Bills & Reminders - Recurring bills
- Savings Goals - Target-based savings

**Use Case:**
- Monthly budget planning
- Bill payment tracking
- Goal progress monitoring

---

### 📊 Investments & Assets
Long-term wealth building.

**Menu Items:**
- Investments - Stocks, ETFs, Crypto
- Fixed Deposits - Time-locked deposits

**Use Case:**
- Portfolio tracking
- Investment performance
- Asset allocation

---

### 📑 Reports & Analytics
Financial insights and analysis.

**Menu Items:**
- Reports - Analytics and trends

**Use Case:**
- Net worth tracking
- Spending analysis
- Financial trends

---

### ⚙️ Settings
Application configuration.

**Menu Items:**
- Settings - App preferences

**Use Case:**
- Currency settings
- Dark mode
- Data export

---

## Visual Grouping

### Expanded Sidebar (220px)
```
┌─────────────────────┐
│ 💰 Expense Tracker │
├─────────────────────┤
│ 📊 Dashboard        │
│                     │
│ ACCOUNTS & TRANS... │
│ 🏦 Accounts         │
│ 💳 Transactions     │
│ 💳 Credit Cards     │
│                     │
│ PLANNING & BUDGETS  │
│ 📈 Budgets          │
│ 📄 Bills & Reminders│
│ 🎯 Savings Goals    │
│                     │
│ INVESTMENTS & ASS...│
│ 📊 Investments      │
│ 💰 Fixed Deposits   │
│                     │
│ REPORTS & ANALYTICS │
│ 📑 Reports          │
│                     │
│ SETTINGS            │
│ ⚙️ Settings         │
└─────────────────────┘
```

### Collapsed Sidebar (60px)
```
┌────┐
│ 💰 │
├────┤
│ 📊 │
├────┤
│ 🏦 │
│ 💳 │
│ 💳 │
├────┤
│ 📈 │
│ 📄 │
│ 🎯 │
├────┤
│ 📊 │
│ 💰 │
├────┤
│ 📑 │
├────┤
│ ⚙️ │
└────┘
```

- Section titles hidden
- Divider lines shown
- Icons only
- Tooltips on hover

---

## Design Principles

### 1. **Logical Grouping**
Related features grouped together for intuitive navigation.

### 2. **Visual Hierarchy**
- Section titles: Uppercase, muted
- Menu items: Regular weight
- Active state: Highlighted background

### 3. **Progressive Disclosure**
- Expanded: Full labels and context
- Collapsed: Icons with tooltips

### 4. **Consistent with Industry**
Follows patterns from popular finance apps:
- **Mint**: Dashboard → Transactions → Budgets → Goals → Trends
- **YNAB**: Budget → Accounts → Reports
- **Personal Capital**: Dashboard → Accounts → Investing → Planning

---

## Styling Details

### Section Titles (Expanded Mode)
```scss
.nav-section-title {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.6875rem; // 11px
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 1rem 1rem 0.375rem 1rem;
  margin-top: 0.5rem;
}
```

### Dividers (Collapsed Mode)
```scss
.nav-divider {
  height: 1px;
  background-color: rgba(255, 255, 255, 0.2);
  margin: 0.625rem 0.75rem;
}
```

### Menu Links
```scss
.nav-link {
  color: white;
  padding: 0.625rem 0.875rem;
  border-radius: 0.375rem;
  margin: 0.125rem 0.5rem;
  font-size: 0.875rem;

  &:hover, &.active {
    background-color: #8b5cf6; // Lighter purple
  }
}
```

---

## Scrollable Sidebar

The sidebar is now scrollable with custom purple-themed scrollbar:

**Features:**
- Auto-scroll when many menu items
- Thin 6px scrollbar
- Semi-transparent styling
- Hover effect on scrollbar

**CSS:**
```scss
.sidebar {
  height: calc(100vh - 56px);
  overflow-y: auto;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }
}
```

---

## Responsive Behavior

### Desktop (> 768px)
- Full sidebar with labels (220px)
- Section titles visible
- Can toggle to icon-only (60px)

### Tablet (768px - 1024px)
- Icon-only by default (60px)
- Expands on toggle
- Dividers instead of titles

### Mobile (< 768px)
- Icon-only (60px)
- Overlay when expanded
- Touch-friendly tap targets

---

## Future Enhancements

### Possible Additions

**Accounts & Transactions:**
- Add "All Transactions" submenu
- Add "Categories" submenu

**Planning & Budgets:**
- Add "Debt Payoff" menu item
- Add "Spending Plan" menu item

**Investments & Assets:**
- Add "Real Estate" menu item
- Add "Retirement" menu item

**Reports & Analytics:**
- Add "Net Worth" submenu
- Add "Cash Flow" submenu
- Add "Tax Reports" submenu

### Collapsible Sections
Could add expand/collapse functionality for each section:

```
▼ ACCOUNTS & TRANSACTIONS
  🏦 Accounts
  💳 Transactions
  💳 Credit Cards
```

### Badge/Counter Support
Show counts on menu items:

```
💳 Transactions (24 new)
📄 Bills & Reminders (3 due)
```

---

## Implementation

### Toggle Sidebar
```vue
<button @click="toggleSidebar()">☰</button>
```

### Active State
```vue
<router-link
  to="/budgets"
  class="nav-link"
  active-class="active"
>
  <i>📈</i>
  <span>Budgets</span>
</router-link>
```

### Conditional Rendering
```vue
<!-- Show title when expanded -->
<li v-if="showSidebar" class="nav-section-title">
  Accounts & Transactions
</li>

<!-- Show divider when collapsed -->
<li v-if="!showSidebar" class="nav-divider"></li>
```

---

## Comparison with Other Apps

| App | Structure | Our Implementation |
|-----|-----------|-------------------|
| **Mint** | Overview → Transactions → Budgets → Goals → Trends | ✅ Similar grouping |
| **YNAB** | Budget → Accounts → Reports | ✅ Budget-focused option available |
| **Personal Capital** | Dashboard → Net Worth → Investing → Planning | ✅ Investment tracking included |
| **Quicken** | Home → Accounts → Bills → Budgets → Investing | ✅ All features covered |

---

## User Benefits

### 1. **Faster Navigation**
Related features grouped together reduces clicks.

### 2. **Better Mental Model**
Clear sections match how users think about finances.

### 3. **Cleaner Interface**
Organized structure reduces cognitive load.

### 4. **Industry Standard**
Familiar to users of other finance apps.

### 5. **Scalable**
Easy to add new features to appropriate sections.

---

## Summary

The new navigation structure:
- ✅ Groups related features logically
- ✅ Follows industry best practices
- ✅ Supports both expanded and collapsed modes
- ✅ Maintains purple theme consistency
- ✅ Provides clear visual hierarchy
- ✅ Scrollable for extensibility
- ✅ Responsive across devices

Perfect for a professional personal finance tracker! 🎯
