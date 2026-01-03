<template>
  <div class="settings-view fade-in">
    <h1 class="text-purple mb-4">Settings</h1>

    <div class="row g-3">
      <!-- Profile Settings -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Profile Information</h5>
          </div>
          <div class="card-body">
            <form @submit.prevent="updateProfile">
              <div class="mb-3">
                <label class="form-label">Username</label>
                <input
                  type="text"
                  class="form-control"
                  :value="currentUser?.username"
                  disabled
                  readonly
                />
                <small class="text-muted">Username cannot be changed</small>
              </div>

              <div class="mb-3">
                <label class="form-label">Full Name</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="profileForm.fullName"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Email</label>
                <input
                  type="email"
                  class="form-control"
                  v-model="profileForm.email"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Role</label>
                <input
                  type="text"
                  class="form-control"
                  :value="currentUser?.role || 'user'"
                  disabled
                  readonly
                />
              </div>

              <button type="submit" class="btn btn-primary" :disabled="profileLoading">
                <span v-if="profileLoading" class="spinner-border spinner-border-sm me-2"></span>
                {{ profileLoading ? 'Updating...' : 'Update Profile' }}
              </button>
            </form>
          </div>
        </div>
      </div>

      <!-- Change Password -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Change Password</h5>
          </div>
          <div class="card-body">
            <form @submit.prevent="changePassword">
              <div class="mb-3">
                <label class="form-label">Current Password</label>
                <input
                  type="password"
                  class="form-control"
                  v-model="passwordForm.currentPassword"
                  required
                  autocomplete="current-password"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">New Password</label>
                <input
                  type="password"
                  class="form-control"
                  v-model="passwordForm.newPassword"
                  required
                  minlength="6"
                  autocomplete="new-password"
                />
                <small class="text-muted">Minimum 6 characters</small>
              </div>

              <div class="mb-3">
                <label class="form-label">Confirm New Password</label>
                <input
                  type="password"
                  class="form-control"
                  v-model="passwordForm.confirmPassword"
                  required
                  minlength="6"
                  autocomplete="new-password"
                />
              </div>

              <button type="submit" class="btn btn-primary" :disabled="passwordLoading">
                <span v-if="passwordLoading" class="spinner-border spinner-border-sm me-2"></span>
                {{ passwordLoading ? 'Changing...' : 'Change Password' }}
              </button>
            </form>
          </div>
        </div>
      </div>

      <!-- General Settings -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">General Settings</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <label class="form-label">Default Currency</label>
              <select class="form-select" v-model="settings.currency" @change="saveSettings">
                <option v-for="currency in currencies" :key="currency.code" :value="currency.code">
                  {{ currency.code }} - {{ currency.name }} ({{ currency.symbol }})
                </option>
              </select>
            </div>

            <div class="mb-3">
              <label class="form-label">Date Format</label>
              <select class="form-select" v-model="settings.dateFormat" @change="saveSettings">
                <option value="MM/DD/YYYY">MM/DD/YYYY</option>
                <option value="DD/MM/YYYY">DD/MM/YYYY</option>
                <option value="YYYY-MM-DD">YYYY-MM-DD</option>
              </select>
            </div>

            <div class="mb-3">
              <label class="form-label">First Day of Week</label>
              <select class="form-select" v-model.number="settings.firstDayOfWeek" @change="saveSettings">
                <option :value="0">Sunday</option>
                <option :value="1">Monday</option>
              </select>
            </div>

            <div class="mb-3">
              <label class="form-label">Language</label>
              <select class="form-select" v-model="settings.language" @change="saveSettings">
                <option value="en">English</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- Appearance -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Appearance</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <div class="form-check form-switch">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="darkMode"
                  v-model="settings.darkMode"
                  @change="saveSettings"
                />
                <label class="form-check-label" for="darkMode">
                  Dark Mode
                </label>
              </div>
              <small class="text-muted">Enable dark theme for better viewing in low light</small>
            </div>
          </div>
        </div>

        <!-- Notifications -->
        <div class="card mt-3">
          <div class="card-header">
            <h5 class="mb-0">Notifications</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <div class="form-check form-switch">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="pushNotifications"
                  v-model="settings.notifications.push"
                  @change="saveSettings"
                />
                <label class="form-check-label" for="pushNotifications">
                  Push Notifications
                </label>
              </div>
            </div>

            <div class="mb-3">
              <div class="form-check form-switch">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="emailNotifications"
                  v-model="settings.notifications.email"
                  @change="saveSettings"
                />
                <label class="form-check-label" for="emailNotifications">
                  Email Notifications
                </label>
              </div>
            </div>

            <div class="mb-3">
              <div class="form-check form-switch">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="budgetAlerts"
                  v-model="settings.notifications.budgetAlerts"
                  @change="saveSettings"
                />
                <label class="form-check-label" for="budgetAlerts">
                  Budget Alerts
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Data Management -->
      <div class="col-12">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Data Management</h5>
          </div>
          <div class="card-body">
            <div class="row g-3">
              <div class="col-12">
                <h6 class="mb-3">Export Your Data</h6>
                <p class="text-muted">Download your financial data in CSV or JSON format</p>
              </div>

              <!-- Export Options -->
              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Transactions</h6>
                    <p class="card-text small text-muted">Export all your income, expenses, and transfers</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('transactions', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('transactions', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Accounts</h6>
                    <p class="card-text small text-muted">Export all your accounts and balances</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('accounts', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('accounts', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Budgets</h6>
                    <p class="card-text small text-muted">Export all your budget configurations</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('budgets', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('budgets', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Goals</h6>
                    <p class="card-text small text-muted">Export all your savings and investment goals</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('goals', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('goals', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Categories</h6>
                    <p class="card-text small text-muted">Export all your income and expense categories</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('categories', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('categories', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100">
                  <div class="card-body">
                    <h6 class="card-title">Assets</h6>
                    <p class="card-text small text-muted">Export all your assets and warranty information</p>
                    <div class="btn-group w-100" role="group">
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('assets', 'csv')">
                        CSV
                      </button>
                      <button class="btn btn-sm btn-outline-primary" @click="exportData('assets', 'json')">
                        JSON
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="col-12 col-md-6 col-lg-4">
                <div class="card h-100 border-primary">
                  <div class="card-body">
                    <h6 class="card-title text-primary">All Data</h6>
                    <p class="card-text small text-muted">Complete export of all your financial data</p>
                    <button class="btn btn-sm btn-primary w-100" @click="exportData('all', 'json')">
                      Export Everything
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Database Backups -->
      <div class="col-12">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">Database Backups</h5>
            <button class="btn btn-primary btn-sm" @click="createBackup" :disabled="backupLoading">
              <span v-if="backupLoading" class="spinner-border spinner-border-sm me-2"></span>
              {{ backupLoading ? 'Creating...' : 'Create New Backup' }}
            </button>
          </div>
          <div class="card-body">
            <div v-if="backups.length === 0" class="text-center text-muted py-4">
              No backups yet. Click "Create New Backup" to get started.
            </div>
            <div v-else class="table-responsive">
              <table class="table table-hover">
                <thead>
                  <tr>
                    <th>File Name</th>
                    <th>Created</th>
                    <th>Size</th>
                    <th>Status</th>
                    <th class="text-end">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="backup in backups" :key="backup.id">
                    <td>
                      <span class="fw-semibold">{{ backup.fileName }}</span>
                    </td>
                    <td>{{ formatDate(backup.createdAt) }}</td>
                    <td>{{ formatFileSize(backup.fileSize) }}</td>
                    <td>
                      <span class="badge" :class="getStatusClass(backup.status)">
                        {{ backup.status }}
                      </span>
                    </td>
                    <td class="text-end">
                      <button
                        v-if="backup.status === 'completed'"
                        class="btn btn-sm btn-success me-2"
                        @click="downloadBackup(backup.id)"
                        title="Download backup"
                      >
                        Download
                      </button>
                      <button
                        class="btn btn-sm btn-danger"
                        @click="deleteBackup(backup.id)"
                        title="Delete backup"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useAuthStore } from '@/stores/auth'
import { useBackupsStore } from '@/stores/backups'
import { useNotification } from '@/composables/useNotification'

const settingsStore = useSettingsStore()
const authStore = useAuthStore()
const backupsStore = useBackupsStore()
const { info, success, error } = useNotification()

const settings = computed(() => settingsStore.settings)
const currencies = computed(() => settingsStore.currencies)
const currentUser = computed(() => authStore.currentUser)
const backups = computed(() => backupsStore.allBackups)
const backupLoading = computed(() => backupsStore.loading)

// Profile form
const profileForm = ref({
  fullName: '',
  email: ''
})
const profileLoading = ref(false)

// Password form
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const passwordLoading = ref(false)

// Initialize profile form with current user data
const initializeProfileForm = () => {
  if (currentUser.value) {
    profileForm.value = {
      fullName: currentUser.value.fullName || '',
      email: currentUser.value.email || ''
    }
  }
}

const updateProfile = async () => {
  profileLoading.value = true
  try {
    await authStore.updateProfile(profileForm.value)
    success('Profile updated successfully')
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error updating profile')
  } finally {
    profileLoading.value = false
  }
}

const changePassword = async () => {
  // Validate passwords match
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    error('New passwords do not match')
    return
  }

  // Validate password length
  if (passwordForm.value.newPassword.length < 6) {
    error('Password must be at least 6 characters')
    return
  }

  passwordLoading.value = true
  try {
    await authStore.changePassword(
      passwordForm.value.currentPassword,
      passwordForm.value.newPassword
    )
    success('Password changed successfully')
    // Clear the form
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error changing password')
  } finally {
    passwordLoading.value = false
  }
}

const saveSettings = async () => {
  try {
    await settingsStore.updateSettings(settings.value)
    success('Settings saved successfully')
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving settings')
  }
}

const exportData = async (type, format) => {
  try {
    info(`Preparing ${type} export...`)

    // Get the auth token
    const token = localStorage.getItem('auth_token')
    if (!token) {
      error('You must be logged in to export data')
      return
    }

    // Build the API URL
    const apiURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
    let url = `${apiURL}/export?type=${type}&format=${format}`

    // For transactions, add date range (last year by default)
    if (type === 'transactions') {
      const endDate = new Date().toISOString().split('T')[0]
      const startDate = new Date(new Date().setFullYear(new Date().getFullYear() - 1)).toISOString().split('T')[0]
      url += `&start_date=${startDate}&end_date=${endDate}`
    }

    // Fetch the file
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      throw new Error(`Export failed: ${response.statusText}`)
    }

    // Get the filename from the Content-Disposition header
    const contentDisposition = response.headers.get('Content-Disposition')
    let filename = `${type}.${format}`
    if (contentDisposition) {
      const matches = /filename=([^;]+)/.exec(contentDisposition)
      if (matches && matches[1]) {
        filename = matches[1].replace(/['"]/g, '')
      }
    }

    // Download the file
    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)

    success(`Successfully exported ${type} as ${format.toUpperCase()}`)
  } catch (err) {
    console.error('Export error:', err)
    error(err.message || 'Failed to export data')
  }
}

// Backup management
const createBackup = async () => {
  try {
    info('Creating database backup...')
    await backupsStore.createBackup()
    success('Backup initiated successfully')
  } catch (err) {
    error(err.response?.data?.error || err.message || 'Failed to create backup')
  }
}

const downloadBackup = async (backupId) => {
  try {
    info('Downloading backup...')
    await backupsStore.downloadBackup(backupId)
    success('Backup downloaded successfully')
  } catch (err) {
    error(err.response?.data?.error || err.message || 'Failed to download backup')
  }
}

const deleteBackup = async (backupId) => {
  if (!confirm('Are you sure you want to delete this backup? This action cannot be undone.')) {
    return
  }

  try {
    await backupsStore.deleteBackup(backupId)
    success('Backup deleted successfully')
  } catch (err) {
    error(err.response?.data?.error || err.message || 'Failed to delete backup')
  }
}

const formatFileSize = (bytes) => {
  return backupsStore.formatFileSize(bytes)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusClass = (status) => {
  switch (status) {
    case 'completed':
      return 'bg-success'
    case 'pending':
      return 'bg-warning'
    case 'failed':
      return 'bg-danger'
    default:
      return 'bg-secondary'
  }
}

onMounted(async () => {
  await settingsStore.loadSettings()
  initializeProfileForm()
  await backupsStore.fetchBackups()
})
</script>
