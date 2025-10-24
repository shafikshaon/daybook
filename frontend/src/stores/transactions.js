import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'
import { useAccountsStore } from './accounts'
import { useSavingsGoalsStore } from './savingsGoals'

export const useTransactionsStore = defineStore('transactions', {
  state: () => ({
    transactions: [],
    categories: [
      // Income categories
      { id: 'salary', name: 'Salary', type: 'income', icon: '💼', color: '#10b981' },
      { id: 'freelance', name: 'Freelance', type: 'income', icon: '💻', color: '#10b981' },
      { id: 'investment_income', name: 'Investment Income', type: 'income', icon: '📈', color: '#10b981' },
      { id: 'other_income', name: 'Other Income', type: 'income', icon: '💰', color: '#10b981' },

      // Expense categories
      { id: 'food', name: 'Food & Dining', type: 'expense', icon: '🍔', color: '#ef4444' },
      { id: 'transport', name: 'Transportation', type: 'expense', icon: '🚗', color: '#ef4444' },
      { id: 'shopping', name: 'Shopping', type: 'expense', icon: '🛍️', color: '#ef4444' },
      { id: 'entertainment', name: 'Entertainment', type: 'expense', icon: '🎬', color: '#ef4444' },
      { id: 'utilities', name: 'Utilities', type: 'expense', icon: '💡', color: '#ef4444' },
      { id: 'healthcare', name: 'Healthcare', type: 'expense', icon: '🏥', color: '#ef4444' },
      { id: 'education', name: 'Education', type: 'expense', icon: '📚', color: '#ef4444' },
      { id: 'housing', name: 'Housing', type: 'expense', icon: '🏠', color: '#ef4444' },
      { id: 'insurance', name: 'Insurance', type: 'expense', icon: '🛡️', color: '#ef4444' },
      { id: 'subscriptions', name: 'Subscriptions', type: 'expense', icon: '📱', color: '#ef4444' },
      { id: 'credit_card_payment', name: 'Credit Card Payment', type: 'expense', icon: '💳', color: '#ef4444' },
      { id: 'other_expense', name: 'Other Expense', type: 'expense', icon: '💸', color: '#ef4444' },

      // Transfer category
      { id: 'transfer', name: 'Transfer', type: 'transfer', icon: '🔄', color: '#3b82f6' }
    ],
    tags: [],
    recurringTransactions: []
  }),

  getters: {
    allTransactions: (state) => {
      return state.transactions.sort((a, b) => new Date(b.date) - new Date(a.date))
    },

    incomeTransactions: (state) => {
      return state.transactions.filter(t => t.type === 'income')
    },

    expenseTransactions: (state) => {
      return state.transactions.filter(t => t.type === 'expense')
    },

    transferTransactions: (state) => {
      return state.transactions.filter(t => t.type === 'transfer')
    },

    transactionsByDateRange: (state) => (startDate, endDate) => {
      return state.transactions.filter(t => {
        const date = new Date(t.date)
        return date >= new Date(startDate) && date <= new Date(endDate)
      })
    },

    transactionsByCategory: (state) => (categoryId) => {
      return state.transactions.filter(t => t.categoryId === categoryId)
    },

    transactionsByAccount: (state) => (accountId) => {
      return state.transactions.filter(t => t.accountId === accountId)
    },

    totalIncome: (state) => (startDate = null, endDate = null) => {
      let transactions = state.transactions.filter(t => t.type === 'income')

      if (startDate && endDate) {
        transactions = transactions.filter(t => {
          const date = new Date(t.date)
          return date >= new Date(startDate) && date <= new Date(endDate)
        })
      }

      return transactions.reduce((total, t) => total + t.amount, 0)
    },

    totalExpense: (state) => (startDate = null, endDate = null) => {
      let transactions = state.transactions.filter(t => t.type === 'expense')

      if (startDate && endDate) {
        transactions = transactions.filter(t => {
          const date = new Date(t.date)
          return date >= new Date(startDate) && date <= new Date(endDate)
        })
      }

      return transactions.reduce((total, t) => total + t.amount, 0)
    },

    categoryBreakdown: (state) => (type = 'expense') => {
      const transactions = state.transactions.filter(t => t.type === type)
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

      return Object.values(breakdown)
    },

    incomeCategories: (state) => state.categories.filter(c => c.type === 'income'),
    expenseCategories: (state) => state.categories.filter(c => c.type === 'expense'),
    transferCategories: (state) => state.categories.filter(c => c.type === 'transfer'),

    getCategoryById: (state) => (id) => {
      return state.categories.find(c => c.id === id)
    }
  },

  actions: {
    async fetchTransactions() {
      try {
        const response = await apiService.get('transactions')
        this.transactions = response.data || []
      } catch (error) {
        console.error('Error fetching transactions:', error)
        throw error
      }
    },

    async createTransaction(transactionData) {
      try {
        const response = await apiService.post('transactions', transactionData)
        this.transactions.push(response.data)

        // Refresh account balances from backend (backend handles balance updates)
        const accountsStore = useAccountsStore()
        await accountsStore.fetchAccounts()

        return response.data
      } catch (error) {
        console.error('Error creating transaction:', error)
        throw error
      }
    },

    async updateTransaction(id, transactionData) {
      try {
        const response = await apiService.put('transactions', id, transactionData)

        const index = this.transactions.findIndex(t => t.id === id)
        if (index !== -1) {
          this.transactions[index] = response.data
        }

        // Refresh account balances from backend (backend handles balance updates)
        const accountsStore = useAccountsStore()
        await accountsStore.fetchAccounts()

        return response.data
      } catch (error) {
        console.error('Error updating transaction:', error)
        throw error
      }
    },

    async deleteTransaction(id) {
      try {
        await apiService.delete('transactions', id)

        // Refresh account balances from backend (backend handles balance updates)
        const accountsStore = useAccountsStore()
        await accountsStore.fetchAccounts()

        this.transactions = this.transactions.filter(t => t.id !== id)
      } catch (error) {
        console.error('Error deleting transaction:', error)
        throw error
      }
    },

    async transferFunds(fromAccountId, toAccountId, amount, description = 'Transfer between accounts', date = new Date().toISOString()) {
      try {
        const transactionData = {
          type: 'transfer',
          amount,
          categoryId: 'transfer',
          accountId: fromAccountId,
          toAccountId: toAccountId,
          date,
          description,
          tags: [],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }

        return await this.createTransaction(transactionData)
      } catch (error) {
        console.error('Error transferring funds:', error)
        throw error
      }
    },

    async transferToSavingsGoal(fromAccountId, savingsGoalId, amount, description = 'Transfer to savings goal', date = new Date().toISOString()) {
      try {
        // Deduct from source account
        const accountsStore = useAccountsStore()
        await accountsStore.updateBalance(fromAccountId, amount, 'subtract')

        // Add contribution to savings goal
        const savingsGoalsStore = useSavingsGoalsStore()
        await savingsGoalsStore.addContribution(savingsGoalId, amount, date)

        // Create transaction record
        const transactionData = {
          type: 'expense',
          amount,
          categoryId: 'other_expense',
          accountId: fromAccountId,
          date,
          description,
          tags: ['savings_goal', savingsGoalId],
          savingsGoalId: savingsGoalId,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }

        const response = await apiService.post('transactions', transactionData)
        this.transactions.push(response.data)

        return response.data
      } catch (error) {
        console.error('Error transferring to savings goal:', error)
        throw error
      }
    },

    async bulkImport(transactions) {
      try {
        const response = await apiService.bulkCreate('transactions', transactions)
        this.transactions.push(...response.data)

        // Update account balances
        const accountsStore = useAccountsStore()
        for (const transaction of response.data) {
          if (transaction.type === 'income') {
            await accountsStore.updateBalance(transaction.accountId, transaction.amount, 'add')
          } else if (transaction.type === 'transfer') {
            await accountsStore.updateBalance(transaction.accountId, transaction.amount, 'subtract')
            if (transaction.toAccountId) {
              await accountsStore.updateBalance(transaction.toAccountId, transaction.amount, 'add')
            }
          } else {
            await accountsStore.updateBalance(transaction.accountId, transaction.amount, 'subtract')
          }
        }

        return response.data
      } catch (error) {
        console.error('Error bulk importing transactions:', error)
        throw error
      }
    },

    async searchTransactions(filters) {
      try {
        const response = await apiService.query('transactions', filters)
        return response.data
      } catch (error) {
        console.error('Error searching transactions:', error)
        throw error
      }
    },

    // Recurring transactions
    async fetchRecurringTransactions() {
      try {
        const response = await apiService.get('recurring-transactions')
        this.recurringTransactions = response.data || []
      } catch (error) {
        console.error('Error fetching recurring transactions:', error)
        throw error
      }
    },

    async createRecurringTransaction(recurringData) {
      try {
        const response = await apiService.post('recurring-transactions', recurringData)
        this.recurringTransactions.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating recurring transaction:', error)
        throw error
      }
    },

    async processRecurringTransactions() {
      const today = new Date()

      for (const recurring of this.recurringTransactions) {
        if (!recurring.enabled) continue

        const lastProcessed = recurring.lastProcessed ? new Date(recurring.lastProcessed) : null
        const nextDate = this.calculateNextDate(lastProcessed || new Date(recurring.startDate), recurring.frequency)

        if (nextDate <= today) {
          await this.createTransaction({
            ...recurring.transactionTemplate,
            date: nextDate.toISOString(),
            recurringId: recurring.id
          })

          await apiService.put('recurring-transactions', recurring.id, {
            ...recurring,
            lastProcessed: nextDate.toISOString()
          })
        }
      }
    },

    calculateNextDate(lastDate, frequency) {
      const next = new Date(lastDate)

      switch (frequency) {
        case 'daily':
          next.setDate(next.getDate() + 1)
          break
        case 'weekly':
          next.setDate(next.getDate() + 7)
          break
        case 'biweekly':
          next.setDate(next.getDate() + 14)
          break
        case 'monthly':
          next.setMonth(next.getMonth() + 1)
          break
        case 'quarterly':
          next.setMonth(next.getMonth() + 3)
          break
        case 'yearly':
          next.setFullYear(next.getFullYear() + 1)
          break
      }

      return next
    },

    // Tags management
    async fetchTags() {
      try {
        const response = await apiService.get('tags')
        this.tags = response.data || []
      } catch (error) {
        console.error('Error fetching tags:', error)
        throw error
      }
    },

    async createTag(tagName) {
      try {
        const tag = { name: tagName, color: this.getRandomColor() }
        const response = await apiService.post('tags', tag)
        this.tags.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating tag:', error)
        throw error
      }
    },

    getRandomColor() {
      const colors = ['#6f42c1', '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#ec4899']
      return colors[Math.floor(Math.random() * colors.length)]
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
