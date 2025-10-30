<template>
  <div class="accounts-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Accounts</h1>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
        </svg>
        Add Account
      </button>
    </div>

    <!-- Summary Cards -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(totalBalance) }}</div>
          <div class="stat-label">Total Balance</div>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon blue">🏦</div>
          <div class="stat-value">{{ accounts.length }}</div>
          <div class="stat-label">Total Accounts</div>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon green">💵</div>
          <div class="stat-value">{{ formatCurrency(cashBalance) }}</div>
          <div class="stat-label">Cash Accounts</div>
        </div>
      </div>

      <div class="col-12 col-md-6 col-lg-3">
        <div class="stat-card">
          <div class="stat-icon orange">🏦</div>
          <div class="stat-value">{{ formatCurrency(bankBalance) }}</div>
          <div class="stat-label">Bank Accounts</div>
        </div>
      </div>
    </div>

    <!-- Accounts List -->
    <div class="card">
      <div class="card-header">
        <h5 class="mb-0">All Accounts</h5>
      </div>
      <div class="card-body p-0">
        <div v-if="accounts.length === 0" class="p-4 text-center text-muted">
          <p>No accounts yet. Create your first account to get started!</p>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Account Name</th>
                <th>Type</th>
                <th>Status</th>
                <th>Description</th>
                <th>Currency</th>
                <th class="text-end">Balance</th>
                <th class="text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="account in accounts" :key="account.id">
                <td>
                  <span class="fw-semibold">{{ account.name }}</span>
                </td>
                <td>
                  <span class="badge bg-secondary text-uppercase">{{ account.type.replace('_', ' ') }}</span>
                </td>
                <td>
                  {{account.active ? 'Active' : 'Inactive' }}
                </td>
                <td>
                  <span class="text-muted">{{ account.description || '-' }}</span>
                </td>
                <td>{{ account.currency || settingsStore.settings.currency }}</td>
                <td class="text-end fw-bold" :title="`Initial: ${formatCurrency(account.initialBalance || account.balance)}`">
                  {{ formatCurrency(account.balance) }}
                  <small v-if="account.initialBalance && account.initialBalance !== account.balance" class="text-muted d-block" style="font-size: 0.75rem;">
                    ({{ account.balance >= account.initialBalance ? '+' : '' }}{{ formatCurrency(account.balance - (account.initialBalance || 0)) }})
                  </small>
                </td>
                <td class="text-center">
                  <button
                    class="btn btn-sm btn-info me-1"
                    @click="showReconcileModal(account)"
                    title="Reconcile Account"
                  >
                    Reconcile
                  </button>
                  <button
                    class="btn btn-sm btn-secondary me-1"
                    @click="showHistoryModal(account)"
                    title="View Reconciliation History"
                  >
                    History
                  </button>
                  <button
                    class="btn btn-sm btn-primary me-1"
                    @click="editAccount(account)"
                  >
                    Edit
                  </button>
                  <button
                    class="btn btn-sm btn-danger"
                    @click="confirmDelete(account)"
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

    <!-- Add/Edit Account Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      tabindex="-1"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Account' : 'Add New Account' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveAccount">
              <div class="mb-3">
                <label class="form-label">Account Name *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.name"
                  required
                  placeholder="e.g., Main Checking"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Account Type *</label>
                <select class="form-select" v-model="form.type" required>
                  <option value="">Select type...</option>
                  <option v-for="type in accountTypes" :key="type.value" :value="type.value">
                    {{ type.label }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Initial Balance (Opening Balance) *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="form.balance"
                  required
                  placeholder="0.00"
                  :disabled="showEditModal"
                />
                <small class="text-muted">
                  <span v-if="!showEditModal">Starting balance when creating the account. This will not change.</span>
                  <span v-else>Initial balance cannot be changed. Current balance updates with transactions.</span>
                </small>
              </div>

              <div class="mb-3">
                <label class="form-label">Currency</label>
                <select class="form-select" v-model="form.currency">
                  <option v-for="currency in currencies" :key="currency.code" :value="currency.code">
                    {{ currency.code }} - {{ currency.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea
                  class="form-control"
                  v-model="form.description"
                  rows="2"
                  placeholder="Optional description or notes"
                ></textarea>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  {{ showEditModal ? 'Update' : 'Create' }} Account
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Reconciliation Modal -->
    <ReconciliationModal
      :show="showReconciliation"
      :account="reconciliationAccount"
      @close="closeReconciliationModal"
      @reconciled="handleReconciled"
    />

    <!-- Reconciliation History Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showHistory }"
      tabindex="-1"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showHistory"
    >
      <div class="modal-dialog modal-xl modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Reconciliation History - {{ reconciliationAccount?.name }}</h5>
            <button type="button" class="btn-close" @click="showHistory = false"></button>
          </div>
          <div class="modal-body">
            <ReconciliationHistory
              v-if="reconciliationAccount"
              :accountId="reconciliationAccount.id"
              @refresh="refreshAccounts"
              ref="historyComponent"
            />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showHistory = false">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'
import ReconciliationModal from '@/components/ReconciliationModal.vue'
import ReconciliationHistory from '@/components/ReconciliationHistory.vue'

const accountsStore = useAccountsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingAccount = ref(null)
const showReconciliation = ref(false)
const showHistory = ref(false)
const reconciliationAccount = ref(null)
const historyComponent = ref(null)

const form = ref({
  name: '',
  type: '',
  balance: 0,
  currency: 'BDT',
  description: ''
})

const accounts = computed(() => accountsStore.allAccounts)
const totalBalance = computed(() => accountsStore.totalBalance)
const accountTypes = computed(() => accountsStore.accountTypes)
const currencies = computed(() => settingsStore.currencies)

const cashBalance = computed(() => {
  return accountsStore.cashAccounts.reduce((sum, acc) => sum + acc.balance, 0)
})

const bankBalance = computed(() => {
  return accountsStore.bankAccounts.reduce((sum, acc) => sum + acc.balance, 0)
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const editAccount = (account) => {
  editingAccount.value = account
  form.value = { ...account }
  showEditModal.value = true
}

const confirmDelete = async (account) => {
  const confirmed = await confirm({
    title: 'Delete Account',
    message: `Are you sure you want to delete "${account.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await accountsStore.deleteAccount(account.id)
      success('Account deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting account')
    }
  }
}

const saveAccount = async () => {
  try {
    if (showEditModal.value) {
      await accountsStore.updateAccount(editingAccount.value.id, form.value)
      success('Account updated successfully')
    } else {
      await accountsStore.createAccount(form.value)
      success('Account created successfully')
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving account')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingAccount.value = null
  form.value = {
    name: '',
    type: '',
    balance: 0,
    currency: settingsStore.settings.currency,
    description: ''
  }
}

const showReconcileModal = (account) => {
  reconciliationAccount.value = account
  showReconciliation.value = true
}

const closeReconciliationModal = () => {
  showReconciliation.value = false
  reconciliationAccount.value = null
}

const showHistoryModal = (account) => {
  reconciliationAccount.value = account
  showHistory.value = true
}

const handleReconciled = async () => {
  success('Account reconciled successfully')
  await refreshAccounts()
  // Close reconciliation modal - don't show history automatically
  closeReconciliationModal()
}

const refreshAccounts = async () => {
  await accountsStore.fetchAccounts()
}

onMounted(async () => {
  await Promise.all([
    accountsStore.fetchAccounts(),
    accountsStore.fetchAccountTypes()
  ])
  form.value.currency = settingsStore.settings.currency
})
</script>
