import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    isAuthenticated: false,
    token: null
  }),

  getters: {
    currentUser: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated,
    userRole: (state) => state.user?.role || 'user'
  },

  actions: {
    async initializeAuth() {
      try {
        // Check if user is logged in from localStorage
        const token = localStorage.getItem('auth_token')
        const userStr = localStorage.getItem('auth_user')

        if (token && userStr) {
          this.token = token
          this.user = JSON.parse(userStr)
          this.isAuthenticated = true

          // Validate token with backend
          try {
            const profileResponse = await apiService.auth.getProfile()
            this.user = profileResponse.data
            localStorage.setItem('auth_user', JSON.stringify(this.user))
          } catch (error) {
            // Token is invalid, logout
            console.error('Token validation failed:', error)
            this.logout()
          }
        }
      } catch (error) {
        console.error('Error initializing auth:', error)
        this.logout()
      }
    },

    async login(username, password) {
      try {
        const response = await apiService.auth.login({ username, password })
        const { token, user } = response.data

        // Store auth data
        this.token = token
        this.user = user
        this.isAuthenticated = true

        // Persist to localStorage
        localStorage.setItem('auth_token', token)
        localStorage.setItem('auth_user', JSON.stringify(user))

        return { success: true, user: this.user }
      } catch (error) {
        console.error('Login error:', error)
        throw error
      }
    },

    async signup(userData) {
      try {
        const response = await apiService.auth.signup(userData)
        const { token, user } = response.data

        // Store auth data
        this.token = token
        this.user = user
        this.isAuthenticated = true

        // Persist to localStorage
        localStorage.setItem('auth_token', token)
        localStorage.setItem('auth_user', JSON.stringify(user))

        return { success: true, user: this.user }
      } catch (error) {
        console.error('Signup error:', error)
        throw error
      }
    },

    async logout() {
      // Clear state
      this.user = null
      this.token = null
      this.isAuthenticated = false

      // Clear localStorage
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    },

    async changePassword(currentPassword, newPassword) {
      try {
        if (!this.user) {
          throw new Error('Not authenticated')
        }

        await apiService.auth.changePassword({
          currentPassword,
          newPassword
        })

        return { success: true }
      } catch (error) {
        console.error('Change password error:', error)
        throw error
      }
    },

    async updateProfile(profileData) {
      try {
        if (!this.user) {
          throw new Error('Not authenticated')
        }

        const response = await apiService.auth.updateProfile(profileData)
        this.user = response.data

        // Update localStorage
        localStorage.setItem('auth_user', JSON.stringify(this.user))

        return { success: true, user: this.user }
      } catch (error) {
        console.error('Update profile error:', error)
        throw error
      }
    }
  }
})
