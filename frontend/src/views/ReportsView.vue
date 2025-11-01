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
          <div class="col-12 col-md-4">
            <label class="form-label">Start Date</label>
            <input type="date" class="form-control" v-model="dateRange.start" />
          </div>
          <div class="col-12 col-md-4">
            <label class="form-label">End Date</label>
            <input type="date" class="form-control" v-model="dateRange.end" />
          </div>
          <div class="col-12 col-md-4">
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
                <span class="text-success">💰 Income</span>
                <span class="fw-bold text-success">{{ formatCurrency(periodIncome) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-success" :style="{ width: '100%' }">100%</div>
              </div>

              <div class="d-flex justify-content-between mb-2">
                <span class="text-danger">🛍️ Regular Expenses</span>
                <span class="fw-bold text-danger">{{ formatCurrency(periodRegularExpense) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-danger" :style="{ width: calcPercentage(periodRegularExpense, periodIncome) + '%' }">
                  {{ calcPercentage(periodRegularExpense, periodIncome) }}%
                </div>
              </div>

              <div class="d-flex justify-content-between mb-2">
                <span class="text-primary">💎 Savings & Investments</span>
                <span class="fw-bold text-primary">{{ formatCurrency(periodSavings) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-primary" :style="{ width: calcPercentage(periodSavings, periodIncome) + '%' }">
                  {{ calcPercentage(periodSavings, periodIncome) }}%
                </div>
              </div>
            </div>
            <div class="text-center p-3 bg-light rounded">
              <div class="fw-bold">Net Cash Flow</div>
              <div :class="periodIncome - periodRegularExpense - periodSavings >= 0 ? 'text-success' : 'text-danger'" style="font-size: 1.5rem;">
                {{ formatCurrency(periodIncome - periodRegularExpense - periodSavings) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Savings Rate Chart -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Money Allocation</h5>
          </div>
          <div class="card-body">
            <div class="text-center mb-3">
              <div style="font-size: 3rem; font-weight: bold;" :class="savingsRate >= 20 ? 'text-success' : savingsRate >= 10 ? 'text-warning' : 'text-danger'">
                {{ savingsRate }}%
              </div>
              <div class="text-muted">Savings Rate</div>
              <small class="text-muted">
                {{ savingsRate >= 20 ? '🎉 Excellent!' : savingsRate >= 10 ? '👍 Good' : '⚠️ Consider saving more' }}
              </small>
            </div>
            <div class="mb-3">
              <div class="d-flex justify-content-between mb-2">
                <span>Savings & Investments</span>
                <span class="fw-bold text-primary">{{ calcPercentage(periodSavings, periodIncome) }}%</span>
              </div>
              <div class="d-flex justify-content-between mb-2">
                <span>Regular Expenses</span>
                <span class="fw-bold text-danger">{{ calcPercentage(periodRegularExpense, periodIncome) }}%</span>
              </div>
              <div class="d-flex justify-content-between mb-2">
                <span>Remaining</span>
                <span class="fw-bold text-muted">
                  {{ Math.max(0, 100 - calcPercentage(periodSavings, periodIncome) - calcPercentage(periodRegularExpense, periodIncome)) }}%
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Regular Expenses Breakdown -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">🛍️ Regular Expenses by Category</h5>
          </div>
          <div class="card-body">
            <div v-if="regularExpenseBreakdown.length > 0">
              <div v-for="breakdown in regularExpenseBreakdown" :key="breakdown.categoryId" class="mb-3">
                <div class="d-flex justify-content-between mb-1">
                  <span>{{ getCategoryName(breakdown.categoryId) }}</span>
                  <span class="fw-bold">{{ formatCurrency(breakdown.amount) }}</span>
                </div>
                <div class="progress" style="height: 8px;">
                  <div
                    class="progress-bar"
                    :style="{
                      width: (breakdown.amount / totalRegularExpenses * 100) + '%',
                      backgroundColor: getCategoryColor(breakdown.categoryId)
                    }"
                  ></div>
                </div>
                <small class="text-muted">{{ Math.round(breakdown.amount / totalRegularExpenses * 100) }}% of regular expenses</small>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No regular expenses in this period</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Savings & Investments Breakdown -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">💎 Savings & Investments by Category</h5>
          </div>
          <div class="card-body">
            <div v-if="savingsBreakdown.length > 0">
              <div v-for="breakdown in savingsBreakdown" :key="breakdown.categoryId" class="mb-3">
                <div class="d-flex justify-content-between mb-1">
                  <span>{{ getCategoryName(breakdown.categoryId) }}</span>
                  <span class="fw-bold">{{ formatCurrency(breakdown.amount) }}</span>
                </div>
                <div class="progress" style="height: 8px;">
                  <div
                    class="progress-bar"
                    :style="{
                      width: (breakdown.amount / totalSavings * 100) + '%',
                      backgroundColor: getCategoryColor(breakdown.categoryId)
                    }"
                  ></div>
                </div>
                <small class="text-muted">{{ Math.round(breakdown.amount / totalSavings * 100) }}% of savings</small>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              <p>No savings or investments in this period</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Monthly Trend -->
      <div class="col-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Monthly Trend (Last 6 Months)</h5>
          </div>
          <div class="card-body">
            <div class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>Month</th>
                    <th class="text-end">Income</th>
                    <th class="text-end">Regular Expenses</th>
                    <th class="text-end">Savings</th>
                    <th class="text-end">Net</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="month in monthlyTrend" :key="month.month">
                    <td>{{ month.month }}</td>
                    <td class="text-end text-success">{{ formatCurrency(month.income) }}</td>
                    <td class="text-end text-danger">{{ formatCurrency(month.regularExpenses) }}</td>
                    <td class="text-end text-primary">{{ formatCurrency(month.savings) }}</td>
                    <td class="text-end" :class="month.net >= 0 ? 'text-success fw-bold' : 'text-danger fw-bold'">
                      {{ formatCurrency(month.net) }}
                    </td>
                  </tr>
                </tbody>
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

const dateRange = ref({
  start: new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString().split('T')[0],
  end: new Date().toISOString().split('T')[0]
})

const periodIncome = ref(0)
const periodRegularExpense = ref(0)
const periodSavings = ref(0)

const netWorth = computed(() => {
  return accountsStore.totalBalance - creditCardsStore.totalOutstanding
})

const totalAssets = computed(() => {
  return accountsStore.totalBalance
})

const totalLiabilities = computed(() => {
  return creditCardsStore.totalOutstanding
})

const savingsRate = computed(() => {
  if (periodIncome.value === 0) return 0
  return Math.round((periodSavings.value / periodIncome.value) * 100)
})

const regularExpenseBreakdown = computed(() => {
  return transactionsStore.groupBreakdown('expense', dateRange.value.start, dateRange.value.end)
})

const savingsBreakdown = computed(() => {
  return transactionsStore.groupBreakdown('savings', dateRange.value.start, dateRange.value.end)
})

const totalRegularExpenses = computed(() => {
  return regularExpenseBreakdown.value.reduce((sum, cat) => sum + cat.amount, 0)
})

const totalSavings = computed(() => {
  return savingsBreakdown.value.reduce((sum, cat) => sum + cat.amount, 0)
})

const calcPercentage = (amount, total) => {
  if (total === 0) return 0
  return Math.round((amount / total) * 100)
}

const monthlyTrend = computed(() => {
  const trends = []
  const now = new Date()

  for (let i = 5; i >= 0; i--) {
    const monthDate = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const startOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1).toISOString().split('T')[0]
    const endOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0).toISOString().split('T')[0]

    const income = transactionsStore.totalByGroup('income', startOfMonth, endOfMonth)
    const regularExpenses = transactionsStore.totalByGroup('expense', startOfMonth, endOfMonth)
    const savings = transactionsStore.totalByGroup('savings', startOfMonth, endOfMonth)

    trends.push({
      month: monthDate.toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
      income,
      regularExpenses,
      savings,
      net: income - regularExpenses - savings
    })
  }

  return trends
})

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)

const getCategoryName = (id) => {
  const cat = transactionsStore.getCategoryById(id)
  return cat ? cat.name : id
}

const getCategoryColor = (id) => {
  const cat = transactionsStore.getCategoryById(id)
  return cat ? cat.color : '#6b7280'
}

const applyDateRange = () => {
  periodIncome.value = transactionsStore.totalByGroup('income', dateRange.value.start, dateRange.value.end)
  periodRegularExpense.value = transactionsStore.totalByGroup('expense', dateRange.value.start, dateRange.value.end)
  periodSavings.value = transactionsStore.totalByGroup('savings', dateRange.value.start, dateRange.value.end)
}

onMounted(async () => {
  await Promise.all([
    accountsStore.fetchAccounts(),
    transactionsStore.fetchTransactions(),
    creditCardsStore.fetchCreditCards()
  ])
  applyDateRange()
})
</script>
