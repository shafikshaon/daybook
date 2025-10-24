import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useInvestmentsStore = defineStore('investments', {
  state: () => ({
    investments: [],
    portfolios: [],
    dividends: [],
    assetTypes: [
      { value: 'stocks', label: 'Stocks', icon: '📈' },
      { value: 'bonds', label: 'Bonds', icon: '📊' },
      { value: 'mutual_funds', label: 'Mutual Funds', icon: '💼' },
      { value: 'etf', label: 'ETFs', icon: '📉' },
      { value: 'crypto', label: 'Cryptocurrency', icon: '₿' },
      { value: 'real_estate', label: 'Real Estate', icon: '🏢' },
      { value: 'commodities', label: 'Commodities', icon: '🥇' },
      { value: 'other', label: 'Other', icon: '💰' }
    ]
  }),

  getters: {
    allInvestments: (state) => state.investments,

    investmentsByType: (state) => (type) => {
      return state.investments.filter(inv => inv.assetType === type)
    },

    investmentsByPortfolio: (state) => (portfolioId) => {
      return state.investments.filter(inv => inv.portfolioId === portfolioId)
    },

    totalInvested: (state) => {
      return state.investments.reduce((sum, inv) => {
        return sum + (inv.quantity * inv.costBasis)
      }, 0)
    },

    totalCurrentValue: (state) => {
      return state.investments.reduce((sum, inv) => {
        return sum + (inv.quantity * inv.currentPrice)
      }, 0)
    },

    totalGainLoss: (state) => {
      return state.totalCurrentValue - state.totalInvested
    },

    totalGainLossPercentage: (state) => {
      if (state.totalInvested === 0) return 0
      return ((state.totalCurrentValue - state.totalInvested) / state.totalInvested) * 100
    },

    assetAllocation: (state) => {
      const allocation = {}
      const total = state.totalCurrentValue

      state.assetTypes.forEach(type => {
        const typeInvestments = state.investments.filter(inv => inv.assetType === type.value)
        const typeValue = typeInvestments.reduce((sum, inv) => {
          return sum + (inv.quantity * inv.currentPrice)
        }, 0)

        if (typeValue > 0) {
          allocation[type.value] = {
            label: type.label,
            value: typeValue,
            percentage: total > 0 ? (typeValue / total) * 100 : 0,
            color: state.getColorForAssetType(type.value)
          }
        }
      })

      return allocation
    },

    getColorForAssetType: () => (assetType) => {
      const colors = {
        stocks: '#3b82f6',
        bonds: '#10b981',
        mutual_funds: '#6f42c1',
        etf: '#8b5cf6',
        crypto: '#f59e0b',
        real_estate: '#ef4444',
        commodities: '#ec4899',
        other: '#6b7280'
      }
      return colors[assetType] || '#6b7280'
    },

    portfolioPerformance: (state) => (portfolioId) => {
      const investments = state.investmentsByPortfolio(portfolioId)

      const invested = investments.reduce((sum, inv) => {
        return sum + (inv.quantity * inv.costBasis)
      }, 0)

      const currentValue = investments.reduce((sum, inv) => {
        return sum + (inv.quantity * inv.currentPrice)
      }, 0)

      const gainLoss = currentValue - invested
      const gainLossPercentage = invested > 0 ? (gainLoss / invested) * 100 : 0

      return {
        invested,
        currentValue,
        gainLoss,
        gainLossPercentage
      }
    },

    topPerformers: (state) => (limit = 5) => {
      return [...state.investments]
        .map(inv => ({
          ...inv,
          gainLoss: (inv.currentPrice - inv.costBasis) * inv.quantity,
          gainLossPercentage: ((inv.currentPrice - inv.costBasis) / inv.costBasis) * 100
        }))
        .sort((a, b) => b.gainLossPercentage - a.gainLossPercentage)
        .slice(0, limit)
    },

    bottomPerformers: (state) => (limit = 5) => {
      return [...state.investments]
        .map(inv => ({
          ...inv,
          gainLoss: (inv.currentPrice - inv.costBasis) * inv.quantity,
          gainLossPercentage: ((inv.currentPrice - inv.costBasis) / inv.costBasis) * 100
        }))
        .sort((a, b) => a.gainLossPercentage - b.gainLossPercentage)
        .slice(0, limit)
    },

    totalDividendsEarned: (state) => {
      return state.dividends.reduce((sum, div) => sum + div.amount, 0)
    },

    dividendsByInvestment: (state) => (investmentId) => {
      return state.dividends.filter(div => div.investmentId === investmentId)
    }
  },

  actions: {
    async fetchInvestments() {
      try {
        const response = await apiService.get('investments')
        this.investments = response.data || []
      } catch (error) {
        console.error('Error fetching investments:', error)
        throw error
      }
    },

    async createInvestment(investmentData) {
      try {
        const response = await apiService.post('investments', investmentData)
        this.investments.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating investment:', error)
        throw error
      }
    },

    async updateInvestment(id, investmentData) {
      try {
        const response = await apiService.put('investments', id, investmentData)
        const index = this.investments.findIndex(inv => inv.id === id)
        if (index !== -1) {
          this.investments[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating investment:', error)
        throw error
      }
    },

    async deleteInvestment(id) {
      try {
        await apiService.delete('investments', id)
        this.investments = this.investments.filter(inv => inv.id !== id)
      } catch (error) {
        console.error('Error deleting investment:', error)
        throw error
      }
    },

    async updatePrice(id, newPrice) {
      const investment = this.investments.find(inv => inv.id === id)
      if (!investment) throw new Error('Investment not found')

      await this.updateInvestment(id, {
        ...investment,
        currentPrice: newPrice,
        lastUpdated: new Date().toISOString()
      })
    },

    async buyShares(id, quantity, price) {
      try {
        const response = await apiService.buyShares(id, { quantity, price })

        // Refresh investments
        await this.fetchInvestments()

        return response.data
      } catch (error) {
        console.error('Error buying shares:', error)
        throw error
      }
    },

    async sellShares(id, quantity, price) {
      try {
        const response = await apiService.sellShares(id, { quantity, price })

        // Refresh investments
        await this.fetchInvestments()

        return response.data
      } catch (error) {
        console.error('Error selling shares:', error)
        throw error
      }
    },

    // Portfolios
    async fetchPortfolios() {
      try {
        const response = await apiService.get('portfolios')
        this.portfolios = response.data || []
      } catch (error) {
        console.error('Error fetching portfolios:', error)
        throw error
      }
    },

    async createPortfolio(portfolioData) {
      try {
        const response = await apiService.post('portfolios', portfolioData)
        this.portfolios.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating portfolio:', error)
        throw error
      }
    },

    // Dividends
    async fetchDividends() {
      try {
        const response = await apiService.get('dividends')
        this.dividends = response.data || []
      } catch (error) {
        console.error('Error fetching dividends:', error)
        throw error
      }
    },

    async recordDividend(dividendData) {
      try {
        const response = await apiService.post('dividends', dividendData)
        this.dividends.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error recording dividend:', error)
        throw error
      }
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
