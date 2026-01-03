import { defineStore } from 'pinia'
import apiService, { api } from '@/services/api-backend'

export const useBackupsStore = defineStore('backups', {
  state: () => ({
    backups: [],
    loading: false
  }),

  getters: {
    allBackups: (state) => state.backups,

    completedBackups: (state) => {
      return state.backups.filter(backup => backup.status === 'completed')
    },

    pendingBackups: (state) => {
      return state.backups.filter(backup => backup.status === 'pending')
    },

    failedBackups: (state) => {
      return state.backups.filter(backup => backup.status === 'failed')
    },

    getBackupById: (state) => (id) => {
      return state.backups.find(backup => backup.id === id)
    }
  },

  actions: {
    async fetchBackups() {
      try {
        this.loading = true
        const response = await apiService.get('backups')
        this.backups = response.data || []
      } catch (error) {
        console.error('Error fetching backups:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createBackup() {
      try {
        this.loading = true
        const response = await apiService.post('backups')
        // Add the new backup to the list
        if (response.data) {
          this.backups.unshift(response.data)
        }
        // Refresh the list to get updated status
        setTimeout(() => this.fetchBackups(), 2000)
        return response.data
      } catch (error) {
        console.error('Error creating backup:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async downloadBackup(backupId) {
      try {
        // Use axios instance to get the proper base URL and auth headers
        const response = await api.get(`/backups/${backupId}/download`, {
          responseType: 'blob'
        })

        // Get filename from Content-Disposition header
        const contentDisposition = response.headers['content-disposition']
        let filename = 'backup.sql'
        if (contentDisposition) {
          const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition)
          if (matches != null && matches[1]) {
            filename = matches[1].replace(/['"]/g, '')
          }
        }

        // Create blob and download
        const blob = new Blob([response.data])
        const downloadUrl = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = downloadUrl
        link.download = filename
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(downloadUrl)
      } catch (error) {
        console.error('Error downloading backup:', error)
        throw error
      }
    },

    async deleteBackup(backupId) {
      try {
        await apiService.delete('backups', backupId)
        // Remove from local state
        this.backups = this.backups.filter(backup => backup.id !== backupId)
      } catch (error) {
        console.error('Error deleting backup:', error)
        throw error
      }
    },

    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
    }
  }
})
