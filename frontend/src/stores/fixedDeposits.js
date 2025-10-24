import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useFixedDepositsStore = defineStore('fixedDeposits', {
  state: () => ({
    fixedDeposits: []
  }),

  getters: {
    allFixedDeposits: (state) => state.fixedDeposits,

    activeFixedDeposits: (state) => {
      const today = new Date()
      return state.fixedDeposits.filter(fd => new Date(fd.maturityDate) > today)
    },

    maturedFixedDeposits: (state) => {
      const today = new Date()
      return state.fixedDeposits.filter(fd => new Date(fd.maturityDate) <= today && !fd.withdrawn)
    },

    upcomingMaturities: (state) => {
      const today = new Date()
      const nextMonth = new Date(today)
      nextMonth.setMonth(nextMonth.getMonth() + 1)

      return state.fixedDeposits.filter(fd => {
        const maturityDate = new Date(fd.maturityDate)
        return maturityDate > today && maturityDate <= nextMonth && !fd.withdrawn
      }).sort((a, b) => new Date(a.maturityDate) - new Date(b.maturityDate))
    },

    totalInvested: (state) => {
      return state.fixedDeposits
        .filter(fd => !fd.withdrawn)
        .reduce((sum, fd) => sum + fd.principal, 0)
    },

    totalMaturityValue: (state) => {
      return state.fixedDeposits
        .filter(fd => !fd.withdrawn)
        .reduce((sum, fd) => {
          return sum + state.calculateMaturityAmount(fd.id)
        }, 0)
    },

    totalInterestEarned: (state) => {
      return state.totalMaturityValue - state.totalInvested
    },

    getFixedDepositById: (state) => (id) => {
      return state.fixedDeposits.find(fd => fd.id === id)
    }
  },

  actions: {
    async fetchFixedDeposits() {
      try {
        const response = await apiService.get('fixed-deposits')
        this.fixedDeposits = response.data || []
      } catch (error) {
        console.error('Error fetching fixed deposits:', error)
        throw error
      }
    },

    async createFixedDeposit(fdData) {
      try {
        const maturityAmount = this.calculateFutureValue(
          fdData.principal,
          fdData.interestRate,
          fdData.tenureMonths
        )

        const response = await apiService.post('fixed-deposits', {
          ...fdData,
          maturityAmount,
          withdrawn: false
        })

        this.fixedDeposits.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating fixed deposit:', error)
        throw error
      }
    },

    async updateFixedDeposit(id, fdData) {
      try {
        const response = await apiService.put('fixed-deposits', id, fdData)
        const index = this.fixedDeposits.findIndex(fd => fd.id === id)
        if (index !== -1) {
          this.fixedDeposits[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating fixed deposit:', error)
        throw error
      }
    },

    async deleteFixedDeposit(id) {
      try {
        await apiService.delete('fixed-deposits', id)
        this.fixedDeposits = this.fixedDeposits.filter(fd => fd.id !== id)
      } catch (error) {
        console.error('Error deleting fixed deposit:', error)
        throw error
      }
    },

    async withdrawFixedDeposit(id) {
      const fd = this.getFixedDepositById(id)
      if (!fd) throw new Error('Fixed deposit not found')

      const maturityAmount = this.calculateMaturityAmount(id)

      await this.updateFixedDeposit(id, {
        ...fd,
        withdrawn: true,
        withdrawnDate: new Date().toISOString(),
        actualMaturityAmount: maturityAmount
      })

      return maturityAmount
    },

    calculateMaturityAmount(fdId) {
      const fd = this.getFixedDepositById(fdId)
      if (!fd) return 0

      const today = new Date()
      const startDate = new Date(fd.startDate)
      const maturityDate = new Date(fd.maturityDate)

      const currentDate = fd.withdrawn ? new Date(fd.withdrawnDate) : today
      const monthsElapsed = this.getMonthsDifference(startDate, currentDate)

      if (fd.compounding === 'simple') {
        const interest = (fd.principal * fd.interestRate * monthsElapsed) / (100 * 12)
        return fd.principal + interest
      } else {
        const compoundingFrequency = this.getCompoundingFrequency(fd.compounding)
        const years = monthsElapsed / 12
        const amount = fd.principal * Math.pow(
          (1 + (fd.interestRate / 100) / compoundingFrequency),
          compoundingFrequency * years
        )
        return amount
      }
    },

    calculateFutureValue(principal, interestRate, tenureMonths, compounding = 'monthly') {
      if (compounding === 'simple') {
        const interest = (principal * interestRate * tenureMonths) / (100 * 12)
        return principal + interest
      } else {
        const compoundingFrequency = this.getCompoundingFrequency(compounding)
        const years = tenureMonths / 12
        const amount = principal * Math.pow(
          (1 + (interestRate / 100) / compoundingFrequency),
          compoundingFrequency * years
        )
        return amount
      }
    },

    getCompoundingFrequency(compounding) {
      const frequencies = {
        daily: 365,
        monthly: 12,
        quarterly: 4,
        'semi-annually': 2,
        annually: 1
      }
      return frequencies[compounding] || 12
    },

    getMonthsDifference(startDate, endDate) {
      const start = new Date(startDate)
      const end = new Date(endDate)

      const months = (end.getFullYear() - start.getFullYear()) * 12 +
        (end.getMonth() - start.getMonth())

      return Math.max(0, months)
    },

    getDaysUntilMaturity(fdId) {
      const fd = this.getFixedDepositById(fdId)
      if (!fd) return 0

      const today = new Date()
      const maturityDate = new Date(fd.maturityDate)
      const difference = maturityDate - today

      return Math.ceil(difference / (1000 * 60 * 60 * 24))
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
