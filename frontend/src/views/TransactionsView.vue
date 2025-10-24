<template>
  <div class="transactions-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Transactions</h1>
      <div class="d-flex gap-2">
        <button class="btn btn-outline-primary" @click="showTransferModal = true">
          🔄 Transfer Funds
        </button>
        <button class="btn btn-primary" @click="showAddModal = true">
          + Add Transaction
        </button>
      </div>
    </div>

    <!-- Summary Stats -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon green">📈</div>
          <div class="stat-value">{{ formatCurrency(totalIncome) }}</div>
          <div class="stat-label">Total Income</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(totalExpense) }}</div>
          <div class="stat-label">Total Expenses</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(totalIncome - totalExpense) }}</div>
          <div class="stat-label">Net Cash Flow</div>
          <div :class="totalIncome - totalExpense >= 0 ? 'stat-change positive' : 'stat-change negative'">
            {{ totalIncome - totalExpense >= 0 ? '↑' : '↓' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-12 col-md-3">
            <label class="form-label">Type</label>
            <select class="form-select" v-model="filters.type">
              <option value="">All</option>
              <option value="income">Income</option>
              <option value="expense">Expense</option>
              <option value="transfer">Transfer</option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Category</label>
            <select class="form-select" v-model="filters.category">
              <option value="">All Categories</option>
              <option v-for="cat in allCategories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Account</label>
            <select class="form-select" v-model="filters.account">
              <option value="">All Accounts</option>
              <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                {{ acc.name }}
              </option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Search</label>
            <input
              type="text"
              class="form-control"
              v-model="filters.search"
              placeholder="Search transactions..."
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Transactions Table -->
    <div class="card">
      <div class="card-header">
        <h5 class="mb-0">Transaction History</h5>
      </div>
      <div class="card-body p-0">
        <div v-if="filteredTransactions.length === 0" class="p-4 text-center text-muted">
          No transactions found
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Category</th>
                <th>Account</th>
                <th>Type</th>
                <th class="text-end">Amount</th>
                <th class="text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="transaction in paginatedTransactions" :key="transaction.id">
                <td>{{ formatDate(transaction.date) }}</td>
                <td>{{ transaction.description || '-' }}</td>
                <td>
                  <span class="badge" :style="{ backgroundColor: getCategoryColor(transaction.categoryId) }">
                    {{ getCategoryName(transaction.categoryId) }}
                  </span>
                </td>
                <td>
                  <span v-if="transaction.type === 'transfer'">
                    {{ getAccountName(transaction.accountId) }} → {{ getAccountName(transaction.toAccountId) }}
                  </span>
                  <span v-else>{{ getAccountName(transaction.accountId) }}</span>
                </td>
                <td>
                  <span
                    class="badge"
                    :class="{
                      'bg-success': transaction.type === 'income',
                      'bg-danger': transaction.type === 'expense',
                      'bg-primary': transaction.type === 'transfer'
                    }"
                  >
                    {{ transaction.type }}
                  </span>
                </td>
                <td class="text-end">
                  <span
                    :class="{
                      'text-success fw-bold': transaction.type === 'income',
                      'text-danger fw-bold': transaction.type === 'expense',
                      'text-primary fw-bold': transaction.type === 'transfer'
                    }"
                  >
                    {{ transaction.type === 'income' ? '+' : transaction.type === 'transfer' ? '↔' : '-' }}{{ formatCurrency(transaction.amount) }}
                  </span>
                </td>
                <td class="text-center">
                  <button
                    class="btn btn-sm btn-outline-primary me-1"
                    @click="editTransaction(transaction)"
                  >
                    Edit
                  </button>
                  <button
                    class="btn btn-sm btn-outline-danger"
                    @click="confirmDelete(transaction)"
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

    <!-- Add/Edit Transaction Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Transaction' : 'Add Transaction' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTransaction">
              <div class="mb-3">
                <label class="form-label">Type *</label>
                <select class="form-select" v-model="form.type" required>
                  <option value="income">Income</option>
                  <option value="expense">Expense</option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="form.amount"
                  required
                  placeholder="0.00"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Category *</label>
                <select class="form-select" v-model="form.categoryId" required>
                  <option value="">Select category...</option>
                  <option
                    v-for="cat in filteredCategories"
                    :key="cat.id"
                    :value="cat.id"
                  >
                    {{ cat.icon }} {{ cat.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Account *</label>
                <select class="form-select" v-model="form.accountId" required>
                  <option value="">Select account...</option>
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="form.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.description"
                  placeholder="Optional description"
                />
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  {{ showEditModal ? 'Update' : 'Create' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Transfer Funds Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showTransferModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showTransferModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Transfer Funds</h5>
            <button type="button" class="btn-close" @click="closeTransferModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTransfer">
              <div class="mb-3">
                <label class="form-label">From Account *</label>
                <select class="form-select" v-model="transferForm.fromAccountId" required>
                  <option value="">Select account...</option>
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.name }} ({{ formatCurrency(acc.balance) }})
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Transfer To *</label>
                <select class="form-select" v-model="transferForm.destinationType" required>
                  <option value="account">Another Account</option>
                  <option value="savings_goal">Savings Goal</option>
                  <option value="fixed_deposit">Fixed Deposit (Create New)</option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'account'">
                <label class="form-label">To Account *</label>
                <select class="form-select" v-model="transferForm.toAccountId" required>
                  <option value="">Select account...</option>
                  <option
                    v-for="acc in accounts"
                    :key="acc.id"
                    :value="acc.id"
                    :disabled="acc.id === transferForm.fromAccountId"
                  >
                    {{ acc.name }} ({{ formatCurrency(acc.balance) }})
                  </option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'savings_goal'">
                <label class="form-label">Savings Goal *</label>
                <select class="form-select" v-model="transferForm.savingsGoalId" required>
                  <option value="">Select savings goal...</option>
                  <option v-for="goal in savingsGoals" :key="goal.id" :value="goal.id">
                    {{ goal.name }} ({{ formatCurrency(goal.currentAmount) }} / {{ formatCurrency(goal.targetAmount) }})
                  </option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Fixed Deposit Name *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="transferForm.fdName"
                  required
                  placeholder="e.g., 1-Year FD"
                />
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Interest Rate (% per annum) *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="transferForm.fdInterestRate"
                  required
                  placeholder="e.g., 5.5"
                />
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Tenure (months) *</label>
                <input
                  type="number"
                  class="form-control"
                  v-model.number="transferForm.fdTenureMonths"
                  required
                  placeholder="e.g., 12"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="transferForm.amount"
                  required
                  placeholder="0.00"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="transferForm.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="transferForm.description"
                  placeholder="Optional description"
                />
              </div>

              <div class="alert alert-info">
                <small v-if="transferForm.destinationType === 'account'">
                  This will transfer {{ formatCurrency(transferForm.amount || 0) }} from {{ getAccountName(transferForm.fromAccountId) || 'source' }} to {{ getAccountName(transferForm.toAccountId) || 'destination' }}.
                </small>
                <small v-else-if="transferForm.destinationType === 'savings_goal'">
                  This will transfer {{ formatCurrency(transferForm.amount || 0) }} from {{ getAccountName(transferForm.fromAccountId) || 'source' }} to your savings goal.
                </small>
                <small v-else-if="transferForm.destinationType === 'fixed_deposit'">
                  This will create a new fixed deposit of {{ formatCurrency(transferForm.amount || 0) }} for {{ transferForm.fdTenureMonths || 0 }} months at {{ transferForm.fdInterestRate || 0 }}% interest.
                </small>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeTransferModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  Transfer
                </button>
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
import { useSavingsGoalsStore } from '@/stores/savingsGoals'
import { useFixedDepositsStore } from '@/stores/fixedDeposits'
import { useSettingsStore } from '@/stores/settings'

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const savingsGoalsStore = useSavingsGoalsStore()
const fixedDepositsStore = useFixedDepositsStore()
const settingsStore = useSettingsStore()

const showAddModal = ref(false)
const showEditModal = ref(false)
const showTransferModal = ref(false)
const editingTransaction = ref(null)

const filters = ref({
  type: '',
  category: '',
  account: '',
  search: ''
})

const form = ref({
  type: 'expense',
  amount: 0,
  categoryId: '',
  accountId: '',
  date: new Date().toISOString().split('T')[0],
  description: '',
  tags: []
})

const transferForm = ref({
  fromAccountId: '',
  destinationType: 'account',
  toAccountId: '',
  savingsGoalId: '',
  fdName: '',
  fdInterestRate: 0,
  fdTenureMonths: 12,
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  description: 'Transfer between accounts'
})

const transactions = computed(() => transactionsStore.allTransactions)
const accounts = computed(() => accountsStore.allAccounts)
const savingsGoals = computed(() => savingsGoalsStore.activeSavingsGoals)
const allCategories = computed(() => transactionsStore.categories)

const filteredCategories = computed(() => {
  return transactionsStore.categories.filter(c => c.type === form.value.type)
})

const totalIncome = computed(() => transactionsStore.totalIncome())
const totalExpense = computed(() => transactionsStore.totalExpense())

const filteredTransactions = computed(() => {
  let result = transactions.value

  if (filters.value.type) {
    result = result.filter(t => t.type === filters.value.type)
  }

  if (filters.value.category) {
    result = result.filter(t => t.categoryId === filters.value.category)
  }

  if (filters.value.account) {
    result = result.filter(t => t.accountId === filters.value.account)
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(t =>
      (t.description || '').toLowerCase().includes(search) ||
      getCategoryName(t.categoryId).toLowerCase().includes(search)
    )
  }

  return result
})

const paginatedTransactions = computed(() => {
  return filteredTransactions.value.slice(0, 50)
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

const getCategoryName = (categoryId) => {
  const category = transactionsStore.getCategoryById(categoryId)
  return category ? category.name : categoryId
}

const getCategoryColor = (categoryId) => {
  const category = transactionsStore.getCategoryById(categoryId)
  return category ? category.color : '#6b7280'
}

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : accountId
}

const editTransaction = (transaction) => {
  editingTransaction.value = transaction
  form.value = {
    ...transaction,
    date: new Date(transaction.date).toISOString().split('T')[0]
  }
  showEditModal.value = true
}

const confirmDelete = async (transaction) => {
  if (confirm('Are you sure you want to delete this transaction?')) {
    await transactionsStore.deleteTransaction(transaction.id)
  }
}

const saveTransaction = async () => {
  try {
    const transactionData = {
      ...form.value,
      date: new Date(form.value.date).toISOString()
    }

    if (showEditModal.value) {
      await transactionsStore.updateTransaction(editingTransaction.value.id, transactionData)
    } else {
      await transactionsStore.createTransaction(transactionData)
    }

    closeModal()
  } catch (error) {
    alert('Error saving transaction: ' + error.message)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingTransaction.value = null
  form.value = {
    type: 'expense',
    amount: 0,
    categoryId: '',
    accountId: '',
    date: new Date().toISOString().split('T')[0],
    description: '',
    tags: []
  }
}

const saveTransfer = async () => {
  try {
    const dateISO = new Date(transferForm.value.date).toISOString()

    if (transferForm.value.destinationType === 'account') {
      if (transferForm.value.fromAccountId === transferForm.value.toAccountId) {
        alert('Source and destination accounts must be different')
        return
      }

      await transactionsStore.transferFunds(
        transferForm.value.fromAccountId,
        transferForm.value.toAccountId,
        transferForm.value.amount,
        transferForm.value.description,
        dateISO
      )
    } else if (transferForm.value.destinationType === 'savings_goal') {
      if (!transferForm.value.savingsGoalId) {
        alert('Please select a savings goal')
        return
      }

      await transactionsStore.transferToSavingsGoal(
        transferForm.value.fromAccountId,
        transferForm.value.savingsGoalId,
        transferForm.value.amount,
        transferForm.value.description,
        dateISO
      )

      await savingsGoalsStore.fetchSavingsGoals()
    } else if (transferForm.value.destinationType === 'fixed_deposit') {
      if (!transferForm.value.fdName || !transferForm.value.fdInterestRate || !transferForm.value.fdTenureMonths) {
        alert('Please fill in all fixed deposit details')
        return
      }

      // Calculate maturity date
      const maturityDate = new Date(transferForm.value.date)
      maturityDate.setMonth(maturityDate.getMonth() + transferForm.value.fdTenureMonths)

      // Deduct from source account
      await accountsStore.updateBalance(transferForm.value.fromAccountId, transferForm.value.amount, 'subtract')

      // Create fixed deposit
      await fixedDepositsStore.createFixedDeposit({
        name: transferForm.value.fdName,
        principal: transferForm.value.amount,
        interestRate: transferForm.value.fdInterestRate,
        tenureMonths: transferForm.value.fdTenureMonths,
        startDate: dateISO,
        maturityDate: maturityDate.toISOString(),
        bank: 'N/A'
      })

      // Create transaction record
      await transactionsStore.createTransaction({
        type: 'expense',
        amount: transferForm.value.amount,
        categoryId: 'other_expense',
        accountId: transferForm.value.fromAccountId,
        date: dateISO,
        description: `Fixed Deposit: ${transferForm.value.fdName}`,
        tags: ['fixed_deposit']
      })

      await fixedDepositsStore.fetchFixedDeposits()
    }

    closeTransferModal()
    // Refresh accounts to show updated balances
    await accountsStore.fetchAccounts()
  } catch (error) {
    alert('Error transferring funds: ' + error.message)
  }
}

const closeTransferModal = () => {
  showTransferModal.value = false
  transferForm.value = {
    fromAccountId: '',
    destinationType: 'account',
    toAccountId: '',
    savingsGoalId: '',
    fdName: '',
    fdInterestRate: 0,
    fdTenureMonths: 12,
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    description: 'Transfer between accounts'
  }
}

onMounted(async () => {
  await Promise.all([
    transactionsStore.fetchTransactions(),
    accountsStore.fetchAccounts(),
    savingsGoalsStore.fetchSavingsGoals(),
    fixedDepositsStore.fetchFixedDeposits()
  ])
})
</script>
