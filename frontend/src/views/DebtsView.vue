<template>
  <div class="debts-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Debts</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Debt</button>
    </div>

    <!-- Debt Summary -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon red">💳</div>
          <div class="stat-value">{{ formatCurrency(debtsStore.totalDebtAmount) }}</div>
          <div class="stat-label">Total Owed</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">📊</div>
          <div class="stat-value">{{ formatCurrency(debtsStore.totalOriginalDebt) }}</div>
          <div class="stat-label">Original Debt</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">✅</div>
          <div class="stat-value">{{ formatCurrency(debtsStore.totalPaidAmount) }}</div>
          <div class="stat-label">Total Paid</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">📝</div>
          <div class="stat-value">{{ debtsStore.activeDebts.length }}</div>
          <div class="stat-label">Active Debts</div>
        </div>
      </div>
    </div>

    <!-- Overdue Alerts -->
    <div v-if="debtsStore.overdueDebts.length > 0" class="alert alert-danger mb-4">
      <h6 class="alert-heading">⚠️ Overdue Debts</h6>
      <ul class="mb-0">
        <li v-for="debt in debtsStore.overdueDebts" :key="debt.id">
          {{ debt.creditorName }}: {{ formatCurrency(debt.remainingAmount) }} - Due {{ formatDate(debt.dueDate) }}
        </li>
      </ul>
    </div>

    <!-- Filter Tabs -->
    <ul class="nav nav-tabs mb-3">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'all' }" @click="filter = 'all'" href="javascript:void(0)">
          All ({{ debtsStore.allDebts.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'active' }" @click="filter = 'active'" href="javascript:void(0)">
          Active ({{ debtsStore.activeDebts.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'partially_paid' }" @click="filter = 'partially_paid'" href="javascript:void(0)">
          Partial ({{ debtsStore.partiallyPaidDebts.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'fully_paid' }" @click="filter = 'fully_paid'" href="javascript:void(0)">
          Paid ({{ debtsStore.fullyPaidDebts.length }})
        </a>
      </li>
    </ul>

    <!-- Debts List -->
    <div class="row g-3">
      <div v-for="debt in filteredDebts" :key="debt.id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-2">
              <h5 class="card-title mb-0">{{ debt.creditorName }}</h5>
              <span class="badge" :class="getStatusClass(debt.status)">{{ formatStatus(debt.status) }}</span>
            </div>

            <p class="text-muted mb-2" v-if="debt.description">{{ debt.description }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Remaining</span>
                <strong>{{ formatCurrency(debt.remainingAmount) }}</strong>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Original</span>
                <span>{{ formatCurrency(debt.originalAmount) }}</span>
              </div>
              <div v-if="debt.accountName" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Account</span>
                <span>{{ debt.accountName }}</span>
              </div>
              <div v-if="debt.dueDate" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Due Date</span>
                <span :class="{ 'text-danger': isOverdue(debt.dueDate) }">
                  {{ formatDate(debt.dueDate) }}
                </span>
              </div>
              <div v-if="debt.interestRate" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Interest Rate</span>
                <span>{{ debt.interestRate }}%</span>
              </div>
            </div>

            <div class="progress mb-3" style="height: 8px;">
              <div
                class="progress-bar bg-success"
                :style="{ width: getProgress(debt) + '%' }"
              ></div>
            </div>

            <div class="d-flex justify-content-between gap-2">
              <button
                v-if="debt.status !== 'fully_paid'"
                class="btn btn-sm btn-success flex-grow-1"
                @click="openPaymentModal(debt)"
              >
                Pay
              </button>
              <button class="btn btn-sm btn-outline-primary" @click="editDebt(debt)">Edit</button>
              <button class="btn btn-sm btn-outline-danger" @click="deleteDebt(debt.id)">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Debt' : 'Add Debt' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveDebt">
              <div class="mb-3">
                <label class="form-label">Creditor Name *</label>
                <input type="text" class="form-control" v-model="form.creditorName" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.originalAmount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Borrowed Date *</label>
                <input type="date" class="form-control" v-model="form.borrowedDate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Due Date</label>
                <input type="date" class="form-control" v-model="form.dueDate" />
              </div>
              <div class="mb-3">
                <label class="form-label">Account</label>
                <select class="form-select" v-model="form.accountId">
                  <option value="">None (doesn't affect balance)</option>
                  <option v-for="account in accounts" :key="account.id" :value="account.id">
                    {{ account.name }} ({{ formatCurrency(account.balance) }})
                  </option>
                </select>
                <div class="form-text">Select an account if this debt affects your balance</div>
              </div>
              <div class="mb-3" v-if="!showEditModal && form.accountId">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" v-model="form.isInitial" id="isInitial">
                  <label class="form-check-label" for="isInitial">
                    This is a pre-existing debt (doesn't affect account balance)
                  </label>
                </div>
              </div>
              <div class="mb-3">
                <label class="form-label">Interest Rate (%)</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.interestRate" />
              </div>
              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea class="form-control" v-model="form.description" rows="3"></textarea>
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
                <button type="submit" class="btn btn-primary">{{ showEditModal ? 'Update' : 'Create' }}</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Payment Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showPaymentModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showPaymentModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Record Payment</h5>
            <button type="button" class="btn-close" @click="closePaymentModal"></button>
          </div>
          <div class="modal-body">
            <div class="alert alert-info mb-3">
              <strong>{{ selectedDebt?.creditorName }}</strong><br>
              Remaining: {{ formatCurrency(selectedDebt?.remainingAmount) }}
            </div>
            <form @submit.prevent="recordPayment">
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="paymentForm.amount"
                  :max="selectedDebt?.remainingAmount"
                  required
                />
              </div>
              <div class="mb-3">
                <label class="form-label">Account *</label>
                <select class="form-select" v-model="paymentForm.accountId" required>
                  <option value="">Select account...</option>
                  <option v-for="account in accounts" :key="account.id" :value="account.id">
                    {{ account.name }} ({{ formatCurrency(account.balance) }})
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Payment Date *</label>
                <input type="date" class="form-control" v-model="paymentForm.paymentDate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Description</label>
                <input type="text" class="form-control" v-model="paymentForm.description" />
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closePaymentModal">Cancel</button>
                <button type="submit" class="btn btn-success">Record Payment</button>
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
import { useDebtsStore } from '@/stores/debts'
import { useAccountsStore } from '@/stores/accounts'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'
import { formatStatus } from '@/utils/textUtils'

const debtsStore = useDebtsStore()
const accountsStore = useAccountsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const filter = ref('all')
const showAddModal = ref(false)
const showEditModal = ref(false)
const showPaymentModal = ref(false)
const editingDebt = ref(null)
const selectedDebt = ref(null)

const form = ref({
  creditorName: '',
  originalAmount: 0,
  borrowedDate: new Date().toISOString().split('T')[0],
  dueDate: '',
  accountId: '',
  interestRate: null,
  description: '',
  isInitial: false
})

const paymentForm = ref({
  amount: 0,
  accountId: '',
  paymentDate: new Date().toISOString().split('T')[0],
  description: ''
})

const accounts = computed(() => accountsStore.allAccounts)

const filteredDebts = computed(() => {
  if (filter.value === 'all') return debtsStore.allDebts
  if (filter.value === 'active') return debtsStore.activeDebts
  if (filter.value === 'partially_paid') return debtsStore.partiallyPaidDebts
  if (filter.value === 'fully_paid') return debtsStore.fullyPaidDebts
  return debtsStore.allDebts
})

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => debtsStore.formatDate(date)
const isOverdue = (date) => debtsStore.isOverdue(date)

const getStatusClass = (status) => {
  if (status === 'active') return 'bg-danger'
  if (status === 'partially_paid') return 'bg-warning'
  if (status === 'fully_paid') return 'bg-success'
  return 'bg-secondary'
}

const getProgress = (debt) => {
  const paid = debt.originalAmount - debt.remainingAmount
  return Math.round((paid / debt.originalAmount) * 100)
}

const editDebt = (debt) => {
  editingDebt.value = debt
  form.value = {
    creditorName: debt.creditorName,
    originalAmount: debt.originalAmount,
    borrowedDate: debt.borrowedDate ? debt.borrowedDate.split('T')[0] : '',
    dueDate: debt.dueDate ? debt.dueDate.split('T')[0] : '',
    interestRate: debt.interestRate,
    description: debt.description
  }
  showEditModal.value = true
}

const deleteDebt = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Debt',
    message: 'Are you sure you want to delete this debt? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await debtsStore.deleteDebt(id)
      success('Debt deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting debt')
    }
  }
}

const saveDebt = async () => {
  try {
    const debtData = { ...form.value }

    // Clean up empty values
    if (!debtData.dueDate) delete debtData.dueDate
    if (!debtData.accountId) delete debtData.accountId
    if (!debtData.interestRate) delete debtData.interestRate

    if (showEditModal.value) {
      await debtsStore.updateDebt(editingDebt.value.id, debtData)
      success('Debt updated successfully')
    } else {
      await debtsStore.createDebt(debtData)
      success('Debt created successfully')
      await accountsStore.fetchAccounts() // Refresh accounts to update balance
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving debt')
  }
}

const openPaymentModal = (debt) => {
  selectedDebt.value = debt
  paymentForm.value = {
    amount: debt.remainingAmount,
    accountId: '',
    paymentDate: new Date().toISOString().split('T')[0],
    description: `Payment to ${debt.creditorName}`
  }
  showPaymentModal.value = true
}

const recordPayment = async () => {
  try {
    await debtsStore.recordPayment(selectedDebt.value.id, paymentForm.value)
    success('Payment recorded successfully')
    await accountsStore.fetchAccounts() // Refresh accounts to update balance
    closePaymentModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error recording payment')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  form.value = {
    creditorName: '',
    originalAmount: 0,
    borrowedDate: new Date().toISOString().split('T')[0],
    dueDate: '',
    accountId: '',
    interestRate: null,
    description: '',
    isInitial: false
  }
  editingDebt.value = null
}

const closePaymentModal = () => {
  showPaymentModal.value = false
  selectedDebt.value = null
  paymentForm.value = {
    amount: 0,
    accountId: '',
    paymentDate: new Date().toISOString().split('T')[0],
    description: ''
  }
}

onMounted(async () => {
  await Promise.all([
    debtsStore.fetchDebts(),
    accountsStore.fetchAccounts()
  ])
})
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.stat-card {
  padding: 1.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.stat-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #2c3e50;
}

.stat-label {
  color: #7f8c8d;
  font-size: 0.9rem;
}

.text-purple { color: #6f42c1; }
.stat-icon.purple { color: #6f42c1; }
.stat-icon.red { color: #dc3545; }
.stat-icon.green { color: #28a745; }
.stat-icon.blue { color: #007bff; }

.card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.nav-link {
  cursor: pointer;
}
</style>
