<template>
  <div class="fixed-deposits-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Fixed Deposits</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add FD</button>
    </div>

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

    <div class="row g-3">
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
            <div class="d-flex gap-2 align-items-center">
              <button v-if="!fd.withdrawn" class="btn btn-sm btn-outline-danger" @click="withdrawFD(fd.id)">
                Withdraw
              </button>
              <span v-else class="badge bg-secondary">Withdrawn</span>
              <span
                v-if="fd.attachments && fd.attachments.length > 0"
                class="badge bg-info"
                style="cursor: pointer;"
                @click="viewAttachments(fd)"
              >
                📎 {{ fd.attachments.length }}
              </span>
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
            <h5 class="modal-title">Add Fixed Deposit</h5>
            <button type="button" class="btn-close" @click="showAddModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveFD">
              <div class="mb-3">
                <label class="form-label">FD Name *</label>
                <input type="text" class="form-control" v-model="form.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Bank *</label>
                <input type="text" class="form-control" v-model="form.bank" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Principal Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.principal" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Interest Rate (%) *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.interestRate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Start Date *</label>
                <input type="date" class="form-control" v-model="form.startDate" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Tenure (Months) *</label>
                <input type="number" class="form-control" v-model.number="form.tenureMonths" required @change="calculateMaturityDate" />
              </div>
              <div class="mb-3">
                <label class="form-label">Maturity Date</label>
                <input type="date" class="form-control" v-model="form.maturityDate" readonly />
              </div>
              <div class="mb-3">
                <FileUpload
                  v-model="fdAttachments"
                  label="FD Certificate & Documents"
                  :multiple="true"
                  :max-files="3"
                  :max-size="10485760"
                  accepted-types=".pdf,.jpg,.jpeg,.png"
                />
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useFixedDepositsStore } from '@/stores/fixedDeposits'
import { useSettingsStore } from '@/stores/settings'
import { FileUpload } from '@/components'

const fixedDepositsStore = useFixedDepositsStore()
const settingsStore = useSettingsStore()
const showAddModal = ref(false)
const fdAttachments = ref([])
const showAttachmentsModal = ref(false)
const viewingFD = ref(null)

const form = ref({
  name: '', bank: '', principal: 0, interestRate: 0,
  startDate: new Date().toISOString().split('T')[0],
  tenureMonths: 12, maturityDate: '', compounding: 'monthly',
  attachments: []
})

const fixedDeposits = computed(() => fixedDepositsStore.allFixedDeposits)
const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })

const calculateMaturityDate = () => {
  const start = new Date(form.value.startDate)
  start.setMonth(start.getMonth() + form.value.tenureMonths)
  form.value.maturityDate = start.toISOString().split('T')[0]
}

const withdrawFD = async (id) => {
  if (confirm('Withdraw this FD?')) {
    await fixedDepositsStore.withdrawFixedDeposit(id)
  }
}

const saveFD = async () => {
  const data = {
    ...form.value,
    attachments: fdAttachments.value.map(f => f.fileUrl)
  }

  await fixedDepositsStore.createFixedDeposit(data)
  showAddModal.value = false
  form.value = {
    name: '', bank: '', principal: 0, interestRate: 0,
    startDate: new Date().toISOString().split('T')[0],
    tenureMonths: 12, maturityDate: '', compounding: 'monthly',
    attachments: []
  }
  fdAttachments.value = []
}

const viewAttachments = (fd) => {
  viewingFD.value = fd
  showAttachmentsModal.value = true
}

const isImageUrl = (url) => {
  return /\.(jpg|jpeg|png|gif|webp)$/i.test(url)
}

const getFileNameFromUrl = (url) => {
  return url.split('/').pop()
}

const openAttachment = (url) => {
  window.open(url, '_blank')
}

onMounted(() => {
  fixedDepositsStore.fetchFixedDeposits()
  calculateMaturityDate()
})
</script>
