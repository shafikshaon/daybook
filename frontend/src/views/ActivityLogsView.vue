<template>
  <div class="activity-logs-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Activity Logs</h1>
      <div class="d-flex gap-2">
        <button class="btn btn-outline-success" @click="showRefillModal = true">
          🔄 Refill All Logs
        </button>
        <button class="btn btn-outline-danger" @click="showCleanupModal = true">
          🗑️ Cleanup Old Logs
        </button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div v-if="summary" class="row g-3 mb-4">
      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon purple">📊</div>
          <div class="stat-value">{{ summary.totalActivities }}</div>
          <div class="stat-label">Total Activities</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            Last {{ summaryDays }} days
          </small>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon blue">🔥</div>
          <div class="stat-value">{{ summary.topModules?.[0]?.module || 'N/A' }}</div>
          <div class="stat-label">Most Active Module</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            {{ summary.topModules?.[0]?.count || 0 }} actions
          </small>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon green">📝</div>
          <div class="stat-value">{{ summary.actionCounts?.length || 0 }}</div>
          <div class="stat-label">Action Types</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            Create, Update, Delete, etc.
          </small>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon orange">⚡</div>
          <div class="stat-value">{{ summary.recentActivities?.length || 0 }}</div>
          <div class="stat-label">Recent Activities</div>
          <small class="text-muted d-block mt-1" style="font-size: 0.75rem;">
            Latest actions
          </small>
        </div>
      </div>
    </div>

    <!-- Activity Summary Charts -->
    <div v-if="summary" class="row g-3 mb-4">
      <div class="col-12 col-md-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Actions Breakdown</h5>
          </div>
          <div class="card-body">
            <div v-if="summary.actionCounts?.length > 0">
              <div v-for="action in summary.actionCounts" :key="action.action" class="mb-2">
                <div class="d-flex justify-content-between align-items-center">
                  <span class="text-capitalize">{{ action.action }}</span>
                  <span class="badge bg-primary">{{ action.count }}</span>
                </div>
                <div class="progress mt-1" style="height: 6px;">
                  <div
                    class="progress-bar"
                    :style="{ width: (action.count / summary.totalActivities * 100) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
            <div v-else class="text-muted text-center py-3">
              No action data available
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-md-6">
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">Module Activity</h5>
          </div>
          <div class="card-body">
            <div v-if="summary.moduleCounts?.length > 0">
              <div v-for="module in summary.moduleCounts" :key="module.module" class="mb-2">
                <div class="d-flex justify-content-between align-items-center">
                  <span class="text-capitalize">{{ module.module.replace('_', ' ') }}</span>
                  <span class="badge bg-info">{{ module.count }}</span>
                </div>
                <div class="progress mt-1" style="height: 6px;">
                  <div
                    class="progress-bar bg-info"
                    :style="{ width: (module.count / summary.totalActivities * 100) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
            <div v-else class="text-muted text-center py-3">
              No module data available
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="card mb-3">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-3">
            <label class="form-label">Module</label>
            <select v-model="localFilters.module" class="form-select" @change="applyFilters">
              <option :value="null">All Modules</option>
              <option value="auth">Authentication</option>
              <option value="account">Accounts</option>
              <option value="transaction">Transactions</option>
              <option value="budget">Budgets</option>
              <option value="credit_card">Credit Cards</option>
              <option value="debt">Debts</option>
              <option value="lend">Lends</option>
              <option value="asset">Assets</option>
              <option value="goal">Goals</option>
              <option value="settings">Settings</option>
            </select>
          </div>

          <div class="col-md-3">
            <label class="form-label">Action</label>
            <select v-model="localFilters.action" class="form-select" @change="applyFilters">
              <option :value="null">All Actions</option>
              <option value="create">Create</option>
              <option value="update">Update</option>
              <option value="delete">Delete</option>
              <option value="view">View</option>
              <option value="login">Login</option>
              <option value="logout">Logout</option>
              <option value="export">Export</option>
              <option value="import">Import</option>
            </select>
          </div>

          <div class="col-md-3">
            <label class="form-label">Start Date</label>
            <input v-model="localFilters.startDate" type="date" class="form-control" @change="applyFilters">
          </div>

          <div class="col-md-3">
            <label class="form-label">End Date</label>
            <input v-model="localFilters.endDate" type="date" class="form-control" @change="applyFilters">
          </div>
        </div>

        <div class="mt-3 d-flex justify-content-end gap-2">
          <button class="btn btn-outline-secondary" @click="clearFilters">
            Clear Filters
          </button>
          <button class="btn btn-primary" @click="refreshLogs">
            🔄 Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Activity Logs List -->
    <div class="card">
      <div class="card-header">
        <h5 class="mb-0">Activity History</h5>
      </div>
      <div class="card-body p-0">
        <div v-if="loading" class="p-4 text-center">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <div v-else-if="activityLogs.length === 0" class="p-4 text-center text-muted">
          <p>No activity logs found.</p>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Date & Time</th>
                <th>Action</th>
                <th>Module</th>
                <th>Description</th>
                <th>IP Address</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="log in activityLogs" :key="log.id">
                <td>
                  <small>{{ formatDateTime(log.createdAt) }}</small>
                </td>
                <td>
                  <span class="badge" :class="getActionBadgeClass(log.action)">
                    {{ log.action }}
                  </span>
                </td>
                <td>
                  <span class="badge bg-secondary text-capitalize">
                    {{ log.module.replace('_', ' ') }}
                  </span>
                </td>
                <td>
                  <span class="text-muted">{{ log.description || '-' }}</span>
                </td>
                <td>
                  <small class="text-muted">{{ log.ipAddress || '-' }}</small>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.totalPages > 1" class="card-footer">
        <div class="d-flex justify-content-between align-items-center">
          <div class="text-muted">
            Showing {{ activityLogs.length }} of {{ pagination.total }} logs
          </div>
          <nav>
            <ul class="pagination mb-0">
              <li class="page-item" :class="{ disabled: !hasPrevPage }">
                <button class="page-link" @click="prevPage" :disabled="!hasPrevPage">
                  Previous
                </button>
              </li>
              <li class="page-item active">
                <span class="page-link">
                  {{ pagination.currentPage }} / {{ pagination.totalPages }}
                </span>
              </li>
              <li class="page-item" :class="{ disabled: !hasNextPage }">
                <button class="page-link" @click="nextPage" :disabled="!hasNextPage">
                  Next
                </button>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>

    <!-- Cleanup Modal -->
    <div
      class="modal fade"
      id="cleanupModal"
      tabindex="-1"
      :class="{ show: showCleanupModal }"
      :style="{ display: showCleanupModal ? 'block' : 'none' }"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Cleanup Old Activity Logs</h5>
            <button type="button" class="btn-close" @click="showCleanupModal = false"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">Delete logs older than (days):</label>
              <input v-model.number="cleanupDays" type="number" min="30" class="form-control">
              <small class="text-muted">Minimum: 30 days</small>
            </div>
            <div class="alert alert-warning">
              <strong>Warning:</strong> This action cannot be undone. All activity logs older than {{ cleanupDays }} days will be permanently deleted.
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showCleanupModal = false">
              Cancel
            </button>
            <button type="button" class="btn btn-danger" @click="cleanupLogs">
              Delete Old Logs
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showCleanupModal" class="modal-backdrop fade show"></div>

    <!-- Refill Modal -->
    <div
      class="modal fade"
      id="refillModal"
      tabindex="-1"
      :class="{ show: showRefillModal }"
      :style="{ display: showRefillModal ? 'block' : 'none' }"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Refill Activity Logs</h5>
            <button type="button" class="btn-close" @click="showRefillModal = false"></button>
          </div>
          <div class="modal-body">
            <div class="alert alert-info">
              <strong>Info:</strong> This will generate activity logs for all existing historical data in your account. This process may take a few moments.
            </div>
            <p class="mb-0">
              The refill process will create activity logs for:
            </p>
            <ul class="mt-2">
              <li>All existing accounts</li>
              <li>All existing transactions</li>
              <li>All existing budgets</li>
              <li>All existing credit cards</li>
              <li>All existing debts and lends</li>
              <li>All existing assets and goals</li>
            </ul>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showRefillModal = false">
              Cancel
            </button>
            <button type="button" class="btn btn-success" @click="refillLogs" :disabled="refilling">
              <span v-if="refilling" class="spinner-border spinner-border-sm me-2" role="status"></span>
              {{ refilling ? 'Refilling...' : 'Refill Logs' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showRefillModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useActivityLogsStore } from '@/stores/activityLogs'
import { format } from 'date-fns'

const activityLogsStore = useActivityLogsStore()

const showCleanupModal = ref(false)
const showRefillModal = ref(false)
const cleanupDays = ref(90)
const summaryDays = ref(30)
const refilling = ref(false)

const localFilters = ref({
  module: null,
  action: null,
  startDate: null,
  endDate: null
})

const activityLogs = computed(() => activityLogsStore.allActivityLogs)
const summary = computed(() => activityLogsStore.summary)
const loading = computed(() => activityLogsStore.loading)
const pagination = computed(() => activityLogsStore.pagination)
const hasNextPage = computed(() => activityLogsStore.hasNextPage)
const hasPrevPage = computed(() => activityLogsStore.hasPrevPage)

onMounted(async () => {
  await refreshLogs()
  await activityLogsStore.fetchActivitySummary(summaryDays.value)
})

const applyFilters = async () => {
  activityLogsStore.setFilters(localFilters.value)
  await activityLogsStore.fetchActivityLogs()
}

const clearFilters = async () => {
  localFilters.value = {
    module: null,
    action: null,
    startDate: null,
    endDate: null
  }
  activityLogsStore.clearFilters()
  await activityLogsStore.fetchActivityLogs()
}

const refreshLogs = async () => {
  await activityLogsStore.fetchActivityLogs()
}

const nextPage = async () => {
  await activityLogsStore.nextPage()
}

const prevPage = async () => {
  await activityLogsStore.prevPage()
}

const cleanupLogs = async () => {
  if (cleanupDays.value < 30) {
    alert('Minimum cleanup period is 30 days')
    return
  }

  if (confirm(`Are you sure you want to delete all activity logs older than ${cleanupDays.value} days? This action cannot be undone.`)) {
    try {
      await activityLogsStore.cleanupOldLogs(cleanupDays.value)
      showCleanupModal.value = false
      alert('Old activity logs deleted successfully')
      await refreshLogs()
      await activityLogsStore.fetchActivitySummary(summaryDays.value)
    } catch (error) {
      alert('Failed to cleanup logs: ' + error.message)
    }
  }
}

const refillLogs = async () => {
  try {
    refilling.value = true
    const result = await activityLogsStore.backfillActivityLogs({})
    showRefillModal.value = false

    // Show success message with details
    const message = `Activity logs refilled successfully!\n\nTotal Records: ${result.totalRecords || 0}\nLogs Created: ${result.logsCreated || 0}\nSkipped: ${result.skipped || 0}`
    alert(message)

    await refreshLogs()
    await activityLogsStore.fetchActivitySummary(summaryDays.value)
  } catch (error) {
    alert('Failed to refill activity logs: ' + error.message)
  } finally {
    refilling.value = false
  }
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  try {
    return format(new Date(dateString), 'MMM dd, yyyy HH:mm:ss')
  } catch (error) {
    return dateString
  }
}

const getActionBadgeClass = (action) => {
  const classes = {
    create: 'bg-success',
    update: 'bg-primary',
    delete: 'bg-danger',
    view: 'bg-info',
    login: 'bg-success',
    logout: 'bg-secondary',
    export: 'bg-warning',
    import: 'bg-warning'
  }
  return classes[action] || 'bg-secondary'
}
</script>

<style scoped>
.activity-logs-view {
  animation: fadeIn 0.3s ease-in;
}

.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: bold;
  color: #333;
}

.stat-label {
  color: #666;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-icon.purple {
  color: #8B5CF6;
}

.stat-icon.blue {
  color: #3B82F6;
}

.stat-icon.green {
  color: #10B981;
}

.stat-icon.orange {
  color: #F59E0B;
}

.text-purple {
  color: #8B5CF6;
}

.table th {
  background-color: #f8f9fa;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  color: #6c757d;
}

.table tbody tr {
  transition: background-color 0.2s;
}

.table tbody tr:hover {
  background-color: #f8f9fa;
}

.modal.show {
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
