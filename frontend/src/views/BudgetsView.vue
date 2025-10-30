<template>
  <div class="budgets-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Budgets</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Budget</button>
    </div>

    <!-- Budget Summary -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalBudgeted) }}</div>
          <div class="stat-label">Total Budgeted</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalSpent) }}</div>
          <div class="stat-label">Total Spent</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon green">💵</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalBudgeted - budgetsStore.totalSpent) }}</div>
          <div class="stat-label">Remaining</div>
        </div>
      </div>
    </div>

    <!-- Budget Alerts -->
    <div v-if="budgetAlerts.length > 0" class="alert alert-warning mb-4">
      <h6 class="alert-heading">⚠️ Budget Alerts</h6>
      <ul class="mb-0">
        <li v-for="alert in budgetAlerts" :key="alert.id">
          {{ alert.message }}: {{ formatCurrency(alert.amount) }} / {{ formatCurrency(alert.budget) }}
        </li>
      </ul>
    </div>

    <!-- Budgets List -->
    <div class="row g-3">
      <div v-for="budget in budgets" :key="budget.id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ getCategoryName(budget.categoryId) }}</h5>
            <p class="text-muted mb-3">{{ budget.period }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span>{{ formatCurrency(getProgress(budget.id)?.spent || 0) }}</span>
                <span>{{ formatCurrency(budget.amount) }}</span>
              </div>
              <div class="progress" style="height: 12px;">
                <div
                  class="progress-bar"
                  :class="getProgressClass(getProgress(budget.id)?.status)"
                  :style="{ width: Math.min(getProgress(budget.id)?.percentage || 0, 100) + '%' }"
                ></div>
              </div>
              <small class="text-muted">
                {{ Math.round(getProgress(budget.id)?.percentage || 0) }}% used
              </small>
            </div>

            <div class="d-flex justify-content-between">
              <button class="btn btn-sm btn-outline-primary" @click="editBudget(budget)">Edit</button>
              <button class="btn btn-sm btn-danger" @click="deleteBudget(budget.id)">Delete</button>
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
            <h5 class="modal-title">{{ showEditModal ? 'Edit Budget' : 'Add Budget' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveBudget">
              <div class="mb-3">
                <label class="form-label">Category *</label>
                <select class="form-select" v-model="form.categoryId" required>
                  <option value="">Select category...</option>
                  <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.id">
                    {{ cat.icon }} {{ cat.name }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.amount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Period *</label>
                <select class="form-select" v-model="form.period" required>
                  <option v-for="period in budgetsStore.budgetPeriods" :key="period.value" :value="period.value">
                    {{ period.label }}
                  </option>
                </select>
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBudgetsStore } from '@/stores/budgets'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const budgetsStore = useBudgetsStore()
const transactionsStore = useTransactionsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingBudget = ref(null)
const form = ref({ categoryId: '', amount: 0, period: 'monthly' })

const budgets = computed(() => budgetsStore.activeBudgets)
const expenseCategories = computed(() => transactionsStore.expenseCategories)
const budgetAlerts = computed(() => budgetsStore.budgetAlerts)

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const getCategoryName = (id) => transactionsStore.getCategoryById(id)?.name || id
const getProgress = (id) => budgetsStore.budgetProgress(id)
const getProgressClass = (status) => {
  return status === 'danger' ? 'bg-danger' : status === 'warning' ? 'bg-warning' : 'bg-success'
}

const editBudget = (budget) => {
  editingBudget.value = budget
  form.value = { ...budget }
  showEditModal.value = true
}

const deleteBudget = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Budget',
    message: 'Are you sure you want to delete this budget? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await budgetsStore.deleteBudget(id)
      success('Budget deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting budget')
    }
  }
}

const saveBudget = async () => {
  try {
    if (showEditModal.value) {
      await budgetsStore.updateBudget(editingBudget.value.id, form.value)
      success('Budget updated successfully')
    } else {
      await budgetsStore.createBudget(form.value)
      success('Budget created successfully')
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving budget')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  form.value = { categoryId: '', amount: 0, period: 'monthly' }
}

onMounted(async () => {
  await Promise.all([
    budgetsStore.fetchBudgets(),
    transactionsStore.fetchTransactions()
  ])
})
</script>
