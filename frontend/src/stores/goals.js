import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { toISOString } from '@/utils/dateUtils'
import { useSettingsStore } from './settings'

export const useGoalsStore = defineStore('goals', {
  state: () => ({
    goals: [],
    holdingTypes: null
  }),

  getters: {
    allGoals: (state) => state.goals,

    activeGoals: (state) => {
      return state.goals.filter(goal => goal.status === 'active')
    },

    achievedGoals: (state) => {
      return state.goals.filter(goal => goal.status === 'achieved')
    },

    pausedGoals: (state) => {
      return state.goals.filter(goal => goal.status === 'paused')
    },

    archivedGoals: (state) => {
      return state.goals.filter(goal => goal.status === 'archived')
    },

    goalsByCategory: (state) => (category) => {
      return state.goals.filter(goal => goal.category === category)
    },

    goalsByPriority: (state) => (priority) => {
      return state.goals.filter(goal => goal.priority === priority)
    },

    highPriorityGoals: (state) => {
      return state.activeGoals.filter(goal => goal.priority === 'high')
    },

    getGoalById: (state) => (id) => {
      return state.goals.find(goal => goal.id === id)
    },

    totalTargetAmount: (state) => {
      return state.activeGoals.reduce((sum, goal) => sum + goal.targetAmount, 0)
    },

    totalCurrentAmount: (state) => {
      return state.activeGoals.reduce((sum, goal) => sum + goal.currentAmount, 0)
    },

    totalProgress: (state) => {
      const total = state.totalTargetAmount
      if (total === 0) return 0
      return (state.totalCurrentAmount / total) * 100
    },

    goalProgress: (state) => (goalId) => {
      const goal = state.getGoalById(goalId)
      if (!goal) return null

      const percentage = goal.targetAmount > 0
        ? (goal.currentAmount / goal.targetAmount) * 100
        : 0
      const remaining = goal.targetAmount - goal.currentAmount

      return {
        percentage: Math.min(percentage, 100),
        remaining: Math.max(remaining, 0),
        achieved: percentage >= 100
      }
    },

    projectedCompletion: (state) => (goalId) => {
      const goal = state.getGoalById(goalId)
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

    goalHoldings: (state) => (goalId) => {
      const goal = state.getGoalById(goalId)
      return goal?.holdings || []
    },

    goalContributions: (state) => (goalId) => {
      const goal = state.getGoalById(goalId)
      return goal?.contributions || []
    },

    totalHoldingsByType: (state) => (goalId) => {
      const holdings = state.goalHoldings(goalId)
      const totals = {}

      holdings.forEach(holding => {
        if (!totals[holding.type]) {
          totals[holding.type] = 0
        }
        totals[holding.type] += holding.currentValue || holding.amount
      })

      return totals
    },

    holdingGainLoss: (state) => (holdingId) => {
      for (const goal of state.goals) {
        const holding = goal.holdings?.find(h => h.id === holdingId)
        if (holding) {
          const gain = (holding.currentValue || holding.amount) - holding.amount
          const gainPercent = holding.amount > 0
            ? (gain / holding.amount) * 100
            : 0
          return { gain, gainPercent }
        }
      }
      return { gain: 0, gainPercent: 0 }
    }
  },

  actions: {
    async fetchGoals(filters = {}) {
      try {
        const params = new URLSearchParams()
        if (filters.status) params.append('status', filters.status)
        if (filters.category) params.append('category', filters.category)
        if (filters.priority) params.append('priority', filters.priority)

        const url = params.toString() ? `goals?${params.toString()}` : 'goals'
        const response = await apiService.get(url)
        this.goals = response.data || []
      } catch (error) {
        console.error('Error fetching goals:', error)
        throw error
      }
    },

    async fetchGoal(id) {
      try {
        const response = await apiService.get(`goals/${id}`)
        const goal = response.data

        // Update in local state if exists, otherwise add
        const index = this.goals.findIndex(g => g.id === id)
        if (index !== -1) {
          this.goals[index] = goal
        } else {
          this.goals.push(goal)
        }

        return goal
      } catch (error) {
        console.error('Error fetching goal:', error)
        throw error
      }
    },

    async createGoal(goalData) {
      try {
        const response = await apiService.post('goals', {
          ...goalData,
          currentAmount: goalData.currentAmount || 0,
          targetDate: goalData.targetDate ? toISOString(goalData.targetDate) : null,
          status: goalData.status || 'active',
          achieved: false
        })

        this.goals.push(response.data)
        return response.data
      } catch (error) {
        console.error('Error creating goal:', error)
        throw error
      }
    },

    async updateGoal(id, goalData) {
      try {
        const dataToSend = { ...goalData }
        if (dataToSend.targetDate) {
          dataToSend.targetDate = toISOString(dataToSend.targetDate)
        }
        if (dataToSend.achievedDate) {
          dataToSend.achievedDate = toISOString(dataToSend.achievedDate)
        }

        const response = await apiService.put('goals', id, dataToSend)
        const index = this.goals.findIndex(goal => goal.id === id)
        if (index !== -1) {
          this.goals[index] = response.data
        }
        return response.data
      } catch (error) {
        console.error('Error updating goal:', error)
        throw error
      }
    },

    async deleteGoal(id) {
      try {
        await apiService.delete('goals', id)
        this.goals = this.goals.filter(goal => goal.id !== id)
      } catch (error) {
        console.error('Error deleting goal:', error)
        throw error
      }
    },

    async addHolding(goalId, holdingData) {
      try {
        const dataToSend = {
          ...holdingData,
          purchaseDate: toISOString(holdingData.purchaseDate),
          maturityDate: holdingData.maturityDate ? toISOString(holdingData.maturityDate) : null
        }

        const response = await apiService.post(`goals/${goalId}/holdings`, dataToSend)

        // Refresh the goal to get updated holdings
        await this.fetchGoal(goalId)

        return response.data
      } catch (error) {
        console.error('Error adding holding:', error)
        throw error
      }
    },

    async updateHolding(holdingId, holdingData) {
      try {
        const dataToSend = { ...holdingData }
        if (dataToSend.purchaseDate) {
          dataToSend.purchaseDate = toISOString(dataToSend.purchaseDate)
        }
        if (dataToSend.maturityDate) {
          dataToSend.maturityDate = toISOString(dataToSend.maturityDate)
        }

        const response = await apiService.put(`goals/holdings/${holdingId}`, dataToSend)

        // Find the goal containing this holding and refresh it
        const goal = this.goals.find(g =>
          g.holdings?.some(h => h.id === holdingId)
        )
        if (goal) {
          await this.fetchGoal(goal.id)
        }

        return response.data
      } catch (error) {
        console.error('Error updating holding:', error)
        throw error
      }
    },

    async removeHolding(holdingId, accountId) {
      try {
        const response = await apiService.delete(`goals/holdings/${holdingId}`, {
          params: { accountId }
        })

        // Find the goal containing this holding and refresh it
        const goal = this.goals.find(g =>
          g.holdings?.some(h => h.id === holdingId)
        )
        if (goal) {
          await this.fetchGoal(goal.id)
        }

        return response.data
      } catch (error) {
        console.error('Error removing holding:', error)
        throw error
      }
    },

    async fetchHoldingTypes() {
      try {
        if (this.holdingTypes) {
          return this.holdingTypes
        }

        const response = await apiService.get('goals/holding-types')
        this.holdingTypes = response.data
        return this.holdingTypes
      } catch (error) {
        console.error('Error fetching holding types:', error)
        throw error
      }
    },

    async pauseGoal(goalId) {
      const goal = this.getGoalById(goalId)
      if (!goal) throw new Error('Goal not found')

      await this.updateGoal(goalId, {
        ...goal,
        status: 'paused'
      })
    },

    async resumeGoal(goalId) {
      const goal = this.getGoalById(goalId)
      if (!goal) throw new Error('Goal not found')

      await this.updateGoal(goalId, {
        ...goal,
        status: 'active'
      })
    },

    async archiveGoal(goalId) {
      const goal = this.getGoalById(goalId)
      if (!goal) throw new Error('Goal not found')

      await this.updateGoal(goalId, {
        ...goal,
        status: 'archived'
      })
    },

    calculateRequiredMonthlyContribution(goalId) {
      const goal = this.getGoalById(goalId)
      if (!goal || !goal.targetDate) return 0

      const remaining = goal.targetAmount - goal.currentAmount
      if (remaining <= 0) return 0

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

    formatAmount(amount) {
      const settingsStore = useSettingsStore()
      return settingsStore.formatCurrency(amount)
    },

    // Helper to get holding type label
    getHoldingTypeLabel(type) {
      if (!this.holdingTypes) return type

      for (const category of Object.values(this.holdingTypes)) {
        const found = category.find(t => t.value === type)
        if (found) return found.label
      }

      return type
    }
  }
})
