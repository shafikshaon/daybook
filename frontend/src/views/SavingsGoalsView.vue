<template>
  <div class="savings-goals-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Savings Goals</h1>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
        </svg>
        Add Goal
      </button>
    </div>

    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon purple">🎯</div>
          <div class="stat-value">{{ savingsGoals.length }}</div>
          <div class="stat-label">Active Goals</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(savingsGoalsStore.totalSavedAmount) }}</div>
          <div class="stat-label">Total Saved</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon green">📊</div>
          <div class="stat-value">{{ Math.round(savingsGoalsStore.totalProgress) }}%</div>
          <div class="stat-label">Overall Progress</div>
        </div>
      </div>
    </div>

    <div class="row g-3">
      <div v-for="goal in savingsGoals" :key="goal.id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-3">
              <div>
                <span style="font-size: 2rem;">{{ goal.icon }}</span>
                <h5 class="mt-2">{{ goal.name }}</h5>
              </div>
            </div>
            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span class="fw-bold">{{ formatCurrency(goal.currentAmount) }}</span>
                <span class="text-muted">{{ formatCurrency(goal.targetAmount) }}</span>
              </div>
              <div class="progress" style="height: 12px;">
                <div
                  class="progress-bar progress-bar-professional"
                  :style="{ width: Math.min((goal.currentAmount / goal.targetAmount) * 100, 100) + '%' }"
                ></div>
              </div>
              <div class="d-flex justify-content-between mt-1">
                <small class="text-muted">{{ Math.round((goal.currentAmount / goal.targetAmount) * 100) }}%</small>
                <small class="text-muted">{{ formatCurrency(goal.targetAmount - goal.currentAmount) }} remaining</small>
              </div>
            </div>
            <div class="mb-3">
              <small class="text-muted">Monthly: {{ formatCurrency(goal.monthlyContribution || 0) }}</small>
              <br>
              <small class="text-muted">Target: {{ formatDate(goal.targetDate) }}</small>
            </div>
            <div class="d-flex gap-2">
              <button class="btn btn-sm btn-add-contribution flex-fill" @click="addContribution(goal.id)">
                Add Contribution
              </button>
              <button class="btn btn-sm btn-delete-goal" @click="deleteGoal(goal.id)">
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Savings Goal</h5>
            <button type="button" class="btn-close" @click="showAddModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveGoal">
              <div class="mb-3">
                <label class="form-label">Goal Name *</label>
                <input type="text" class="form-control" v-model="form.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Icon</label>
                <input type="text" class="form-control" v-model="form.icon" placeholder="e.g., 🎯" />
              </div>
              <div class="mb-3">
                <label class="form-label">Target Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.targetAmount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Current Amount</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.currentAmount" />
              </div>
              <div class="mb-3">
                <label class="form-label">Monthly Contribution</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.monthlyContribution" />
              </div>
              <div class="mb-3">
                <label class="form-label">Target Date</label>
                <input type="date" class="form-control" v-model="form.targetDate" />
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showAddModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Create</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Contribution Modal -->
    <div class="modal fade" :class="{ 'show d-block': showContributionModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showContributionModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Contribution</h5>
            <button type="button" class="btn-close" @click="showContributionModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveContribution">
              <div class="mb-3">
                <label class="form-label">From Account *</label>
                <select class="form-select" v-model="contributionForm.accountId" required>
                  <option value="">Select account...</option>
                  <option v-for="account in accounts" :key="account.id" :value="account.id">
                    {{ account.name }} ({{ formatCurrency(account.balance) }})
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="contributionForm.amount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input type="date" class="form-control" v-model="contributionForm.date" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Description</label>
                <input type="text" class="form-control" v-model="contributionForm.description" placeholder="Optional note" />
              </div>
              <div class="alert alert-info">
                <small>
                  This will deduct {{ formatCurrency(contributionForm.amount || 0) }} from
                  {{ getAccountName(contributionForm.accountId) || 'selected account' }} and add it to your savings goal.
                  A transaction record will be created.
                </small>
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showContributionModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Add Contribution</button>
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
import { useSavingsGoalsStore } from '@/stores/savingsGoals'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const savingsGoalsStore = useSavingsGoalsStore()
const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()
const showAddModal = ref(false)
const showContributionModal = ref(false)
const selectedGoalId = ref(null)
const form = ref({ name: '', icon: '🎯', targetAmount: 0, currentAmount: 0, monthlyContribution: 0, targetDate: '' })

const contributionForm = ref({
  accountId: '',
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  description: ''
})

const savingsGoals = computed(() => savingsGoalsStore.activeSavingsGoals)
const accounts = computed(() => accountsStore.allAccounts)
const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : ''
}

const addContribution = (goalId) => {
  selectedGoalId.value = goalId
  contributionForm.value = {
    accountId: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    description: ''
  }
  showContributionModal.value = true
}

const saveContribution = async () => {
  try {
    if (!contributionForm.value.accountId) {
      error('Please select an account')
      return
    }

    if (contributionForm.value.amount <= 0) {
      error('Amount must be greater than 0')
      return
    }

    // Get the selected account
    const account = accountsStore.getAccountById(contributionForm.value.accountId)
    if (!account) {
      error('Selected account not found')
      return
    }

    // Check balance (ensure both are numbers for comparison)
    // Parse balance by removing any currency symbols and converting to number
    let accountBalance = account.balance

    // If balance is a string with currency formatting, extract the numeric value
    if (typeof accountBalance === 'string') {
      // Remove currency symbols and commas, then parse
      accountBalance = parseFloat(accountBalance.replace(/[^\d.-]/g, ''))
    } else {
      accountBalance = Number(accountBalance)
    }

    const contributionAmount = Number(contributionForm.value.amount)

    console.log('Balance Check:', {
      rawBalance: account.balance,
      accountBalance,
      contributionAmount,
      comparison: `${accountBalance} < ${contributionAmount} = ${accountBalance < contributionAmount}`
    })

    if (isNaN(accountBalance) || isNaN(contributionAmount)) {
      error('Invalid amount or balance value')
      return
    }

    if (accountBalance < contributionAmount) {
      error(`Insufficient balance. Available: ${formatCurrency(accountBalance)}, Requested: ${formatCurrency(contributionAmount)}`)
      return
    }

    // Get the savings goal to include in description
    const goal = savingsGoals.value.find(g => g.id === selectedGoalId.value)
    const notes = contributionForm.value.description || `Contribution to ${goal?.name || 'Savings Goal'}`

    // Add contribution to savings goal (backend will create transaction and update account)
    await savingsGoalsStore.addContribution(
      selectedGoalId.value,
      contributionForm.value.amount,
      contributionForm.value.accountId,
      contributionForm.value.date,
      notes
    )

    success('Contribution added successfully')
    showContributionModal.value = false

    // Refresh data
    await Promise.all([
      accountsStore.fetchAccounts(),
      savingsGoalsStore.fetchSavingsGoals(),
      transactionsStore.fetchTransactions(1, 20) // Refresh transactions to show new entry
    ])
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error adding contribution')
  }
}

const deleteGoal = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Savings Goal',
    message: 'Are you sure you want to delete this savings goal? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await savingsGoalsStore.deleteSavingsGoal(id)
      success('Savings goal deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting savings goal')
    }
  }
}

const saveGoal = async () => {
  try {
    await savingsGoalsStore.createSavingsGoal(form.value)
    success('Savings goal created successfully')
    showAddModal.value = false
    form.value = { name: '', icon: '🎯', targetAmount: 0, currentAmount: 0, monthlyContribution: 0, targetDate: '' }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error creating savings goal')
  }
}

onMounted(async () => {
  await Promise.all([
    savingsGoalsStore.fetchSavingsGoals(),
    accountsStore.fetchAccounts()
  ])
})
</script>

<style scoped>
/* Professional progress bar */
.progress-bar-professional {
  background-color: #3b82f6;
  background-image: linear-gradient(45deg, rgba(255, 255, 255, 0.15) 25%, transparent 25%, transparent 50%, rgba(255, 255, 255, 0.15) 50%, rgba(255, 255, 255, 0.15) 75%, transparent 75%, transparent);
  background-size: 1rem 1rem;
}

/* Professional button styles */
.btn-add-contribution {
  color: #1e40af;
  border-color: #3b82f6;
  background-color: #eff6ff;
}

.btn-add-contribution:hover {
  color: #ffffff;
  background-color: #3b82f6;
  border-color: #2563eb;
}

.btn-delete-goal {
  color: #991b1b;
  border-color: #ef4444;
  background-color: #fef2f2;
}

.btn-delete-goal:hover {
  color: #ffffff;
  background-color: #ef4444;
  border-color: #dc2626;
}

/* Dark mode support */
.dark-mode .progress-bar-professional {
  background-color: #3b82f6;
}

.dark-mode .btn-add-contribution {
  color: #93c5fd;
  border-color: #3b82f6;
  background-color: #1e3a5f;
}

.dark-mode .btn-add-contribution:hover {
  color: #ffffff;
  background-color: #3b82f6;
  border-color: #60a5fa;
}

.dark-mode .btn-delete-goal {
  color: #fca5a5;
  border-color: #ef4444;
  background-color: #5f1e1e;
}

.dark-mode .btn-delete-goal:hover {
  color: #ffffff;
  background-color: #ef4444;
  border-color: #f87171;
}
</style>
