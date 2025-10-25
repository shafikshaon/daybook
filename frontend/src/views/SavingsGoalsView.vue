<template>
  <div class="savings-goals-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Savings Goals</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Goal</button>
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
                  class="progress-bar bg-success"
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
              <button class="btn btn-sm btn-outline-success flex-fill" @click="addContribution(goal.id)">
                Add Contribution
              </button>
              <button class="btn btn-sm btn-outline-danger" @click="deleteGoal(goal.id)">Delete</button>
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
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="contributionAmount" required />
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showContributionModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Add</button>
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
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const savingsGoalsStore = useSavingsGoalsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()
const showAddModal = ref(false)
const showContributionModal = ref(false)
const selectedGoalId = ref(null)
const contributionAmount = ref(0)
const form = ref({ name: '', icon: '🎯', targetAmount: 0, currentAmount: 0, monthlyContribution: 0, targetDate: '' })

const savingsGoals = computed(() => savingsGoalsStore.activeSavingsGoals)
const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })

const addContribution = (goalId) => {
  selectedGoalId.value = goalId
  contributionAmount.value = 0
  showContributionModal.value = true
}

const saveContribution = async () => {
  try {
    await savingsGoalsStore.addContribution(selectedGoalId.value, contributionAmount.value)
    success('Contribution added successfully')
    showContributionModal.value = false
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

onMounted(() => savingsGoalsStore.fetchSavingsGoals())
</script>
