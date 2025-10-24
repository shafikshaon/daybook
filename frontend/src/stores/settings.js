import { defineStore } from 'pinia'
import apiService from '@/services/api-backend'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {
      currency: 'USD',
      darkMode: false,
      notifications: {
        push: true,
        email: true,
        budgetAlerts: true,
        billReminders: true
      },
      dateFormat: 'MM/DD/YYYY',
      firstDayOfWeek: 0, // 0 = Sunday
      language: 'en'
    },
    currencies: [
      { code: 'USD', symbol: '$', name: 'US Dollar' },
      { code: 'EUR', symbol: '€', name: 'Euro' },
      { code: 'GBP', symbol: '£', name: 'British Pound' },
      { code: 'JPY', symbol: '¥', name: 'Japanese Yen' },
      { code: 'INR', symbol: '₹', name: 'Indian Rupee' },
      { code: 'AUD', symbol: 'A$', name: 'Australian Dollar' },
      { code: 'CAD', symbol: 'C$', name: 'Canadian Dollar' },
      { code: 'CHF', symbol: 'CHF', name: 'Swiss Franc' },
      { code: 'CNY', symbol: '¥', name: 'Chinese Yuan' },
      { code: 'BDT', symbol: '৳', name: 'Bangladeshi Taka' }
    ]
  }),

  getters: {
    currentCurrency: (state) => {
      return state.currencies.find(c => c.code === state.settings.currency) || state.currencies[0]
    },

    isDarkMode: (state) => state.settings.darkMode,

    currencySymbol: (state) => {
      const currency = state.currencies.find(c => c.code === state.settings.currency)
      return currency ? currency.symbol : '$'
    }
  },

  actions: {
    async loadSettings() {
      try {
        // Check if user is authenticated before loading settings
        const token = localStorage.getItem('auth_token')
        if (!token) {
          console.log('No auth token, using default settings')
          return
        }

        const response = await apiService.get('settings')
        if (response.data) {
          this.settings = { ...this.settings, ...response.data }
        }

        // Apply dark mode
        if (this.settings.darkMode) {
          document.body.classList.add('dark-mode')
        } else {
          document.body.classList.remove('dark-mode')
        }
      } catch (error) {
        console.error('Error loading settings:', error)
        // Use default settings if endpoint doesn't exist yet
        // Don't throw error, just use defaults
      }
    },

    async updateSettings(newSettings) {
      try {
        this.settings = { ...this.settings, ...newSettings }

        // Settings endpoint doesn't use ID, pass null to use /settings directly
        await apiService.put('settings', null, this.settings)

        // Apply dark mode
        if (this.settings.darkMode) {
          document.body.classList.add('dark-mode')
        } else {
          document.body.classList.remove('dark-mode')
        }
      } catch (error) {
        console.error('Error updating settings:', error)
        // If PUT fails, try POST for creating settings
        try {
          await apiService.post('settings', this.settings)
          if (this.settings.darkMode) {
            document.body.classList.add('dark-mode')
          } else {
            document.body.classList.remove('dark-mode')
          }
        } catch (postError) {
          console.error('Error creating settings:', postError)
          throw postError
        }
      }
    },

    toggleDarkMode() {
      this.updateSettings({ darkMode: !this.settings.darkMode })
    },

    setCurrency(currencyCode) {
      this.updateSettings({ currency: currencyCode })
    },

    formatCurrency(amount, currencyCode = null) {
      const currency = currencyCode || this.settings.currency
      const currencyObj = this.currencies.find(c => c.code === currency)
      const symbol = currencyObj ? currencyObj.symbol : '$'

      return `${symbol}${Math.abs(amount).toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })}`
    }
  }
})
