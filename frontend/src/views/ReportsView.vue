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
          <div class="stat-value">{{ formatCurrency(cashFlow) }}</div>
          <div class="stat-label">Cash Flow (Monthly)</div>
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
      <!-- Income vs Expenses -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Income vs Expenses</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <div class="d-flex justify-content-between mb-2">
                <span class="text-success">Income</span>
                <span class="fw-bold text-success">{{ formatCurrency(periodIncome) }}</span>
              </div>
              <div class="progress mb-3" style="height: 20px;">
                <div class="progress-bar bg-success" :style="{ width: '50%' }"></div>
              </div>

              <div class="d-flex justify-content-between mb-2">
                <span class="text-danger">Expenses</span>
                <span class="fw-bold text-danger">{{ formatCurrency(periodExpense) }}</span>
              </div>
              <div class="progress" style="height: 20px;">
                <div class="progress-bar bg-danger" :style="{ width: '50%' }"></div>
              </div>
            </div>
            <div class="text-center p-3 bg-light rounded">
              <div class="fw-bold">Net Cash Flow</div>
              <div :class="periodIncome - periodExpense >= 0 ? 'text-success' : 'text-danger'" style="font-size: 1.5rem;">
                {{ formatCurrency(periodIncome - periodExpense) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Category Breakdown -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Expense by Category</h5>
          </div>
          <div class="card-body">
            <div v-for="breakdown in categoryBreakdown" :key="breakdown.categoryId" class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span>{{ getCategoryName(breakdown.categoryId) }}</span>
                <span class="fw-bold">{{ formatCurrency(breakdown.amount) }}</span>
              </div>
              <div class="progress" style="height: 8px;">
                <div
                  class="progress-bar"
                  :style="{
                    width: (breakdown.amount / totalCategoryExpenses * 100) + '%',
                    backgroundColor: getCategoryColor(breakdown.categoryId)
                  }"
                ></div>
              </div>
              <small class="text-muted">{{ Math.round(breakdown.amount / totalCategoryExpenses * 100) }}% of total</small>
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
                    <th class="text-end">Expenses</th>
                    <th class="text-end">Net</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="month in monthlyTrend" :key="month.month">
                    <td>{{ month.month }}</td>
                    <td class="text-end text-success">{{ formatCurrency(month.income) }}</td>
                    <td class="text-end text-danger">{{ formatCurrency(month.expenses) }}</td>
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
import { useInvestmentsStore } from '@/stores/investments'
import { useCreditCardsStore } from '@/stores/creditCards'
import { useSettingsStore } from '@/stores/settings'

const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const investmentsStore = useInvestmentsStore()
const creditCardsStore = useCreditCardsStore()
const settingsStore = useSettingsStore()

const dateRange = ref({
  start: new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString().split('T')[0],
  end: new Date().toISOString().split('T')[0]
})

const periodIncome = ref(0)
const periodExpense = ref(0)

const netWorth = computed(() => {
  return accountsStore.totalBalance + investmentsStore.totalCurrentValue - creditCardsStore.totalOutstanding
})

const totalAssets = computed(() => {
  return accountsStore.totalBalance + investmentsStore.totalCurrentValue
})

const totalLiabilities = computed(() => {
  return creditCardsStore.totalOutstanding
})

const cashFlow = computed(() => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  const income = transactionsStore.totalIncome(startOfMonth, endOfMonth)
  const expense = transactionsStore.totalExpense(startOfMonth, endOfMonth)
  return income - expense
})

const categoryBreakdown = computed(() => {
  return transactionsStore.categoryBreakdown('expense')
})

const totalCategoryExpenses = computed(() => {
  return categoryBreakdown.value.reduce((sum, cat) => sum + cat.amount, 0)
})

const monthlyTrend = computed(() => {
  const trends = []
  const now = new Date()

  for (let i = 5; i >= 0; i--) {
    const monthDate = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const startOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth(), 1)
    const endOfMonth = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0)

    const income = transactionsStore.totalIncome(startOfMonth, endOfMonth)
    const expenses = transactionsStore.totalExpense(startOfMonth, endOfMonth)

    trends.push({
      month: monthDate.toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
      income,
      expenses,
      net: income - expenses
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
  periodIncome.value = transactionsStore.totalIncome(dateRange.value.start, dateRange.value.end)
  periodExpense.value = transactionsStore.totalExpense(dateRange.value.start, dateRange.value.end)
}

onMounted(async () => {
  await Promise.all([
    accountsStore.fetchAccounts(),
    transactionsStore.fetchTransactions(),
    investmentsStore.fetchInvestments(),
    creditCardsStore.fetchCreditCards()
  ])
  applyDateRange()
})
</script>
