import axios from 'axios'

// Create axios instance with backend API baseURL
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
api.interceptors.response.use(
  response => {
    // Extract data from standardized backend response
    if (response.data && response.data.data !== undefined) {
      return { ...response, data: response.data.data }
    }
    return response
  },
  error => {
    if (error.response?.status === 401) {
      // Unauthorized - clear auth and redirect to login
      // But don't redirect if we're already on login/signup page or if this is a login request
      const isAuthPage = window.location.pathname === '/login' || window.location.pathname === '/signup'
      const isAuthRequest = error.config?.url?.includes('/auth/login') || error.config?.url?.includes('/auth/signup')

      if (!isAuthPage && !isAuthRequest) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        window.location.href = '/login'
      }
    }

    // Extract error message from backend response
    const errorMessage = error.response?.data?.error || error.message || 'An error occurred'
    return Promise.reject(new Error(errorMessage))
  }
)

// API service methods
const apiService = {
  // Authentication
  auth: {
    async signup(userData) {
      const response = await api.post('/auth/signup', userData)
      return response
    },

    async login(credentials) {
      const response = await api.post('/auth/login', credentials)
      return response
    },

    async getProfile() {
      const response = await api.get('/auth/me')
      return response
    },

    async updateProfile(profileData) {
      const response = await api.put('/auth/profile', profileData)
      return response
    },

    async changePassword(passwordData) {
      const response = await api.put('/auth/change-password', passwordData)
      return response
    }
  },

  // Generic CRUD operations
  async get(endpoint, params = null) {
    let url = `/${endpoint}`

    // If params is a string, treat it as an ID
    if (typeof params === 'string') {
      url = `/${endpoint}/${params}`
    }
    // If params is an object, treat it as query parameters
    else if (params && typeof params === 'object') {
      const queryParams = new URLSearchParams()
      Object.keys(params).forEach(key => {
        if (params[key] !== undefined && params[key] !== null) {
          queryParams.append(key, params[key])
        }
      })
      const queryString = queryParams.toString()
      if (queryString) {
        url = `/${endpoint}?${queryString}`
      }
    }

    const response = await api.get(url)
    return response
  },

  async post(endpoint, payload) {
    const response = await api.post(`/${endpoint}`, payload)
    return response
  },

  async put(endpoint, id, payload) {
    const url = id ? `/${endpoint}/${id}` : `/${endpoint}`
    const response = await api.put(url, payload)
    return response
  },

  async patch(endpoint, id, payload) {
    const response = await api.patch(`/${endpoint}/${id}`, payload)
    return response
  },

  async delete(endpoint, id) {
    const response = await api.delete(`/${endpoint}/${id}`)
    return response
  },

  // Bulk operations
  async bulkCreate(endpoint, items) {
    const response = await api.post(`/${endpoint}/bulk`, items)
    return response
  },

  // Query with filters
  async query(endpoint, filters = {}) {
    const params = new URLSearchParams()
    Object.keys(filters).forEach(key => {
      if (filters[key] !== undefined && filters[key] !== null) {
        params.append(key, filters[key])
      }
    })

    const url = params.toString() ? `/${endpoint}?${params.toString()}` : `/${endpoint}`
    const response = await api.get(url)
    return response
  },

  // Transaction-specific
  async getTransactionStats(filters = {}) {
    const params = new URLSearchParams()
    Object.keys(filters).forEach(key => {
      if (filters[key]) params.append(key, filters[key])
    })

    const url = params.toString() ? `/transactions/stats?${params.toString()}` : '/transactions/stats'
    const response = await api.get(url)
    return response
  },

  // Credit card payment
  async recordCreditCardPayment(cardId, paymentData) {
    const response = await api.post(`/credit-cards/${cardId}/payment`, paymentData)
    return response
  },

  // Investment buy/sell
  async buyShares(investmentId, data) {
    const response = await api.post(`/investments/${investmentId}/buy`, data)
    return response
  },

  async sellShares(investmentId, data) {
    const response = await api.post(`/investments/${investmentId}/sell`, data)
    return response
  },

  // Bill payment
  async payBill(billId, paymentData) {
    const response = await api.post(`/bills/${billId}/pay`, paymentData)
    return response
  },

  // Budget progress
  async getBudgetProgress(budgetId) {
    const response = await api.get(`/budgets/${budgetId}/progress`)
    return response
  },

  // Savings goal actions
  async contributeToGoal(goalId, contributionData) {
    const response = await api.post(`/savings-goals/${goalId}/contribute`, contributionData)
    return response
  },

  async withdrawFromGoal(goalId, withdrawData) {
    const response = await api.post(`/savings-goals/${goalId}/withdraw`, withdrawData)
    return response
  },

  // Fixed deposit withdrawal
  async withdrawFixedDeposit(depositId) {
    const response = await api.post(`/fixed-deposits/${depositId}/withdraw`)
    return response
  },

  // Account balance update
  async updateAccountBalance(accountId, balanceData) {
    const response = await api.patch(`/accounts/${accountId}/balance`, balanceData)
    return response
  },

  // Reconciliation
  reconciliation: {
    // Get all reconciliations (optionally filtered by accountId)
    async getAll(accountId = null) {
      const url = accountId ? `/reconciliations?accountId=${accountId}` : '/reconciliations'
      const response = await api.get(url)
      return response
    },

    // Get reconciliations for a specific account
    async getByAccount(accountId) {
      const response = await api.get(`/accounts/${accountId}/reconciliations`)
      return response
    },

    // Get a single reconciliation by ID
    async getById(reconciliationId) {
      const response = await api.get(`/reconciliations/${reconciliationId}`)
      return response
    },

    // Create a new reconciliation
    async create(reconciliationData) {
      const response = await api.post('/reconciliations', reconciliationData)
      return response
    },

    // Update a reconciliation
    async update(reconciliationId, reconciliationData) {
      const response = await api.put(`/reconciliations/${reconciliationId}`, reconciliationData)
      return response
    },

    // Delete a reconciliation
    async delete(reconciliationId) {
      const response = await api.delete(`/reconciliations/${reconciliationId}`)
      return response
    },

    // Get unreconciled transactions for an account
    async getUnreconciledTransactions(accountId) {
      const response = await api.get(`/accounts/${accountId}/unreconciled-transactions`)
      return response
    },

    // Get reconciliation statistics for an account
    async getStats(accountId) {
      const response = await api.get(`/accounts/${accountId}/reconciliations/stats`)
      return response
    }
  },

  // Utility methods
  generateId() {
    // Note: Backend generates UUIDs, so this is not used with real backend
    return Date.now().toString(36) + Math.random().toString(36).substring(2)
  }
}

export default apiService
export { api }
