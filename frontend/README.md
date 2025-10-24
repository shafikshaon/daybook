# Personal Expense Tracker

A comprehensive personal finance management application built with Vue.js, Bootstrap, and Pinia. Track your income, expenses, budgets, investments, and more - all in one place!

## Features

### 💰 Account Management
- Multiple account types (Cash, Checking, Savings, Credit Card)
- Multi-currency support
- Balance tracking and reconciliation
- Account-level transaction history

### 📊 Transaction Management
- Income and expense tracking
- Category-based organization
- Tag support for advanced filtering
- Recurring transactions (subscriptions, bills)
- Split transactions across categories
- Bulk import capability
- Advanced search and filters

### 📈 Budget System
- Category-based budgets
- Multiple budget periods (weekly, monthly, quarterly, yearly, custom)
- Budget alerts and notifications
- Spending trends and forecasting
- Budget rollover support

### 💳 Credit Card Management
- Multiple credit card tracking
- Statement tracking with due dates
- Interest calculation
- Rewards and cashback tracking
- Payment reminders
- Credit utilization monitoring

### 📊 Investment Portfolio
- Support for stocks, bonds, mutual funds, ETFs, crypto, and more
- Portfolio performance tracking
- Cost basis and capital gains calculation
- Dividend and interest tracking
- Asset allocation visualization
- Top/bottom performers analysis

### 💎 Fixed Deposits & Savings Goals
- FD tracking with maturity dates and interest
- Savings goals with progress tracking
- Goal projections based on current savings rate
- Automated savings rules
- Multiple compounding frequency options

### 📑 Analytics & Reporting
- Net worth tracking over time
- Cash flow analysis (income vs expenses)
- Category-wise spending breakdown
- Monthly and yearly comparisons
- Customizable date ranges
- Export capabilities

### 🔔 Bill Payment Reminders
- Recurring bill tracking
- Due date reminders
- Payment history
- Overdue bill alerts

### 🎨 Additional Features
- **Dark Mode** - Easy on the eyes
- **Multi-currency** - Support for 10+ major currencies
- **Responsive Design** - Works on desktop, tablet, and mobile
- **Data Export** - Export all your data as JSON
- **Sample Data** - Pre-populated demo data to explore features
- **Purple Theme** - Beautiful, consistent purple-themed UI

## Technology Stack

- **Vue 3** - Progressive JavaScript framework
- **Vite** - Next generation frontend tooling
- **Pinia** - State management
- **Bootstrap 5** - CSS framework
- **Axios** - HTTP client (localStorage adapter for now)
- **SCSS** - Custom styling

## Project Structure

```
src/
├── assets/
│   └── styles/
│       └── custom.scss         # Custom Bootstrap theme
├── components/                 # Reusable components (expandable)
├── router/
│   └── index.js               # Vue Router configuration
├── services/
│   └── api.js                 # API service with localStorage adapter
├── stores/
│   ├── accounts.js            # Account management store
│   ├── bills.js               # Bills and reminders store
│   ├── budgets.js             # Budget management store
│   ├── creditCards.js         # Credit card store
│   ├── fixedDeposits.js       # Fixed deposits store
│   ├── investments.js         # Investment portfolio store
│   ├── savingsGoals.js        # Savings goals store
│   ├── settings.js            # App settings store
│   └── transactions.js        # Transaction management store
├── utils/
│   └── seedData.js            # Sample data generator
├── views/
│   ├── AccountsView.vue       # Accounts page
│   ├── BillsView.vue          # Bills & reminders page
│   ├── BudgetsView.vue        # Budgets page
│   ├── CreditCardsView.vue    # Credit cards page
│   ├── DashboardView.vue      # Main dashboard
│   ├── FixedDepositsView.vue  # Fixed deposits page
│   ├── InvestmentsView.vue    # Investments page
│   ├── ReportsView.vue        # Analytics & reports page
│   ├── SavingsGoalsView.vue   # Savings goals page
│   ├── SettingsView.vue       # Settings page
│   └── TransactionsView.vue   # Transactions page
├── App.vue                    # Root component with navigation
└── main.js                    # App entry point
```

## Getting Started

### Prerequisites

- Node.js 16+ and npm

### Installation

1. Clone or navigate to the project directory:

```bash
cd /Users/shafikshaon/workplace/development/projects/daybook/frontend
```

2. Install dependencies:

```bash
npm install
```

3. Start the development server:

```bash
npm run dev
```

4. Open your browser and navigate to `http://localhost:3000`

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

### Preview Production Build

```bash
npm run preview
```

## Usage

### First Time Setup

1. **Load Sample Data**: Click the "Load Sample Data" button on the dashboard to populate the app with example data
2. **Or Create Your Own**: Start fresh by creating your accounts, transactions, budgets, etc.

### Key Workflows

#### Adding a Transaction
1. Navigate to Transactions
2. Click "+ Add Transaction"
3. Fill in the details (type, amount, category, account, date)
4. Click "Create"

#### Setting Up a Budget
1. Navigate to Budgets
2. Click "+ Add Budget"
3. Select category, set amount, and choose period
4. Click "Create"
5. Monitor your spending against the budget

#### Tracking Investments
1. Navigate to Investments
2. Click "+ Add Investment"
3. Enter symbol, name, asset type, quantity, and prices
4. Track performance over time

#### Creating Savings Goals
1. Navigate to Savings Goals
2. Click "+ Add Goal"
3. Set target amount, monthly contribution, and target date
4. Add contributions regularly to track progress

### Settings

Access Settings to:
- Change default currency
- Toggle dark mode
- Configure notifications
- Export all data
- Clear all data (careful!)

## Data Storage

All data is stored locally in your browser's localStorage. This means:
- ✅ Your data never leaves your device
- ✅ No server required
- ✅ Complete privacy
- ⚠️ Data is browser-specific (won't sync across devices)
- ⚠️ Clearing browser data will delete your records

**Important**: Use the Export feature regularly to backup your data!

## Future Backend Integration

The app is designed to easily switch from localStorage to a real backend API:

1. The `api.js` service uses Axios with interceptors
2. All data operations go through the API service
3. To switch to a real backend:
   - Update the `baseURL` in `api.js`
   - Remove the localStorage interceptor
   - Implement real API endpoints
   - No changes needed in components or stores!

## Customization

### Changing the Theme Color

Edit `src/assets/styles/custom.scss`:

```scss
$primary: #your-color-here;
$sidebar-bg: #your-color-here;
```

### Adding New Categories

Edit `src/stores/transactions.js` and add to the `categories` array.

### Adding New Currency

Edit `src/stores/settings.js` and add to the `currencies` array.

## License

MIT License - feel free to use this project for personal or commercial purposes!

## Support

For issues or questions, please open an issue on the repository.

## Acknowledgments

Built with ❤️ using Vue.js, Bootstrap, and Pinia.
