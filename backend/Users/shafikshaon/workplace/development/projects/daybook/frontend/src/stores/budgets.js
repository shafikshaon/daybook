  getters: {
    allBudgets: (state) => state.budgets,

    activeBudgets: (state) => {
      console.log('activeBudgets - all budgets:', state.budgets)
      return state.budgets.filter(item => {
        // Backend returns BudgetProgress objects with nested budget
        const budget = item.budget || item
        const isEnabled = budget.enabled !== false
        console.log('Budget:', budget.id, 'enabled:', isEnabled)
        return isEnabled
      })
    },

    budgetsByPeriod: (state) => (period) => {
      return state.budgets.filter(item => {
        const budget = item.budget || item
        return budget.period === period
      })
    },

    getBudgetById: (state) => (id) => {
      return state.budgets.find(item => {
        const budget = item.budget || item
        return budget.id === id
      })
    },

    // Backend now provides progress data, no need to calculate
    budgetProgress: (state) => (budgetId) => {
      const item = state.budgets.find(b => {
        const budget = b.budget || b
        return budget.id === budgetId
      })

      if (!item) return null

      // If this is a BudgetProgress object from backend
      if (item.totalSpent !== undefined) {
        return {
          spent: item.totalSpent,
          remaining: item.remaining,
          percentage: item.percentageUsed,
          isOverBudget: item.isOverBudget,
          status: item.isOverBudget ? 'danger' : 
                  item.alertTriggered ? 'warning' : 'success',
          startDate: item.startDate,
          endDate: item.endDate
        }
      }

      // Fallback
      return {
        spent: 0,
        remaining: (item.budget || item).amount || 0,
        percentage: 0,
        isOverBudget: false,
        status: 'success'
      }
    },

    budgetAlerts: (state) => {
      const alerts = []

      state.budgets.forEach(item => {
        const budget = item.budget || item

        if (!budget.enabled) return

        const transactionsStore = useTransactionsStore()
        const category = transactionsStore.getCategoryById(budget.categoryId)

        // Use backend-provided alert status
        if (item.isOverBudget) {
          alerts.push({
            id: budget.id,
            type: 'danger',
            message: `Budget exceeded for ${category?.name || 'Unknown'}`,
            amount: item.totalSpent,
            budget: budget.amount
          })
        } else if (item.alertTriggered) {
          alerts.push({
            id: budget.id,
            type: 'warning',
            message: `Approaching budget limit for ${category?.name || 'Unknown'}`,
            amount: item.totalSpent,
            budget: budget.amount
          })
        }
      })

      return alerts
    },

    totalBudgeted: (state) => {
      return state.budgets
        .filter(item => {
          const budget = item.budget || item
          return budget.enabled
        })
        .reduce((sum, item) => {
          const budget = item.budget || item
          return sum + budget.amount
        }, 0)
    },

    totalSpent: (state) => {
      return state.budgets
        .filter(item => {
          const budget = item.budget || item
          return budget.enabled
        })
        .reduce((sum, item) => {
          return sum + (item.totalSpent || 0)
        }, 0)
    }
  },
