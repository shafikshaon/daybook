<template>
  <div class="lends-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-blue">Lends</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Lend</button>
    </div>

    <!-- Lend Summary -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(lendsStore.totalLendAmount) }}</div>
          <div class="stat-label">Total Owed to Me</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">📊</div>
          <div class="stat-value">{{ formatCurrency(lendsStore.totalOriginalLend) }}</div>
          <div class="stat-label">Original Lent</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">✅</div>
          <div class="stat-value">{{ formatCurrency(lendsStore.totalReceivedAmount) }}</div>
          <div class="stat-label">Total Received</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon orange">📝</div>
          <div class="stat-value">{{ lendsStore.activeLends.length }}</div>
          <div class="stat-label">Active Lends</div>
        </div>
      </div>
    </div>

    <!-- Overdue Alerts -->
    <div v-if="lendsStore.overdueLends.length > 0" class="alert alert-warning mb-4">
      <h6 class="alert-heading">⚠️ Overdue Lends</h6>
      <ul class="mb-0">
        <li v-for="lend in lendsStore.overdueLends" :key="lend.id">
          {{ lend.debtorName }}: {{ formatCurrency(lend.remainingAmount) }} - Due {{ formatDate(lend.dueDate) }}
        </li>
      </ul>
    </div>

    <!-- Filter Tabs -->
    <ul class="nav nav-tabs mb-3">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'all' }" @click="filter = 'all'" href="javascript:void(0)">
          All ({{ lendsStore.allLends.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'active' }" @click="filter = 'active'" href="javascript:void(0)">
          Active ({{ lendsStore.activeLends.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'partially_received' }" @click="filter = 'partially_received'" href="javascript:void(0)">
          Partial ({{ lendsStore.partiallyReceivedLends.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'fully_received' }" @click="filter = 'fully_received'" href="javascript:void(0)">
          Received ({{ lendsStore.fullyReceivedLends.length }})
        </a>
      </li>
    </ul>

    <!-- Lends List -->
    <div class="row g-3">
      <div v-for="lend in filteredLends" :key="lend.id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-2">
              <h5 class="card-title mb-0">{{ lend.debtorName }}</h5>
              <span class="badge" :class="getStatusClass(lend.status)">{{ lend.status }}</span>
            </div>

            <p class="text-muted mb-2" v-if="lend.description">{{ lend.description }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Remaining</span>
                <strong>{{ formatCurrency(lend.remainingAmount) }}</strong>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Original</span>
                <span>{{ formatCurrency(lend.originalAmount) }}</span>
              </div>
              <div v-if="lend.accountName" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Account</span>
                <span>{{ lend.accountName }}</span>
              </div>
              <div v-if="lend.dueDate" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Due Date</span>
                <span :class="{ 'text-warning': isOverdue(lend.dueDate) }">
                  {{ formatDate(lend.dueDate) }}
                </span>
              </div>
              <div v-if="lend.interestRate" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Interest Rate</span>
                <span>{{ lend.interestRate }}%</span>
              </div>
            </div>

            <div class="progress mb-3" style="height: 8px;">
              <div
                class="progress-bar bg-primary"
                :style="{ width: getProgress(lend) + '%' }"
              ></div>
            </div>

            <div class="d-flex justify-content-between gap-2">
              <button
                v-if="lend.status !== 'fully_received'"
                class="btn btn-sm btn-success flex-grow-1"
                @click="openPaymentModal(lend)"
              >
                Receive
              </button>
              <button class="btn btn-sm btn-outline-primary" @click="editLend(lend)">Edit</button>
              <button class="btn btn-sm btn-outline-danger" @click="deleteLend(lend.id)">Delete</button>
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
            <h5 class="modal-title">{{ showEditModal ? 'Edit Lend' : 'Add Lend' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveLend">
              <div class="mb-3">
                <label class="form-label">Debtor Name *</label>
                <input type="text" class="form-control" v-model="form.debtorName" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.originalAmount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Lent Date *</label>
                <input type="date" class="form-control" v-model="form.lentDate" required />
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
                <div class="form-text">Select an account if this lend affects your balance</div>
              </div>
              <div class="mb-3" v-if="!showEditModal && form.accountId">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" v-model="form.isInitial" id="isInitial">
                  <label class="form-check-label" for="isInitial">
                    This is a pre-existing lend (doesn't affect account balance)
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
            <h5 class="modal-title">Record Payment Received</h5>
            <button type="button" class="btn-close" @click="closePaymentModal"></button>
          </div>
          <div class="modal-body">
            <div class="alert alert-info mb-3">
              <strong>{{ selectedLend?.debtorName }}</strong><br>
              Remaining: {{ formatCurrency(selectedLend?.remainingAmount) }}
            </div>
            <form @submit.prevent="recordPayment">
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="paymentForm.amount"
                  :max="selectedLend?.remainingAmount"
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
import { useLendsStore } from '@/stores/lends'
import { useAccountsStore } from '@/stores/accounts'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const lendsStore = useLendsStore()
const accountsStore = useAccountsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const filter = ref('all')
const showAddModal = ref(false)
const showEditModal = ref(false)
const showPaymentModal = ref(false)
const editingLend = ref(null)
const selectedLend = ref(null)

const form = ref({
  debtorName: '',
  originalAmount: 0,
  lentDate: new Date().toISOString().split('T')[0],
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

const filteredLends = computed(() => {
  if (filter.value === 'all') return lendsStore.allLends
  if (filter.value === 'active') return lendsStore.activeLends
  if (filter.value === 'partially_received') return lendsStore.partiallyReceivedLends
  if (filter.value === 'fully_received') return lendsStore.fullyReceivedLends
  return lendsStore.allLends
})

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => lendsStore.formatDate(date)
const isOverdue = (date) => lendsStore.isOverdue(date)

const getStatusClass = (status) => {
  if (status === 'active') return 'bg-primary'
  if (status === 'partially_received') return 'bg-warning'
  if (status === 'fully_received') return 'bg-success'
  return 'bg-secondary'
}

const getProgress = (lend) => {
  const received = lend.originalAmount - lend.remainingAmount
  return Math.round((received / lend.originalAmount) * 100)
}

const editLend = (lend) => {
  editingLend.value = lend
  form.value = {
    debtorName: lend.debtorName,
    dueDate: lend.dueDate ? lend.dueDate.split('T')[0] : '',
    interestRate: lend.interestRate,
    description: lend.description
  }
  showEditModal.value = true
}

const deleteLend = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Lend',
    message: 'Are you sure you want to delete this lend? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await lendsStore.deleteLend(id)
      success('Lend deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting lend')
    }
  }
}

const saveLend = async () => {
  try {
    const lendData = { ...form.value }

    // Clean up empty values
    if (!lendData.dueDate) delete lendData.dueDate
    if (!lendData.accountId) delete lendData.accountId
    if (!lendData.interestRate) delete lendData.interestRate

    if (showEditModal.value) {
      await lendsStore.updateLend(editingLend.value.id, lendData)
      success('Lend updated successfully')
    } else {
      await lendsStore.createLend(lendData)
      success('Lend created successfully')
      await accountsStore.fetchAccounts() // Refresh accounts to update balance
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving lend')
  }
}

const openPaymentModal = (lend) => {
  selectedLend.value = lend
  paymentForm.value = {
    amount: lend.remainingAmount,
    accountId: '',
    paymentDate: new Date().toISOString().split('T')[0],
    description: `Payment from ${lend.debtorName}`
  }
  showPaymentModal.value = true
}

const recordPayment = async () => {
  try {
    await lendsStore.recordPayment(selectedLend.value.id, paymentForm.value)
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
    debtorName: '',
    originalAmount: 0,
    lentDate: new Date().toISOString().split('T')[0],
    dueDate: '',
    accountId: '',
    interestRate: null,
    description: '',
    isInitial: false
  }
  editingLend.value = null
}

const closePaymentModal = () => {
  showPaymentModal.value = false
  selectedLend.value = null
  paymentForm.value = {
    amount: 0,
    accountId: '',
    paymentDate: new Date().toISOString().split('T')[0],
    description: ''
  }
}

onMounted(async () => {
  await Promise.all([
    lendsStore.fetchLends(),
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

.text-blue { color: #007bff; }
.stat-icon.blue { color: #007bff; }
.stat-icon.purple { color: #6f42c1; }
.stat-icon.green { color: #28a745; }
.stat-icon.orange { color: #fd7e14; }

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
