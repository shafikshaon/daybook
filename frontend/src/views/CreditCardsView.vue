<template>
  <div class="credit-cards-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Credit Cards</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Card</button>
    </div>

    <!-- Summary -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">💳</div>
          <div class="stat-value">{{ creditCards.length }}</div>
          <div class="stat-label">Total Cards</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(creditCardsStore.totalCreditLimit) }}</div>
          <div class="stat-label">Credit Limit</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(creditCardsStore.totalOutstanding) }}</div>
          <div class="stat-label">Outstanding</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon orange">📊</div>
          <div class="stat-value">{{ Math.round(creditCardsStore.creditUtilization) }}%</div>
          <div class="stat-label">Utilization</div>
        </div>
      </div>
    </div>

    <!-- Cards Grid -->
    <div class="row g-3">
      <div v-for="card in creditCards" :key="card.id" class="col-12 col-md-6 col-lg-4">
        <div class="card credit-card-item" style="background: linear-gradient(135deg, #6f42c1 0%, #8b5cf6 100%); color: white;">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-3">
              <div>
                <h5 class="card-title mb-1">{{ card.name }}</h5>
                <p class="mb-0 small">**** **** **** {{ card.lastFourDigits }}</p>
              </div>
              <div class="dropdown">
                <button class="btn btn-sm btn-light dropdown-toggle" type="button" data-bs-toggle="dropdown">
                  •••
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" @click="viewCardDetails(card)">View Details</a></li>
                  <li><a class="dropdown-item" @click="editCard(card)">Edit Card</a></li>
                  <li><hr class="dropdown-divider"></li>
                  <li><a class="dropdown-item text-danger" @click="deleteCard(card.id)">Delete</a></li>
                </ul>
              </div>
            </div>

            <div class="mb-3">
              <small>Current Balance</small>
              <h3 class="mb-2">{{ formatCurrency(card.currentBalance) }}</h3>
              <div class="progress" style="height: 6px; background-color: rgba(255,255,255,0.3);">
                <div
                  class="progress-bar"
                  :class="getUtilizationClass(card.currentBalance / card.creditLimit * 100)"
                  :style="{ width: Math.min((card.currentBalance / card.creditLimit * 100), 100) + '%' }"
                ></div>
              </div>
              <small class="d-flex justify-content-between mt-1">
                <span>{{ Math.round((card.currentBalance / card.creditLimit * 100)) }}% used</span>
                <span>Limit: {{ formatCurrency(card.creditLimit) }}</span>
              </small>
            </div>

            <div class="d-flex justify-content-between text-sm mb-3">
              <span>APR: {{ card.apr }}%</span>
              <span v-if="card.dueDate">Due: {{ formatDate(card.dueDate) }}</span>
            </div>

            <div class="d-flex gap-2">
              <button
                class="btn btn-light flex-fill btn-sm"
                @click="openTransactionModal(card)"
              >
                📝 Add Expense
              </button>
              <button
                class="btn btn-success flex-fill btn-sm"
                @click="openPaymentModal(card)"
                :disabled="card.currentBalance <= 0"
              >
                💳 Pay Bill
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Card Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Credit Card</h5>
            <button type="button" class="btn-close" @click="showAddModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveCard">
              <div class="mb-3">
                <label class="form-label">Card Name *</label>
                <input type="text" class="form-control" v-model="form.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Last 4 Digits *</label>
                <input type="text" class="form-control" v-model="form.lastFourDigits" maxlength="4" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Card Network</label>
                <select class="form-select" v-model="form.cardNetwork">
                  <option value="">Select network...</option>
                  <option value="Visa">Visa</option>
                  <option value="Mastercard">Mastercard</option>
                  <option value="American Express">American Express</option>
                  <option value="Discover">Discover</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Credit Limit *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.creditLimit" required />
              </div>
              <div class="mb-3">
                <label class="form-label">APR (%) *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.apr" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Current Balance</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.currentBalance" />
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

    <!-- Add Transaction Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showTransactionModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showTransactionModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Credit Card Expense</h5>
            <button type="button" class="btn-close" @click="closeTransactionModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTransaction">
              <div class="mb-3">
                <label class="form-label">Card</label>
                <input type="text" class="form-control" :value="transactionForm.cardName" disabled />
              </div>

              <div class="mb-3">
                <label class="form-label">Type *</label>
                <select class="form-select" v-model="transactionForm.type" required>
                  <option value="purchase">Purchase</option>
                  <option value="fee">Fee</option>
                  <option value="interest">Interest Charge</option>
                  <option value="refund">Refund</option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="transactionForm.amount"
                  required
                  placeholder="0.00"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Category</label>
                <select class="form-select" v-model="transactionForm.categoryId">
                  <option value="">Select category...</option>
                  <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.id">
                    {{ cat.icon }} {{ cat.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Merchant</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="transactionForm.merchant"
                  placeholder="Where did you spend?"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="transactionForm.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea
                  class="form-control"
                  v-model="transactionForm.description"
                  rows="2"
                  placeholder="Additional details..."
                ></textarea>
              </div>

              <div class="alert alert-info">
                <small>This will add {{ formatCurrency(transactionForm.amount || 0) }} to your credit card balance.</small>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeTransactionModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  Record Transaction
                </button>
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
            <h5 class="modal-title">Make Credit Card Payment</h5>
            <button type="button" class="btn-close" @click="closePaymentModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="savePayment">
              <div class="mb-3">
                <label class="form-label">Card</label>
                <input type="text" class="form-control" :value="paymentForm.cardName" disabled />
              </div>

              <div class="mb-3">
                <label class="form-label">Current Balance</label>
                <input type="text" class="form-control" :value="formatCurrency(paymentForm.currentBalance)" disabled />
              </div>

              <div class="mb-3">
                <label class="form-label">Payment Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="paymentForm.amount"
                  required
                  placeholder="0.00"
                  :max="paymentForm.currentBalance"
                />
                <small class="text-muted">Maximum: {{ formatCurrency(paymentForm.currentBalance) }}</small>
              </div>

              <div class="mb-3">
                <label class="form-label">Pay From Account *</label>
                <select class="form-select" v-model="paymentForm.paymentAccountId" required>
                  <option value="">Select account...</option>
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.name }} ({{ formatCurrency(acc.balance) }})
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Payment Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="paymentForm.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="paymentForm.description"
                  placeholder="Optional note..."
                />
              </div>

              <div class="alert alert-info">
                <small>This will deduct {{ formatCurrency(paymentForm.amount || 0) }} from {{ getAccountName(paymentForm.paymentAccountId) || 'selected account' }} and reduce your credit card balance.</small>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closePaymentModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  Pay {{ formatCurrency(paymentForm.amount || 0) }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Card Details Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showDetailsModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showDetailsModal"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ selectedCard?.name }} - Transactions</h5>
            <button type="button" class="btn-close" @click="closeDetailsModal"></button>
          </div>
          <div class="modal-body">
            <!-- Card Info -->
            <div class="card mb-3" style="background: linear-gradient(135deg, #6f42c1 0%, #8b5cf6 100%); color: white;">
              <div class="card-body">
                <div class="row">
                  <div class="col-md-6">
                    <small>Current Balance</small>
                    <h4>{{ formatCurrency(selectedCard?.currentBalance) }}</h4>
                  </div>
                  <div class="col-md-6">
                    <small>Available Credit</small>
                    <h4>{{ formatCurrency((selectedCard?.creditLimit || 0) - (selectedCard?.currentBalance || 0)) }}</h4>
                  </div>
                </div>
              </div>
            </div>

            <!-- Transactions List -->
            <h6 class="mb-3">Transaction History</h6>
            <div v-if="cardTransactions.length === 0" class="text-center text-muted py-4">
              No transactions yet
            </div>
            <div v-else class="list-group">
              <div v-for="transaction in cardTransactions" :key="transaction.id" class="list-group-item">
                <div class="d-flex justify-content-between align-items-start">
                  <div>
                    <h6 class="mb-1">
                      <span class="badge" :class="getTransactionTypeBadge(transaction.type)">
                        {{ transaction.type }}
                      </span>
                      {{ transaction.merchant || transaction.description || 'Transaction' }}
                    </h6>
                    <small class="text-muted">{{ formatDate(transaction.date) }}</small>
                    <p class="mb-0 small" v-if="transaction.description">{{ transaction.description }}</p>
                  </div>
                  <div class="text-end">
                    <strong
                      :class="{
                        'text-danger': transaction.type === 'purchase' || transaction.type === 'fee' || transaction.type === 'interest',
                        'text-success': transaction.type === 'payment' || transaction.type === 'refund'
                      }"
                    >
                      {{ transaction.type === 'payment' || transaction.type === 'refund' ? '-' : '+' }}{{ formatCurrency(transaction.amount) }}
                    </strong>
                    <button
                      class="btn btn-sm btn-outline-danger mt-2"
                      @click="deleteCardTransaction(transaction.id)"
                      v-if="transaction.type !== 'payment'"
                    >
                      Delete
                    </button>
                  </div>
                </div>
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
import { useCreditCardsStore } from '@/stores/creditCards'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const creditCardsStore = useCreditCardsStore()
const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const settingsStore = useSettingsStore()
const { success, error, confirm } = useNotification()

const showAddModal = ref(false)
const showPaymentModal = ref(false)
const showTransactionModal = ref(false)
const showDetailsModal = ref(false)
const selectedCard = ref(null)
const cardTransactions = ref([])

const form = ref({ name: '', lastFourDigits: '', cardNetwork: '', creditLimit: 0, apr: 0, currentBalance: 0 })

const transactionForm = ref({
  cardId: '',
  cardName: '',
  type: 'purchase',
  amount: 0,
  categoryId: '',
  merchant: '',
  date: new Date().toISOString().split('T')[0],
  description: ''
})

const paymentForm = ref({
  cardId: '',
  cardName: '',
  currentBalance: 0,
  amount: 0,
  paymentAccountId: '',
  date: new Date().toISOString().split('T')[0],
  description: ''
})

const creditCards = computed(() => creditCardsStore.allCreditCards)
const accounts = computed(() => accountsStore.allAccounts)
const expenseCategories = computed(() => transactionsStore.expenseCategories)

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : accountId
}

const getUtilizationClass = (utilization) => {
  if (utilization >= 90) return 'bg-danger'
  if (utilization >= 70) return 'bg-warning'
  return 'bg-success'
}

const getTransactionTypeBadge = (type) => {
  const badges = {
    purchase: 'bg-primary',
    fee: 'bg-warning',
    interest: 'bg-danger',
    refund: 'bg-success',
    payment: 'bg-info'
  }
  return badges[type] || 'bg-secondary'
}

const openTransactionModal = (card) => {
  transactionForm.value = {
    cardId: card.id,
    cardName: card.name,
    type: 'purchase',
    amount: 0,
    categoryId: '',
    merchant: '',
    date: new Date().toISOString().split('T')[0],
    description: ''
  }
  showTransactionModal.value = true
}

const closeTransactionModal = () => {
  showTransactionModal.value = false
  transactionForm.value = {
    cardId: '',
    cardName: '',
    type: 'purchase',
    amount: 0,
    categoryId: '',
    merchant: '',
    date: new Date().toISOString().split('T')[0],
    description: ''
  }
}

const saveTransaction = async () => {
  try {
    await creditCardsStore.recordTransaction(transactionForm.value.cardId, {
      type: transactionForm.value.type,
      amount: transactionForm.value.amount,
      categoryId: transactionForm.value.categoryId,
      merchant: transactionForm.value.merchant,
      date: new Date(transactionForm.value.date).toISOString(),
      description: transactionForm.value.description
    })
    success('Transaction recorded successfully')
    closeTransactionModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error recording transaction')
  }
}

const openPaymentModal = (card) => {
  paymentForm.value = {
    cardId: card.id,
    cardName: card.name,
    currentBalance: card.currentBalance,
    amount: card.currentBalance,
    paymentAccountId: '',
    date: new Date().toISOString().split('T')[0],
    description: `Credit card payment - ${card.name}`
  }
  showPaymentModal.value = true
}

const closePaymentModal = () => {
  showPaymentModal.value = false
  paymentForm.value = {
    cardId: '',
    cardName: '',
    currentBalance: 0,
    amount: 0,
    paymentAccountId: '',
    date: new Date().toISOString().split('T')[0],
    description: ''
  }
}

const savePayment = async () => {
  try {
    if (paymentForm.value.amount > paymentForm.value.currentBalance) {
      error('Payment amount cannot exceed current balance')
      return
    }

    if (!paymentForm.value.paymentAccountId) {
      error('Please select a payment account')
      return
    }

    await creditCardsStore.recordPayment(
      paymentForm.value.cardId,
      paymentForm.value.amount,
      paymentForm.value.paymentAccountId,
      new Date(paymentForm.value.date).toISOString(),
      paymentForm.value.description
    )

    closePaymentModal()
    await accountsStore.fetchAccounts()
    success('Payment recorded successfully')
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error recording payment')
  }
}

const saveCard = async () => {
  try {
    await creditCardsStore.createCreditCard(form.value)
    success('Credit card added successfully')
    showAddModal.value = false
    form.value = { name: '', lastFourDigits: '', cardNetwork: '', creditLimit: 0, apr: 0, currentBalance: 0 }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error adding credit card')
  }
}

const viewCardDetails = async (card) => {
  selectedCard.value = card
  try {
    const transactions = await creditCardsStore.fetchTransactions(card.id)
    cardTransactions.value = transactions
    showDetailsModal.value = true
  } catch (err) {
    error('Error loading card transactions')
  }
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedCard.value = null
  cardTransactions.value = []
}

const deleteCardTransaction = async (transactionId) => {
  const confirmed = await confirm({
    title: 'Delete Transaction',
    message: 'Are you sure you want to delete this transaction? This will adjust your card balance.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await creditCardsStore.deleteTransaction(selectedCard.value.id, transactionId)
      cardTransactions.value = cardTransactions.value.filter(t => t.id !== transactionId)
      success('Transaction deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting transaction')
    }
  }
}

const editCard = (card) => {
  form.value = { ...card }
  showAddModal.value = true
}

const deleteCard = async (cardId) => {
  const confirmed = await confirm({
    title: 'Delete Credit Card',
    message: 'Are you sure you want to delete this credit card? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await creditCardsStore.deleteCreditCard(cardId)
      success('Credit card deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting credit card')
    }
  }
}

onMounted(async () => {
  await Promise.all([
    creditCardsStore.fetchCreditCards(),
    accountsStore.fetchAccounts(),
    transactionsStore.fetchTransactions()
  ])
})
</script>

<style scoped>
.credit-card-item {
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.credit-card-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
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
</style>
