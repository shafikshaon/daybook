import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'

export const useAccountTypesStore = defineStore('accountTypes', {
  state: () => ({
    accountTypes: [],
    loading: false,
    error: null
  }),

  getters: {
    allAccountTypes: (state) => state.accountTypes,

    getTypeById: (state) => (id) => {
      return state.accountTypes.find(type => type.id === id)
    }
  },

  actions: {
    async fetchAccountTypes() {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('account-types')
        this.accountTypes = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching account types:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createAccountType(typeData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('account-types', typeData)
        this.accountTypes.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating account type:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateAccountType(id, typeData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('account-types', id, typeData)
        const index = this.accountTypes.findIndex(t => t.id === id)
        if (index !== -1) {
          this.accountTypes[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error updating account type:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAccountType(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('account-types', id)
        this.accountTypes = this.accountTypes.filter(t => t.id !== id)
      } catch (error) {
        this.error = error.message
        console.error('Error deleting account type:', error)
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
