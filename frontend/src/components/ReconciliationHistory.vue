<template>
  <div class="reconciliation-history">
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h5 class="mb-0">Reconciliation History</h5>
      <button
        class="btn btn-sm btn-outline-primary"
        @click="loadReconciliations"
        :disabled="loading"
      >
        <span v-if="loading">Loading...</span>
        <span v-else>Refresh</span>
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="row g-3 mb-4" v-if="stats">
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h3 class="mb-1">{{ stats.totalReconciliations }}</h3>
            <small class="text-muted">Total</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h3 class="mb-1 text-success">{{ stats.completedReconciliations }}</h3>
            <small class="text-muted">Completed</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h3 class="mb-1 text-warning">{{ stats.discrepancyReconciliations }}</h3>
            <small class="text-muted">Discrepancies</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card">
          <div class="card-body text-center">
            <h3 class="mb-1 text-info">{{ formatCurrency(stats.averageDifference || 0) }}</h3>
            <small class="text-muted">Avg. Difference</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Reconciliation List -->
    <div v-if="loading && reconciliations.length === 0" class="text-center py-4">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-else-if="reconciliations.length === 0" class="text-center text-muted py-4">
      <p>No reconciliation history yet</p>
    </div>

    <div v-else class="table-responsive">
      <table class="table table-hover">
        <thead>
          <tr>
            <th>Date</th>
            <th>Statement Balance</th>
            <th>Book Balance</th>
            <th>Difference</th>
            <th>Status</th>
            <th>Notes</th>
            <th class="text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="reconciliation in reconciliations" :key="reconciliation.id">
            <td>{{ formatDate(reconciliation.reconciliationDate) }}</td>
            <td>{{ formatCurrency(reconciliation.statementBalance) }}</td>
            <td>{{ formatCurrency(reconciliation.bookBalance) }}</td>
            <td :class="getDifferenceClass(reconciliation.difference)">
              {{ formatCurrency(Math.abs(reconciliation.difference)) }}
              <span v-if="reconciliation.difference !== 0">
                {{ reconciliation.difference > 0 ? '↑' : '↓' }}
              </span>
            </td>
            <td>
              <span class="badge" :class="getStatusBadgeClass(reconciliation.status)">
                {{ formatStatus(reconciliation.status) }}
              </span>
            </td>
            <td>
              <small class="text-muted">{{ reconciliation.notes || '-' }}</small>
            </td>
            <td class="text-center">
              <button
                class="btn btn-sm btn-outline-primary me-1"
                @click="viewDetails(reconciliation)"
                title="View Details"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z"/>
                  <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0"/>
                </svg>
              </button>
              <button
                class="btn btn-sm btn-outline-danger"
                @click="deleteReconciliation(reconciliation)"
                title="Delete"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                  <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Details Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showDetailsModal }"
      tabindex="-1"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showDetailsModal"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Reconciliation Details</h5>
            <button type="button" class="btn-close" @click="showDetailsModal = false"></button>
          </div>
          <div class="modal-body" v-if="selectedReconciliation">
            <div class="row mb-4">
              <div class="col-md-6">
                <strong>Date:</strong> {{ formatDate(selectedReconciliation.reconciliationDate) }}<br>
                <strong>Statement Balance:</strong> {{ formatCurrency(selectedReconciliation.statementBalance) }}<br>
                <strong>Book Balance:</strong> {{ formatCurrency(selectedReconciliation.bookBalance) }}<br>
                <strong>Difference:</strong>
                <span :class="getDifferenceClass(selectedReconciliation.difference)">
                  {{ formatCurrency(Math.abs(selectedReconciliation.difference)) }}
                </span>
              </div>
              <div class="col-md-6">
                <strong>Status:</strong>
                <span class="badge ms-2" :class="getStatusBadgeClass(selectedReconciliation.status)">
                  {{ formatStatus(selectedReconciliation.status) }}
                </span><br>
                <strong>Notes:</strong> {{ selectedReconciliation.notes || '-' }}
              </div>
            </div>

            <h6>Reconciled Transactions ({{ selectedReconciliation.transactions?.length || 0 }})</h6>
            <div v-if="!selectedReconciliation.transactions || selectedReconciliation.transactions.length === 0" class="text-muted">
              No transactions linked to this reconciliation
            </div>
            <div v-else class="table-responsive">
              <table class="table table-sm">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Description</th>
                    <th>Category</th>
                    <th class="text-end">Amount</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="rt in selectedReconciliation.transactions" :key="rt.id">
                    <td>{{ formatDate(rt.transaction?.date) }}</td>
                    <td>{{ rt.transaction?.description || '-' }}</td>
                    <td>
                      <span class="badge bg-secondary">{{ rt.transaction?.categoryId }}</span>
                    </td>
                    <td class="text-end" :class="rt.transaction?.type === 'expense' ? 'text-danger' : 'text-success'">
                      {{ rt.transaction?.type === 'expense' ? '-' : '+' }}{{ formatCurrency(rt.transaction?.amount || 0) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showDetailsModal = false">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import apiService from '@/services/api-backend'
import { useSettingsStore } from '@/stores/settings'

const props = defineProps({
  accountId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['refresh'])

const settingsStore = useSettingsStore()

const loading = ref(false)
const reconciliations = ref([])
const stats = ref(null)
const showDetailsModal = ref(false)
const selectedReconciliation = ref(null)

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const formatStatus = (status) => {
  const statusMap = {
    'pending': 'Pending',
    'completed': 'Completed',
    'discrepancy': 'Discrepancy'
  }
  return statusMap[status] || status
}

const getStatusBadgeClass = (status) => {
  const classMap = {
    'pending': 'bg-secondary',
    'completed': 'bg-success',
    'discrepancy': 'bg-warning'
  }
  return classMap[status] || 'bg-secondary'
}

const getDifferenceClass = (difference) => {
  if (difference === 0) return 'text-success fw-bold'
  if (Math.abs(difference) < 1) return 'text-success'
  return 'text-warning'
}

const loadReconciliations = async () => {
  loading.value = true
  try {
    const [reconciliationsResponse, statsResponse] = await Promise.all([
      apiService.reconciliation.getByAccount(props.accountId),
      apiService.reconciliation.getStats(props.accountId)
    ])
    reconciliations.value = reconciliationsResponse.data || []
    stats.value = statsResponse.data || null
  } catch (error) {
    console.error('Failed to load reconciliations:', error)
  } finally {
    loading.value = false
  }
}

const viewDetails = async (reconciliation) => {
  try {
    const response = await apiService.reconciliation.getById(reconciliation.id)
    selectedReconciliation.value = response.data
    showDetailsModal.value = true
  } catch (error) {
    console.error('Failed to load reconciliation details:', error)
    alert('Failed to load reconciliation details')
  }
}

const deleteReconciliation = async (reconciliation) => {
  if (!confirm(`Are you sure you want to delete this reconciliation from ${formatDate(reconciliation.reconciliationDate)}?`)) {
    return
  }

  try {
    await apiService.reconciliation.delete(reconciliation.id)
    await loadReconciliations()
    emit('refresh')
  } catch (error) {
    console.error('Failed to delete reconciliation:', error)
    alert('Failed to delete reconciliation')
  }
}

onMounted(() => {
  loadReconciliations()
})

// Expose method for parent component to refresh
defineExpose({
  loadReconciliations
})
</script>

<style scoped>
.card {
  border: 1px solid #dee2e6;
}

.card-body {
  padding: 1rem;
}

.card-body h3 {
  margin-bottom: 0.25rem;
  font-size: 1.5rem;
}
</style>
