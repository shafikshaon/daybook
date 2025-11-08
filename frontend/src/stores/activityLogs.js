import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'

export const useActivityLogsStore = defineStore('activityLogs', {
  state: () => ({
    activityLogs: [],
    summary: null,
    pagination: {
      currentPage: 1,
      limit: 50,
      total: 0,
      totalPages: 0
    },
    filters: {
      module: null,
      action: null,
      startDate: null,
      endDate: null
    },
    loading: false
  }),

  getters: {
    allActivityLogs: (state) => state.activityLogs,

    logsByModule: (state) => (module) => {
      return state.activityLogs.filter(log => log.module === module)
    },

    logsByAction: (state) => (action) => {
      return state.activityLogs.filter(log => log.action === action)
    },

    recentLogs: (state) => (limit = 10) => {
      return state.activityLogs.slice(0, limit)
    },

    hasNextPage: (state) => {
      return state.pagination.currentPage < state.pagination.totalPages
    },

    hasPrevPage: (state) => {
      return state.pagination.currentPage > 1
    }
  },

  actions: {
    async fetchActivityLogs(filters = {}, page = 1) {
      this.loading = true
      try {
        const params = {
          page,
          limit: this.pagination.limit,
          ...this.filters,
          ...filters
        }

        // Remove null/undefined values
        Object.keys(params).forEach(key => {
          if (params[key] === null || params[key] === undefined || params[key] === '') {
            delete params[key]
          }
        })

        const response = await apiService.activityLogs.getAll(params)

        if (response.data) {
          this.activityLogs = response.data.data || []
          this.pagination = {
            currentPage: response.data.page || 1,
            limit: response.data.limit || 50,
            total: response.data.total || 0,
            totalPages: response.data.totalPages || 0
          }
        }
      } catch (error) {
        console.error('Error fetching activity logs:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchActivitySummary(days = 30) {
      this.loading = true
      try {
        const response = await apiService.activityLogs.getSummary(days)
        this.summary = response.data
        return response.data
      } catch (error) {
        console.error('Error fetching activity summary:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchActivityLogById(id) {
      try {
        const response = await apiService.activityLogs.getById(id)
        return response.data
      } catch (error) {
        console.error('Error fetching activity log:', error)
        throw error
      }
    },

    async cleanupOldLogs(days = 90) {
      try {
        const response = await apiService.activityLogs.cleanup(days)
        // Refresh the list after cleanup
        await this.fetchActivityLogs()
        return response.data
      } catch (error) {
        console.error('Error cleaning up activity logs:', error)
        throw error
      }
    },

    async backfillActivityLogs(options = {}) {
      try {
        const response = await apiService.activityLogs.backfill(options)
        // Refresh the list and summary after backfill
        await this.fetchActivityLogs()
        await this.fetchActivitySummary()
        return response.data
      } catch (error) {
        console.error('Error backfilling activity logs:', error)
        throw error
      }
    },

    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
    },

    clearFilters() {
      this.filters = {
        module: null,
        action: null,
        startDate: null,
        endDate: null
      }
    },

    async nextPage() {
      if (this.hasNextPage) {
        await this.fetchActivityLogs(this.filters, this.pagination.currentPage + 1)
      }
    },

    async prevPage() {
      if (this.hasPrevPage) {
        await this.fetchActivityLogs(this.filters, this.pagination.currentPage - 1)
      }
    },

    async goToPage(page) {
      if (page >= 1 && page <= this.pagination.totalPages) {
        await this.fetchActivityLogs(this.filters, page)
      }
    }
  }
})
