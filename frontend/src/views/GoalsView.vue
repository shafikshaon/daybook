<template>
  <div class="goals-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Financial Goals</h1>
      <button class="btn btn-primary" @click="showAddGoalModal = true">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
        </svg>
        Add Goal
      </button>
    </div>

    <!-- Filters -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <select class="form-select" v-model="filters.status" @change="applyFilters">
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="achieved">Achieved</option>
          <option value="paused">Paused</option>
          <option value="archived">Archived</option>
        </select>
      </div>
      <div class="col-12 col-md-3">
        <select class="form-select" v-model="filters.category" @change="applyFilters">
          <option value="">All Categories</option>
          <option value="emergency_fund">Emergency Fund</option>
          <option value="vacation">Vacation</option>
          <option value="retirement">Retirement</option>
          <option value="home">Home</option>
          <option value="education">Education</option>
          <option value="car">Car</option>
          <option value="wedding">Wedding</option>
          <option value="business">Business</option>
          <option value="other">Other</option>
        </select>
      </div>
      <div class="col-12 col-md-3">
        <select class="form-select" v-model="filters.priority" @change="applyFilters">
          <option value="">All Priorities</option>
          <option value="high">High Priority</option>
          <option value="medium">Medium Priority</option>
          <option value="low">Low Priority</option>
        </select>
      </div>
    </div>

    <!-- Stats -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">🎯</div>
          <div class="stat-value">{{ activeGoals.length }}</div>
          <div class="stat-label">Active Goals</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(goalsStore.totalCurrentAmount) }}</div>
          <div class="stat-label">Total Value</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">📊</div>
          <div class="stat-value">{{ Math.round(goalsStore.totalProgress) }}%</div>
          <div class="stat-label">Overall Progress</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon orange">🔥</div>
          <div class="stat-value">{{ goalsStore.highPriorityGoals.length }}</div>
          <div class="stat-label">High Priority</div>
        </div>
      </div>
    </div>

    <!-- Goals List -->
    <div class="row g-3">
      <div v-for="goal in goals" :key="goal.id" class="col-12 col-md-6 col-lg-4">
        <div class="card goal-card" @click="viewGoalDetails(goal.id)">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-3">
              <div>
                <span style="font-size: 2rem;">{{ goal.icon || '🎯' }}</span>
                <h5 class="mt-2">{{ goal.name }}</h5>
                <div class="d-flex gap-2 mt-1">
                  <span class="badge" :class="categoryBadgeClass(goal.category)">
                    {{ formatCategory(goal.category) }}
                  </span>
                  <span class="badge" :class="priorityBadgeClass(goal.priority)">
                    {{ formatPriority(goal.priority) }}
                  </span>
                  <span class="badge" :class="statusBadgeClass(goal.status)">
                    {{ formatStatus(goal.status) }}
                  </span>
                </div>
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
                <small class="text-muted">{{ formatCurrency(Math.max(0, goal.targetAmount - goal.currentAmount)) }} remaining</small>
              </div>
            </div>

            <div class="mb-3">
              <small class="text-muted">
                <strong>{{ goal.holdings?.length || 0 }}</strong> Holdings
              </small>
              <br>
              <small class="text-muted" v-if="goal.targetDate">
                Target: {{ formatDate(goal.targetDate) }}
              </small>
            </div>

            <div class="d-flex gap-2" @click.stop>
              <button class="btn btn-sm btn-add-holding flex-fill" @click="addHolding(goal.id)">
                Add Holding
              </button>
              <button class="btn btn-sm btn-view-goal" @click="viewGoalDetails(goal.id)">
                Details
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="goals.length === 0" class="col-12">
        <div class="card">
          <div class="card-body text-center py-5">
            <div style="font-size: 4rem; opacity: 0.5;">🎯</div>
            <h5 class="mt-3">No goals found</h5>
            <p class="text-muted">Create your first financial goal to start tracking your progress</p>
            <button class="btn btn-primary mt-3" @click="showAddGoalModal = true">
              Create Goal
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Goal Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddGoalModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddGoalModal">
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Create Financial Goal</h5>
            <button type="button" class="btn-close" @click="showAddGoalModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveGoal">
              <div class="row">
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Goal Name *</label>
                  <input type="text" class="form-control" v-model="goalForm.name" required />
                </div>
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Icon</label>
                  <input type="text" class="form-control" v-model="goalForm.icon" placeholder="e.g., 🏠" />
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea class="form-control" v-model="goalForm.description" rows="2"></textarea>
              </div>

              <div class="row">
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Category *</label>
                  <select class="form-select" v-model="goalForm.category" required>
                    <option value="">Select...</option>
                    <option value="emergency_fund">Emergency Fund</option>
                    <option value="vacation">Vacation</option>
                    <option value="retirement">Retirement</option>
                    <option value="home">Home</option>
                    <option value="education">Education</option>
                    <option value="car">Car</option>
                    <option value="wedding">Wedding</option>
                    <option value="business">Business</option>
                    <option value="other">Other</option>
                  </select>
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Priority *</label>
                  <select class="form-select" v-model="goalForm.priority" required>
                    <option value="">Select...</option>
                    <option value="high">High</option>
                    <option value="medium">Medium</option>
                    <option value="low">Low</option>
                  </select>
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Color</label>
                  <input type="color" class="form-control form-control-color" v-model="goalForm.color" />
                </div>
              </div>

              <div class="row">
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Target Amount *</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="goalForm.targetAmount" required />
                </div>
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Target Date</label>
                  <input type="date" class="form-control" v-model="goalForm.targetDate" />
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Monthly Contribution Target</label>
                <input type="number" step="0.01" class="form-control" v-model.number="goalForm.monthlyContribution" />
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showAddGoalModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Create Goal</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Goal Detail Modal -->
    <div class="modal fade" :class="{ 'show d-block': showDetailModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showDetailModal">
      <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              {{ selectedGoal?.icon }} {{ selectedGoal?.name }}
            </h5>
            <button type="button" class="btn-close" @click="showDetailModal = false"></button>
          </div>
          <div class="modal-body" v-if="selectedGoal">
            <!-- Goal Progress -->
            <div class="mb-4">
              <div class="d-flex justify-content-between mb-2">
                <span class="fw-bold">Progress</span>
                <span class="text-muted">{{ Math.round((selectedGoal.currentAmount / selectedGoal.targetAmount) * 100) }}%</span>
              </div>
              <div class="progress" style="height: 20px;">
                <div
                  class="progress-bar progress-bar-professional"
                  :style="{ width: Math.min((selectedGoal.currentAmount / selectedGoal.targetAmount) * 100, 100) + '%' }"
                >
                  {{ formatCurrency(selectedGoal.currentAmount) }} / {{ formatCurrency(selectedGoal.targetAmount) }}
                </div>
              </div>
            </div>

            <!-- Goal Info -->
            <div class="row mb-4">
              <div class="col-6 col-md-3">
                <small class="text-muted">Category</small>
                <div class="fw-bold">{{ formatCategory(selectedGoal.category) }}</div>
              </div>
              <div class="col-6 col-md-3">
                <small class="text-muted">Priority</small>
                <div class="fw-bold">{{ formatPriority(selectedGoal.priority) }}</div>
              </div>
              <div class="col-6 col-md-3">
                <small class="text-muted">Status</small>
                <div class="fw-bold">{{ formatStatus(selectedGoal.status) }}</div>
              </div>
              <div class="col-6 col-md-3">
                <small class="text-muted">Target Date</small>
                <div class="fw-bold">{{ selectedGoal.targetDate ? formatDate(selectedGoal.targetDate) : 'Not set' }}</div>
              </div>
            </div>

            <!-- Holdings List -->
            <div class="d-flex justify-content-between align-items-center mb-3">
              <h6>Holdings ({{ selectedGoal.holdings?.length || 0 }})</h6>
              <button class="btn btn-sm btn-primary" @click="addHolding(selectedGoal.id)">
                Add Holding
              </button>
            </div>

            <div v-if="selectedGoal.holdings && selectedGoal.holdings.length > 0" class="table-responsive">
              <table class="table">
                <thead>
                  <tr>
                    <th>Type</th>
                    <th>Name</th>
                    <th>Amount</th>
                    <th>Current Value</th>
                    <th>Gain/Loss</th>
                    <th>Date</th>
                    <th>Status</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="holding in selectedGoal.holdings" :key="holding.id">
                    <td>
                      <span class="badge bg-secondary">
                        {{ goalsStore.getHoldingTypeLabel(holding.type) }}
                      </span>
                    </td>
                    <td>{{ holding.name }}</td>
                    <td>{{ formatCurrency(holding.amount) }}</td>
                    <td>{{ formatCurrency(holding.currentValue) }}</td>
                    <td :class="getGainLossClass(holding)">
                      {{ formatGainLoss(holding) }}
                    </td>
                    <td>{{ formatDate(holding.purchaseDate) }}</td>
                    <td>
                      <span class="badge" :class="holdingStatusBadgeClass(holding.status)">
                        {{ holding.status }}
                      </span>
                    </td>
                    <td>
                      <button class="btn btn-sm btn-outline-danger" @click="removeHolding(holding.id)">
                        Remove
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="text-center py-4 text-muted">
              <p>No holdings yet. Add your first investment or savings to this goal.</p>
            </div>

            <!-- Actions -->
            <div class="d-flex gap-2 mt-4">
              <button class="btn btn-warning" v-if="selectedGoal.status === 'active'" @click="pauseGoal(selectedGoal.id)">
                Pause Goal
              </button>
              <button class="btn btn-success" v-if="selectedGoal.status === 'paused'" @click="resumeGoal(selectedGoal.id)">
                Resume Goal
              </button>
              <button class="btn btn-secondary" @click="archiveGoal(selectedGoal.id)">
                Archive Goal
              </button>
              <button class="btn btn-danger ms-auto" @click="deleteGoal(selectedGoal.id)">
                Delete Goal
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Holding Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddHoldingModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddHoldingModal">
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Holding to Goal</h5>
            <button type="button" class="btn-close" @click="showAddHoldingModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveHolding">
              <div class="mb-3">
                <label class="form-label">From Account *</label>
                <select class="form-select" v-model="holdingForm.accountId" required>
                  <option value="">Select account...</option>
                  <option v-for="account in accounts" :key="account.id" :value="account.id">
                    {{ account.name }} ({{ formatCurrency(account.balance) }})
                  </option>
                </select>
              </div>

              <div class="row">
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Holding Type *</label>
                  <select class="form-select" v-model="holdingForm.type" required @change="onHoldingTypeChange">
                    <option value="">Select type...</option>
                    <optgroup v-for="(types, category) in holdingTypes" :key="category" :label="category">
                      <option v-for="type in types" :key="type.value" :value="type.value">
                        {{ type.label }}
                      </option>
                    </optgroup>
                  </select>
                </div>
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Holding Name *</label>
                  <input type="text" class="form-control" v-model="holdingForm.name" required />
                </div>
              </div>

              <div class="row">
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Amount *</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="holdingForm.amount" required />
                </div>
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Purchase Date & Time *</label>
                  <input type="datetime-local" class="form-control" v-model="holdingForm.purchaseDate" required />
                </div>
              </div>

              <!-- Market Instruments Fields -->
              <div v-if="isMarketInstrument" class="row">
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Symbol</label>
                  <input type="text" class="form-control" v-model="holdingForm.symbol" placeholder="e.g., AAPL" />
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Quantity</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="holdingForm.quantity" />
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Cost Per Unit</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="holdingForm.costBasis" />
                </div>
              </div>

              <!-- Bank Products Fields -->
              <div v-if="isBankProduct" class="row">
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Institution</label>
                  <input type="text" class="form-control" v-model="holdingForm.institution" />
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Interest Rate (%)</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="holdingForm.interestRate" />
                </div>
                <div class="col-12 col-md-4 mb-3">
                  <label class="form-label">Tenure (months)</label>
                  <input type="number" class="form-control" v-model.number="holdingForm.tenureMonths" />
                </div>
              </div>

              <div v-if="isBankProduct" class="row">
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Maturity Date</label>
                  <input type="date" class="form-control" v-model="holdingForm.maturityDate" />
                </div>
                <div class="col-12 col-md-6 mb-3">
                  <label class="form-label">Maturity Amount</label>
                  <input type="number" step="0.01" class="form-control" v-model.number="holdingForm.maturityAmount" />
                </div>
              </div>

              <div class="alert alert-info">
                <small>
                  This will deduct {{ formatCurrency(holdingForm.amount || 0) }} from
                  {{ getAccountName(holdingForm.accountId) || 'selected account' }} and add it as a holding to your goal.
                  A transaction record will be created.
                </small>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showAddHoldingModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Add Holding</button>
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
import { useGoalsStore } from '@/stores/goals'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'

const goalsStore = useGoalsStore()
const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const showAddGoalModal = ref(false)
const showDetailModal = ref(false)
const showAddHoldingModal = ref(false)
const selectedGoal = ref(null)
const selectedGoalIdForHolding = ref(null)
const holdingTypes = ref(null)

const filters = ref({
  status: '',
  category: '',
  priority: ''
})

const goalForm = ref({
  name: '',
  description: '',
  icon: '🎯',
  color: '#3b82f6',
  category: '',
  priority: '',
  targetAmount: 0,
  targetDate: '',
  monthlyContribution: 0
})

const holdingForm = ref({
  accountId: '',
  type: '',
  name: '',
  amount: 0,
  purchaseDate: new Date().toISOString().slice(0, 16), // Format: YYYY-MM-DDTHH:mm
  // Market instruments
  symbol: '',
  quantity: null,
  costBasis: null,
  // Bank products
  institution: '',
  interestRate: null,
  tenureMonths: null,
  maturityDate: '',
  maturityAmount: null
})

const goals = computed(() => goalsStore.allGoals)
const activeGoals = computed(() => goalsStore.activeGoals)
const accounts = computed(() => accountsStore.allAccounts)

const isMarketInstrument = computed(() => {
  const marketTypes = ['stocks', 'mutual_fund', 'etf', 'index_fund', 'bonds', 'cryptocurrency', 'commodities', 'gold']
  return marketTypes.includes(holdingForm.value.type)
})

const isBankProduct = computed(() => {
  const bankTypes = ['fixed_deposit', 'dps', 'recurring_deposit', 'savings_bond', 'ppf', 'nsc']
  return bankTypes.includes(holdingForm.value.type)
})

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleString('en-US', {
  year: 'numeric',
  month: 'short',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  hour12: false
})

const formatCategory = (category) => {
  const map = {
    emergency_fund: 'Emergency Fund',
    vacation: 'Vacation',
    retirement: 'Retirement',
    home: 'Home',
    education: 'Education',
    car: 'Car',
    wedding: 'Wedding',
    business: 'Business',
    other: 'Other'
  }
  return map[category] || category
}

const formatPriority = (priority) => {
  const map = { high: 'High', medium: 'Medium', low: 'Low' }
  return map[priority] || priority
}

const formatStatus = (status) => {
  const map = { active: 'Active', achieved: 'Achieved', paused: 'Paused', archived: 'Archived' }
  return map[status] || status
}

const categoryBadgeClass = (category) => 'bg-info'
const priorityBadgeClass = (priority) => {
  const map = { high: 'bg-danger', medium: 'bg-warning', low: 'bg-secondary' }
  return map[priority] || 'bg-secondary'
}
const statusBadgeClass = (status) => {
  const map = { active: 'bg-success', achieved: 'bg-primary', paused: 'bg-warning', archived: 'bg-secondary' }
  return map[status] || 'bg-secondary'
}

const holdingStatusBadgeClass = (status) => {
  const map = { active: 'bg-success', matured: 'bg-primary', sold: 'bg-info', closed: 'bg-secondary', withdrawn: 'bg-warning' }
  return map[status] || 'bg-secondary'
}

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : ''
}

const formatGainLoss = (holding) => {
  const gain = (holding.currentValue || holding.amount) - holding.amount
  const gainPercent = holding.amount > 0 ? (gain / holding.amount) * 100 : 0
  const sign = gain >= 0 ? '+' : ''
  return `${sign}${formatCurrency(gain)} (${sign}${gainPercent.toFixed(2)}%)`
}

const getGainLossClass = (holding) => {
  const gain = (holding.currentValue || holding.amount) - holding.amount
  return gain >= 0 ? 'text-success' : 'text-danger'
}

const applyFilters = async () => {
  await goalsStore.fetchGoals(filters.value)
}

const viewGoalDetails = async (goalId) => {
  try {
    const goal = await goalsStore.fetchGoal(goalId)
    selectedGoal.value = goal
    showDetailModal.value = true
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error fetching goal details')
  }
}

const addHolding = (goalId) => {
  selectedGoalIdForHolding.value = goalId
  holdingForm.value = {
    accountId: '',
    type: '',
    name: '',
    amount: 0,
    purchaseDate: new Date().toISOString().slice(0, 16), // Format: YYYY-MM-DDTHH:mm
    symbol: '',
    quantity: null,
    costBasis: null,
    institution: '',
    interestRate: null,
    tenureMonths: null,
    maturityDate: '',
    maturityAmount: null
  }
  showAddHoldingModal.value = true
}

const onHoldingTypeChange = () => {
  // Auto-fill name based on type if empty
  if (!holdingForm.value.name && holdingForm.value.type) {
    holdingForm.value.name = goalsStore.getHoldingTypeLabel(holdingForm.value.type)
  }
}

const saveGoal = async () => {
  try {
    await goalsStore.createGoal(goalForm.value)
    success('Goal created successfully')
    showAddGoalModal.value = false
    goalForm.value = {
      name: '',
      description: '',
      icon: '🎯',
      color: '#3b82f6',
      category: '',
      priority: '',
      targetAmount: 0,
      targetDate: '',
      monthlyContribution: 0
    }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error creating goal')
  }
}

const saveHolding = async () => {
  try {
    if (!holdingForm.value.accountId) {
      error('Please select an account')
      return
    }

    if (holdingForm.value.amount <= 0) {
      error('Amount must be greater than 0')
      return
    }

    // Check account balance
    const account = accountsStore.getAccountById(holdingForm.value.accountId)
    if (!account) {
      error('Selected account not found')
      return
    }

    let accountBalance = typeof account.balance === 'string'
      ? parseFloat(account.balance.replace(/[^\d.-]/g, ''))
      : Number(account.balance)

    if (accountBalance < holdingForm.value.amount) {
      error(`Insufficient balance. Available: ${formatCurrency(accountBalance)}`)
      return
    }

    // Prepare holding data
    const holdingData = {
      ...holdingForm.value,
      currentValue: holdingForm.value.amount // Initial value equals amount
    }

    await goalsStore.addHolding(selectedGoalIdForHolding.value, holdingData)

    success('Holding added successfully')
    showAddHoldingModal.value = false

    // Refresh data
    await Promise.all([
      accountsStore.fetchAccounts(),
      goalsStore.fetchGoals(),
      transactionsStore.fetchTransactions(1, 20)
    ])

    // Refresh detail modal if open
    if (showDetailModal.value && selectedGoal.value) {
      await viewGoalDetails(selectedGoal.value.id)
    }
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error adding holding')
  }
}

const removeHolding = async (holdingId) => {
  const confirmed = await confirm({
    title: 'Remove Holding',
    message: 'This will liquidate/sell this holding and credit the money back to your account. Continue?',
    confirmText: 'Remove',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    // Show account selector
    const accountId = await selectAccount('Select account to receive funds')
    if (!accountId) return

    try {
      await goalsStore.removeHolding(holdingId, accountId)
      success('Holding removed successfully')

      // Refresh data
      await Promise.all([
        accountsStore.fetchAccounts(),
        goalsStore.fetchGoals(),
        transactionsStore.fetchTransactions(1, 20)
      ])

      // Refresh detail modal
      if (showDetailModal.value && selectedGoal.value) {
        await viewGoalDetails(selectedGoal.value.id)
      }
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error removing holding')
    }
  }
}

const selectAccount = async (title) => {
  // Simple implementation - in production, use a proper modal
  const accountId = prompt(`${title}\n\nEnter account ID from: ${accounts.value.map(a => `${a.name} (${a.id})`).join(', ')}`)
  return accountId
}

const pauseGoal = async (goalId) => {
  try {
    await goalsStore.pauseGoal(goalId)
    success('Goal paused')
    await viewGoalDetails(goalId)
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error pausing goal')
  }
}

const resumeGoal = async (goalId) => {
  try {
    await goalsStore.resumeGoal(goalId)
    success('Goal resumed')
    await viewGoalDetails(goalId)
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error resuming goal')
  }
}

const archiveGoal = async (goalId) => {
  const confirmed = await confirm({
    title: 'Archive Goal',
    message: 'Archive this goal? You can restore it later.',
    confirmText: 'Archive',
    cancelText: 'Cancel'
  })

  if (confirmed) {
    try {
      await goalsStore.archiveGoal(goalId)
      success('Goal archived')
      showDetailModal.value = false
      await goalsStore.fetchGoals()
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error archiving goal')
    }
  }
}

const deleteGoal = async (goalId) => {
  const confirmed = await confirm({
    title: 'Delete Goal',
    message: 'Are you sure you want to delete this goal? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await goalsStore.deleteGoal(goalId)
      success('Goal deleted successfully')
      showDetailModal.value = false
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting goal')
    }
  }
}

onMounted(async () => {
  await Promise.all([
    goalsStore.fetchGoals(),
    accountsStore.fetchAccounts(),
    goalsStore.fetchHoldingTypes().then(types => {
      holdingTypes.value = types
    })
  ])
})
</script>

<style scoped>
.goal-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.goal-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.progress-bar-professional {
  background-color: #3b82f6;
  background-image: linear-gradient(45deg, rgba(255, 255, 255, 0.15) 25%, transparent 25%, transparent 50%, rgba(255, 255, 255, 0.15) 50%, rgba(255, 255, 255, 0.15) 75%, transparent 75%, transparent);
  background-size: 1rem 1rem;
}

.btn-add-holding {
  color: #1e40af;
  border-color: #3b82f6;
  background-color: #eff6ff;
}

.btn-add-holding:hover {
  color: #ffffff;
  background-color: #3b82f6;
  border-color: #2563eb;
}

.btn-view-goal {
  color: #059669;
  border-color: #10b981;
  background-color: #ecfdf5;
}

.btn-view-goal:hover {
  color: #ffffff;
  background-color: #10b981;
  border-color: #059669;
}

.dark-mode .progress-bar-professional {
  background-color: #3b82f6;
}

.dark-mode .btn-add-holding {
  color: #93c5fd;
  border-color: #3b82f6;
  background-color: #1e3a5f;
}

.dark-mode .btn-add-holding:hover {
  color: #ffffff;
  background-color: #3b82f6;
  border-color: #60a5fa;
}

.dark-mode .btn-view-goal {
  color: #6ee7b7;
  border-color: #10b981;
  background-color: #064e3b;
}

.dark-mode .btn-view-goal:hover {
  color: #ffffff;
  background-color: #10b981;
  border-color: #34d399;
}
</style>
