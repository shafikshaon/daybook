<template>
  <div class="investments-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Investments & Fixed Deposits</h1>
      <button class="btn btn-primary" @click="handleAddClick">+ Add {{ activeTab === 'stocks' ? 'Investment' : 'Fixed Deposit' }}</button>
    </div>

    <!-- Tabs -->
    <ul class="nav nav-tabs mb-4">
      <li class="nav-item">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'stocks' }"
          @click="activeTab = 'stocks'"
        >
          📈 Stocks & Investments
        </button>
      </li>
      <li class="nav-item">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'fixed-deposits' }"
          @click="activeTab = 'fixed-deposits'"
        >
          💰 Fixed Deposits
        </button>
      </li>
    </ul>

    <!-- Stocks & Investments Tab -->
    <div v-if="activeTab === 'stocks'">
      <!-- Portfolio Summary -->
      <div class="row g-3 mb-4">
        <div class="col-12 col-md-3">
          <div class="stat-card">
            <div class="stat-icon purple">💎</div>
            <div class="stat-value">{{ formatCurrency(investmentsStore.totalCurrentValue) }}</div>
            <div class="stat-label">Portfolio Value</div>
          </div>
        </div>
        <div class="col-12 col-md-3">
          <div class="stat-card">
            <div class="stat-icon blue">📊</div>
            <div class="stat-value">{{ formatCurrency(investmentsStore.totalInvested) }}</div>
            <div class="stat-label">Total Invested</div>
          </div>
        </div>
        <div class="col-12 col-md-3">
          <div class="stat-card">
            <div class="stat-icon green">📈</div>
            <div class="stat-value">{{ formatCurrency(investmentsStore.totalGainLoss) }}</div>
            <div class="stat-label">Gain/Loss</div>
            <div :class="investmentsStore.totalGainLoss >= 0 ? 'stat-change positive' : 'stat-change negative'">
              {{ Math.round(investmentsStore.totalGainLossPercentage) }}%
            </div>
          </div>
        </div>
        <div class="col-12 col-md-3">
          <div class="stat-card">
            <div class="stat-icon orange">💰</div>
            <div class="stat-value">{{ formatCurrency(investmentsStore.totalDividendsEarned) }}</div>
            <div class="stat-label">Total Dividends</div>
          </div>
        </div>
      </div>

      <!-- Investments Table -->
      <div class="card">
        <div class="card-header">
          <h5 class="mb-0">Holdings</h5>
        </div>
        <div class="card-body p-0">
          <div v-if="investments.length === 0" class="p-4 text-center text-muted">No investments yet</div>
          <div v-else class="table-responsive">
            <table class="table table-hover mb-0">
              <thead>
                <tr>
                  <th>Symbol</th>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Quantity</th>
                  <th>Cost Basis</th>
                  <th>Current Price</th>
                  <th class="text-end">Value</th>
                  <th class="text-end">Gain/Loss</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="inv in investments" :key="inv.id">
                  <td class="fw-bold">{{ inv.symbol }}</td>
                  <td>{{ inv.name }}</td>
                  <td><span class="badge bg-secondary">{{ inv.assetType }}</span></td>
                  <td>{{ inv.quantity }}</td>
                  <td>{{ formatCurrency(inv.costBasis) }}</td>
                  <td>{{ formatCurrency(inv.currentPrice) }}</td>
                  <td class="text-end fw-bold">{{ formatCurrency(inv.quantity * inv.currentPrice) }}</td>
                  <td class="text-end">
                    <span :class="(inv.currentPrice - inv.costBasis) >= 0 ? 'text-success' : 'text-danger'">
                      {{ formatCurrency((inv.currentPrice - inv.costBasis) * inv.quantity) }}
                      ({{ Math.round(((inv.currentPrice - inv.costBasis) / inv.costBasis) * 100) }}%)
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Fixed Deposits Tab -->
    <div v-if="activeTab === 'fixed-deposits'">
      <!-- FD Summary -->
      <div class="row g-3 mb-4">
        <div class="col-12 col-md-4">
          <div class="stat-card">
            <div class="stat-icon purple">💰</div>
            <div class="stat-value">{{ formatCurrency(fixedDepositsStore.totalInvested) }}</div>
            <div class="stat-label">Total Invested</div>
          </div>
        </div>
        <div class="col-12 col-md-4">
          <div class="stat-card">
            <div class="stat-icon green">📈</div>
            <div class="stat-value">{{ formatCurrency(fixedDepositsStore.totalMaturityValue) }}</div>
            <div class="stat-label">Maturity Value</div>
          </div>
        </div>
        <div class="col-12 col-md-4">
          <div class="stat-card">
            <div class="stat-icon blue">💵</div>
            <div class="stat-value">{{ formatCurrency(fixedDepositsStore.totalInterestEarned) }}</div>
            <div class="stat-label">Interest Earned</div>
          </div>
        </div>
      </div>

      <!-- FD Cards -->
      <div v-if="fixedDeposits.length === 0" class="card">
        <div class="card-body text-center text-muted p-5">
          No fixed deposits yet
        </div>
      </div>
      <div v-else class="row g-3">
        <div v-for="fd in fixedDeposits" :key="fd.id" class="col-12 col-md-6">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">{{ fd.name }}</h5>
              <p class="text-muted">{{ fd.bank }}</p>
              <div class="mb-3">
                <div class="d-flex justify-content-between mb-2">
                  <span>Principal:</span>
                  <span class="fw-bold">{{ formatCurrency(fd.principal) }}</span>
                </div>
                <div class="d-flex justify-content-between mb-2">
                  <span>Interest Rate:</span>
                  <span class="fw-bold">{{ fd.interestRate }}%</span>
                </div>
                <div class="d-flex justify-content-between mb-2">
                  <span>Maturity Value:</span>
                  <span class="fw-bold text-success">{{ formatCurrency(fixedDepositsStore.calculateMaturityAmount(fd.id)) }}</span>
                </div>
                <div class="d-flex justify-content-between mb-2">
                  <span>Maturity Date:</span>
                  <span>{{ formatDate(fd.maturityDate) }}</span>
                </div>
                <div class="d-flex justify-content-between">
                  <span>Days Until Maturity:</span>
                  <span class="fw-bold">{{ fixedDepositsStore.getDaysUntilMaturity(fd.id) }} days</span>
                </div>
              </div>
              <button v-if="!fd.withdrawn" class="btn btn-sm btn-outline-danger" @click="withdrawFD(fd.id)">
                Withdraw
              </button>
              <span v-else class="badge bg-secondary">Withdrawn</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Investment Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddInvestmentModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddInvestmentModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Investment</h5>
            <button type="button" class="btn-close" @click="showAddInvestmentModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveInvestment">
              <div class="mb-3">
                <label class="form-label">Symbol *</label>
                <input type="text" class="form-control" v-model="investmentForm.symbol" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Name *</label>
                <input type="text" class="form-control" v-model="investmentForm.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Asset Type *</label>
                <select class="form-select" v-model="investmentForm.assetType" required>
                  <option v-for="type in investmentsStore.assetTypes" :key="type.value" :value="type.value">
                    {{ type.icon }} {{ type.label }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Quantity *</label>
                <input type="number" step="0.00001" class="form-control" v-model.number="investmentForm.quantity" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Cost Basis *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="investmentForm.costBasis" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Current Price *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="investmentForm.currentPrice" required />
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showAddInvestmentModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Create</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Fixed Deposit Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddFDModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddFDModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Fixed Deposit</h5>
            <button type="button" class="btn-close" @click="showAddFDModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveFD">
              <div class="mb-3">
                <label class="form-label">FD Name *</label>
                <input type="text" class="form-control" v-model="fdForm.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Bank *</label>
                <input type="text" class="form-control" v-model="fdForm.bank" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Principal Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="fdForm.principal" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Interest Rate (%) *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="fdForm.interestRate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Start Date *</label>
                <input type="date" class="form-control" v-model="fdForm.startDate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Tenure (Months) *</label>
                <input type="number" class="form-control" v-model.number="fdForm.tenureMonths" required @change="calculateMaturityDate" />
              </div>
              <div class="mb-3">
                <label class="form-label">Maturity Date</label>
                <input type="date" class="form-control" v-model="fdForm.maturityDate" readonly />
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="showAddFDModal = false">Cancel</button>
                <button type="submit" class="btn btn-primary">Create</button>
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
import { useInvestmentsStore } from '@/stores/investments'
import { useFixedDepositsStore } from '@/stores/fixedDeposits'
import { useSettingsStore } from '@/stores/settings'

const investmentsStore = useInvestmentsStore()
const fixedDepositsStore = useFixedDepositsStore()
const settingsStore = useSettingsStore()

const activeTab = ref('stocks')
const showAddInvestmentModal = ref(false)
const showAddFDModal = ref(false)

const investmentForm = ref({ symbol: '', name: '', assetType: 'stocks', quantity: 0, costBasis: 0, currentPrice: 0 })
const fdForm = ref({
  name: '', bank: '', principal: 0, interestRate: 0,
  startDate: new Date().toISOString().split('T')[0],
  tenureMonths: 12, maturityDate: '', compounding: 'monthly'
})

const investments = computed(() => investmentsStore.allInvestments)
const fixedDeposits = computed(() => fixedDepositsStore.allFixedDeposits)
const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })

const handleAddClick = () => {
  if (activeTab.value === 'stocks') {
    showAddInvestmentModal.value = true
  } else {
    showAddFDModal.value = true
  }
}

const saveInvestment = async () => {
  await investmentsStore.createInvestment(investmentForm.value)
  showAddInvestmentModal.value = false
  investmentForm.value = { symbol: '', name: '', assetType: 'stocks', quantity: 0, costBasis: 0, currentPrice: 0 }
}

const calculateMaturityDate = () => {
  const start = new Date(fdForm.value.startDate)
  start.setMonth(start.getMonth() + fdForm.value.tenureMonths)
  fdForm.value.maturityDate = start.toISOString().split('T')[0]
}

const withdrawFD = async (id) => {
  if (confirm('Withdraw this FD?')) {
    await fixedDepositsStore.withdrawFixedDeposit(id)
  }
}

const saveFD = async () => {
  await fixedDepositsStore.createFixedDeposit(fdForm.value)
  showAddFDModal.value = false
  fdForm.value = {
    name: '', bank: '', principal: 0, interestRate: 0,
    startDate: new Date().toISOString().split('T')[0],
    tenureMonths: 12, maturityDate: '', compounding: 'monthly'
  }
}

onMounted(async () => {
  await Promise.all([
    investmentsStore.fetchInvestments(),
    investmentsStore.fetchDividends(),
    fixedDepositsStore.fetchFixedDeposits()
  ])
  calculateMaturityDate()
})
</script>

<style scoped>
.nav-tabs {
  border-bottom: 2px solid #e3e8ee;
}

.nav-tabs .nav-link {
  border: none;
  background: none;
  color: #64748b;
  font-weight: 500;
  padding: 0.75rem 1.5rem;
  border-bottom: 3px solid transparent;
  transition: all 0.2s;
}

.nav-tabs .nav-link:hover {
  color: #1e293b;
  border-bottom-color: #cbd5e1;
}

.nav-tabs .nav-link.active {
  color: #635bff;
  border-bottom-color: #635bff;
  background: none;
}
</style>
