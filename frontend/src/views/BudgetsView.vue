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

    <!-- Debug Info -->
    <div v-if="budgets.length === 0" class="alert alert-info mb-4">
      <p class="mb-0">No budgets found. Click "Add Budget" to create your first budget.</p>
    </div>

    <!-- Budgets List -->
    <div class="row g-3">
      <div v-for="item in budgets" :key="(item.budget || item).id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ getCategoryName((item.budget || item).categoryId) }}</h5>
            <p class="text-muted mb-3">{{ (item.budget || item).period }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span>{{ formatCurrency(item.totalSpent || 0) }}</span>
                <span>{{ formatCurrency((item.budget || item).amount) }}</span>
              </div>
              <div class="progress" style="height: 12px;">
                <div
                  class="progress-bar"
                  :class="getProgressClass(item)"
                  :style="{ width: Math.min(item.percentageUsed || 0, 100) + '%' }"
                ></div>
              </div>
              <small class="text-muted">
                {{ Math.round(item.percentageUsed || 0) }}% used
              </small>
            </div>

            <div class="mb-2">
              <small class="text-muted">
                Remaining: {{ formatCurrency(item.remaining || 0) }}
              </small>
            </div>

            <div v-if="item.alertTriggered || item.isOverBudget" class="mb-2">
              <span v-if="item.isOverBudget" class="badge bg-danger">
                ⛔ Over Budget
              </span>
              <span v-else-if="item.alertTriggered" class="badge bg-warning text-dark">
                ⚠️ Alert
              </span>
            </div>

            <div class="d-flex justify-content-between">
              <button class="btn btn-sm btn-outline-primary" @click="editBudget(item.budget || item)">Edit</button>
              <button class="btn btn-sm btn-danger" @click="deleteBudget((item.budget || item).id)">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Budget Modal -->
    <div v-if="showAddModal || showEditModal" class="modal show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ editingBudget ? 'Edit Budget' : 'Add Budget' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveBudget">
              <div class="mb-3">
                <label for="category" class="form-label">Category</label>
                <select v-model="budgetForm.categoryId" class="form-select" id="category" required>
                  <option value="">Select a category</option>
                  <option v-for="category in categories" :key="category.id" :value="category.id">
                    {{ category.name }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label for="amount" class="form-label">Budget Amount</label>
                <input v-model.number="budgetForm.amount" type="number" class="form-control" id="amount" min="0" step="0.01" required>
              </div>
              <div class="mb-3">
                <label for="period" class="form-label">Budget Period</label>
                <select v-model="budgetForm.period" class="form-select" id="period" required>
                  <option value="monthly">Monthly</option>
                  <option value="weekly">Weekly</option>
                  <option value="yearly">Yearly</option>
                </select>
              </div>
              <div class="mb-3">
                <label for="alertThreshold" class="form-label">Alert Threshold (%)</label>
                <input v-model.number="budgetForm.alertThreshold" type="number" class="form-control" id="alertThreshold" min="0" max="100" placeholder="80">
              </div>
              <div class="form-check mb-3">
                <input v-model="budgetForm.enabled" class="form-check-input" type="checkbox" id="enabled">
                <label class="form-check-label" for="enabled">
                  Enable this budget
                </label>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="button" class="btn btn-primary" @click="saveBudget">
              {{ editingBudget ? 'Update' : 'Create' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useBudgetsStore } from '@/stores/budgets'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'

export default {
  name: 'BudgetsView',
  setup() {
    const budgetsStore = useBudgetsStore()
    const transactionsStore = useTransactionsStore()
    const settingsStore = useSettingsStore()

    const showAddModal = ref(false)
    const showEditModal = ref(false)
    const editingBudget = ref(null)

    const budgetForm = ref({
      categoryId: '',
      amount: 0,
      period: 'monthly',
      alertThreshold: 80,
      enabled: true
    })

    const budgets = computed(() => budgetsStore.activeBudgets)
    const budgetAlerts = computed(() => budgetsStore.budgetAlerts)
    const categories = computed(() => transactionsStore.expenseCategories)

    const formatCurrency = (amount) => {
      return settingsStore.formatCurrency(amount)
    }

    const getCategoryName = (categoryId) => {
      const category = transactionsStore.getCategoryById(categoryId)
      return category?.name || 'Unknown'
    }

    const getProgressClass = (item) => {
      if (item.isOverBudget) return 'bg-danger'
      if (item.alertTriggered) return 'bg-warning'
      return 'bg-success'
    }

    const editBudget = (budget) => {
      editingBudget.value = budget
      budgetForm.value = {
        categoryId: budget.categoryId,
        amount: budget.amount,
        period: budget.period,
        alertThreshold: budget.alertThreshold || 80,
        enabled: budget.enabled !== false
      }
      showEditModal.value = true
    }

    const deleteBudget = async (budgetId) => {
      if (confirm('Are you sure you want to delete this budget?')) {
        try {
          await budgetsStore.deleteBudget(budgetId)
        } catch (error) {
          console.error('Error deleting budget:', error)
        }
      }
    }

    const saveBudget = async () => {
      try {
        if (editingBudget.value) {
          await budgetsStore.updateBudget(editingBudget.value.id, budgetForm.value)
        } else {
          await budgetsStore.createBudget(budgetForm.value)
        }
        closeModal()
      } catch (error) {
        console.error('Error saving budget:', error)
      }
    }

    const closeModal = () => {
      showAddModal.value = false
      showEditModal.value = false
      editingBudget.value = null
      budgetForm.value = {
        categoryId: '',
        amount: 0,
        period: 'monthly',
        alertThreshold: 80,
        enabled: true
      }
    }

    onMounted(async () => {
      try {
        await Promise.all([
          budgetsStore.fetchBudgets(),
          transactionsStore.fetchCategories()
        ])
      } catch (error) {
        console.error('Error loading data:', error)
      }
    })

    return {
      budgets,
      budgetAlerts,
      categories,
      budgetsStore,
      showAddModal,
      showEditModal,
      editingBudget,
      budgetForm,
      formatCurrency,
      getCategoryName,
      getProgressClass,
      editBudget,
      deleteBudget,
      saveBudget,
      closeModal
    }
  }
}
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  border-left: 4px solid var(--bs-primary);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.stat-icon.purple { background: rgba(111, 66, 193, 0.1); }
.stat-icon.red { background: rgba(239, 68, 68, 0.1); }
.stat-icon.green { background: rgba(16, 185, 129, 0.1); }

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--bs-dark);
  margin-bottom: 0.25rem;
}

.stat-label {
  color: var(--bs-secondary);
  font-size: 0.875rem;
  font-weight: 500;
}

.card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.progress {
  border-radius: 6px;
}

.text-purple {
  color: var(--bs-primary) !important;
}
</style>