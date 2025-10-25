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
        <div class="card" style="background: linear-gradient(135deg, #6f42c1 0%, #8b5cf6 100%); color: white;">
          <div class="card-body">
            <h5 class="card-title">{{ card.name }}</h5>
            <p class="mb-2">**** **** **** {{ card.lastFourDigits }}</p>
            <div class="mb-3">
              <small>Balance</small>
              <h3>{{ formatCurrency(card.currentBalance) }}</h3>
              <div class="progress" style="height: 6px; background-color: rgba(255,255,255,0.3);">
                <div
                  class="progress-bar bg-warning"
                  :style="{ width: (card.currentBalance / card.creditLimit * 100) + '%' }"
                ></div>
              </div>
              <small>Limit: {{ formatCurrency(card.creditLimit) }}</small>
            </div>
            <div class="d-flex justify-content-between text-sm mb-3">
              <span>APR: {{ card.apr }}%</span>
              <span>Due: {{ formatDate(card.dueDate) }}</span>
            </div>
            <button
              class="btn btn-light w-100"
              @click="openPaymentModal(card)"
              :disabled="card.currentBalance <= 0"
            >
              💳 Make Payment
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Modal -->
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
                <label class="form-label">Credit Limit *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.creditLimit" required />
              </div>
              <div class="mb-3">
                <label class="form-label">APR (%) *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.apr" required />
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCreditCardsStore } from '@/stores/creditCards'
import { useAccountsStore } from '@/stores/accounts'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const creditCardsStore = useCreditCardsStore()
const accountsStore = useAccountsStore()
const settingsStore = useSettingsStore()
const { success, error } = useNotification()

const showAddModal = ref(false)
const showPaymentModal = ref(false)
const form = ref({ name: '', lastFourDigits: '', creditLimit: 0, apr: 0, currentBalance: 0 })

const paymentForm = ref({
  cardId: '',
  cardName: '',
  currentBalance: 0,
  amount: 0,
  paymentAccountId: '',
  date: new Date().toISOString().split('T')[0]
})

const creditCards = computed(() => creditCardsStore.allCreditCards)
const accounts = computed(() => accountsStore.allAccounts)

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : accountId
}

const openPaymentModal = (card) => {
  paymentForm.value = {
    cardId: card.id,
    cardName: card.name,
    currentBalance: card.currentBalance,
    amount: card.currentBalance,
    paymentAccountId: '',
    date: new Date().toISOString().split('T')[0]
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
    date: new Date().toISOString().split('T')[0]
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
      paymentForm.value.date,
      paymentForm.value.paymentAccountId
    )

    closePaymentModal()
    // Refresh accounts to show updated balances
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
    form.value = { name: '', lastFourDigits: '', creditLimit: 0, apr: 0, currentBalance: 0 }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error adding credit card')
  }
}

onMounted(async () => {
  await Promise.all([
    creditCardsStore.fetchCreditCards(),
    accountsStore.fetchAccounts()
  ])
})
</script>
