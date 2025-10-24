import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useSavingsGoalsStore = defineStore('savingsGoals', {
  state: () => ({
    savingsGoals: [],
    automatedRules: []
  }),

  getters: {
    allSavingsGoals: (state) => state.savingsGoals,

    activeSavingsGoals: (state) => {
      return state.savingsGoals.filter(goal => !goal.achieved && !goal.archived)
    },

    achievedSavingsGoals: (state) => {
      return state.savingsGoals.filter(goal => goal.achieved)
    },

    totalTargetAmount: (state) => {
      return state.activeSavingsGoals.reduce((sum, goal) => sum + goal.targetAmount, 0)
    },

    totalSavedAmount: (state) => {
      return state.activeSavingsGoals.reduce((sum, goal) => sum + goal.currentAmount, 0)
    },

    totalProgress: (state) => {
      if (state.totalTargetAmount === 0) return 0
      return (state.totalSavedAmount / state.totalTargetAmount) * 100
    },

    getSavingsGoalById: (state) => (id) => {
      return state.savingsGoals.find(goal => goal.id === id)
    },

    goalProgress: (state) => (goalId) => {
      const goal = state.getSavingsGoalById(goalId)
      if (!goal) return null

      const percentage = (goal.currentAmount / goal.targetAmount) * 100
      const remaining = goal.targetAmount - goal.currentAmount

      return {
        percentage: Math.min(percentage, 100),
        remaining,
        achieved: percentage >= 100
      }
    },

    projectedCompletion: (state) => (goalId) => {
      const goal = state.getSavingsGoalById(goalId)
      if (!goal || goal.currentAmount >= goal.targetAmount) return null

      const remaining = goal.targetAmount - goal.currentAmount

      if (goal.monthlyContribution && goal.monthlyContribution > 0) {
        const monthsNeeded = Math.ceil(remaining / goal.monthlyContribution)
        const completionDate = new Date()
        completionDate.setMonth(completionDate.getMonth() + monthsNeeded)

        return {
          monthsNeeded,
          completionDate: completionDate.toISOString(),
          onTrack: goal.targetDate ? new Date(goal.targetDate) >= completionDate : true
        }
      }

      return null
    },

    activeAutomatedRules: (state) => {
      return state.automatedRules.filter(rule => rule.enabled)
    }
  },

  actions: {
    async fetchSavingsGoals() {
      try {
        const response = await apiService.get('savings-goals')
        this.savingsGoals = response.data || []
      } catch (error) {
        console.error('Error fetching savings goals:', error)
        throw error
      }
    },

    async createSavingsGoal(goalData) {
      try {
        const response = await apiService.post('savings-goals', {
          ...goalData,
          currentAmount: goalData.currentAmount || 0,
          achieved: false,
          archived: false
        })

        this.savingsGoals.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating savings goal:', error)
        throw error
      }
    },

    async updateSavingsGoal(id, goalData) {
      try {
        const response = await apiService.put('savings-goals', id, goalData)
        const index = this.savingsGoals.findIndex(goal => goal.id === id)
        if (index !== -1) {
          this.savingsGoals[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating savings goal:', error)
        throw error
      }
    },

    async deleteSavingsGoal(id) {
      try {
        await apiService.delete('savings-goals', id)
        this.savingsGoals = this.savingsGoals.filter(goal => goal.id !== id)
      } catch (error) {
        console.error('Error deleting savings goal:', error)
        throw error
      }
    },

    async addContribution(goalId, amount, date = null) {
      try {
        const contribution = {
          amount,
          date: date || new Date().toISOString(),
          notes: 'Contribution to savings goal'
        }

        const response = await apiService.contributeToGoal(goalId, contribution)

        // Refresh savings goals
        await this.fetchSavingsGoals()

        return response.data
      } catch (error) {
        console.error('Error adding contribution:', error)
        throw error
      }
    },

    async withdrawFromGoal(goalId, amount) {
      try {
        const withdrawal = {
          amount,
          notes: 'Withdrawal from savings goal'
        }

        const response = await apiService.withdrawFromGoal(goalId, withdrawal)

        // Refresh savings goals
        await this.fetchSavingsGoals()

        return response.data
      } catch (error) {
        console.error('Error withdrawing from goal:', error)
        throw error
      }
    },

    async archiveGoal(goalId) {
      const goal = this.getSavingsGoalById(goalId)
      if (!goal) throw new Error('Savings goal not found')

      await this.updateSavingsGoal(goalId, {
        ...goal,
        archived: true,
        archivedDate: new Date().toISOString()
      })
    },

    calculateRequiredMonthlyContribution(goalId) {
      const goal = this.getSavingsGoalById(goalId)
      if (!goal || !goal.targetDate) return 0

      const remaining = goal.targetAmount - goal.currentAmount
      const today = new Date()
      const targetDate = new Date(goal.targetDate)

      const monthsRemaining = Math.max(1, this.getMonthsDifference(today, targetDate))

      return remaining / monthsRemaining
    },

    getMonthsDifference(startDate, endDate) {
      const start = new Date(startDate)
      const end = new Date(endDate)

      const months = (end.getFullYear() - start.getFullYear()) * 12 +
        (end.getMonth() - start.getMonth())

      return Math.max(0, months)
    },

    // Automated savings rules
    async fetchAutomatedRules() {
      try {
        const response = await apiService.get('automated-rules')
        this.automatedRules = response.data || []
      } catch (error) {
        console.error('Error fetching automated rules:', error)
        throw error
      }
    },

    async createAutomatedRule(ruleData) {
      try {
        const response = await apiService.post('automated-rules', {
          ...ruleData,
          enabled: true
        })

        this.automatedRules.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating automated rule:', error)
        throw error
      }
    },

    async toggleAutomatedRule(ruleId) {
      const rule = this.automatedRules.find(r => r.id === ruleId)
      if (!rule) throw new Error('Rule not found')

      const response = await apiService.put('automated-rules', ruleId, {
        ...rule,
        enabled: !rule.enabled
      })

      const index = this.automatedRules.findIndex(r => r.id === ruleId)
      if (index !== -1) {
        this.automatedRules[index] = response.data
      }
    },

    async deleteAutomatedRule(ruleId) {
      try {
        await apiService.delete('automated-rules', ruleId)
        this.automatedRules = this.automatedRules.filter(r => r.id !== ruleId)
      } catch (error) {
        console.error('Error deleting automated rule:', error)
        throw error
      }
    },

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    }
  }
})
