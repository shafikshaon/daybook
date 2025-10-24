import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'
import { useTransactionsStore } from './transactions'

export const useCreditCardsStore = defineStore('creditCards', {
  state: () => ({
    creditCards: [],
    statements: [],
    rewards: []
  }),

  getters: {
    allCreditCards: (state) => state.creditCards,

    activeCreditCards: (state) => {
      return state.creditCards.filter(cc => cc.active)
    },

    totalCreditLimit: (state) => {
      return state.creditCards.reduce((sum, cc) => sum + cc.creditLimit, 0)
    },

    totalOutstanding: (state) => {
      return state.creditCards.reduce((sum, cc) => sum + cc.currentBalance, 0)
    },

    totalAvailableCredit: (state) => {
      return state.creditCards.reduce((sum, cc) => sum + (cc.creditLimit - cc.currentBalance), 0)
    },

    creditUtilization: (state) => {
      const total = state.totalCreditLimit
      return total > 0 ? (state.totalOutstanding / total) * 100 : 0
    },

    getCreditCardById: (state) => (id) => {
      return state.creditCards.find(cc => cc.id === id)
    },

    upcomingPayments: (state) => {
      const upcoming = []
      const today = new Date()
      const nextMonth = new Date(today)
      nextMonth.setMonth(nextMonth.getMonth() + 1)

      state.creditCards.forEach(cc => {
        if (cc.dueDate) {
          const dueDate = new Date(cc.dueDate)
          if (dueDate >= today && dueDate <= nextMonth) {
            upcoming.push({
              cardId: cc.id,
              cardName: cc.name,
              dueDate: cc.dueDate,
              minimumPayment: cc.minimumPayment,
              currentBalance: cc.currentBalance,
              daysUntilDue: Math.ceil((dueDate - today) / (1000 * 60 * 60 * 24))
            })
          }
        }
      })

      return upcoming.sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate))
    },

    statementsByCard: (state) => (cardId) => {
      return state.statements
        .filter(s => s.cardId === cardId)
        .sort((a, b) => new Date(b.statementDate) - new Date(a.statementDate))
    },

    rewardsByCard: (state) => (cardId) => {
      return state.rewards.filter(r => r.cardId === cardId)
    },

    totalRewardsEarned: (state) => (cardId = null) => {
      const relevantRewards = cardId
        ? state.rewards.filter(r => r.cardId === cardId)
        : state.rewards

      return relevantRewards.reduce((sum, r) => sum + (r.amount || 0), 0)
    }
  },

  actions: {
    async fetchCreditCards() {
      try {
        const response = await apiService.get('credit-cards')
        this.creditCards = response.data || []
      } catch (error) {
        console.error('Error fetching credit cards:', error)
        throw error
      }
    },

    async createCreditCard(cardData) {
      try {
        const response = await apiService.post('credit-cards', {
          ...cardData,
          active: true,
          currentBalance: cardData.currentBalance || 0
        })
        this.creditCards.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating credit card:', error)
        throw error
      }
    },

    async updateCreditCard(id, cardData) {
      try {
        const response = await apiService.put('credit-cards', id, cardData)
        const index = this.creditCards.findIndex(cc => cc.id === id)
        if (index !== -1) {
          this.creditCards[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating credit card:', error)
        throw error
      }
    },

    async deleteCreditCard(id) {
      try {
        await apiService.delete('credit-cards', id)
        this.creditCards = this.creditCards.filter(cc => cc.id !== id)
      } catch (error) {
        console.error('Error deleting credit card:', error)
        throw error
      }
    },

    calculateInterest(cardId) {
      const card = this.getCreditCardById(cardId)
      if (!card || !card.apr) return 0

      const monthlyRate = card.apr / 100 / 12
      const monthlyInterest = card.currentBalance * monthlyRate

      return monthlyInterest
    },

    calculateMinimumPayment(cardId) {
      const card = this.getCreditCardById(cardId)
      if (!card) return 0

      const interest = this.calculateInterest(cardId)
      const percentageBased = card.currentBalance * 0.02
      const minimumFloor = 25

      return Math.max(interest + percentageBased, minimumFloor)
    },

    async recordPayment(cardId, amount, paymentDate, paymentAccountId) {
      try {
        const card = this.getCreditCardById(cardId)
        if (!card) throw new Error('Card not found')

        const newBalance = card.currentBalance - amount

        // Update credit card balance
        await this.updateCreditCard(cardId, {
          ...card,
          currentBalance: newBalance,
          lastPaymentDate: paymentDate,
          lastPaymentAmount: amount
        })

        // Create transaction record if payment account is specified
        if (paymentAccountId) {
          const transactionsStore = useTransactionsStore()
          await transactionsStore.createTransaction({
            type: 'expense',
            amount,
            categoryId: 'credit_card_payment',
            accountId: paymentAccountId,
            date: paymentDate,
            description: `Credit card payment - ${card.name}`,
            creditCardId: cardId,
            tags: ['credit_card_payment'],
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          })
        }

        return newBalance
      } catch (error) {
        console.error('Error recording payment:', error)
        throw error
      }
    },

    async addCharge(cardId, amount, description) {
      try {
        const card = this.getCreditCardById(cardId)
        if (!card) throw new Error('Card not found')

        const newBalance = card.currentBalance + amount

        if (newBalance > card.creditLimit) {
          throw new Error('Transaction exceeds credit limit')
        }

        await this.updateCreditCard(cardId, {
          ...card,
          currentBalance: newBalance
        })

        return newBalance
      } catch (error) {
        console.error('Error adding charge:', error)
        throw error
      }
    },

    // Statements
    async fetchStatements(cardId = null) {
      try {
        const endpoint = cardId ? `statements?cardId=${cardId}` : 'statements'
        const response = await apiService.get('statements')
        this.statements = response.data || []
      } catch (error) {
        console.error('Error fetching statements:', error)
        throw error
      }
    },

    async createStatement(statementData) {
      try {
        const response = await apiService.post('statements', statementData)
        this.statements.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating statement:', error)
        throw error
      }
    },

    // Rewards
    async fetchRewards(cardId = null) {
      try {
        const response = await apiService.get('rewards')
        this.rewards = response.data || []
      } catch (error) {
        console.error('Error fetching rewards:', error)
        throw error
      }
    },

    async recordReward(rewardData) {
      try {
        const response = await apiService.post('rewards', rewardData)
        this.rewards.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error recording reward:', error)
        throw error
      }
    },

    async redeemReward(rewardId) {
      try {
        const reward = this.rewards.find(r => r.id === rewardId)
        if (!reward) throw new Error('Reward not found')

        const response = await apiService.put('rewards', rewardId, {
          ...reward,
          redeemed: true,
          redeemedAt: new Date().toISOString()
        })

        const index = this.rewards.findIndex(r => r.id === rewardId)
        if (index !== -1) {
          this.rewards[index] = response.data
        }

        return response.data
      } catch (error) {
        console.error('Error redeeming reward:', error)
        throw error
      }
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
