<template>
  <div class="goods-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-blue">My Goods</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Good</button>
    </div>

    <!-- Summary Stats -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">📦</div>
          <div class="stat-value">{{ goodsStore.activeGoods.length }}</div>
          <div class="stat-label">Active Goods</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(goodsStore.totalValue) }}</div>
          <div class="stat-label">Total Value</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon orange">🔧</div>
          <div class="stat-value">{{ formatCurrency(goodsStore.totalServiceCost) }}</div>
          <div class="stat-label">Total Service Cost</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">📜</div>
          <div class="stat-value">{{ goodsStore.goodsUnderWarranty.length }}</div>
          <div class="stat-label">Under Warranty</div>
        </div>
      </div>
    </div>

    <!-- Warranty Alerts -->
    <div v-if="goodsStore.goodsWarrantyExpiringSoon.length > 0" class="alert alert-warning mb-4">
      <h6 class="alert-heading">⚠️ Warranty Expiring Soon</h6>
      <ul class="mb-0">
        <li v-for="good in goodsStore.goodsWarrantyExpiringSoon" :key="good.id">
          {{ good.name }}: {{ good.warrantyDaysRemaining }} days remaining
        </li>
      </ul>
    </div>

    <!-- Filter Tabs -->
    <ul class="nav nav-tabs mb-3">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'all' }" @click="filter = 'all'" href="javascript:void(0)">
          All ({{ goodsStore.allGoods.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'active' }" @click="filter = 'active'" href="javascript:void(0)">
          Active ({{ goodsStore.activeGoods.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'warranty' }" @click="filter = 'warranty'" href="javascript:void(0)">
          Under Warranty ({{ goodsStore.goodsUnderWarranty.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'archived' }" @click="filter = 'archived'" href="javascript:void(0)">
          Archived ({{ goodsStore.archivedGoods.length }})
        </a>
      </li>
    </ul>

    <!-- Goods List -->
    <div v-if="loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status"></div>
    </div>

    <div v-else class="row g-3">
      <div v-for="good in filteredGoods" :key="good.id" class="col-12 col-md-6 col-lg-4">
        <div class="card h-100">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-2">
              <h5 class="card-title mb-0">{{ good.name }}</h5>
              <span class="badge" :class="getStatusClass(good.status)">{{ formatStatus(good.status) }}</span>
            </div>

            <p class="text-muted small mb-2" v-if="good.category">{{ good.category }}</p>
            <p class="text-muted mb-2" v-if="good.description">{{ good.description }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Purchase Price</span>
                <strong>{{ formatCurrency(good.purchasePrice) }}</strong>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Purchase Date</span>
                <span>{{ formatDate(good.purchaseDate) }}</span>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Days Owned</span>
                <span>{{ good.daysOwned || 0 }} days</span>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Price/Day</span>
                <span>{{ formatCurrency(good.pricePerDay || 0) }}</span>
              </div>
              <div v-if="good.totalServiceCost > 0" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Service Cost</span>
                <span class="text-warning">{{ formatCurrency(good.totalServiceCost) }}</span>
              </div>
            </div>

            <!-- Warranty Info -->
            <div v-if="good.warrantyStatus !== 'no_warranty'" class="mb-3 p-2 rounded" :class="getWarrantyBgClass(good.warrantyStatus)">
              <div class="d-flex justify-content-between align-items-center mb-1">
                <span class="small">{{ getWarrantyStatusText(good.warrantyStatus) }}</span>
                <span class="badge" :class="getWarrantyBadgeClass(good.warrantyStatus)">
                  {{ good.warrantyStatus === 'active' ? good.warrantyDaysRemaining + ' days' : 'Expired' }}
                </span>
              </div>
              <div v-if="good.warrantyStatus === 'active'" class="progress" style="height: 4px;">
                <div
                  class="progress-bar"
                  :class="good.warrantyDaysRemaining <= 30 ? 'bg-warning' : 'bg-success'"
                  :style="{ width: getWarrantyProgress(good) + '%' }"
                ></div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="d-flex flex-wrap gap-1 mb-2">
              <button class="btn btn-sm btn-outline-primary flex-grow-1" @click="viewDetails(good)">
                View Details
              </button>
              <button class="btn btn-sm btn-outline-success" @click="openServiceModal(good)">
                Service
              </button>
              <button class="btn btn-sm btn-outline-info" @click="openAttachmentsModal(good)">
                Files
              </button>
            </div>
            <div class="d-flex gap-1">
              <button class="btn btn-sm btn-outline-secondary flex-grow-1" @click="editGood(good)">
                Edit
              </button>
              <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(good.id)">
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="filteredGoods.length === 0" class="col-12">
        <div class="text-center py-5 text-muted">
          <p class="h5">No goods found</p>
          <p>Add your first good to start tracking!</p>
        </div>
      </div>
    </div>

    <!-- Add/Edit Good Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Good' : 'Add Good' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveGood">
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Name *</label>
                  <input v-model="form.name" type="text" class="form-control" required />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Category</label>
                  <input v-model="form.category" type="text" class="form-control"
                         list="categories" placeholder="e.g., Electronics, Furniture" />
                  <datalist id="categories">
                    <option v-for="cat in goodsStore.categories" :key="cat" :value="cat"></option>
                  </datalist>
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Brand</label>
                  <input v-model="form.brand" type="text" class="form-control" />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Model</label>
                  <input v-model="form.model" type="text" class="form-control" />
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Serial Number</label>
                <input v-model="form.serialNumber" type="text" class="form-control" />
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Purchase Date *</label>
                  <input v-model="form.purchaseDate" type="date" class="form-control" required />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Purchase Price *</label>
                  <input v-model.number="form.purchasePrice" type="number" step="0.01" class="form-control" required />
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Purchase Location</label>
                <input v-model="form.purchaseLocation" type="text" class="form-control" placeholder="Store name or website" />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea v-model="form.description" class="form-control" rows="2"></textarea>
              </div>

              <h6 class="mt-4 mb-3">Warranty Information</h6>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Warranty Start Date</label>
                  <input v-model="form.warrantyStartDate" type="date" class="form-control" />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Warranty End Date</label>
                  <input v-model="form.warrantyEndDate" type="date" class="form-control" />
                </div>
              </div>

              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Warranty Provider</label>
                  <input v-model="form.warrantyProvider" type="text" class="form-control" placeholder="Manufacturer, retailer, etc." />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Warranty Type</label>
                  <select v-model="form.warrantyType" class="form-select">
                    <option value="">Select type</option>
                    <option value="manufacturer">Manufacturer</option>
                    <option value="extended">Extended</option>
                    <option value="lifetime">Lifetime</option>
                  </select>
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Status</label>
                <select v-model="form.status" class="form-select">
                  <option value="active">Active</option>
                  <option value="archived">Archived</option>
                  <option value="sold">Sold</option>
                  <option value="disposed">Disposed</option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Notes</label>
                <textarea v-model="form.notes" class="form-control" rows="2"></textarea>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
                <button type="submit" class="btn btn-primary" :disabled="loading">
                  {{ loading ? 'Saving...' : 'Save' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Service Records Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showServiceModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showServiceModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Service Records - {{ selectedGood?.name }}</h5>
            <button type="button" class="btn-close" @click="closeServiceModal"></button>
          </div>
          <div class="modal-body">
            <!-- Add Service Form -->
            <form @submit.prevent="addService" class="mb-4 p-3 bg-light rounded">
              <h6 class="mb-3">Add Service Record</h6>
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Service Date *</label>
                  <input v-model="serviceForm.serviceDate" type="date" class="form-control" required />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Service Type *</label>
                  <select v-model="serviceForm.serviceType" class="form-select" required>
                    <option value="">Select type</option>
                    <option value="repair">Repair</option>
                    <option value="maintenance">Maintenance</option>
                    <option value="inspection">Inspection</option>
                    <option value="replacement">Replacement</option>
                  </select>
                </div>
              </div>
              <div class="row">
                <div class="col-md-6 mb-3">
                  <label class="form-label">Service Provider</label>
                  <input v-model="serviceForm.serviceProvider" type="text" class="form-control" />
                </div>
                <div class="col-md-6 mb-3">
                  <label class="form-label">Cost *</label>
                  <input v-model.number="serviceForm.cost" type="number" step="0.01" class="form-control" required />
                </div>
              </div>
              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea v-model="serviceForm.description" class="form-control" rows="2"></textarea>
              </div>
              <div class="form-check mb-3">
                <input v-model="serviceForm.warrantyCovered" type="checkbox" class="form-check-input" id="warrantyCovered" />
                <label class="form-check-label" for="warrantyCovered">Covered by warranty</label>
              </div>
              <button type="submit" class="btn btn-primary btn-sm" :disabled="loading">
                {{ loading ? 'Adding...' : 'Add Service Record' }}
              </button>
            </form>

            <!-- Service Records List -->
            <h6 class="mb-3">Service History</h6>
            <div v-if="currentServiceRecords.length === 0" class="text-muted text-center py-3">
              No service records yet
            </div>
            <div v-else>
              <div v-for="service in currentServiceRecords" :key="service.id" class="card mb-2">
                <div class="card-body p-3">
                  <div class="d-flex justify-content-between align-items-start">
                    <div>
                      <div class="d-flex align-items-center gap-2 mb-1">
                        <strong>{{ service.serviceType }}</strong>
                        <span v-if="service.warrantyCovered" class="badge bg-success">Warranty</span>
                      </div>
                      <div class="text-muted small">{{ formatDate(service.serviceDate) }}</div>
                      <div v-if="service.serviceProvider" class="small">{{ service.serviceProvider }}</div>
                      <div v-if="service.description" class="small mt-1">{{ service.description }}</div>
                    </div>
                    <div class="text-end">
                      <div class="fw-bold">{{ formatCurrency(service.cost) }}</div>
                      <button class="btn btn-sm btn-outline-danger mt-1" @click="deleteService(service.id)">
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

    <!-- Attachments Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAttachmentsModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAttachmentsModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Attachments - {{ selectedGood?.name }}</h5>
            <button type="button" class="btn-close" @click="closeAttachmentsModal"></button>
          </div>
          <div class="modal-body">
            <p class="text-muted small mb-3">Upload receipts, warranty documents, photos, and other files related to this good.</p>

            <!-- File Upload Info -->
            <div class="alert alert-info small">
              To upload files: Use the file upload endpoint at <code>/api/v1/uploads</code> first, then link them here.
            </div>

            <!-- Attachments List -->
            <h6 class="mb-3">Current Attachments</h6>
            <div v-if="currentAttachments.length === 0" class="text-muted text-center py-3">
              No attachments yet
            </div>
            <div v-else class="row g-2">
              <div v-for="attachment in currentAttachments" :key="attachment.id" class="col-12 col-md-6">
                <div class="card">
                  <div class="card-body p-2">
                    <div class="d-flex justify-content-between align-items-center">
                      <div class="flex-grow-1">
                        <div class="small fw-bold text-truncate">{{ attachment.originalName }}</div>
                        <div class="text-muted small">{{ attachment.attachmentType || 'file' }}</div>
                      </div>
                      <button class="btn btn-sm btn-outline-danger" @click="deleteAttachment(attachment.id)">
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

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useGoodsStore } from '@/stores/goods'
import { useSettingsStore } from '@/stores/settings'

const goodsStore = useGoodsStore()
const settingsStore = useSettingsStore()

const loading = ref(false)
const filter = ref('all')
const showAddModal = ref(false)
const showEditModal = ref(false)
const showServiceModal = ref(false)
const showAttachmentsModal = ref(false)
const selectedGood = ref(null)
const editingGoodId = ref(null)

const form = ref({
  name: '',
  description: '',
  category: '',
  brand: '',
  model: '',
  serialNumber: '',
  purchaseDate: '',
  purchasePrice: 0,
  purchaseLocation: '',
  warrantyStartDate: null,
  warrantyEndDate: null,
  warrantyProvider: '',
  warrantyType: '',
  status: 'active',
  notes: ''
})

const serviceForm = ref({
  serviceDate: '',
  serviceType: '',
  serviceProvider: '',
  cost: 0,
  description: '',
  warrantyCovered: false
})

const filteredGoods = computed(() => {
  if (filter.value === 'all') return goodsStore.allGoods
  if (filter.value === 'active') return goodsStore.activeGoods
  if (filter.value === 'archived') return goodsStore.archivedGoods
  if (filter.value === 'warranty') return goodsStore.goodsUnderWarranty
  return goodsStore.allGoods
})

const currentServiceRecords = computed(() => {
  if (!selectedGood.value) return []
  return goodsStore.getServiceRecordsForGood(selectedGood.value.id)
})

const currentAttachments = computed(() => {
  if (!selectedGood.value) return []
  return goodsStore.getAttachmentsForGood(selectedGood.value.id)
})

onMounted(async () => {
  loading.value = true
  try {
    await goodsStore.fetchGoods()
    await goodsStore.fetchStats()
  } catch (error) {
    console.error('Error loading goods:', error)
  } finally {
    loading.value = false
  }
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString()
}

const formatStatus = (status) => {
  return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const getStatusClass = (status) => {
  const classes = {
    active: 'bg-success',
    archived: 'bg-secondary',
    sold: 'bg-info',
    disposed: 'bg-dark'
  }
  return classes[status] || 'bg-secondary'
}

const getWarrantyStatusText = (status) => {
  if (status === 'active') return '✓ Warranty Active'
  if (status === 'expired') return '✗ Warranty Expired'
  return 'No Warranty'
}

const getWarrantyBadgeClass = (status) => {
  return status === 'active' ? 'bg-success' : 'bg-danger'
}

const getWarrantyBgClass = (status) => {
  return status === 'active' ? 'bg-light-success' : 'bg-light-danger'
}

const getWarrantyProgress = (good) => {
  if (!good.warrantyDaysTotal || good.warrantyDaysTotal === 0) return 0
  const passed = good.warrantyDaysPassed || 0
  return Math.min(100, (passed / good.warrantyDaysTotal) * 100)
}

const editGood = (good) => {
  editingGoodId.value = good.id
  form.value = {
    name: good.name,
    description: good.description || '',
    category: good.category || '',
    brand: good.brand || '',
    model: good.model || '',
    serialNumber: good.serialNumber || '',
    purchaseDate: good.purchaseDate?.split('T')[0] || '',
    purchasePrice: good.purchasePrice,
    purchaseLocation: good.purchaseLocation || '',
    warrantyStartDate: good.warrantyStartDate?.split('T')[0] || null,
    warrantyEndDate: good.warrantyEndDate?.split('T')[0] || null,
    warrantyProvider: good.warrantyProvider || '',
    warrantyType: good.warrantyType || '',
    status: good.status,
    notes: good.notes || ''
  }
  showEditModal.value = true
}

const saveGood = async () => {
  loading.value = true
  try {
    if (showEditModal.value) {
      await goodsStore.updateGood(editingGoodId.value, form.value)
    } else {
      await goodsStore.createGood(form.value)
    }
    closeModal()
  } catch (error) {
    console.error('Error saving good:', error)
    alert('Error saving good. Please try again.')
  } finally {
    loading.value = false
  }
}

const confirmDelete = async (id) => {
  if (confirm('Are you sure you want to delete this good? This will also delete all service records and attachments.')) {
    loading.value = true
    try {
      await goodsStore.deleteGood(id)
    } catch (error) {
      console.error('Error deleting good:', error)
      alert('Error deleting good. Please try again.')
    } finally {
      loading.value = false
    }
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingGoodId.value = null
  form.value = {
    name: '',
    description: '',
    category: '',
    brand: '',
    model: '',
    serialNumber: '',
    purchaseDate: '',
    purchasePrice: 0,
    purchaseLocation: '',
    warrantyStartDate: null,
    warrantyEndDate: null,
    warrantyProvider: '',
    warrantyType: '',
    status: 'active',
    notes: ''
  }
}

const viewDetails = async (good) => {
  selectedGood.value = good
  await goodsStore.fetchServiceRecords(good.id)
  await goodsStore.fetchAttachments(good.id)
  showServiceModal.value = true
}

const openServiceModal = async (good) => {
  selectedGood.value = good
  await goodsStore.fetchServiceRecords(good.id)
  showServiceModal.value = true
}

const closeServiceModal = () => {
  showServiceModal.value = false
  selectedGood.value = null
  serviceForm.value = {
    serviceDate: '',
    serviceType: '',
    serviceProvider: '',
    cost: 0,
    description: '',
    warrantyCovered: false
  }
}

const addService = async () => {
  if (!selectedGood.value) return

  loading.value = true
  try {
    await goodsStore.createServiceRecord(selectedGood.value.id, serviceForm.value)
    serviceForm.value = {
      serviceDate: '',
      serviceType: '',
      serviceProvider: '',
      cost: 0,
      description: '',
      warrantyCovered: false
    }
  } catch (error) {
    console.error('Error adding service record:', error)
    alert('Error adding service record. Please try again.')
  } finally {
    loading.value = false
  }
}

const deleteService = async (serviceId) => {
  if (!selectedGood.value) return
  if (!confirm('Are you sure you want to delete this service record?')) return

  loading.value = true
  try {
    await goodsStore.deleteServiceRecord(selectedGood.value.id, serviceId)
  } catch (error) {
    console.error('Error deleting service record:', error)
    alert('Error deleting service record. Please try again.')
  } finally {
    loading.value = false
  }
}

const openAttachmentsModal = async (good) => {
  selectedGood.value = good
  await goodsStore.fetchAttachments(good.id)
  showAttachmentsModal.value = true
}

const closeAttachmentsModal = () => {
  showAttachmentsModal.value = false
  selectedGood.value = null
}

const deleteAttachment = async (attachmentId) => {
  if (!selectedGood.value) return
  if (!confirm('Are you sure you want to delete this attachment?')) return

  loading.value = true
  try {
    await goodsStore.deleteAttachment(selectedGood.value.id, attachmentId)
  } catch (error) {
    console.error('Error deleting attachment:', error)
    alert('Error deleting attachment. Please try again.')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.bg-light-success {
  background-color: #d4edda;
}

.bg-light-danger {
  background-color: #f8d7da;
}
</style>
