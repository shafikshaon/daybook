import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'

export const useCategoriesStore = defineStore('categories', {
  state: () => ({
    categories: [],
    loading: false,
    error: null
  }),

  getters: {
    allCategories: (state) => state.categories,

    incomeCategories: (state) => state.categories.filter(cat => cat.type === 'income'),

    expenseCategories: (state) => state.categories.filter(cat => cat.type === 'expense'),

    getCategoryById: (state) => (id) => {
      return state.categories.find(cat => cat.id === id)
    }
  },

  actions: {
    async fetchCategories() {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.get('categories')
        this.categories = response.data || []
      } catch (error) {
        this.error = error.message
        console.error('Error fetching categories:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async createCategory(categoryData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.post('categories', categoryData)
        this.categories.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error creating category:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateCategory(id, categoryData) {
      this.loading = true
      this.error = null
      try {
        const response = await apiService.put('categories', id, categoryData)
        const index = this.categories.findIndex(c => c.id === id)
        if (index !== -1) {
          this.categories[index] = response.data
        }
        return response.data
      } catch (error) {
        this.error = error.message
        console.error('Error updating category:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteCategory(id) {
      this.loading = true
      this.error = null
      try {
        await apiService.delete('categories', id)
        this.categories = this.categories.filter(c => c.id !== id)
      } catch (error) {
        this.error = error.message
        console.error('Error deleting category:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async reorderCategories(categoryOrders) {
      this.loading = true
      this.error = null
      try {
        await apiService.put('categories/reorder', null, { categories: categoryOrders })

        // Update local state with new orders
        categoryOrders.forEach(({ id, order }) => {
          const category = this.categories.find(c => c.id === id)
          if (category) {
            category.order = order
          }
        })

        // Re-sort categories by type and order
        this.categories.sort((a, b) => {
          if (a.type !== b.type) return a.type.localeCompare(b.type)
          return (a.order || 0) - (b.order || 0)
        })
      } catch (error) {
        this.error = error.message
        console.error('Error reordering categories:', error)
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
