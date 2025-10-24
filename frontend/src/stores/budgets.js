import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useTransactionsStore } from './transactions'
import { useSettingsStore } from './settings'

export const useBudgetsStore = defineStore('budgets', {
  state: () => ({
    budgets: [],
    budgetPeriods: [
      { value: 'weekly', label: 'Weekly' },
      { value: 'monthly', label: 'Monthly' },
      { value: 'quarterly', label: 'Quarterly' },
      { value: 'yearly', label: 'Yearly' },
      { value: 'custom', label: 'Custom' }
    ]
  }),

  getters: {
    allBudgets: (state) => state.budgets,

    activeBudgets: (state) => {
      return state.budgets.filter(b => b.enabled)
    },

    budgetsByPeriod: (state) => (period) => {
      return state.budgets.filter(b => b.period === period)
    },

    getBudgetById: (state) => (id) => {
      return state.budgets.find(b => b.id === id)
    },

    budgetProgress: (state) => (budgetId) => {
      const budget = state.budgets.find(b => b.id === budgetId)
      if (!budget) return null

      const transactionsStore = useTransactionsStore()
      const { startDate, endDate } = state.getBudgetDateRange(budget)

      const spent = transactionsStore.transactions
        .filter(t => {
          const tDate = new Date(t.date)
          return t.type === 'expense' &&
                 t.categoryId === budget.categoryId &&
                 tDate >= startDate &&
                 tDate <= endDate
        })
        .reduce((sum, t) => sum + t.amount, 0)

      const percentage = (spent / budget.amount) * 100
      const remaining = budget.amount - spent

      return {
        spent,
        remaining,
        percentage,
        isOverBudget: spent > budget.amount,
        status: percentage >= 100 ? 'danger' : percentage >= 80 ? 'warning' : 'success'
      }
    },

    getBudgetDateRange: () => (budget) => {
      const today = new Date()
      let startDate, endDate

      switch (budget.period) {
        case 'weekly':
          startDate = new Date(today.getFullYear(), today.getMonth(), today.getDate() - today.getDay())
          endDate = new Date(startDate)
          endDate.setDate(endDate.getDate() + 6)
          break

        case 'monthly':
          startDate = new Date(today.getFullYear(), today.getMonth(), 1)
          endDate = new Date(today.getFullYear(), today.getMonth() + 1, 0)
          break

        case 'quarterly':
          const quarter = Math.floor(today.getMonth() / 3)
          startDate = new Date(today.getFullYear(), quarter * 3, 1)
          endDate = new Date(today.getFullYear(), (quarter + 1) * 3, 0)
          break

        case 'yearly':
          startDate = new Date(today.getFullYear(), 0, 1)
          endDate = new Date(today.getFullYear(), 11, 31)
          break

        case 'custom':
          startDate = new Date(budget.customStartDate)
          endDate = new Date(budget.customEndDate)
          break

        default:
          startDate = today
          endDate = today
      }

      return { startDate, endDate }
    },

    budgetAlerts: (state) => {
      const alerts = []

      state.budgets.forEach(budget => {
        if (!budget.enabled) return

        const progress = state.budgetProgress(budget.id)
        if (!progress) return

        const transactionsStore = useTransactionsStore()
        const category = transactionsStore.getCategoryById(budget.categoryId)

        if (progress.percentage >= 100) {
          alerts.push({
            id: budget.id,
            type: 'danger',
            message: `Budget exceeded for ${category?.name || 'Unknown'}`,
            amount: progress.spent,
            budget: budget.amount
          })
        } else if (progress.percentage >= 80) {
          alerts.push({
            id: budget.id,
            type: 'warning',
            message: `Approaching budget limit for ${category?.name || 'Unknown'}`,
            amount: progress.spent,
            budget: budget.amount
          })
        }
      })

      return alerts
    },

    totalBudgeted: (state) => {
      return state.budgets
        .filter(b => b.enabled)
        .reduce((sum, b) => sum + b.amount, 0)
    },

    totalSpent: (state) => {
      let total = 0

      state.budgets.forEach(budget => {
        if (!budget.enabled) return
        const progress = state.budgetProgress(budget.id)
        if (progress) {
          total += progress.spent
        }
      })

      return total
    }
  },

  actions: {
    async fetchBudgets() {
      try {
        const response = await apiService.get('budgets')
        this.budgets = response.data || []
      } catch (error) {
        console.error('Error fetching budgets:', error)
        throw error
      }
    },

    async createBudget(budgetData) {
      try {
        const response = await apiService.post('budgets', {
          ...budgetData,
          enabled: true,
          rollover: budgetData.rollover || false
        })
        this.budgets.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating budget:', error)
        throw error
      }
    },

    async updateBudget(id, budgetData) {
      try {
        const response = await apiService.put('budgets', id, budgetData)
        const index = this.budgets.findIndex(b => b.id === id)
        if (index !== -1) {
          this.budgets[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating budget:', error)
        throw error
      }
    },

    async deleteBudget(id) {
      try {
        await apiService.delete('budgets', id)
        this.budgets = this.budgets.filter(b => b.id !== id)
      } catch (error) {
        console.error('Error deleting budget:', error)
        throw error
      }
    },

    async toggleBudget(id) {
      const budget = this.getBudgetById(id)
      if (budget) {
        await this.updateBudget(id, { ...budget, enabled: !budget.enabled })
      }
    },

    calculateForecast(budgetId, monthsAhead = 3) {
      const budget = this.getBudgetById(budgetId)
      if (!budget) return []

      const transactionsStore = useTransactionsStore()
      const transactions = transactionsStore.transactionsByCategory(budget.categoryId)
        .filter(t => t.type === 'expense')

      if (transactions.length === 0) return []

      const avgMonthlySpending = transactions.reduce((sum, t) => sum + t.amount, 0) /
        Math.max(1, this.getMonthsSpan(transactions))

      const forecast = []
      for (let i = 1; i <= monthsAhead; i++) {
        const projectedSpending = avgMonthlySpending * i
        forecast.push({
          month: i,
          projected: projectedSpending,
          budgetLimit: budget.amount * i,
          overBudget: projectedSpending > (budget.amount * i)
        })
      }

      return forecast
    },

    getMonthsSpan(transactions) {
      if (transactions.length === 0) return 1

      const dates = transactions.map(t => new Date(t.date))
      const earliest = new Date(Math.min(...dates))
      const latest = new Date(Math.max(...dates))

      const months = (latest.getFullYear() - earliest.getFullYear()) * 12 +
        (latest.getMonth() - earliest.getMonth())

      return Math.max(1, months || 1)
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
