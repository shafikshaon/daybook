<template>
  <div class="reports-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple mb-0">📊 Reports & Analytics</h1>
      <div class="d-flex gap-2">
        <select class="form-select form-select-sm" v-model="selectedPeriod" @change="handlePeriodChange" style="width: auto;">
          <option value="this_month">This Month</option>
          <option value="last_month">Last Month</option>
          <option value="this_quarter">This Quarter</option>
          <option value="last_quarter">Last Quarter</option>
          <option value="this_year">This Year</option>
          <option value="last_year">Last Year</option>
          <option value="custom">Custom Range</option>
        </select>
        <button class="btn btn-sm btn-outline-primary" @click="refreshAllData">
          <span v-if="loading">⏳</span>
          <span v-else>🔄</span> Refresh
        </button>
      </div>
    </div>

    <!-- Custom Date Range Modal (if custom selected) -->
    <div v-if="selectedPeriod === 'custom'" class="card mb-4">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-4">
            <label class="form-label">Start Date</label>
            <input type="date" class="form-control" v-model="customRange.start" />
          </div>
          <div class="col-md-4">
            <label class="form-label">End Date</label>
            <input type="date" class="form-control" v-model="customRange.end" />
          </div>
          <div class="col-md-4 d-flex align-items-end">
            <button class="btn btn-primary w-100" @click="applyCustomRange">Apply Range</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="initialLoading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading reports...</span>
      </div>
      <p class="mt-3 text-muted">Loading comprehensive reports...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger">
      <h5>Error Loading Reports</h5>
      <p>{{ error }}</p>
      <button class="btn btn-sm btn-danger" @click="refreshAllData">Retry</button>
    </div>

    <!-- Dashboard Content -->
    <div v-else>
      <!-- Key Metrics Cards -->
      <div class="row g-3 mb-4">
        <div class="col-6 col-md-3">
          <div class="metric-card purple-gradient">
            <div class="metric-icon">💎</div>
            <div class="metric-value">{{ formatCurrency(dashboardData?.netWorth?.netWorth || 0) }}</div>
            <div class="metric-label">Net Worth</div>
          </div>
        </div>
        <div class="col-6 col-md-3">
          <div class="metric-card green-gradient">
            <div class="metric-icon">💰</div>
            <div class="metric-value">{{ formatCurrency(dashboardData?.currentMonth?.totalIncome || 0) }}</div>
            <div class="metric-label">Current Month Income</div>
          </div>
        </div>
        <div class="col-6 col-md-3">
          <div class="metric-card red-gradient">
            <div class="metric-icon">💸</div>
            <div class="metric-value">{{ formatCurrency(dashboardData?.currentMonth?.totalExpense || 0) }}</div>
            <div class="metric-label">Current Month Expense</div>
          </div>
        </div>
        <div class="col-6 col-md-3">
          <div class="metric-card blue-gradient">
            <div class="metric-icon">📈</div>
            <div class="metric-value">{{ formatCurrency(dashboardData?.currentMonth?.netSavings || 0) }}</div>
            <div class="metric-label">Current Month Savings</div>
            <div class="metric-change" :class="savingsChangeClass">
              {{ savingsChangeText }}
            </div>
          </div>
        </div>
      </div>

      <!-- Income vs Expense Trend -->
      <div class="card mb-4">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">📊 Income vs Expense Trend</h5>
          <div class="btn-group btn-group-sm">
            <button
              v-for="period in ['day', 'week', 'month', 'year']"
              :key="period"
              class="btn"
              :class="trendGroupBy === period ? 'btn-primary' : 'btn-outline-primary'"
              @click="changeTrendGrouping(period)"
            >
              {{ period.charAt(0).toUpperCase() + period.slice(1) }}
            </button>
          </div>
        </div>
        <div class="card-body">
          <div v-if="loadingTrend" class="text-center py-4">
            <div class="spinner-border spinner-border-sm text-primary"></div>
          </div>
          <Line v-else-if="incomeTrendData.labels.length > 0" :data="incomeTrendData" :options="trendChartOptions" />
          <div v-else class="text-center text-muted py-4">No data available for selected period</div>
        </div>
      </div>

      <!-- Category Analysis and Account Distribution Row -->
      <div class="row g-3 mb-4">
        <!-- Expense Categories Breakdown -->
        <div class="col-lg-6">
          <div class="card h-100">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h5 class="mb-0">🏷️ Expense Categories</h5>
              <select class="form-select form-select-sm" v-model="categoryLimit" @change="loadCategoryData" style="width: auto;">
                <option :value="5">Top 5</option>
                <option :value="10">Top 10</option>
                <option :value="15">Top 15</option>
              </select>
            </div>
            <div class="card-body">
              <div v-if="loadingCategories" class="text-center py-4">
                <div class="spinner-border spinner-border-sm text-primary"></div>
              </div>
              <Doughnut v-else-if="categoryData.labels.length > 0" :data="categoryData" :options="categoryChartOptions" />
              <div v-else class="text-center text-muted py-4">No expense data available</div>
            </div>
          </div>
        </div>

        <!-- Account Balances -->
        <div class="col-lg-6">
          <div class="card h-100">
            <div class="card-header">
              <h5 class="mb-0">🏦 Account Distribution</h5>
            </div>
            <div class="card-body">
              <div v-if="loadingAccounts" class="text-center py-4">
                <div class="spinner-border spinner-border-sm text-primary"></div>
              </div>
              <div v-else-if="accountData.labels.length > 0">
                <Pie :data="accountData" :options="accountChartOptions" />
                <div class="mt-3">
                  <div class="d-flex justify-content-between align-items-center">
                    <strong>Total Balance:</strong>
                    <strong class="text-success">{{ formatCurrency(accountsReport?.totalBalance || 0) }}</strong>
                  </div>
                </div>
              </div>
              <div v-else class="text-center text-muted py-4">No account data available</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Cash Flow Chart -->
      <div class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">💵 Cash Flow Analysis</h5>
        </div>
        <div class="card-body">
          <div v-if="loadingCashFlow" class="text-center py-4">
            <div class="spinner-border spinner-border-sm text-primary"></div>
          </div>
          <Bar v-else-if="cashFlowData.labels.length > 0" :data="cashFlowData" :options="cashFlowChartOptions" />
          <div v-else class="text-center text-muted py-4">No cash flow data available</div>
        </div>
      </div>

      <!-- Net Worth Trend -->
      <div class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">💎 Net Worth Trend</h5>
        </div>
        <div class="card-body">
          <div v-if="loadingNetWorth" class="text-center py-4">
            <div class="spinner-border spinner-border-sm text-primary"></div>
          </div>
          <Line v-else-if="netWorthData.labels.length > 0" :data="netWorthData" :options="netWorthChartOptions" />
          <div v-else class="text-center text-muted py-4">No net worth history available</div>
        </div>
      </div>

      <!-- Budget Performance (if available) -->
      <div v-if="budgetReport && budgetReport.budgets && budgetReport.budgets.length > 0" class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">🎯 Budget Performance</h5>
        </div>
        <div class="card-body">
          <div class="row g-3">
            <div v-for="budget in budgetReport.budgets" :key="budget.categoryId" class="col-md-6 col-lg-4">
              <div class="budget-card">
                <div class="d-flex justify-content-between align-items-center mb-2">
                  <strong>{{ budget.categoryName }}</strong>
                  <span class="badge" :class="getBudgetStatusClass(budget.status)">
                    {{ budget.status }}
                  </span>
                </div>
                <div class="progress mb-2" style="height: 24px;">
                  <div
                    class="progress-bar"
                    :class="getBudgetProgressClass(budget.percentageUsed)"
                    :style="{ width: Math.min(budget.percentageUsed, 100) + '%' }"
                  >
                    {{ Math.round(budget.percentageUsed) }}%
                  </div>
                </div>
                <div class="d-flex justify-content-between small text-muted">
                  <span>{{ formatCurrency(budget.spent) }} / {{ formatCurrency(budget.budgeted) }}</span>
                  <span>{{ formatCurrency(budget.remaining) }} left</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Expense Categories List -->
      <div v-if="categoryReport && categoryReport.breakdown && categoryReport.breakdown.length > 0" class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">📋 Top Expense Categories (Detailed)</h5>
        </div>
        <div class="card-body">
          <div class="table-responsive">
            <table class="table table-hover">
              <thead>
                <tr>
                  <th>Category</th>
                  <th class="text-end">Amount</th>
                  <th class="text-center">Transactions</th>
                  <th class="text-end">Percentage</th>
                  <th>Distribution</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="category in categoryReport.breakdown.slice(0, 10)" :key="category.categoryId">
                  <td><strong>{{ category.categoryName }}</strong></td>
                  <td class="text-end">{{ formatCurrency(category.amount) }}</td>
                  <td class="text-center">{{ category.count }}</td>
                  <td class="text-end">{{ category.percentage.toFixed(1) }}%</td>
                  <td>
                    <div class="progress" style="height: 8px;">
                      <div class="progress-bar bg-danger" :style="{ width: category.percentage + '%' }"></div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Account Balances List -->
      <div v-if="accountsReport && accountsReport.accounts && accountsReport.accounts.length > 0" class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">💳 All Accounts</h5>
        </div>
        <div class="card-body">
          <div class="table-responsive">
            <table class="table table-hover">
              <thead>
                <tr>
                  <th>Account</th>
                  <th>Type</th>
                  <th class="text-end">Balance</th>
                  <th class="text-end">Percentage</th>
                  <th>Distribution</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="account in accountsReport.accounts" :key="account.accountId">
                  <td><strong>{{ account.accountName }}</strong></td>
                  <td>{{ account.accountType }}</td>
                  <td class="text-end" :class="account.balance >= 0 ? 'text-success' : 'text-danger'">
                    {{ formatCurrency(account.balance) }}
                  </td>
                  <td class="text-end">{{ account.percentage.toFixed(1) }}%</td>
                  <td>
                    <div class="progress" style="height: 8px;">
                      <div class="progress-bar bg-primary" :style="{ width: account.percentage + '%' }"></div>
                    </div>
                  </td>
                </tr>
              </tbody>
              <tfoot>
                <tr class="table-active">
                  <td colspan="2"><strong>Total Balance</strong></td>
                  <td class="text-end"><strong>{{ formatCurrency(accountsReport.totalBalance) }}</strong></td>
                  <td colspan="2"></td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>

      <!-- Monthly Summary Comparison -->
      <div v-if="dashboardData?.currentMonth && dashboardData?.previousMonth" class="card">
        <div class="card-header">
          <h5 class="mb-0">📅 Month-over-Month Comparison</h5>
        </div>
        <div class="card-body">
          <div class="row g-3">
            <div class="col-md-4">
              <div class="comparison-card">
                <div class="text-muted mb-1">Total Income</div>
                <div class="d-flex justify-content-between align-items-center">
                  <div>
                    <div class="h5 mb-0">{{ formatCurrency(dashboardData.currentMonth.totalIncome) }}</div>
                    <small class="text-muted">This Month</small>
                  </div>
                  <div class="text-end">
                    <div class="small text-muted">{{ formatCurrency(dashboardData.previousMonth.totalIncome) }}</div>
                    <div :class="getChangeClass(dashboardData.currentMonth.totalIncome - dashboardData.previousMonth.totalIncome)">
                      {{ getChangeText(dashboardData.currentMonth.totalIncome, dashboardData.previousMonth.totalIncome) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <div class="comparison-card">
                <div class="text-muted mb-1">Total Expense</div>
                <div class="d-flex justify-content-between align-items-center">
                  <div>
                    <div class="h5 mb-0">{{ formatCurrency(dashboardData.currentMonth.totalExpense) }}</div>
                    <small class="text-muted">This Month</small>
                  </div>
                  <div class="text-end">
                    <div class="small text-muted">{{ formatCurrency(dashboardData.previousMonth.totalExpense) }}</div>
                    <div :class="getChangeClass(dashboardData.previousMonth.totalExpense - dashboardData.currentMonth.totalExpense)">
                      {{ getChangeText(dashboardData.currentMonth.totalExpense, dashboardData.previousMonth.totalExpense, true) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <div class="comparison-card">
                <div class="text-muted mb-1">Net Savings</div>
                <div class="d-flex justify-content-between align-items-center">
                  <div>
                    <div class="h5 mb-0">{{ formatCurrency(dashboardData.currentMonth.netSavings) }}</div>
                    <small class="text-muted">This Month</small>
                  </div>
                  <div class="text-end">
                    <div class="small text-muted">{{ formatCurrency(dashboardData.previousMonth.netSavings) }}</div>
                    <div :class="getChangeClass(dashboardData.currentMonth.netSavings - dashboardData.previousMonth.netSavings)">
                      {{ getChangeText(dashboardData.currentMonth.netSavings, dashboardData.previousMonth.netSavings) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { Line, Bar, Pie, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import apiService from '@/services/api-backend'
import { format, startOfMonth, endOfMonth, subMonths, startOfQuarter, endOfQuarter, startOfYear, endOfYear, subYears } from 'date-fns'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

export default {
  name: 'ReportsView',
  components: {
    Line,
    Bar,
    Pie,
    Doughnut
  },
  setup() {
    // State
    const initialLoading = ref(true)
    const loading = ref(false)
    const loadingTrend = ref(false)
    const loadingCategories = ref(false)
    const loadingAccounts = ref(false)
    const loadingCashFlow = ref(false)
    const loadingNetWorth = ref(false)
    const error = ref(null)

    const selectedPeriod = ref('this_month')
    const customRange = ref({
      start: format(startOfMonth(new Date()), 'yyyy-MM-dd'),
      end: format(endOfMonth(new Date()), 'yyyy-MM-dd')
    })
    const dateRange = ref({
      start: format(startOfMonth(new Date()), 'yyyy-MM-dd'),
      end: format(endOfMonth(new Date()), 'yyyy-MM-dd')
    })

    const trendGroupBy = ref('month')
    const categoryLimit = ref(10)

    // Data
    const dashboardData = ref(null)
    const incomeTrendReport = ref(null)
    const categoryReport = ref(null)
    const accountsReport = ref(null)
    const cashFlowReport = ref(null)
    const netWorthReport = ref(null)
    const budgetReport = ref(null)

    // Chart Data
    const incomeTrendData = computed(() => {
      if (!incomeTrendReport.value || !incomeTrendReport.value.trend) {
        return { labels: [], datasets: [] }
      }

      const trend = incomeTrendReport.value.trend
      return {
        labels: trend.map(t => t.period),
        datasets: [
          {
            label: 'Income',
            data: trend.map(t => t.income),
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            tension: 0.4,
            fill: true
          },
          {
            label: 'Expense',
            data: trend.map(t => t.expense),
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            tension: 0.4,
            fill: true
          },
          {
            label: 'Net',
            data: trend.map(t => t.net),
            borderColor: 'rgb(153, 102, 255)',
            backgroundColor: 'rgba(153, 102, 255, 0.2)',
            tension: 0.4,
            fill: false,
            borderDash: [5, 5]
          }
        ]
      }
    })

    const categoryData = computed(() => {
      if (!categoryReport.value || !categoryReport.value.breakdown) {
        return { labels: [], datasets: [] }
      }

      const breakdown = categoryReport.value.breakdown
      return {
        labels: breakdown.map(c => c.categoryName),
        datasets: [{
          data: breakdown.map(c => c.amount),
          backgroundColor: [
            'rgba(255, 99, 132, 0.8)',
            'rgba(54, 162, 235, 0.8)',
            'rgba(255, 206, 86, 0.8)',
            'rgba(75, 192, 192, 0.8)',
            'rgba(153, 102, 255, 0.8)',
            'rgba(255, 159, 64, 0.8)',
            'rgba(199, 199, 199, 0.8)',
            'rgba(83, 102, 255, 0.8)',
            'rgba(255, 99, 255, 0.8)',
            'rgba(99, 255, 132, 0.8)'
          ]
        }]
      }
    })

    const accountData = computed(() => {
      if (!accountsReport.value || !accountsReport.value.accounts) {
        return { labels: [], datasets: [] }
      }

      const accounts = accountsReport.value.accounts
      return {
        labels: accounts.map(a => a.accountName),
        datasets: [{
          data: accounts.map(a => a.balance),
          backgroundColor: [
            'rgba(54, 162, 235, 0.8)',
            'rgba(75, 192, 192, 0.8)',
            'rgba(255, 206, 86, 0.8)',
            'rgba(153, 102, 255, 0.8)',
            'rgba(255, 159, 64, 0.8)',
            'rgba(255, 99, 132, 0.8)'
          ]
        }]
      }
    })

    const cashFlowData = computed(() => {
      if (!cashFlowReport.value || !cashFlowReport.value.monthlyBreakdown) {
        return { labels: [], datasets: [] }
      }

      const breakdown = cashFlowReport.value.monthlyBreakdown
      return {
        labels: breakdown.map(m => m.month),
        datasets: [
          {
            label: 'Inflow',
            data: breakdown.map(m => m.inflow),
            backgroundColor: 'rgba(75, 192, 192, 0.8)'
          },
          {
            label: 'Outflow',
            data: breakdown.map(m => m.outflow),
            backgroundColor: 'rgba(255, 99, 132, 0.8)'
          }
        ]
      }
    })

    const netWorthData = computed(() => {
      if (!netWorthReport.value || !netWorthReport.value.trend) {
        return { labels: [], datasets: [] }
      }

      const trend = netWorthReport.value.trend
      return {
        labels: trend.map(t => t.period),
        datasets: [
          {
            label: 'Net Worth',
            data: trend.map(t => t.netWorth),
            borderColor: 'rgb(153, 102, 255)',
            backgroundColor: 'rgba(153, 102, 255, 0.2)',
            tension: 0.4,
            fill: true
          },
          {
            label: 'Assets',
            data: trend.map(t => t.totalAssets),
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.1)',
            tension: 0.4,
            fill: false
          },
          {
            label: 'Liabilities',
            data: trend.map(t => t.totalLiabilities),
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.1)',
            tension: 0.4,
            fill: false
          }
        ]
      }
    })

    // Chart Options
    const trendChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top'
        },
        tooltip: {
          mode: 'index',
          intersect: false,
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || ''
              if (label) label += ': '
              if (context.parsed.y !== null) {
                label += formatCurrency(context.parsed.y)
              }
              return label
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              return '$' + value.toLocaleString()
            }
          }
        }
      }
    }

    const categoryChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'right'
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              const label = context.label || ''
              const value = context.parsed || 0
              const percentage = context.dataset.data.reduce((a, b) => a + b, 0)
              const percent = ((value / percentage) * 100).toFixed(1)
              return `${label}: ${formatCurrency(value)} (${percent}%)`
            }
          }
        }
      }
    }

    const accountChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom'
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              const label = context.label || ''
              const value = context.parsed || 0
              return `${label}: ${formatCurrency(value)}`
            }
          }
        }
      }
    }

    const cashFlowChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top'
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || ''
              if (label) label += ': '
              if (context.parsed.y !== null) {
                label += formatCurrency(context.parsed.y)
              }
              return label
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              return '$' + value.toLocaleString()
            }
          }
        }
      }
    }

    const netWorthChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top'
        },
        tooltip: {
          mode: 'index',
          intersect: false,
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || ''
              if (label) label += ': '
              if (context.parsed.y !== null) {
                label += formatCurrency(context.parsed.y)
              }
              return label
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              return '$' + value.toLocaleString()
            }
          }
        }
      }
    }

    // Computed
    const savingsChangeText = computed(() => {
      if (!dashboardData.value?.currentMonth || !dashboardData.value?.previousMonth) return ''
      const current = dashboardData.value.currentMonth.netSavings
      const previous = dashboardData.value.previousMonth.netSavings
      return getChangeText(current, previous)
    })

    const savingsChangeClass = computed(() => {
      if (!dashboardData.value?.currentMonth || !dashboardData.value?.previousMonth) return ''
      const current = dashboardData.value.currentMonth.netSavings
      const previous = dashboardData.value.previousMonth.netSavings
      return getChangeClass(current - previous)
    })

    // Methods
    const formatCurrency = (value) => {
      if (value === null || value === undefined) return '$0.00'
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        minimumFractionDigits: 2
      }).format(value)
    }

    const getChangeText = (current, previous, inverse = false) => {
      if (!previous) return 'N/A'
      const change = ((current - previous) / Math.abs(previous)) * 100
      const absChange = Math.abs(change)
      const direction = inverse ? (change > 0 ? '↓' : '↑') : (change > 0 ? '↑' : '↓')
      return `${direction} ${absChange.toFixed(1)}%`
    }

    const getChangeClass = (change) => {
      if (change > 0) return 'text-success small'
      if (change < 0) return 'text-danger small'
      return 'text-muted small'
    }

    const getBudgetStatusClass = (status) => {
      const statusMap = {
        'over': 'bg-danger',
        'on-track': 'bg-success',
        'under': 'bg-warning'
      }
      return statusMap[status] || 'bg-secondary'
    }

    const getBudgetProgressClass = (percentage) => {
      if (percentage >= 100) return 'bg-danger'
      if (percentage >= 80) return 'bg-warning'
      return 'bg-success'
    }

    const handlePeriodChange = () => {
      if (selectedPeriod.value !== 'custom') {
        updateDateRangeFromPeriod()
        refreshAllData()
      }
    }

    const updateDateRangeFromPeriod = () => {
      const now = new Date()
      let start, end

      switch (selectedPeriod.value) {
        case 'this_month':
          start = startOfMonth(now)
          end = endOfMonth(now)
          break
        case 'last_month':
          start = startOfMonth(subMonths(now, 1))
          end = endOfMonth(subMonths(now, 1))
          break
        case 'this_quarter':
          start = startOfQuarter(now)
          end = endOfQuarter(now)
          break
        case 'last_quarter':
          const lastQ = subMonths(now, 3)
          start = startOfQuarter(lastQ)
          end = endOfQuarter(lastQ)
          break
        case 'this_year':
          start = startOfYear(now)
          end = endOfYear(now)
          break
        case 'last_year':
          start = startOfYear(subYears(now, 1))
          end = endOfYear(subYears(now, 1))
          break
        default:
          return
      }

      dateRange.value = {
        start: format(start, 'yyyy-MM-dd'),
        end: format(end, 'yyyy-MM-dd')
      }
    }

    const applyCustomRange = () => {
      dateRange.value = { ...customRange.value }
      refreshAllData()
    }

    const changeTrendGrouping = async (groupBy) => {
      trendGroupBy.value = groupBy
      await loadTrendData()
    }

    const loadDashboardSummary = async () => {
      try {
        const response = await apiService.reports.getDashboardSummary()
        dashboardData.value = response.data
      } catch (err) {
        console.error('Error loading dashboard summary:', err)
        throw err
      }
    }

    const loadTrendData = async () => {
      loadingTrend.value = true
      try {
        const response = await apiService.reports.getIncomeExpense(
          dateRange.value.start,
          dateRange.value.end,
          trendGroupBy.value
        )
        incomeTrendReport.value = response.data
      } catch (err) {
        console.error('Error loading trend data:', err)
      } finally {
        loadingTrend.value = false
      }
    }

    const loadCategoryData = async () => {
      loadingCategories.value = true
      try {
        const response = await apiService.reports.getCategoryAnalysis(
          dateRange.value.start,
          dateRange.value.end,
          'expense',
          categoryLimit.value
        )
        categoryReport.value = response.data
      } catch (err) {
        console.error('Error loading category data:', err)
      } finally {
        loadingCategories.value = false
      }
    }

    const loadAccountsData = async () => {
      loadingAccounts.value = true
      try {
        const response = await apiService.reports.getAccountBalances()
        accountsReport.value = response.data
      } catch (err) {
        console.error('Error loading accounts data:', err)
      } finally {
        loadingAccounts.value = false
      }
    }

    const loadCashFlowData = async () => {
      loadingCashFlow.value = true
      try {
        const response = await apiService.reports.getCashFlow(
          dateRange.value.start,
          dateRange.value.end,
          'month'
        )
        cashFlowReport.value = response.data
      } catch (err) {
        console.error('Error loading cash flow data:', err)
      } finally {
        loadingCashFlow.value = false
      }
    }

    const loadNetWorthData = async () => {
      loadingNetWorth.value = true
      try {
        const response = await apiService.reports.getNetWorth(
          dateRange.value.start,
          dateRange.value.end,
          'month'
        )
        netWorthReport.value = response.data
      } catch (err) {
        console.error('Error loading net worth data:', err)
      } finally {
        loadingNetWorth.value = false
      }
    }

    const loadBudgetData = async () => {
      try {
        const currentMonth = format(new Date(), 'yyyy-MM')
        const response = await apiService.reports.getBudgetPerformance(currentMonth)
        budgetReport.value = response.data
      } catch (err) {
        console.error('Error loading budget data:', err)
      }
    }

    const refreshAllData = async () => {
      loading.value = true
      error.value = null
      try {
        await Promise.all([
          loadDashboardSummary(),
          loadTrendData(),
          loadCategoryData(),
          loadAccountsData(),
          loadCashFlowData(),
          loadNetWorthData(),
          loadBudgetData()
        ])
      } catch (err) {
        error.value = err.message || 'Failed to load reports data'
      } finally {
        loading.value = false
        initialLoading.value = false
      }
    }

    // Lifecycle
    onMounted(async () => {
      updateDateRangeFromPeriod()
      await refreshAllData()
    })

    return {
      // State
      initialLoading,
      loading,
      loadingTrend,
      loadingCategories,
      loadingAccounts,
      loadingCashFlow,
      loadingNetWorth,
      error,
      selectedPeriod,
      customRange,
      dateRange,
      trendGroupBy,
      categoryLimit,

      // Data
      dashboardData,
      incomeTrendReport,
      categoryReport,
      accountsReport,
      cashFlowReport,
      netWorthReport,
      budgetReport,

      // Chart Data
      incomeTrendData,
      categoryData,
      accountData,
      cashFlowData,
      netWorthData,

      // Chart Options
      trendChartOptions,
      categoryChartOptions,
      accountChartOptions,
      cashFlowChartOptions,
      netWorthChartOptions,

      // Computed
      savingsChangeText,
      savingsChangeClass,

      // Methods
      formatCurrency,
      getChangeText,
      getChangeClass,
      getBudgetStatusClass,
      getBudgetProgressClass,
      handlePeriodChange,
      applyCustomRange,
      changeTrendGrouping,
      loadCategoryData,
      refreshAllData
    }
  }
}
</script>

<style scoped>
.reports-view {
  padding: 1.5rem;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.metric-card {
  padding: 1.5rem;
  border-radius: 12px;
  color: white;
  text-align: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.metric-card:hover {
  transform: translateY(-4px);
}

.purple-gradient {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.green-gradient {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.red-gradient {
  background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
}

.blue-gradient {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.metric-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.metric-value {
  font-size: 1.5rem;
  font-weight: bold;
  margin-bottom: 0.25rem;
}

.metric-label {
  font-size: 0.875rem;
  opacity: 0.9;
}

.metric-change {
  font-size: 0.75rem;
  margin-top: 0.25rem;
}

.card {
  border: none;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background: white;
}

.card-header {
  background: white;
  border-bottom: 1px solid #e9ecef;
  border-radius: 12px 12px 0 0 !important;
  padding: 1rem 1.5rem;
}

.budget-card {
  padding: 1rem;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  background: white;
}

.comparison-card {
  padding: 1rem;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  background: #f8f9fa;
}

.text-purple {
  color: #667eea;
}

.fade-in {
  animation: fadeIn 0.5s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Chart containers */
.card-body canvas {
  max-height: 400px;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .metric-card {
    margin-bottom: 1rem;
  }

  .card-header h5 {
    font-size: 1rem;
  }
}
</style>
