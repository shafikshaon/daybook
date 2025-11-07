<template>
  <div class="scheduled-transactions-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Scheduled Transactions</h1>
      <div class="d-flex gap-2">
        <button class="btn btn-success" @click="processScheduled" :disabled="processing">
          <span v-if="processing">Processing...</span>
          <span v-else>Process Scheduled</span>
        </button>
        <button class="btn btn-primary" @click="showAddModal = true">+ Add Scheduled Transaction</button>
      </div>
    </div>

    <!-- Info Alert -->
    <div class="alert alert-info mb-4">
      <h6 class="alert-heading">Automated Transaction Scheduling</h6>
      <p class="mb-0">
        Scheduled transactions are automatically created based on their frequency.
        Click "Process Scheduled" to generate all pending transactions, or they will be created automatically at their scheduled time.
      </p>
    </div>

    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon purple">📅</div>
          <div class="stat-value">{{ activeSchedules.length }}</div>
          <div class="stat-label">Active Schedules</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon blue">🔄</div>
          <div class="stat-value">{{ recurringTransactions.length }}</div>
          <div class="stat-label">Total Schedules</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon orange">⏸️</div>
          <div class="stat-value">{{ disabledSchedules.length }}</div>
          <div class="stat-label">Disabled</div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h5 class="mb-0">All Scheduled Transactions</h5>
      </div>
      <div class="card-body p-0">
        <div v-if="recurringTransactions.length === 0" class="p-4 text-center text-muted">
          No scheduled transactions yet. Click "Add Scheduled Transaction" to create one.
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Status</th>
                <th>Description</th>
                <th>Type</th>
                <th>Amount</th>
                <th>Frequency</th>
                <th>Start Date</th>
                <th>End Date</th>
                <th>Last Processed</th>
                <th class="text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="schedule in recurringTransactions" :key="schedule.id">
                <td>
                  <span class="badge" :class="schedule.enabled ? 'bg-success' : 'bg-secondary'">
                    {{ schedule.enabled ? 'Active' : 'Disabled' }}
                  </span>
                </td>
                <td class="fw-semibold">{{ schedule.transactionTemplate?.description || 'N/A' }}</td>
                <td>
                  <span class="badge" :class="getTypeBadgeClass(schedule.transactionTemplate?.type)">
                    {{ formatType(schedule.transactionTemplate?.type) }}
                  </span>
                </td>
                <td class="fw-bold" :class="getAmountClass(schedule.transactionTemplate?.type)">
                  {{ formatCurrency(schedule.transactionTemplate?.amount) }}
                </td>
                <td>{{ formatFrequency(schedule.frequency) }}</td>
                <td>{{ formatDate(schedule.startDate) }}</td>
                <td>{{ schedule.endDate ? formatDate(schedule.endDate) : 'Never' }}</td>
                <td>{{ schedule.lastProcessed ? formatDate(schedule.lastProcessed) : 'Not yet' }}</td>
                <td class="text-center">
                  <button
                    class="btn btn-sm me-1"
                    :class="schedule.enabled ? 'btn-warning' : 'btn-success'"
                    @click="toggleEnabled(schedule)"
                  >
                    {{ schedule.enabled ? 'Disable' : 'Enable' }}
                  </button>
                  <button class="btn btn-sm btn-info me-1" @click="editSchedule(schedule)">
                    Edit
                  </button>
                  <button class="btn btn-sm btn-delete" @click="deleteSchedule(schedule.id)">
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddModal">
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ isEditing ? 'Edit' : 'Add' }} Scheduled Transaction</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveSchedule">
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Description *</label>
                  <input type="text" class="form-control" v-model="form.description" required />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Type *</label>
                  <select class="form-select" v-model="form.type" required @change="onTypeChange">
                    <option value="income">Income</option>
                    <option value="expense">Expense</option>
                    <option value="transfer">Transfer</option>
                  </select>
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Amount *</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="form.amount" required />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Category *</label>
                  <select class="form-select" v-model="form.categoryId" required>
                    <option v-for="cat in filteredCategories" :key="cat.id" :value="cat.id">
                      {{ cat.icon }} {{ cat.name }}
                    </option>
                  </select>
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Account *</label>
                  <select class="form-select" v-model="form.accountId" required>
                    <option value="">Select Account</option>
                    <option v-for="account in accounts" :key="account.id" :value="account.id">
                      {{ account.name }} ({{ formatCurrency(account.balance) }})
                    </option>
                  </select>
                </div>
                <div class="col-md-6 mb-3" v-if="form.type === 'transfer'">
                  <label class="form-label">To Account *</label>
                  <select class="form-select" v-model="form.toAccountId" :required="form.type === 'transfer'">
                    <option value="">Select Destination Account</option>
                    <option v-for="account in accounts" :key="account.id" :value="account.id">
                      {{ account.name }} ({{ formatCurrency(account.balance) }})
                    </option>
                  </select>
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Frequency *</label>
                  <select class="form-select" v-model="form.frequency" required>
                    <option value="daily">Daily</option>
                    <option value="weekly">Weekly</option>
                    <option value="biweekly">Bi-Weekly (Every 2 weeks)</option>
                    <option value="monthly">Monthly</option>
                    <option value="quarterly">Quarterly (Every 3 months)</option>
                    <option value="yearly">Yearly</option>
                  </select>
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Start Date *</label>
                  <input type="date" class="form-control" v-model="form.startDate" required />
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">End Date (Optional)</label>
                  <input type="date" class="form-control" v-model="form.endDate" />
                  <small class="text-muted">Leave empty for ongoing schedule</small>
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Status</label>
                  <div class="form-check form-switch mt-2">
                    <input class="form-check-input" type="checkbox" v-model="form.enabled" id="enabledSwitch">
                    <label class="form-check-label" for="enabledSwitch">
                      {{ form.enabled ? 'Enabled' : 'Disabled' }}
                    </label>
                  </div>
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Notes (Optional)</label>
                <textarea class="form-control" v-model="form.notes" rows="2"></textarea>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
                <button type="submit" class="btn btn-primary">{{ isEditing ? 'Update' : 'Create' }}</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const processing = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const form = ref({
  description: '',
  type: 'expense',
  amount: 0,
  categoryId: '',
  accountId: '',
  toAccountId: '',
  frequency: 'monthly',
  startDate: new Date().toISOString().split('T')[0],
  endDate: '',
  enabled: true,
  notes: ''
})

const recurringTransactions = computed(() => transactionsStore.recurringTransactions)
const activeSchedules = computed(() => recurringTransactions.value.filter(rt => rt.enabled))
const disabledSchedules = computed(() => recurringTransactions.value.filter(rt => !rt.enabled))
const accounts = computed(() => accountsStore.accounts)
const categories = computed(() => transactionsStore.categories)

const filteredCategories = computed(() => {
  return categories.value.filter(cat => cat.type === form.value.type)
})

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => date ? new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }) : '-'

const formatType = (type) => {
  const types = { income: 'Income', expense: 'Expense', transfer: 'Transfer' }
  return types[type] || type
}

const formatFrequency = (frequency) => {
  const frequencies = {
    daily: 'Daily',
    weekly: 'Weekly',
    biweekly: 'Bi-Weekly',
    monthly: 'Monthly',
    quarterly: 'Quarterly',
    yearly: 'Yearly'
  }
  return frequencies[frequency] || frequency
}

const getTypeBadgeClass = (type) => {
  if (type === 'income') return 'bg-success'
  if (type === 'expense') return 'bg-danger'
  if (type === 'transfer') return 'bg-info'
  return 'bg-secondary'
}

const getAmountClass = (type) => {
  if (type === 'income') return 'text-success'
  if (type === 'expense') return 'text-danger'
  return 'text-info'
}

const onTypeChange = () => {
  // Reset category when type changes
  form.value.categoryId = ''
  // Reset toAccountId if not transfer
  if (form.value.type !== 'transfer') {
    form.value.toAccountId = ''
  }
}

const processScheduled = async () => {
  processing.value = true
  try {
    const result = await transactionsStore.processRecurringTransactions()
    await transactionsStore.fetchRecurringTransactions()
    success(`Successfully created ${result.created || 0} transaction(s)`)
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error processing scheduled transactions')
  } finally {
    processing.value = false
  }
}

const toggleEnabled = async (schedule) => {
  try {
    // Build transaction template without empty UUID fields
    const transactionTemplate = {
      description: schedule.transactionTemplate.description,
      type: schedule.transactionTemplate.type,
      amount: schedule.transactionTemplate.amount,
      categoryId: schedule.transactionTemplate.categoryId,
      accountId: schedule.transactionTemplate.accountId,
      notes: schedule.transactionTemplate.notes || ''
    }

    // Only include toAccountId if it exists and has a value
    if (schedule.transactionTemplate.toAccountId) {
      transactionTemplate.toAccountId = schedule.transactionTemplate.toAccountId
    }

    const updatedData = {
      transactionTemplate,
      frequency: schedule.frequency,
      startDate: schedule.startDate,
      endDate: schedule.endDate,
      enabled: !schedule.enabled
    }
    await transactionsStore.updateRecurringTransaction(schedule.id, updatedData)
    success(`Schedule ${updatedData.enabled ? 'enabled' : 'disabled'} successfully`)
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error updating schedule')
  }
}

const editSchedule = (schedule) => {
  isEditing.value = true
  editingId.value = schedule.id
  form.value = {
    description: schedule.transactionTemplate?.description || '',
    type: schedule.transactionTemplate?.type || 'expense',
    amount: schedule.transactionTemplate?.amount || 0,
    categoryId: schedule.transactionTemplate?.categoryId || '',
    accountId: schedule.transactionTemplate?.accountId || '',
    toAccountId: schedule.transactionTemplate?.toAccountId || '',
    frequency: schedule.frequency,
    startDate: schedule.startDate ? new Date(schedule.startDate).toISOString().split('T')[0] : '',
    endDate: schedule.endDate ? new Date(schedule.endDate).toISOString().split('T')[0] : '',
    enabled: schedule.enabled,
    notes: schedule.transactionTemplate?.notes || ''
  }
  showAddModal.value = true
}

const deleteSchedule = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Scheduled Transaction',
    message: 'Are you sure you want to delete this scheduled transaction? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await transactionsStore.deleteRecurringTransaction(id)
      success('Scheduled transaction deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting scheduled transaction')
    }
  }
}

const saveSchedule = async () => {
  try {
    // Build transaction template without empty UUID fields
    const transactionTemplate = {
      description: form.value.description,
      type: form.value.type,
      amount: form.value.amount,
      categoryId: form.value.categoryId,
      accountId: form.value.accountId,
      notes: form.value.notes
    }

    // Only include toAccountId if it's a transfer and has a value
    if (form.value.type === 'transfer' && form.value.toAccountId) {
      transactionTemplate.toAccountId = form.value.toAccountId
    }

    const scheduleData = {
      transactionTemplate,
      frequency: form.value.frequency,
      startDate: form.value.startDate,
      endDate: form.value.endDate || null,
      enabled: form.value.enabled
    }

    if (isEditing.value) {
      await transactionsStore.updateRecurringTransaction(editingId.value, scheduleData)
      success('Scheduled transaction updated successfully')
    } else {
      await transactionsStore.createRecurringTransaction(scheduleData)
      success('Scheduled transaction created successfully')
    }

    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving scheduled transaction')
  }
}

const closeModal = () => {
  showAddModal.value = false
  isEditing.value = false
  editingId.value = null
  form.value = {
    description: '',
    type: 'expense',
    amount: 0,
    categoryId: '',
    accountId: '',
    toAccountId: '',
    frequency: 'monthly',
    startDate: new Date().toISOString().split('T')[0],
    endDate: '',
    enabled: true,
    notes: ''
  }
}

onMounted(async () => {
  await transactionsStore.fetchRecurringTransactions()
  await accountsStore.fetchAccounts()
})
</script>

<style scoped>
/* Professional button styles */
.btn-delete {
  color: #991b1b;
  border-color: #ef4444;
  background-color: #fef2f2;
}

.btn-delete:hover {
  color: #ffffff;
  background-color: #ef4444;
  border-color: #dc2626;
}

/* Dark mode support */
.dark-mode .btn-delete {
  color: #fca5a5;
  border-color: #ef4444;
  background-color: #5f1e1e;
}

.dark-mode .btn-delete:hover {
  color: #ffffff;
  background-color: #ef4444;
  border-color: #f87171;
}

/* Table styling */
.table tbody tr:hover {
  background-color: rgba(139, 92, 246, 0.05);
}

.dark-mode .table tbody tr:hover {
  background-color: rgba(139, 92, 246, 0.1);
}
</style>
