<template>
  <div class="bills-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Bills & Reminders</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Bill</button>
    </div>

    <!-- Upcoming Bills Alert -->
    <div v-if="upcomingBills.length > 0" class="alert alert-info mb-4">
      <h6 class="alert-heading">📅 Upcoming Bills (Next 30 Days)</h6>
      <ul class="mb-0">
        <li v-for="bill in upcomingBills.slice(0, 5)" :key="bill.id">
          {{ bill.name }} - {{ formatCurrency(bill.amount) }} (Due in {{ bill.daysUntilDue }} days)
        </li>
      </ul>
    </div>

    <!-- Overdue Bills Alert -->
    <div v-if="overdueBills.length > 0" class="alert alert-danger mb-4">
      <h6 class="alert-heading">⚠️ Overdue Bills</h6>
      <ul class="mb-0">
        <li v-for="bill in overdueBills" :key="bill.id">
          {{ bill.name }} - {{ formatCurrency(bill.amount) }} ({{ bill.daysOverdue }} days overdue)
        </li>
      </ul>
    </div>

    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon purple">📄</div>
          <div class="stat-value">{{ activeBills.length }}</div>
          <div class="stat-label">Active Bills</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon red">💰</div>
          <div class="stat-value">{{ formatCurrency(billsStore.totalMonthlyBills) }}</div>
          <div class="stat-label">Monthly Bills</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon orange">🔔</div>
          <div class="stat-value">{{ upcomingBills.length }}</div>
          <div class="stat-label">Upcoming</div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h5 class="mb-0">All Bills</h5>
      </div>
      <div class="card-body p-0">
        <div v-if="activeBills.length === 0" class="p-4 text-center text-muted">No bills yet</div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Bill Name</th>
                <th>Category</th>
                <th>Amount</th>
                <th>Frequency</th>
                <th>Next Due</th>
                <th class="text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="bill in activeBills" :key="bill.id">
                <td class="fw-semibold">{{ bill.name }}</td>
                <td><span class="badge bg-secondary">{{ bill.category }}</span></td>
                <td class="fw-bold">{{ formatCurrency(bill.amount) }}</td>
                <td>{{ bill.frequency }}</td>
                <td>{{ formatDate(billsStore.getNextDueDate(bill)) }}</td>
                <td class="text-center">
                  <button class="btn btn-sm btn-outline-success me-1" @click="markPaid(bill.id)">
                    Mark Paid
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="deleteBill(bill.id)">
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Add Modal -->
    <div class="modal fade" :class="{ 'show d-block': showAddModal }" style="background-color: rgba(0,0,0,0.5);" v-if="showAddModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add Bill</h5>
            <button type="button" class="btn-close" @click="showAddModal = false"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveBill">
              <div class="mb-3">
                <label class="form-label">Bill Name *</label>
                <input type="text" class="form-control" v-model="form.name" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Category *</label>
                <select class="form-select" v-model="form.category" required>
                  <option value="housing">Housing</option>
                  <option value="utilities">Utilities</option>
                  <option value="subscriptions">Subscriptions</option>
                  <option value="insurance">Insurance</option>
                  <option value="healthcare">Healthcare</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input type="number" step="0.01" class="form-control" v-model.number="form.amount" required />
              </div>
              <div class="mb-3">
                <label class="form-label">Frequency *</label>
                <select class="form-select" v-model="form.frequency" required>
                  <option value="monthly">Monthly</option>
                  <option value="weekly">Weekly</option>
                  <option value="quarterly">Quarterly</option>
                  <option value="annually">Annually</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Start Date *</label>
                <input type="date" class="form-control" v-model="form.startDate" required />
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
import { useBillsStore } from '@/stores/bills'
import { useSettingsStore } from '@/stores/settings'

const billsStore = useBillsStore()
const settingsStore = useSettingsStore()
const showAddModal = ref(false)
const form = ref({ name: '', category: 'utilities', amount: 0, frequency: 'monthly', startDate: new Date().toISOString().split('T')[0] })

const activeBills = computed(() => billsStore.activeBills)
const upcomingBills = computed(() => billsStore.upcomingBills)
const overdueBills = computed(() => billsStore.overdueBills)

const formatCurrency = (amount) => settingsStore.formatCurrency(amount)
const formatDate = (date) => date ? new Date(date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) : '-'

const markPaid = async (billId) => {
  await billsStore.markAsPaid(billId)
}

const deleteBill = async (id) => {
  if (confirm('Delete this bill?')) {
    await billsStore.deleteBill(id)
  }
}

const saveBill = async () => {
  await billsStore.createBill(form.value)
  showAddModal.value = false
  form.value = { name: '', category: 'utilities', amount: 0, frequency: 'monthly', startDate: new Date().toISOString().split('T')[0] }
}

onMounted(() => billsStore.fetchBills())
</script>
