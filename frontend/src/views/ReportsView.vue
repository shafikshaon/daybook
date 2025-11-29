<template>
  <div class="reports-view fade-in">
    <h1 class="text-purple mb-4">Reports & Analytics</h1>

    <!-- Net Worth -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">💎</div>
          <div class="stat-value">{{ formatCurrency(netWorth) }}</div>
          <div class="stat-label">Net Worth</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">📈</div>
          <div class="stat-value">{{ formatCurrency(totalAssets) }}</div>
          <div class="stat-label">Total Assets</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(totalLiabilities) }}</div>
          <div class="stat-label">Total Liabilities</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ savingsRate }}%</div>
          <div class="stat-label">Savings Rate</div>
        </div>
      </div>
    </div>

    <!-- Date Range Filter -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-12 col-md-3">
            <label class="form-label">Quick Filter</label>
            <select class="form-select" v-model="quickFilter" @change="applyQuickFilter">
              <option value="custom">Custom Range</option>
              <option value="this_month">This Month</option>
              <option value="last_month">Last Month</option>
              <option value="this_quarter">This Quarter</option>
              <option value="this_year">This Year</option>
              <option value="last_year">Last Year</option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Start Date</label>
            <input type="date" class="form-control" v-model="dateRange.start" @change="quickFilter = 'custom'" />
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">End Date</label>
            <input type="date" class="form-control" v-model="dateRange.end" @change="quickFilter = 'custom'" />
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">&nbsp;</label>
            <button class="btn btn-primary w-100" @click="applyDateRange">Apply</button>
          </div>
        </div>
      </div>
    </div>

    <div class="row g-3">
      <!-- Cash Flow Summary -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Cash Flow Summary</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <div class="d-flex justify-content-between mb-2">
                <span class="text-success">💰 Total Income</span>
                <span class="fw-bold text-success">{{ formatCurrency(periodIncome) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-success" :style="{ width: '100%' }">100%</div>
              </div>

              <div class="d-flex justify-content-between mb-2">
                <span class="text-danger">🛍️ Total Expenses</span>
                <span class="fw-bold text-danger">{{ formatCurrency(periodExpense) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-danger" :style="{ width: calcPercentage(periodExpense, periodIncome) + '%' }">
                  {{ calcPercentage(periodExpense, periodIncome) }}%
                </div>
              </div>
            </div>
            <div class="text-center p-3 rounded" :class="netCashFlow >= 0 ? 'bg-success bg-opacity-10' : 'bg-danger bg-opacity-10'">
              <div class="fw-bold">Net Cash Flow</div>
              <div :class="netCashFlow >= 0 ? 'text-success' : 'text-danger'" style="font-size: 1.5rem;">
                {{ netCashFlow >= 0 ? '+' : '' }}{{ formatCurrency(netCashFlow) }}
              </div>
              <small class="text-muted">{{ dateRangeText }}</small>
            </div>
          </div>
        </div>
      </div>

      <!-- Savings Rate Chart -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Financial Health</h5>
          </div>
          <div class="card-body">
            <div class="text-center mb-3">
              <div style="font-size: 3rem; font-weight: bold;" :class="savingsRate >= 20 ? 'text-success' : savingsRate >= 10 ? 'text-warning' : 'text-danger'">
                {{ savingsRate }}%
              </div>
              <div class="text-muted">Savings Rate</div>
              <small class="text-muted">
                {{ savingsRate >= 20 ? '🎉 Excellent! You\'re saving over 20%' : savingsRate >= 10 ? '👍 Good! Try to save more' : '⚠️ Consider reducing expenses' }}
              </small>
            </div>
            <div class="mb-3">
              <div class="d-flex justify-content-between mb-2">
                <span>Saved</span>
                <span class="fw-bold text-success">{{ formatCurrency(netCashFlow >= 0 ? netCashFlow : 0) }}</span>
              </div>
              <div class="d-flex justify-content-between mb-2">
                <span>Spent</span>
                <span class="fw-bold text-danger">{{ formatCurrency(periodExpense) }}</span>
              </div>
              <div class="d-flex justify-content-between mb-2">
                <span>Earned</span>
                <span class="fw-bold text-primary">{{ formatCurrency(periodIncome) }}</span>
              </div>
            </div>
            <div class="alert alert-info mb-0">
              <small>
                <strong>Tip:</strong> Financial experts recommend saving at least 20% of your income for long-term financial security.
              </small>
            </div>
          </div>
        </div>
      </div>

      <!-- Income Breakdown -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">💰 Income by Category</h5>
            <span class="badge bg-success">{{ formatCurrency(periodIncome) }}</span>
          </div>
          <div class="card-body">
            <div v-if="incomeBreakdown.length > 0">
              <div v-for="breakdown in incomeBreakdown" :key="breakdown.categoryId" class="mb-3">
                <div class="d-flex justify-content-between mb-1">
                  <span>
                    <span class="me-1">{{ getCategoryIcon(breakdown.categoryId) }}</span>
                    {{ getCategoryName(breakdown.categoryId) }}
                  </span>
                  <span class="fw-bold text-success">{{ formatCurrency(breakdown.amount) }}</span>
                </div>
                <div class="progress" style="height: 10px;">
                  <div
                    class="progress-bar bg-success"
                    :style="{
                      width: (breakdown.amount / periodIncome * 100) + '%'
                    }"
                  ></div>
                </div>
                <small class="text-muted">
                  {{ Math.round(breakdown.amount / periodIncome * 100) }}% of total income
                  ({{ breakdown.count }} transaction{{ breakdown.count > 1 ? 's' : '' }})
                </small>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No income in this period</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Expense Breakdown -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">🛍️ Expenses by Category</h5>
            <span class="badge bg-danger">{{ formatCurrency(periodExpense) }}</span>
          </div>
          <div class="card-body">
            <div v-if="expenseBreakdown.length > 0">
              <div v-for="breakdown in expenseBreakdown" :key="breakdown.categoryId" class="mb-3">
                <div class="d-flex justify-content-between mb-1">
                  <span>
                    <span class="me-1">{{ getCategoryIcon(breakdown.categoryId) }}</span>
                    {{ getCategoryName(breakdown.categoryId) }}
                  </span>
                  <span class="fw-bold text-danger">{{ formatCurrency(breakdown.amount) }}</span>
                </div>
                <div class="progress" style="height: 10px;">
                  <div
                    class="progress-bar"
                    :style="{
                      width: (breakdown.amount / periodExpense * 100) + '%',
                      backgroundColor: getCategoryColor(breakdown.categoryId)
                    }"
                  ></div>
                </div>
                <small class="text-muted">
                  {{ Math.round(breakdown.amount / periodExpense * 100) }}% of total expenses
                  ({{ breakdown.count }} transaction{{ breakdown.count > 1 ? 's' : '' }})
                </small>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No expenses in this period</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Spending Categories -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">🔥 Top 5 Spending Categories</h5>
          </div>
          <div class="card-body">
            <div v-if="topExpenseCategories.length > 0">
              <div class="list-group list-group-flush">
                <div
                  v-for="(category, index) in topExpenseCategories"
                  :key="category.categoryId"
                  class="list-group-item d-flex justify-content-between align-items-center px-0"
                >
                  <div class="d-flex align-items-center">
                    <span class="badge bg-secondary me-2">#{{ index + 1 }}</span>
                    <span class="me-2">{{ getCategoryIcon(category.categoryId) }}</span>
                    <span>{{ getCategoryName(category.categoryId) }}</span>
                  </div>
                  <span class="fw-bold text-danger">{{ formatCurrency(category.amount) }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No expense data available</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Income Sources -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">💎 Top 5 Income Sources</h5>
          </div>
          <div class="card-body">
            <div v-if="topIncomeCategories.length > 0">
              <div class="list-group list-group-flush">
                <div
                  v-for="(category, index) in topIncomeCategories"
                  :key="category.categoryId"
                  class="list-group-item d-flex justify-content-between align-items-center px-0"
                >
                  <div class="d-flex align-items-center">
                    <span class="badge bg-secondary me-2">#{{ index + 1 }}</span>
                    <span class="me-2">{{ getCategoryIcon(category.categoryId) }}</span>
                    <span>{{ getCategoryName(category.categoryId) }}</span>
                  </div>
                  <span class="fw-bold text-success">{{ formatCurrency(category.amount) }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No income data available</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Monthly Trend -->
      <div class="col-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">📊 Monthly Trend (Last 6 Months)</h5>
          </div>
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>Month</th>
                    <th class="text-end">Income</th>
                    <th class="text-end">Expenses</th>
                    <th class="text-end">Net</th>
                    <th class="text-end">Savings Rate</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="month in monthlyTrend" :key="month.month">
                    <td class="fw-bold">{{ month.month }}</td>
                    <td class="text-end text-success">{{ formatCurrency(month.income) }}</td>
                    <td class="text-end text-danger">{{ formatCurrency(month.expenses) }}</td>
                    <td class="text-end" :class="month.net >= 0 ? 'text-success fw-bold' : 'text-danger fw-bold'">
                      {{ month.net >= 0 ? '+' : '' }}{{ formatCurrency(month.net) }}
                    </td>
                    <td class="text-end">
                      <span class="badge" :class="month.savingsRate >= 20 ? 'bg-success' : month.savingsRate >= 10 ? 'bg-warning' : 'bg-danger'">
                        {{ month.savingsRate }}%
                      </span>
                    </td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr class="table-light">
                    <td class="fw-bold">Average</td>
                    <td class="text-end text-success fw-bold">{{ formatCurrency(averageMonthlyIncome) }}</td>
                    <td class="text-end text-danger fw-bold">{{ formatCurrency(averageMonthlyExpense) }}</td>
                    <td class="text-end fw-bold" :class="averageMonthlyNet >= 0 ? 'text-success' : 'text-danger'">
                      {{ averageMonthlyNet >= 0 ? '+' : '' }}{{ formatCurrency(averageMonthlyNet) }}
                    </td>
                    <td class="text-end">
                      <span class="badge" :class="averageSavingsRate >= 20 ? 'bg-success' : averageSavingsRate >= 10 ? 'bg-warning' : 'bg-danger'">
                        {{ averageSavingsRate }}%
                      </span>
                    </td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import { useCreditCardsStore } from '@/stores/creditCards'
import { useSettingsStore } from '@/stores/settings'

const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const creditCardsStore = useCreditCardsStore()
const settingsStore = useSettingsStore()

// Initialize date range to current month (from 1st to today)
const today = new Date()
const currentMonthStart = new Date(today.getFullYear(), today.getMonth(), 1)

const dateRange = ref({
  start: currentMonthStart.toISOString().split('T')[0],
  end: today.toISOString().split('T')[0]
})

const quickFilter = ref('this_month')

// Helper to get transactions in date range
const getTransactionsInRange = (startDate, endDate) => {
  return transactionsStore.transactions.filter(t => {
    const transactionDate = new Date(t.date).toISOString().split('T')[0]
    return transactionDate >= startDate && transactionDate <= endDate
  })
}

// Helper to get category breakdown
const getCategoryBreakdown = (type, startDate, endDate) => {
  const transactions = getTransactionsInRange(startDate, endDate).filter(t => t.type === type)
  const breakdown = {}

  transactions.forEach(t => {
    if (!breakdown[t.categoryId]) {
      breakdown[t.categoryId] = {
        categoryId: t.categoryId,
        amount: 0,
        count: 0
      }
    }
    breakdown[t.categoryId].amount += t.amount
    breakdown[t.categoryId].count++
  })

  return Object.values(breakdown).sort((a, b) => b.amount - a.amount)
}

// Period calculations
const periodIncome = computed(() => {
  return transactionsStore.totalIncome(dateRange.value.start, dateRange.value.end)
})

const periodExpense = computed(() => {
  return transactionsStore.totalExpense(dateRange.value.start, dateRange.value.end)
})

const netCashFlow = computed(() => {
  return periodIncome.value - periodExpense.value
})

const savingsRate = computed(() => {
  if (periodIncome.value === 0) return 0
  return Math.round((netCashFlow.value / periodIncome.value) * 100)
})

// Net worth
const netWorth = computed(() => {
  return accountsStore.totalBalance - creditCardsStore.totalOutstanding
})

const totalAssets = computed(() => {
  return accountsStore.totalBalance
})

const totalLiabilities = computed(() => {
  return creditCardsStore.totalOutstanding
})

// Category breakdowns
const incomeBreakdown = computed(() => {
  return getCategoryBreakdown('income', dateRange.value.start, dateRange.value.end)
})

const expenseBreakdown = computed(() => {
  return getCategoryBreakdown('expense', dateRange.value.start, dateRange.value.end)
})

// Top categories
const topExpenseCategories = computed(() => {
  return expenseBreakdown.value.slice(0, 5)
})

const topIncomeCategories = computed(() => {
  return incomeBreakdown.value.slice(0, 5)
})

// Monthly trend
const monthlyTrend = computed(() => {
  const trends = []
  const now = new Date()

  for (let i = 5; i >= 0; i--) {
    const monthDate = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const startOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1).toISOString().split('T')[0]
    const endOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0).toISOString().split('T')[0]

    const income = transactionsStore.totalIncome(startOfMonth, endOfMonth)
    const expenses = transactionsStore.totalExpense(startOfMonth, endOfMonth)
    const net = income - expenses
    const rate = income > 0 ? Math.round((net / income) * 100) : 0

    trends.push({
      month: monthDate.toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
      income,
      expenses,
      net,
      savingsRate: rate
    })
  }

  return trends
})

// Monthly averages
const averageMonthlyIncome = computed(() => {
  if (monthlyTrend.value.length === 0) return 0
  const total = monthlyTrend.value.reduce((sum, m) => sum + m.income, 0)
  return Math.round(total / monthlyTrend.value.length)
})

const averageMonthlyExpense = computed(() => {
  if (monthlyTrend.value.length === 0) return 0
  const total = monthlyTrend.value.reduce((sum, m) => sum + m.expenses, 0)
  return Math.round(total / monthlyTrend.value.length)
})

const averageMonthlyNet = computed(() => {
  return averageMonthlyIncome.value - averageMonthlyExpense.value
})

const averageSavingsRate = computed(() => {
  if (averageMonthlyIncome.value === 0) return 0
  return Math.round((averageMonthlyNet.value / averageMonthlyIncome.value) * 100)
})

// Date range text
const dateRangeText = computed(() => {
  const start = new Date(dateRange.value.start)
  const end = new Date(dateRange.value.end)
  return `${start.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${end.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}`
})

// Helpers
const calcPercentage = (amount, total) => {
  if (total === 0) return 0
  return Math.min(100, Math.round((amount / total) * 100))
}

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)

const getCategoryName = (id) => {
  const cat = transactionsStore.getCategoryById(id)
  return cat ? cat.name : id
}

const getCategoryIcon = (id) => {
  const cat = transactionsStore.getCategoryById(id)
  return cat ? cat.icon : '📦'
}

const getCategoryColor = (id) => {
  const cat = transactionsStore.getCategoryById(id)
  return cat ? cat.color : '#6b7280'
}

const applyQuickFilter = () => {
  const today = new Date()

  switch (quickFilter.value) {
    case 'this_month':
      dateRange.value.start = new Date(today.getFullYear(), today.getMonth(), 1).toISOString().split('T')[0]
      dateRange.value.end = today.toISOString().split('T')[0]
      break
    case 'last_month':
      const lastMonth = new Date(today.getFullYear(), today.getMonth() - 1, 1)
      dateRange.value.start = lastMonth.toISOString().split('T')[0]
      dateRange.value.end = new Date(today.getFullYear(), today.getMonth(), 0).toISOString().split('T')[0]
      break
    case 'this_quarter':
      const quarter = Math.floor(today.getMonth() / 3)
      dateRange.value.start = new Date(today.getFullYear(), quarter * 3, 1).toISOString().split('T')[0]
      dateRange.value.end = today.toISOString().split('T')[0]
      break
    case 'this_year':
      dateRange.value.start = new Date(today.getFullYear(), 0, 1).toISOString().split('T')[0]
      dateRange.value.end = today.toISOString().split('T')[0]
      break
    case 'last_year':
      dateRange.value.start = new Date(today.getFullYear() - 1, 0, 1).toISOString().split('T')[0]
      dateRange.value.end = new Date(today.getFullYear() - 1, 11, 31).toISOString().split('T')[0]
      break
  }
}

const applyDateRange = () => {
  // The computed properties will automatically update
  console.log('=== Report Date Range Applied ===')
  console.log('Date Range:', dateRange.value.start, 'to', dateRange.value.end)
  console.log('Period Income:', periodIncome.value)
  console.log('Period Expense:', periodExpense.value)
  console.log('Net Cash Flow:', netCashFlow.value)
  console.log('================================')
}

onMounted(async () => {
  console.log('Loading report data...')
  await Promise.all([
    accountsStore.fetchAccounts(),
    transactionsStore.fetchCategories(),
    transactionsStore.fetchTransactions(1, 1000), // Get more transactions for reports
    creditCardsStore.fetchCreditCards()
  ])
  console.log('Transactions loaded:', transactionsStore.transactions.length)
  console.log('Categories loaded:', transactionsStore.categories.length)
})
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.stat-icon.green {
  color: #10b981;
}

.stat-icon.red {
  color: #ef4444;
}

.stat-icon.blue {
  color: #3b82f6;
}

.stat-icon.purple {
  color: #8b5cf6;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.stat-label {
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
}

.card {
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  background-color: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
  padding: 1rem 1.5rem;
}

.text-purple {
  color: #8b5cf6;
}

/* Dark mode */
.dark-mode .stat-card {
  background-color: #1f2937;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.dark-mode .stat-value {
  color: #f3f4f6;
}

.dark-mode .stat-label {
  color: #9ca3af;
}

.dark-mode .card {
  background-color: #1f2937;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.dark-mode .card-header {
  background-color: #111827;
  border-bottom-color: #374151;
}

.dark-mode .list-group-item {
  background-color: #1f2937;
  border-color: #374151;
  color: #f3f4f6;
}

.dark-mode .table {
  color: #f3f4f6;
}

.dark-mode .table-light {
  background-color: #374151;
  color: #f3f4f6;
}
</style>
