import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'
import { useSettingsStore } from './settings'

export const useGoodsStore = defineStore('goods', {
  state: () => ({
    goods: [],
    serviceRecords: {},
    attachments: {},
    stats: null,
    loading: false,
    error: null
  }),

  getters: {
    allGoods: (state) => state.goods,

    activeGoods: (state) => {
      return state.goods.filter(g => g.status === 'active')
    },

    archivedGoods: (state) => {
      return state.goods.filter(g => g.status === 'archived')
    },

    soldGoods: (state) => {
      return state.goods.filter(g => g.status === 'sold')
    },

    disposedGoods: (state) => {
      return state.goods.filter(g => g.status === 'disposed')
    },

    getGoodById: (state) => (id) => {
      return state.goods.find(g => g.id === id)
    },

    goodsByCategory: (state) => (category) => {
      return state.goods.filter(g =>
        g.category && g.category.toLowerCase().includes(category.toLowerCase())
      )
    },

    goodsUnderWarranty: (state) => {
      return state.goods.filter(g => g.warrantyStatus === 'active')
    },

    goodsWarrantyExpired: (state) => {
      return state.goods.filter(g => g.warrantyStatus === 'expired')
    },

    goodsWarrantyExpiringSoon: (state) => {
      // Warranty expiring in next 30 days
      return state.goods.filter(g => {
        return g.warrantyStatus === 'active' &&
               g.warrantyDaysRemaining !== null &&
               g.warrantyDaysRemaining <= 30
      })
    },

    totalValue: (state) => {
      return state.goods
        .filter(g => g.status === 'active')
        .reduce((sum, g) => sum + g.purchasePrice, 0)
    },

    totalServiceCost: (state) => {
      return state.goods.reduce((sum, g) => sum + (g.totalServiceCost || 0), 0)
    },

    totalCost: (state) => {
      return state.goods
        .filter(g => g.status === 'active')
        .reduce((sum, g) => sum + (g.totalCost || g.purchasePrice), 0)
    },

    getServiceRecordsForGood: (state) => (goodId) => {
      return state.serviceRecords[goodId] || []
    },

    getAttachmentsForGood: (state) => (goodId) => {
      return state.attachments[goodId] || []
    },

    categories: (state) => {
      const cats = new Set()
      state.goods.forEach(g => {
        if (g.category) cats.add(g.category)
      })
      return Array.from(cats).sort()
    },

    brands: (state) => {
      const brands = new Set()
      state.goods.forEach(g => {
        if (g.brand) brands.add(g.brand)
      })
      return Array.from(brands).sort()
    }
  },

  actions: {
    async fetchGoods(filters = {}) {
      this.loading = true
      this.error = null
      try {
        const params = {}
        if (filters.status) params.status = filters.status
        if (filters.category) params.category = filters.category

        const response = await apiService.query('goods', params)
        this.goods = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching goods:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchGood(id) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('goods', id)
        const index = this.goods.findIndex(g => g.id === id)
        if (index !== -1) {
          this.goods[index] = response.data
        } else {
          this.goods.push(response.data)
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

    async createGood(goodData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('goods', goodData)
        this.goods.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating good:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateGood(id, goodData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('goods', id, goodData)
        const index = this.goods.findIndex(g => g.id === id)
        if (index !== -1) {
          this.goods[index] = response.data
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

    async deleteGood(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('goods', id)
        this.goods = this.goods.filter(g => g.id !== id)
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

    async createServiceRecord(goodId, serviceData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`goods/${goodId}/services`, serviceData)

        // Add to local cache
        if (!this.serviceRecords[goodId]) {
          this.serviceRecords[goodId] = []
        }
        this.serviceRecords[goodId].unshift(response.data)

        // Refresh the good to update stats
        await this.fetchGood(goodId)

        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating service record:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchServiceRecords(goodId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`goods/${goodId}/services`)
        this.serviceRecords[goodId] = response.data || []
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching service records:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteServiceRecord(goodId, serviceId) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete(`goods/${goodId}/services/${serviceId}`)

        // Remove from local cache
        if (this.serviceRecords[goodId]) {
          this.serviceRecords[goodId] = this.serviceRecords[goodId].filter(
            s => s.id !== serviceId
          )
        }

        // Refresh the good to update stats
        await this.fetchGood(goodId)
      } catch (error) {
        this.error = error.message
        console.error('Error deleting service record:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async addAttachment(goodId, attachmentData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post(`goods/${goodId}/attachments`, attachmentData)

        // Add to local cache
        if (!this.attachments[goodId]) {
          this.attachments[goodId] = []
        }
        this.attachments[goodId].unshift(response.data)

        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error adding attachment:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchAttachments(goodId) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get(`goods/${goodId}/attachments`)
        this.attachments[goodId] = response.data || []
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error fetching attachments:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAttachment(goodId, attachmentId) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete(`goods/${goodId}/attachments/${attachmentId}`)

        // Remove from local cache
        if (this.attachments[goodId]) {
          this.attachments[goodId] = this.attachments[goodId].filter(
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
        const response = await apiService.get('goods/stats')
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
