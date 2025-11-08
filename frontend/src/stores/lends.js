import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useLendsStore = defineStore('lends', {
  state: () => ({
    lends: [],
    lendPayments: {},
    loading: false,
    error: null
  }),

  getters: {
    allLends: (state) => state.lends,

    activeLends: (state) => {
      return state.lends.filter(l => l.status === 'active')
    },

    partiallyReceivedLends: (state) => {
      return state.lends.filter(l => l.status === 'partially_received')
    },

    fullyReceivedLends: (state) => {
      return state.lends.filter(l => l.status === 'fully_received')
    },

    getLendById: (state) => (id) => {
      return state.lends.find(l => l.id === id)
    },

    totalLendAmount: (state) => {
      return state.lends
        .filter(l => l.status !== 'fully_received')
        .reduce((sum, l) => sum + l.remainingAmount, 0)
    },

    totalOriginalLend: (state) => {
      return state.lends
        .filter(l => l.status !== 'fully_received')
        .reduce((sum, l) => sum + l.originalAmount, 0)
    },

    totalReceivedAmount: (state) => {
      return state.lends.reduce((sum, l) => sum + (l.originalAmount - l.remainingAmount), 0)
    },

    overdueLends: (state) => {
      const today = new Date()
      return state.lends.filter(l => {
        if (!l.dueDate || l.status === 'fully_received') return false
        return new Date(l.dueDate) < today
      })
    },

    lendsByDebtor: (state) => (debtorName) => {
      return state.lends.filter(l =>
        l.debtorName.toLowerCase().includes(debtorName.toLowerCase())
      )
    },

    getPaymentsForLend: (state) => (lendId) => {
      return state.lendPayments[lendId] || []
    },

    lendProgress: (state) => (lendId) => {
      const lend = state.lends.find(l => l.id === lendId)
      if (!lend) return null

      const receivedAmount = lend.originalAmount - lend.remainingAmount
      const percentage = (receivedAmount / lend.originalAmount) * 100

      return {
        receivedAmount,
        remainingAmount: lend.remainingAmount,
        percentage,
        status: lend.status
      }
    }
  },

  actions: {
    async fetchLends() {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('lends')
        this.lends = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching lends:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchLend(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('lends', id)
        const index = this.lends.findIndex(l => l.id === id)
        if (index !== -1) {
          this.lends[index] = response.data
        } else {
          this.lends.push(response.data)
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching lend:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createLend(lendData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('lends', lendData)
        this.lends.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating lend:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateLend(id, lendData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('lends', id, lendData)
        const index = this.lends.findIndex(l => l.id === id)
        if (index !== -1) {
          this.lends[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error updating lend:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteLend(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('lends', id)
        this.lends = this.lends.filter(l => l.id !== id)
        delete this.lendPayments[id]
      } catch (error) {
        this.error = error.message
        console.error('Error deleting lend:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async recordPayment(lendId, paymentData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`lends/${lendId}/payments`, paymentData)

        // Update the lend record
        if (response.data.lend) {
          const index = this.lends.findIndex(l => l.id === lendId)
          if (index !== -1) {
            this.lends[index] = response.data.lend
          }
        }

        // Add payment to the local payments cache
        if (!this.lendPayments[lendId]) {
          this.lendPayments[lendId] = []
        }
        if (response.data.payment) {
          this.lendPayments[lendId].unshift(response.data.payment)
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

    async fetchPayments(lendId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`lends/${lendId}/payments`)
        this.lendPayments[lendId] = response.data || []
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
