import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useTransactionsStore } from './transactions'
import { useSettingsStore } from './settings'

export const useBudgetsStore = defineStore('budgets', {
  state: () => ({
    budgets: [], // Now stores BudgetProgress objects from backend
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
      return state.budgets.filter(b => b.budget?.enabled || b.enabled)
    },

    budgetsByPeriod: (state) => (period) => {
      return state.budgets.filter(b => {
        const budget = b.budget || b
        return budget.period === period
      })
    },

    getBudgetById: (state) => (id) => {
      return state.budgets.find(b => {
        const budget = b.budget || b
        return budget.id === id
      })
    },

    // Use backend-provided progress data
    budgetProgress: (state) => (budgetId) => {
      const budgetProgress = state.budgets.find(b => {
        const budget = b.budget || b
        return budget.id === budgetId
      })

      if (!budgetProgress) return null

      // If this is already a BudgetProgress object from backend
      if (budgetProgress.totalSpent !== undefined) {
        return {
          spent: budgetProgress.totalSpent,
          remaining: budgetProgress.remaining,
          percentage: budgetProgress.percentageUsed,
          isOverBudget: budgetProgress.isOverBudget,
          status: budgetProgress.isOverBudget ? 'danger' : 
                  budgetProgress.alertTriggered ? 'warning' : 'success',
          startDate: budgetProgress.startDate,
          endDate: budgetProgress.endDate
        }
      }

      // Fallback for old format (shouldn't happen with new API)
      return {
        spent: 0,
        remaining: budgetProgress.amount || 0,
        percentage: 0,
        isOverBudget: false,
        status: 'success'
      }
    },

    budgetAlerts: (state) => {
      const alerts = []

      state.budgets.forEach(budgetProgress => {
        const budget = budgetProgress.budget || budgetProgress

        if (!budget.enabled) return

        const transactionsStore = useTransactionsStore()
        const category = transactionsStore.getCategoryById(budget.categoryId)

        // Use backend-provided alert status
        if (budgetProgress.isOverBudget) {
          alerts.push({
            id: budget.id,
            type: 'danger',
            message: `Budget exceeded for ${category?.name || 'Unknown'}`,
            amount: budgetProgress.totalSpent,
            budget: budget.amount
          })
        } else if (budgetProgress.alertTriggered) {
          alerts.push({
            id: budget.id,
            type: 'warning',
            message: `Approaching budget limit for ${category?.name || 'Unknown'}`,
            amount: budgetProgress.totalSpent,
            budget: budget.amount
          })
        }
      })

      return alerts
    },

    totalBudgeted: (state) => {
      return state.budgets
        .filter(b => {
          const budget = b.budget || b
          return budget.enabled
        })
        .reduce((sum, b) => {
          const budget = b.budget || b
          return sum + budget.amount
        }, 0)
    },

    totalSpent: (state) => {
      return state.budgets
        .filter(b => {
          const budget = b.budget || b
          return budget.enabled
        })
        .reduce((sum, b) => {
          return sum + (b.totalSpent || 0)
        }, 0)
    }
  },

  actions: {
    async fetchBudgets() {
      try {
        // Backend now returns BudgetProgress objects with automatic expense calculation
        const response = await apiService.get('budgets')
        this.budgets = response.data || []
      } catch (error) {
        console.error('Error fetching budgets:', error)
        throw error
      }
    },

    async fetchBudgetProgress(budgetId) {
      try {
        const response = await apiService.getBudgetProgress(budgetId)
        // Update the specific budget in the list
        const index = this.budgets.findIndex(b => {
          const budget = b.budget || b
          return budget.id === budgetId
        })
        if (index !== -1) {
          this.budgets[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error fetching budget progress:', error)
        throw error
      }
    },
