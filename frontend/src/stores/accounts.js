import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useAccountsStore = defineStore('accounts', {
  state: () => ({
    accounts: [],
    accountTypes: [
      // Liquid Assets
      { value: 'cash', label: 'Cash', icon: '💵', category: 'liquid' },
      { value: 'checking', label: 'Checking Account', icon: '🏦', category: 'liquid' },
      { value: 'savings', label: 'Savings Account', icon: '💰', category: 'liquid' },
      { value: 'money_market', label: 'Money Market Account', icon: '📊', category: 'liquid' },

      // Credit Accounts
      { value: 'credit_card', label: 'Credit Card', icon: '💳', category: 'credit' },
      { value: 'line_of_credit', label: 'Line of Credit', icon: '💳', category: 'credit' },

      // Investment Related (Cash in brokerage)
      { value: 'brokerage_cash', label: 'Brokerage Cash', icon: '🏦', category: 'investment' },

      // Other
      { value: 'prepaid', label: 'Prepaid Card', icon: '💳', category: 'liquid' },
      { value: 'other', label: 'Other', icon: '🏦', category: 'other' }
    ]
  }),

  getters: {
    allAccounts: (state) => state.accounts,

    accountsByType: (state) => (type) => {
      return state.accounts.filter(account => account.type === type)
    },

    totalBalance: (state) => {
      return state.accounts.reduce((total, account) => {
        if (account.type === 'credit_card') {
          return total - account.balance
        }
        return total + account.balance
      }, 0)
    },

    cashAccounts: (state) => {
      return state.accounts.filter(a => a.type === 'cash')
    },

    bankAccounts: (state) => {
      return state.accounts.filter(a => ['checking', 'savings'].includes(a.type))
    },

    creditCardAccounts: (state) => {
      return state.accounts.filter(a => a.type === 'credit_card')
    },

    getAccountById: (state) => (id) => {
      return state.accounts.find(account => account.id === id)
    }
  },

  actions: {
    async fetchAccounts() {
      try {
        const response = await apiService.get('accounts')
        this.accounts = response.data || []
      } catch (error) {
        console.error('Error fetching accounts:', error)
        throw error
      }
    },

    async createAccount(accountData) {
      try {
        const response = await apiService.post('accounts', accountData)
        this.accounts.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating account:', error)
        throw error
      }
    },

    async updateAccount(id, accountData) {
      try {
        const response = await apiService.put('accounts', id, accountData)
        const index = this.accounts.findIndex(a => a.id === id)
        if (index !== -1) {
          this.accounts[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating account:', error)
        throw error
      }
    },

    async deleteAccount(id) {
      try {
        await apiService.delete('accounts', id)
        this.accounts = this.accounts.filter(a => a.id !== id)
      } catch (error) {
        console.error('Error deleting account:', error)
        throw error
      }
    },

    async updateBalance(id, amount, operation = 'add') {
      try {
        const account = this.getAccountById(id)
        if (!account) throw new Error('Account not found')

        // Use backend balance update endpoint
        const response = await apiService.updateAccountBalance(id, {
          balance: amount,
          operation: operation === 'add' ? 'add' : 'subtract'
        })

        // Update local state
        const index = this.accounts.findIndex(a => a.id === id)
        if (index !== -1) {
          this.accounts[index] = response.data
        }
      } catch (error) {
        console.error('Error updating balance:', error)
        throw error
      }
    },

    async reconcileAccount(id, actualBalance) {
      try {
        const account = this.getAccountById(id)
        if (!account) throw new Error('Account not found')

        const difference = actualBalance - account.balance

        await this.updateAccount(id, {
          ...account,
          balance: actualBalance,
          lastReconciled: new Date().toISOString(),
          reconciliationDifference: difference
        })

        return difference
      } catch (error) {
        console.error('Error reconciling account:', error)
        throw error
      }
    },

    formatBalance(balance) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(balance)
    }
  }
})
