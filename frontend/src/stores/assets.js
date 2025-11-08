import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useAssetsStore = defineStore('assets', {
  state: () => ({
    assets: [],
    serviceRecords: {},
    attachments: {},
    stats: null,
    loading: false,
    error: null
  }),

  getters: {
    allAssets: (state) => state.assets,

    activeAssets: (state) => {
      return state.assets.filter(g => g.status === 'active')
    },

    archivedAssets: (state) => {
      return state.assets.filter(g => g.status === 'archived')
    },

    soldAssets: (state) => {
      return state.assets.filter(g => g.status === 'sold')
    },

    disposedAssets: (state) => {
      return state.assets.filter(g => g.status === 'disposed')
    },

    getAssetById: (state) => (id) => {
      return state.assets.find(g => g.id === id)
    },

    assetsByCategory: (state) => (category) => {
      return state.assets.filter(g =>
        g.category && g.category.toLowerCase().includes(category.toLowerCase())
      )
    },

    assetsUnderWarranty: (state) => {
      return state.assets.filter(g => g.warrantyStatus === 'active')
    },

    assetsWarrantyExpired: (state) => {
      return state.assets.filter(g => g.warrantyStatus === 'expired')
    },

    assetsWarrantyExpiringSoon: (state) => {
      // Warranty expiring in next 30 days
      return state.assets.filter(g => {
        return g.warrantyStatus === 'active' &&
               g.warrantyDaysRemaining !== null &&
               g.warrantyDaysRemaining <= 30
      })
    },

    totalValue: (state) => {
      return state.assets
        .filter(g => g.status === 'active')
        .reduce((sum, g) => sum + g.purchasePrice, 0)
    },

    totalServiceCost: (state) => {
      return state.assets.reduce((sum, g) => sum + (g.totalServiceCost || 0), 0)
    },

    totalCost: (state) => {
      return state.assets
        .filter(g => g.status === 'active')
        .reduce((sum, g) => sum + (g.totalCost || g.purchasePrice), 0)
    },

    getServiceRecordsForAsset: (state) => (assetId) => {
      return state.serviceRecords[assetId] || []
    },

    getAttachmentsForAsset: (state) => (assetId) => {
      return state.attachments[assetId] || []
    },

    categories: (state) => {
      const cats = new Set()
      state.assets.forEach(g => {
        if (g.category) cats.add(g.category)
      })
      return Array.from(cats).sort()
    },

    brands: (state) => {
      const brands = new Set()
      state.assets.forEach(g => {
        if (g.brand) brands.add(g.brand)
      })
      return Array.from(brands).sort()
    }
  },

  actions: {
    async fetchAssets(filters = {}) {
      this.loading = true
      this.error = null
      try {
        const params = {}
        if (filters.status) params.status = filters.status
        if (filters.category) params.category = filters.category

        const response = await apiService.query('assets', params)
        this.assets = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching assets:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchAsset(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('assets', id)
        const index = this.assets.findIndex(g => g.id === id)
        if (index !== -1) {
          this.assets[index] = response.data
        } else {
          this.assets.push(response.data)
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching good:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createAsset(assetData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('assets', assetData)
        this.assets.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating good:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateAsset(id, assetData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('assets', id, assetData)
        const index = this.assets.findIndex(g => g.id === id)
        if (index !== -1) {
          this.assets[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error updating good:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAsset(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('assets', id)
        this.assets = this.assets.filter(g => g.id !== id)
        delete this.serviceRecords[id]
        delete this.attachments[id]
      } catch (error) {
        this.error = error.message
        console.error('Error deleting good:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createServiceRecord(assetId, serviceData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`assets/${assetId}/services`, serviceData)

        // Add to local cache
        if (!this.serviceRecords[assetId]) {
          this.serviceRecords[assetId] = []
        }
        this.serviceRecords[assetId].unshift(response.data)

        // Refresh the good to update stats
        await this.fetchAsset(assetId)

        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating service record:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchServiceRecords(assetId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`assets/${assetId}/services`)
        this.serviceRecords[assetId] = response.data || []
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching service records:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteServiceRecord(assetId, serviceId) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete(`assets/${assetId}/services/${serviceId}`)

        // Remove from local cache
        if (this.serviceRecords[assetId]) {
          this.serviceRecords[assetId] = this.serviceRecords[assetId].filter(
            s => s.id !== serviceId
          )
        }

        // Refresh the good to update stats
        await this.fetchAsset(assetId)
      } catch (error) {
        this.error = error.message
        console.error('Error deleting service record:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async addAttachment(assetId, attachmentData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`assets/${assetId}/attachments`, attachmentData)

        // Add to local cache
        if (!this.attachments[assetId]) {
          this.attachments[assetId] = []
        }
        this.attachments[assetId].unshift(response.data)

        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error adding attachment:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchAttachments(assetId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`assets/${assetId}/attachments`)
        this.attachments[assetId] = response.data || []
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching attachments:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAttachment(assetId, attachmentId) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete(`assets/${assetId}/attachments/${attachmentId}`)

        // Remove from local cache
        if (this.attachments[assetId]) {
          this.attachments[assetId] = this.attachments[assetId].filter(
            a => a.id !== attachmentId
          )
        }
      } catch (error) {
        this.error = error.message
        console.error('Error deleting attachment:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchStats() {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('assets/stats')
        this.stats = response.data
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching stats:', error)
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

    isWarrantyExpired(warrantyEndDate) {
      if (!warrantyEndDate) return false
      return new Date(warrantyEndDate) < new Date()
    },

    isWarrantyExpiringSoon(warrantyDaysRemaining) {
      return warrantyDaysRemaining !== null && warrantyDaysRemaining <= 30 && warrantyDaysRemaining > 0
    }
  }
})
