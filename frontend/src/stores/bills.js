import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useBillsStore = defineStore('bills', {
  state: () => ({
    bills: [],
    billPayments: []
  }),

  getters: {
    allBills: (state) => state.bills,

    activeBills: (state) => {
      return state.bills.filter(bill => bill.active)
    },

    upcomingBills: (state) => {
      const today = new Date()
      const nextMonth = new Date(today)
      nextMonth.setMonth(nextMonth.getMonth() + 1)

      return state.activeBills
        .filter(bill => {
          const dueDate = state.getNextDueDate(bill)
          return dueDate && dueDate >= today && dueDate <= nextMonth
        })
        .map(bill => ({
          ...bill,
          nextDueDate: state.getNextDueDate(bill),
          daysUntilDue: state.getDaysUntilDue(bill)
        }))
        .sort((a, b) => a.nextDueDate - b.nextDueDate)
    },

    overdueBills: (state) => {
      const today = new Date()

      return state.activeBills
        .filter(bill => {
          const dueDate = state.getNextDueDate(bill)
          return dueDate && dueDate < today && !state.isPaid(bill, dueDate)
        })
        .map(bill => ({
          ...bill,
          nextDueDate: state.getNextDueDate(bill),
          daysOverdue: Math.abs(state.getDaysUntilDue(bill))
        }))
    },

    totalMonthlyBills: (state) => {
      return state.activeBills
        .filter(bill => bill.frequency === 'monthly')
        .reduce((sum, bill) => sum + bill.amount, 0)
    },

    billsByCategory: (state) => {
      const categories = {}

      state.activeBills.forEach(bill => {
        if (!categories[bill.category]) {
          categories[bill.category] = {
            count: 0,
            totalAmount: 0,
            bills: []
          }
        }

        categories[bill.category].count++
        categories[bill.category].totalAmount += bill.amount
        categories[bill.category].bills.push(bill)
      })

      return categories
    },

    getBillById: (state) => (id) => {
      return state.bills.find(bill => bill.id === id)
    },

    getNextDueDate: () => (bill) => {
      const today = new Date()
      const lastPaid = bill.lastPaidDate ? new Date(bill.lastPaidDate) : null

      let nextDate = new Date()

      if (lastPaid) {
        nextDate = new Date(lastPaid)
      } else if (bill.startDate) {
        nextDate = new Date(bill.startDate)
      }

      while (nextDate <= today) {
        switch (bill.frequency) {
          case 'weekly':
            nextDate.setDate(nextDate.getDate() + 7)
            break
          case 'biweekly':
            nextDate.setDate(nextDate.getDate() + 14)
            break
          case 'monthly':
            nextDate.setMonth(nextDate.getMonth() + 1)
            break
          case 'quarterly':
            nextDate.setMonth(nextDate.getMonth() + 3)
            break
          case 'semi-annually':
            nextDate.setMonth(nextDate.getMonth() + 6)
            break
          case 'annually':
            nextDate.setFullYear(nextDate.getFullYear() + 1)
            break
          default:
            return null
        }
      }

      return nextDate
    },

    getDaysUntilDue: (state) => (bill) => {
      const dueDate = state.getNextDueDate(bill)
      if (!dueDate) return null

      const today = new Date()
      const difference = dueDate - today

      return Math.ceil(difference / (1000 * 60 * 60 * 24))
    },

    isPaid: (state) => (bill, dueDate) => {
      const dueDateStr = dueDate.toISOString().split('T')[0]
      return state.billPayments.some(payment =>
        payment.billId === bill.id &&
        payment.paymentDate.split('T')[0] === dueDateStr
      )
    },

    paymentHistory: (state) => (billId) => {
      return state.billPayments
        .filter(payment => payment.billId === billId)
        .sort((a, b) => new Date(b.paymentDate) - new Date(a.paymentDate))
    }
  },

  actions: {
    async fetchBills() {
      try {
        const response = await apiService.get('bills')
        this.bills = response.data || []
      } catch (error) {
        console.error('Error fetching bills:', error)
        throw error
      }
    },

    async createBill(billData) {
      try {
        const response = await apiService.post('bills', {
          ...billData,
          active: true
        })

        this.bills.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating bill:', error)
        throw error
      }
    },

    async updateBill(id, billData) {
      try {
        const response = await apiService.put('bills', id, billData)
        const index = this.bills.findIndex(bill => bill.id === id)
        if (index !== -1) {
          this.bills[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating bill:', error)
        throw error
      }
    },

    async deleteBill(id) {
      try {
        await apiService.delete('bills', id)
        this.bills = this.bills.filter(bill => bill.id !== id)
      } catch (error) {
        console.error('Error deleting bill:', error)
        throw error
      }
    },

    async markAsPaid(billId, paymentDate = null, amount = null, accountId = null) {
      try {
        const bill = this.getBillById(billId)
        if (!bill) throw new Error('Bill not found')

        const payment = {
          amount: amount || bill.amount,
          paymentDate: paymentDate || new Date().toISOString(),
          accountId,
          notes: 'Bill payment'
        }

        const response = await apiService.payBill(billId, payment)

        // Refresh bills and payments
        await this.fetchBills()
        await this.fetchBillPayments()

        return response.data
      } catch (error) {
        console.error('Error marking bill as paid:', error)
        throw error
      }
    },

    async fetchBillPayments() {
      try {
        const response = await apiService.get('bill-payments')
        this.billPayments = response.data || []
      } catch (error) {
        console.error('Error fetching bill payments:', error)
        throw error
      }
    },

    async toggleBillActive(billId) {
      const bill = this.getBillById(billId)
      if (!bill) throw new Error('Bill not found')

      await this.updateBill(billId, {
        ...bill,
        active: !bill.active
      })
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
