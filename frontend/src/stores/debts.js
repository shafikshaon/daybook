import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useDebtsStore = defineStore('debts', {
  state: () => ({
    debts: [],
    debtPayments: {},
    loading: false,
    error: null
  }),

  getters: {
    allDebts: (state) => state.debts,

    activeDebts: (state) => {
      return state.debts.filter(d => d.status === 'active')
    },

    partiallyPaidDebts: (state) => {
      return state.debts.filter(d => d.status === 'partially_paid')
    },

    fullyPaidDebts: (state) => {
      return state.debts.filter(d => d.status === 'fully_paid')
    },

    getDebtById: (state) => (id) => {
      return state.debts.find(d => d.id === id)
    },

    totalDebtAmount: (state) => {
      return state.debts
        .filter(d => d.status !== 'fully_paid')
        .reduce((sum, d) => sum + d.remainingAmount, 0)
    },

    totalOriginalDebt: (state) => {
      return state.debts
        .filter(d => d.status !== 'fully_paid')
        .reduce((sum, d) => sum + d.originalAmount, 0)
    },

    totalPaidAmount: (state) => {
      return state.debts.reduce((sum, d) => sum + (d.originalAmount - d.remainingAmount), 0)
    },

    overdueDebts: (state) => {
      const today = new Date()
      return state.debts.filter(d => {
        if (!d.dueDate || d.status === 'fully_paid') return false
        return new Date(d.dueDate) < today
      })
    },

    debtsByCreditor: (state) => (creditorName) => {
      return state.debts.filter(d =>
        d.creditorName.toLowerCase().includes(creditorName.toLowerCase())
      )
    },

    getPaymentsForDebt: (state) => (debtId) => {
      return state.debtPayments[debtId] || []
    },

    debtProgress: (state) => (debtId) => {
      const debt = state.debts.find(d => d.id === debtId)
      if (!debt) return null

      const paidAmount = debt.originalAmount - debt.remainingAmount
      const percentage = (paidAmount / debt.originalAmount) * 100

      return {
        paidAmount,
        remainingAmount: debt.remainingAmount,
        percentage,
        status: debt.status
      }
    }
  },

  actions: {
    async fetchDebts() {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('debts')
        this.debts = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching debts:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchDebt(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('debts', id)
        const index = this.debts.findIndex(d => d.id === id)
        if (index !== -1) {
          this.debts[index] = response.data
        } else {
          this.debts.push(response.data)
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching debt:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createDebt(debtData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('debts', debtData)
        this.debts.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating debt:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateDebt(id, debtData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('debts', id, debtData)
        const index = this.debts.findIndex(d => d.id === id)
        if (index !== -1) {
          this.debts[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error updating debt:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteDebt(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('debts', id)
        this.debts = this.debts.filter(d => d.id !== id)
        delete this.debtPayments[id]
      } catch (error) {
        this.error = error.message
        console.error('Error deleting debt:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async recordPayment(debtId, paymentData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`debts/${debtId}/payments`, paymentData)

        // Update the debt record
        if (response.data.debt) {
          const index = this.debts.findIndex(d => d.id === debtId)
          if (index !== -1) {
            this.debts[index] = response.data.debt
          }
        }

        // Add payment to the local payments cache
        if (!this.debtPayments[debtId]) {
          this.debtPayments[debtId] = []
        }
        if (response.data.payment) {
          this.debtPayments[debtId].unshift(response.data.payment)
        }

        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error recording payment:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchPayments(debtId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`debts/${debtId}/payments`)
        this.debtPayments[debtId] = response.data || []
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching payments:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    },

    formatDate(date) {
      if (!date) return ''
      return new Date(date).toLocaleDateString()
    },

    isOverdue(dueDate) {
      if (!dueDate) return false
      return new Date(dueDate) < new Date()
    }
  }
})
