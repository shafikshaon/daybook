<template>
  <div
    class="modal fade"
    :class="{ 'show d-block': show }"
    tabindex="-1"
    style="background-color: rgba(0,0,0,0.5);"
    v-if="show"
  >
    <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Reconcile Account: {{ account?.name }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitReconciliation">
            <!-- Current Balance Info -->
            <div class="alert alert-info mb-4">
              <div class="row">
                <div class="col-md-6">
                  <strong>Current Book Balance:</strong><br>
                  <span class="fs-5">{{ formatCurrency(account?.balance || 0) }}</span>
                </div>
                <div class="col-md-6">
                  <strong>Last Reconciled:</strong><br>
                  <span>{{ account?.lastReconciled ? formatDate(account.lastReconciled) : 'Never' }}</span>
                </div>
              </div>
            </div>

            <!-- Reconciliation Form -->
            <div class="mb-3">
              <label class="form-label">Reconciliation Date *</label>
              <input
                type="date"
                class="form-control"
                v-model="form.reconciliationDate"
                required
              />
            </div>

            <div class="mb-3">
              <label class="form-label">Statement Balance *</label>
              <input
                type="number"
                step="0.01"
                class="form-control"
                v-model.number="form.statementBalance"
                required
                placeholder="Enter balance from bank statement"
              />
              <small class="text-muted">Enter the balance shown on your bank statement</small>
            </div>

            <div class="mb-3">
              <label class="form-label">Notes</label>
              <textarea
                class="form-control"
                v-model="form.notes"
                rows="2"
                placeholder="Optional notes about this reconciliation"
              ></textarea>
            </div>

            <!-- Difference Alert -->
            <div v-if="difference !== null" class="alert mb-4" :class="differenceAlertClass">
              <div class="d-flex justify-content-between align-items-center">
                <div>
                  <strong>Difference:</strong>
                  <span class="fs-5 ms-2">{{ formatCurrency(Math.abs(difference)) }}</span>
                  <span v-if="difference > 0" class="ms-2 text-success">
                    (Statement is higher)
                  </span>
                  <span v-else-if="difference < 0" class="ms-2 text-danger">
                    (Book balance is higher)
                  </span>
                  <span v-else class="ms-2 text-success">
                    ✓ Balanced
                  </span>
                </div>
              </div>
            </div>
            <div class="d-flex justify-content-end gap-2">
              <button type="button" class="btn btn-secondary" @click="closeModal">
                Cancel
              </button>
              <button type="submit" class="btn btn-primary" :disabled="submitting">
                <span v-if="submitting">Reconciling...</span>
                <span v-else>Reconcile Account</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import apiService from '@/services/api-backend'
import { useSettingsStore } from '@/stores/settings'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  account: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'reconciled'])

const settingsStore = useSettingsStore()

const loading = ref(false)
const submitting = ref(false)
const unreconciledTransactions = ref([])

const form = ref({
  reconciliationDate: new Date().toISOString().split('T')[0],
  statementBalance: 0,
  notes: '',
  transactionIds: []
})

const difference = computed(() => {
  if (!props.account || form.value.statementBalance === null) return null
  return form.value.statementBalance - (props.account.balance || 0)
})

const differenceAlertClass = computed(() => {
  if (difference.value === null) return 'alert-secondary'
  if (difference.value === 0) return 'alert-success'
  return 'alert-warning'
})

const allTransactionsSelected = computed(() => {
  return unreconciledTransactions.value.length > 0 &&
    unreconciledTransactions.value.every(t => form.value.transactionIds.includes(t.id))
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const loadUnreconciledTransactions = async () => {
  if (!props.account?.id) {
    console.warn('No account ID provided')
    return
  }

  console.log('Loading unreconciled transactions for account:', props.account.id)
  loading.value = true
  try {
    const response = await apiService.reconciliation.getUnreconciledTransactions(props.account.id)
    console.log('Unreconciled transactions response:', response)
    unreconciledTransactions.value = response.data || []
    console.log('Loaded transactions count:', unreconciledTransactions.value.length)
  } catch (error) {
    console.error('Failed to load unreconciled transactions:', error)
    alert('Failed to load transactions: ' + (error.message || 'Unknown error'))
  } finally {
    loading.value = false
  }
}

const toggleTransaction = (transactionId) => {
  const index = form.value.transactionIds.indexOf(transactionId)
  if (index > -1) {
    form.value.transactionIds.splice(index, 1)
  } else {
    form.value.transactionIds.push(transactionId)
  }
}

const toggleAllTransactions = () => {
  if (allTransactionsSelected.value) {
    form.value.transactionIds = []
  } else {
    form.value.transactionIds = unreconciledTransactions.value.map(t => t.id)
  }
}

const submitReconciliation = async () => {
  if (!props.account?.id) return

  submitting.value = true
  try {
    const reconciliationData = {
      accountId: props.account.id,
      reconciliationDate: new Date(form.value.reconciliationDate).toISOString(),
      statementBalance: form.value.statementBalance,
      notes: form.value.notes,
      transactionIds: form.value.transactionIds
    }

    await apiService.reconciliation.create(reconciliationData)
    emit('reconciled')
    closeModal()
  } catch (error) {
    console.error('Failed to create reconciliation:', error)
    alert(error.message || 'Failed to reconcile account')
  } finally {
    submitting.value = false
  }
}

const closeModal = () => {
  form.value = {
    reconciliationDate: new Date().toISOString().split('T')[0],
    statementBalance: 0,
    notes: '',
    transactionIds: []
  }
  unreconciledTransactions.value = []
  emit('close')
}

// Watch for modal show to load transactions
watch(() => props.show, (newVal) => {
  if (newVal && props.account?.id) {
    form.value.statementBalance = props.account.balance || 0
    loadUnreconciledTransactions()
  }
}, { immediate: true })
</script>

<style scoped>
.table thead.sticky-top {
  position: sticky;
  top: 0;
  z-index: 1;
}

.modal-dialog-scrollable .modal-body {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}
</style>
