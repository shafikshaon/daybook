<template>
  <div class="settings-view fade-in">
    <h1 class="text-purple mb-4">Settings</h1>

    <div class="row g-3">
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

            <div class="mb-3">
              <div class="form-check form-switch">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="billReminders"
                  v-model="settings.notifications.billReminders"
                  @change="saveSettings"
                />
                <label class="form-check-label" for="billReminders">
                  Bill Reminders
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
              <div class="col-12 col-md-4">
                <button class="btn btn-outline-primary w-100" @click="exportData">
                  📤 Export Data
                </button>
                <small class="text-muted d-block mt-2">Download all your data as JSON</small>
              </div>
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
import { useNotification } from '@/composables/useNotification'

const settingsStore = useSettingsStore()
const { info, success, error } = useNotification()

const settings = computed(() => settingsStore.settings)
const currencies = computed(() => settingsStore.currencies)

const saveSettings = async () => {
  try {
    await settingsStore.updateSettings(settings.value)
    success('Settings saved successfully')
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving settings')
  }
}

const exportData = () => {
  // Export data functionality would require backend API endpoint
  info('Export functionality will be available in a future update.')
}

onMounted(async () => {
  await settingsStore.loadSettings()
})
</script>
