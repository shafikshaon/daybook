<template>
  <div class="dashboard fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Dashboard</h1>
    </div>

    <!-- Summary Cards -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-sm-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(totalNetWorth) }}</div>
          <div class="stat-label">Total Net Worth</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            Liquid: {{ formatCurrency(totalBalance) }} | Goals: {{ formatCurrency(totalGoalsValue) }}
          </small>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon blue">🎯</div>
          <div class="stat-value">{{ formatCurrency(totalGoalsValue) }}</div>
          <div class="stat-label">Goals & Investments</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            {{ activeGoalsCount }} active goals
          </small>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon green">📈</div>
          <div class="stat-value">{{ formatCurrency(monthlyIncome) }}</div>
          <div class="stat-label">Monthly Income</div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(monthlyExpenses) }}</div>
          <div class="stat-label">Monthly Expenses</div>
        </div>
      </div>
    </div>

    <div class="row g-3">
      <!-- Recent Transactions -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Recent Transactions</h5>
          </div>
          <div class="card-body p-0">
            <div v-if="recentTransactions.length === 0" class="p-4 text-center text-muted">
              No transactions yet
            </div>
            <div v-else class="table-responsive">
              <table class="table table-hover mb-0">
                <tbody>
                  <tr v-for="transaction in recentTransactions.slice(0, 5)" :key="transaction.id">
                    <td>
                      <div class="fw-semibold">{{ getCategoryName(transaction.categoryId) }}</div>
                      <small class="text-muted">{{ formatDate(transaction.date) }}</small>
                    </td>
                    <td class="text-end">
                      <span :class="transaction.type === 'income' ? 'text-success' : 'text-danger'">
                        {{ transaction.type === 'income' ? '+' : '-' }}{{ formatCurrency(transaction.amount) }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Budget Overview -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Budget Overview</h5>
          </div>
          <div class="card-body">
            <div v-if="budgets.length === 0" class="text-center text-muted">
              No budgets set
            </div>
            <div v-else>
              <div v-for="budget in budgets.slice(0, 4)" :key="budget.id" class="mb-3">
                <div class="d-flex justify-content-between mb-1">
                  <span>{{ getCategoryName(budget.categoryId) }}</span>
                  <span class="text-muted">
                    {{ formatCurrency(getBudgetProgress(budget.id)?.spent || 0) }} / {{ formatCurrency(budget.amount) }}
                  </span>
                </div>
                <div class="progress" style="height: 8px;">
                  <div
                    class="progress-bar"
                    :class="getProgressBarClass(getBudgetProgress(budget.id)?.percentage)"
                    :style="{ width: Math.min(getBudgetProgress(budget.id)?.percentage || 0, 100) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Accounts Summary -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Accounts</h5>
          </div>
          <div class="card-body">
            <div v-if="accounts.length === 0" class="text-center text-muted">
              No accounts yet
            </div>
            <div v-else>
              <div v-for="account in accounts" :key="account.id" class="d-flex justify-content-between align-items-center mb-2">
                <div>
                  <span class="fw-semibold">{{ account.name }}</span>
                  <span class="badge bg-secondary ms-2 text-uppercase" style="font-size: 0.7rem;">
                    {{ account.type }}
                  </span>
                </div>
                <div class="fw-bold">{{ formatCurrency(account.balance) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Upcoming Bills -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Upcoming Bills</h5>
          </div>
          <div class="card-body">
            <div v-if="upcomingBills.length === 0" class="text-center text-muted">
              No upcoming bills
            </div>
            <div v-else>
              <div v-for="bill in upcomingBills.slice(0, 5)" :key="bill.id" class="d-flex justify-content-between align-items-center mb-2">
                <div>
                  <div class="fw-semibold">{{ bill.name }}</div>
                  <small class="text-muted">Due in {{ bill.daysUntilDue }} days</small>
                </div>
                <div class="fw-bold">{{ formatCurrency(bill.amount) }}</div>
              </div>
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
import { useBudgetsStore } from '@/stores/budgets'
import { useBillsStore } from '@/stores/bills'
import { useGoalsStore } from '@/stores/goals'
import { useSettingsStore } from '@/stores/settings'

const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const budgetsStore = useBudgetsStore()
const billsStore = useBillsStore()
const goalsStore = useGoalsStore()
const settingsStore = useSettingsStore()

const accounts = computed(() => accountsStore.allAccounts)
const recentTransactions = computed(() => transactionsStore.allTransactions)
const budgets = computed(() => budgetsStore.activeBudgets)
const upcomingBills = computed(() => billsStore.upcomingBills)

const totalBalance = computed(() => accountsStore.totalBalance)
const totalGoalsValue = computed(() => goalsStore.totalCurrentAmount)
const activeGoalsCount = computed(() => goalsStore.activeGoals.length)
const totalNetWorth = computed(() => totalBalance.value + totalGoalsValue.value)

const monthlyIncome = computed(() => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  return transactionsStore.totalIncome(startOfMonth, endOfMonth)
})

const monthlyExpenses = computed(() => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  return transactionsStore.totalExpense(startOfMonth, endOfMonth)
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const getCategoryName = (categoryId) => {
  const category = transactionsStore.getCategoryById(categoryId)
  return category ? category.name : categoryId
}

const getBudgetProgress = (budgetId) => {
  return budgetsStore.budgetProgress(budgetId)
}

const getProgressBarClass = (percentage) => {
  if (percentage >= 100) return 'bg-danger'
  if (percentage >= 80) return 'bg-warning'
  return 'bg-success'
}

const loadData = async () => {
  await Promise.all([
    accountsStore.fetchAccounts(),
    transactionsStore.fetchTransactions(),
    budgetsStore.fetchBudgets(),
    billsStore.fetchBills(),
    goalsStore.fetchGoals()
  ])
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}
</style>
