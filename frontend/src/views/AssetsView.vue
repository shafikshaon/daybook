<template>
  <div class="assets-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-blue">My Assets</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Asset</button>
    </div>

    <!-- Summary Stats -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon blue">📦</div>
          <div class="stat-value">{{ assetsStore.activeAssets.length }}</div>
          <div class="stat-label">Active Goods</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(assetsStore.totalValue) }}</div>
          <div class="stat-label">Total Value</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon orange">🔧</div>
          <div class="stat-value">{{ formatCurrency(assetsStore.totalServiceCost) }}</div>
          <div class="stat-label">Total Service Cost</div>
        </div>
      </div>
      <div class="col-12 col-md-3">
        <div class="stat-card">
          <div class="stat-icon green">📜</div>
          <div class="stat-value">{{ assetsStore.assetsUnderWarranty.length }}</div>
          <div class="stat-label">Under Warranty</div>
        </div>
      </div>
    </div>

    <!-- Warranty Alerts -->
    <div v-if="assetsStore.assetsWarrantyExpiringSoon.length > 0" class="alert alert-warning mb-4">
      <h6 class="alert-heading">⚠️ Warranty Expiring Soon</h6>
      <ul class="mb-0">
        <li v-for="asset in assetsStore.assetsWarrantyExpiringSoon" :key="asset.id">
          {{ asset.name }}: {{ asset.warrantyDaysRemaining }} days remaining
        </li>
      </ul>
    </div>

    <!-- Filter Tabs -->
    <ul class="nav nav-tabs mb-3">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'all' }" @click="filter = 'all'" href="javascript:void(0)">
          All ({{ assetsStore.allAssets.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'active' }" @click="filter = 'active'" href="javascript:void(0)">
          Active ({{ assetsStore.activeAssets.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'warranty' }" @click="filter = 'warranty'" href="javascript:void(0)">
          Under Warranty ({{ assetsStore.assetsUnderWarranty.length }})
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'archived' }" @click="filter = 'archived'" href="javascript:void(0)">
          Archived ({{ assetsStore.archivedAssets.length }})
        </a>
      </li>
    </ul>

    <!-- Goods List -->
    <div v-if="loading" class="text-center my-5">
      <div class="spinner-border text-primary" role="status"></div>
    </div>

    <div v-else class="row g-3">
      <div v-for="asset in filteredAssets" :key="asset.id" class="col-12 col-md-6 col-lg-4">
        <div class="card h-100">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-2">
              <h5 class="card-title mb-0">{{ asset.name }}</h5>
              <span class="badge" :class="getStatusClass(asset.status)">{{ formatStatus(asset.status) }}</span>
            </div>

            <p class="text-muted small mb-2" v-if="asset.category">{{ asset.category }}</p>
            <p class="text-muted mb-2" v-if="asset.description">{{ asset.description }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Purchase Price</span>
                <strong>{{ formatCurrency(asset.purchasePrice) }}</strong>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Purchase Date</span>
                <span>{{ formatDate(asset.purchaseDate) }}</span>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Days Owned</span>
                <span>{{ asset.daysOwned || 0 }} days</span>
              </div>
              <div class="d-flex justify-content-between mb-1">
                <span class="text-muted">Price/Day</span>
                <span>{{ formatCurrency(asset.pricePerDay || 0) }}</span>
              </div>
              <div v-if="asset.totalServiceCost > 0" class="d-flex justify-content-between mb-1">
                <span class="text-muted">Service Cost</span>
                <span class="text-warning">{{ formatCurrency(asset.totalServiceCost) }}</span>
              </div>
            </div>

            <!-- Warranty Info -->
            <div v-if="asset.warrantyStatus !== 'no_warranty'" class="mb-3 p-2 rounded" :class="getWarrantyBgClass(asset.warrantyStatus)">
              <div class="d-flex justify-content-between align-items-center mb-1">
                <span class="small">{{ getWarrantyStatusText(asset.warrantyStatus) }}</span>
                <span class="badge" :class="getWarrantyBadgeClass(asset.warrantyStatus)">
                  {{ asset.warrantyStatus === 'active' ? asset.warrantyDaysRemaining + ' days' : 'Expired' }}
                </span>
              </div>
              <div v-if="asset.warrantyStatus === 'active'" class="progress" style="height: 4px;">
                <div
                  class="progress-bar"
                  :class="asset.warrantyDaysRemaining <= 30 ? 'bg-warning' : 'bg-primary'"
                  :style="{ width: getWarrantyProgress(asset) + '%' }"
                ></div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="d-flex flex-wrap gap-1 mb-2">
              <button class="btn btn-sm btn-outline-primary flex-grow-1" @click="viewDetails(asset)">
                View Details
              </button>
              <button class="btn btn-sm btn-outline-primary" @click="openServiceModal(asset)">
                Service
              </button>
              <button class="btn btn-sm btn-outline-info" @click="openAttachmentsModal(asset)">
                Files
              </button>
            </div>
            <div class="d-flex gap-1">
              <button class="btn btn-sm btn-outline-secondary flex-grow-1" @click="editGood(asset)">
                Edit
              </button>
              <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(asset.id)">
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="filteredAssets.length === 0" class="col-12">
        <div class="text-center py-5 text-muted">
          <p class="h5">No assets found</p>
          <p>Add your first asset to start tracking!</p>
        </div>
      </div>
    </div>

    <!-- Add/Edit Asset Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Asset' : 'Add Asset' }}</h5>
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
                    <option v-for="cat in assetsStore.categories" :key="cat" :value="cat"></option>
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

              <!-- Attachments Section for Edit Mode -->
              <div v-if="showEditModal && editingAssetId" class="mb-4">
                <h6 class="mb-3">Attachments</h6>

                <!-- Upload New Files -->
                <div class="mb-3">
                  <label class="form-label">Upload Files</label>
                  <div class="row mb-2">
                    <div class="col-md-6">
                      <select v-model="attachmentType" class="form-select form-select-sm">
                        <option value="photo">Photo</option>
                        <option value="receipt">Receipt</option>
                        <option value="warranty_document">Warranty Document</option>
                        <option value="manual">Manual</option>
                        <option value="other">Other</option>
                      </select>
                    </div>
                    <div class="col-md-6">
                      <input v-model="attachmentDescription" class="form-control form-control-sm" placeholder="Description (optional)" />
                    </div>
                  </div>
                  <FileUpload
                    label=""
                    :multiple="true"
                    :maxFiles="10"
                    :autoUpload="true"
                    @upload-success="handleEditFormFilesUploaded"
                    ref="editFileUploadRef"
                  />
                </div>

                <!-- Current Attachments -->
                <div v-if="editFormAttachments.length > 0">
                  <label class="form-label">Current Files</label>
                  <div class="list-group">
                    <div v-for="attachment in editFormAttachments" :key="attachment.id" class="list-group-item">
                      <div class="d-flex justify-content-between align-items-center">
                        <div class="flex-grow-1">
                          <div class="fw-bold small">{{ attachment.originalName }}</div>
                          <span class="badge bg-secondary me-1">{{ formatAttachmentType(attachment.attachmentType) }}</span>
                          <span v-if="attachment.description" class="text-muted small">{{ attachment.description }}</span>
                        </div>
                        <div class="d-flex gap-1">
                          <button type="button" class="btn btn-sm btn-outline-primary" @click="viewAttachment(attachment)">👁️</button>
                          <a :href="getAttachmentUrl(attachment)" :download="attachment.originalName" class="btn btn-sm btn-outline-primary">⬇️</a>
                          <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteEditFormAttachment(attachment.id)">🗑️</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
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

    <!-- Asset Details Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showDetailsModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showDetailsModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Asset Details - {{ selectedAsset?.name }}</h5>
            <button type="button" class="btn-close" @click="closeDetailsModal"></button>
          </div>
          <div class="modal-body">
            <div v-if="selectedAsset">
              <!-- Basic Information -->
              <div class="mb-4">
                <h6 class="text-primary mb-3">Basic Information</h6>
                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="text-muted small">Name</label>
                    <div class="fw-bold">{{ selectedAsset.name }}</div>
                  </div>
                  <div class="col-md-6">
                    <label class="text-muted small">Status</label>
                    <div><span class="badge" :class="getStatusClass(selectedAsset.status)">{{ formatStatus(selectedAsset.status) }}</span></div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.category">
                    <label class="text-muted small">Category</label>
                    <div>{{ selectedAsset.category }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.brand">
                    <label class="text-muted small">Brand</label>
                    <div>{{ selectedAsset.brand }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.model">
                    <label class="text-muted small">Model</label>
                    <div>{{ selectedAsset.model }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.serialNumber">
                    <label class="text-muted small">Serial Number</label>
                    <div>{{ selectedAsset.serialNumber }}</div>
                  </div>
                  <div class="col-12" v-if="selectedAsset.description">
                    <label class="text-muted small">Description</label>
                    <div>{{ selectedAsset.description }}</div>
                  </div>
                </div>
              </div>

              <!-- Purchase Information -->
              <div class="mb-4">
                <h6 class="text-primary mb-3">Purchase Information</h6>
                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="text-muted small">Purchase Date</label>
                    <div>{{ formatDate(selectedAsset.purchaseDate) }}</div>
                  </div>
                  <div class="col-md-6">
                    <label class="text-muted small">Purchase Price</label>
                    <div class="fw-bold">{{ formatCurrency(selectedAsset.purchasePrice) }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.purchaseLocation">
                    <label class="text-muted small">Purchase Location</label>
                    <div>{{ selectedAsset.purchaseLocation }}</div>
                  </div>
                  <div class="col-md-6">
                    <label class="text-muted small">Days Owned</label>
                    <div>{{ selectedAsset.daysOwned || 0 }} days</div>
                  </div>
                  <div class="col-md-6">
                    <label class="text-muted small">Price per Day</label>
                    <div>{{ formatCurrency(selectedAsset.pricePerDay || 0) }}</div>
                  </div>
                </div>
              </div>

              <!-- Warranty Information -->
              <div class="mb-4" v-if="selectedAsset.warrantyStatus !== 'no_warranty'">
                <h6 class="text-primary mb-3">Warranty Information</h6>
                <div class="warranty-info-card p-3 rounded mb-3" :class="getWarrantyBgClass(selectedAsset.warrantyStatus)">
                  <div class="d-flex justify-content-between align-items-center mb-2">
                    <span class="fw-bold">{{ getWarrantyStatusText(selectedAsset.warrantyStatus) }}</span>
                    <span class="badge" :class="getWarrantyBadgeClass(selectedAsset.warrantyStatus)">
                      {{ selectedAsset.warrantyStatus === 'active' ? selectedAsset.warrantyDaysRemaining + ' days remaining' : 'Expired' }}
                    </span>
                  </div>
                  <div v-if="selectedAsset.warrantyStatus === 'active'" class="progress" style="height: 6px;">
                    <div
                      class="progress-bar"
                      :class="selectedAsset.warrantyDaysRemaining <= 30 ? 'bg-warning' : 'bg-primary'"
                      :style="{ width: getWarrantyProgress(selectedAsset) + '%' }"
                    ></div>
                  </div>
                </div>
                <div class="row g-3">
                  <div class="col-md-6" v-if="selectedAsset.warrantyStartDate">
                    <label class="text-muted small">Warranty Start Date</label>
                    <div>{{ formatDate(selectedAsset.warrantyStartDate) }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.warrantyEndDate">
                    <label class="text-muted small">Warranty End Date</label>
                    <div>{{ formatDate(selectedAsset.warrantyEndDate) }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.warrantyProvider">
                    <label class="text-muted small">Warranty Provider</label>
                    <div>{{ selectedAsset.warrantyProvider }}</div>
                  </div>
                  <div class="col-md-6" v-if="selectedAsset.warrantyType">
                    <label class="text-muted small">Warranty Type</label>
                    <div>{{ formatStatus(selectedAsset.warrantyType) }}</div>
                  </div>
                </div>
              </div>

              <!-- Service History Summary -->
              <div class="mb-4" v-if="currentServiceRecords.length > 0">
                <h6 class="text-primary mb-3">Service History</h6>
                <div class="alert alert-info">
                  <div class="d-flex justify-content-between align-items-center">
                    <span>Total Service Records: <strong>{{ currentServiceRecords.length }}</strong></span>
                    <span>Total Service Cost: <strong>{{ formatCurrency(selectedAsset.totalServiceCost || 0) }}</strong></span>
                  </div>
                </div>
                <div class="list-group">
                  <div v-for="service in currentServiceRecords.slice(0, 3)" :key="service.id" class="list-group-item">
                    <div class="d-flex justify-content-between align-items-start">
                      <div>
                        <div class="d-flex align-items-center gap-2">
                          <strong>{{ service.serviceType }}</strong>
                          <span v-if="service.warrantyCovered" class="badge bg-primary">Warranty</span>
                        </div>
                        <div class="text-muted small">{{ formatDate(service.serviceDate) }}</div>
                      </div>
                      <div class="fw-bold">{{ formatCurrency(service.cost) }}</div>
                    </div>
                  </div>
                </div>
                <button class="btn btn-sm btn-outline-primary mt-2 w-100" @click="openServiceModal(selectedAsset); closeDetailsModal()">
                  View All Service Records
                </button>
              </div>

              <!-- Attachments Summary -->
              <div class="mb-4" v-if="currentAttachments.length > 0">
                <h6 class="text-primary mb-3">Attachments</h6>
                <div class="alert alert-info">
                  Total Attachments: <strong>{{ currentAttachments.length }}</strong>
                </div>
                <div class="list-group">
                  <div v-for="attachment in currentAttachments.slice(0, 3)" :key="attachment.id" class="list-group-item">
                    <div class="d-flex justify-content-between align-items-center">
                      <div>
                        <div class="fw-bold small">{{ attachment.originalName }}</div>
                        <span class="badge bg-secondary">{{ formatAttachmentType(attachment.attachmentType) }}</span>
                      </div>
                      <button class="btn btn-sm btn-outline-primary" @click="viewAttachment(attachment)">View</button>
                    </div>
                  </div>
                </div>
                <button class="btn btn-sm btn-outline-primary mt-2 w-100" @click="openAttachmentsModal(selectedAsset); closeDetailsModal()">
                  View All Attachments
                </button>
              </div>

              <!-- Notes -->
              <div v-if="selectedAsset.notes">
                <h6 class="text-primary mb-3">Notes</h6>
                <div class="alert alert-secondary">
                  {{ selectedAsset.notes }}
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="d-flex gap-2 mt-4">
                <button class="btn btn-outline-secondary flex-grow-1" @click="editGood(selectedAsset); closeDetailsModal()">
                  Edit Asset
                </button>
                <button class="btn btn-outline-primary" @click="openServiceModal(selectedAsset); closeDetailsModal()">
                  Add Service
                </button>
                <button class="btn btn-outline-info" @click="openAttachmentsModal(selectedAsset); closeDetailsModal()">
                  Add Files
                </button>
              </div>
            </div>
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
            <h5 class="modal-title">Service Records - {{ selectedAsset?.name }}</h5>
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
                        <span v-if="service.warrantyCovered" class="badge bg-primary">Warranty</span>
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
            <h5 class="modal-title">Attachments - {{ selectedAsset?.name }}</h5>
            <button type="button" class="btn-close" @click="closeAttachmentsModal"></button>
          </div>
          <div class="modal-body">
            <p class="text-muted small mb-3">Upload receipts, warranty documents, photos, and other files related to this asset.</p>

            <!-- File Upload Component -->
            <div class="mb-4">
              <h6 class="mb-3">Upload New Files</h6>

              <!-- Attachment Type Selection (before upload) -->
              <div class="row mb-3">
                <div class="col-md-6">
                  <label class="form-label">Attachment Type</label>
                  <select v-model="attachmentType" class="form-select">
                    <option value="photo">Photo</option>
                    <option value="receipt">Receipt</option>
                    <option value="warranty_document">Warranty Document</option>
                    <option value="manual">Manual</option>
                    <option value="other">Other</option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label">Description (Optional)</label>
                  <input v-model="attachmentDescription" class="form-control" placeholder="Add a description..." />
                </div>
              </div>

              <FileUpload
                label=""
                :multiple="true"
                :maxFiles="10"
                :autoUpload="true"
                @upload-success="handleFilesUploaded"
                ref="fileUploadRef"
              />
            </div>

            <!-- Attachments List -->
            <h6 class="mb-3">Current Attachments</h6>
            <div v-if="currentAttachments.length === 0" class="text-muted text-center py-3">
              No attachments yet
            </div>
            <div v-else class="row g-2">
              <div v-for="attachment in currentAttachments" :key="attachment.id" class="col-12">
                <div class="card">
                  <div class="card-body p-3">
                    <div class="d-flex justify-content-between align-items-start">
                      <div class="flex-grow-1">
                        <div class="fw-bold">{{ attachment.originalName }}</div>
                        <div class="text-muted small">
                          <span class="badge bg-secondary me-2">{{ formatAttachmentType(attachment.attachmentType) }}</span>
                          <span v-if="attachment.fileSize">{{ formatFileSize(attachment.fileSize) }}</span>
                        </div>
                        <div v-if="attachment.description" class="text-muted small mt-1">{{ attachment.description }}</div>
                      </div>
                      <div class="d-flex gap-1">
                        <button class="btn btn-sm btn-outline-primary" @click="viewAttachment(attachment)" title="View">
                          👁️
                        </button>
                        <a :href="getAttachmentUrl(attachment)" :download="attachment.originalName" class="btn btn-sm btn-outline-primary" title="Download">
                          ⬇️
                        </a>
                        <button class="btn btn-sm btn-outline-danger" @click="deleteAttachment(attachment.id)" title="Delete">
                          🗑️
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

    <!-- View Attachment Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showViewAttachmentModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showViewAttachmentModal"
    >
      <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ viewingAttachment?.originalName }}</h5>
            <button type="button" class="btn-close" @click="closeViewAttachmentModal"></button>
          </div>
          <div class="modal-body text-center">
            <img v-if="isImageAttachment(viewingAttachment)" :src="getAttachmentUrl(viewingAttachment)" :alt="viewingAttachment.originalName" class="img-fluid" style="max-height: 70vh;" />
            <div v-else class="py-5">
              <p class="text-muted">Preview not available for this file type.</p>
              <a :href="getAttachmentUrl(viewingAttachment)" :download="viewingAttachment.originalName" class="btn btn-primary">
                Download File
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAssetsStore } from '@/stores/assets'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'
import FileUpload from '@/components/FileUpload.vue'

const assetsStore = useAssetsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const loading = ref(false)
const filter = ref('all')
const showAddModal = ref(false)
const showEditModal = ref(false)
const showServiceModal = ref(false)
const showDetailsModal = ref(false)
const showAttachmentsModal = ref(false)
const showViewAttachmentModal = ref(false)
const selectedAsset = ref(null)
const editingAssetId = ref(null)
const viewingAttachment = ref(null)
const fileUploadRef = ref(null)
const editFileUploadRef = ref(null)
const pendingFiles = ref([])
const attachmentType = ref('photo')
const attachmentDescription = ref('')
const uploading = ref(false)
const editFormAttachments = ref([])

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

const filteredAssets = computed(() => {
  if (filter.value === 'all') return assetsStore.allAssets
  if (filter.value === 'active') return assetsStore.activeAssets
  if (filter.value === 'archived') return assetsStore.archivedAssets
  if (filter.value === 'warranty') return assetsStore.assetsUnderWarranty
  return assetsStore.allAssets
})

const currentServiceRecords = computed(() => {
  if (!selectedAsset.value) return []
  return assetsStore.getServiceRecordsForAsset(selectedAsset.value.id)
})

const currentAttachments = computed(() => {
  if (!selectedAsset.value) return []
  return assetsStore.getAttachmentsForAsset(selectedAsset.value.id)
})

onMounted(async () => {
  loading.value = true
  try {
    await assetsStore.fetchAssets()
    await assetsStore.fetchStats()
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
    active: 'bg-primary',
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
  return status === 'active' ? 'bg-primary' : 'bg-danger'
}

const getWarrantyBgClass = (status) => {
  return status === 'active' ? 'bg-light-success' : 'bg-light-danger'
}

const getWarrantyProgress = (asset) => {
  if (!asset.warrantyDaysTotal || asset.warrantyDaysTotal === 0) return 0
  const passed = asset.warrantyDaysPassed || 0
  return Math.min(100, (passed / asset.warrantyDaysTotal) * 100)
}

const editGood = async (asset) => {
  editingAssetId.value = asset.id
  form.value = {
    name: asset.name,
    description: asset.description || '',
    category: asset.category || '',
    brand: asset.brand || '',
    model: asset.model || '',
    serialNumber: asset.serialNumber || '',
    purchaseDate: asset.purchaseDate?.split('T')[0] || '',
    purchasePrice: asset.purchasePrice,
    purchaseLocation: asset.purchaseLocation || '',
    warrantyStartDate: asset.warrantyStartDate?.split('T')[0] || null,
    warrantyEndDate: asset.warrantyEndDate?.split('T')[0] || null,
    warrantyProvider: asset.warrantyProvider || '',
    warrantyType: asset.warrantyType || '',
    status: asset.status,
    notes: asset.notes || ''
  }

  // Fetch attachments for the asset
  try {
    await assetsStore.fetchAttachments(asset.id)
    editFormAttachments.value = assetsStore.getAttachmentsForAsset(asset.id)
  } catch (error) {
    console.error('Error loading attachments:', error)
    editFormAttachments.value = []
  }

  showEditModal.value = true
}

const saveGood = async () => {
  loading.value = true
  try {
    if (showEditModal.value) {
      await assetsStore.updateAsset(editingAssetId.value, form.value)
    } else {
      await assetsStore.createAsset(form.value)
    }
    closeModal()
  } catch (error) {
    console.error('Error saving asset:', error)
    alert('Error saving asset. Please try again.')
  } finally {
    loading.value = false
  }
}

const confirmDelete = async (id) => {
  const confirmed = await confirm({
    title: 'Delete Asset',
    message: 'Are you sure you want to delete this asset? This will also delete all service records and attachments.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    loading.value = true
    try {
      await assetsStore.deleteAsset(id)
      success('Asset deleted successfully')
    } catch (err) {
      console.error('Error deleting asset:', err)
      error(err.response?.data?.message || err.message || 'Error deleting asset. Please try again.')
    } finally {
      loading.value = false
    }
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingAssetId.value = null
  editFormAttachments.value = []
  attachmentType.value = 'photo'
  attachmentDescription.value = ''
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

const viewDetails = async (asset) => {
  selectedAsset.value = asset
  await assetsStore.fetchServiceRecords(asset.id)
  await assetsStore.fetchAttachments(asset.id)
  showDetailsModal.value = true
}

const openServiceModal = async (asset) => {
  selectedAsset.value = asset
  await assetsStore.fetchServiceRecords(asset.id)
  showServiceModal.value = true
}

const closeServiceModal = () => {
  showServiceModal.value = false
  selectedAsset.value = null
  serviceForm.value = {
    serviceDate: '',
    serviceType: '',
    serviceProvider: '',
    cost: 0,
    description: '',
    warrantyCovered: false
  }
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedAsset.value = null
}

const addService = async () => {
  if (!selectedAsset.value) return

  loading.value = true
  try {
    await assetsStore.createServiceRecord(selectedAsset.value.id, serviceForm.value)
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
  if (!selectedAsset.value) return

  const confirmed = await confirm({
    title: 'Delete Service Record',
    message: 'Are you sure you want to delete this service record?',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (!confirmed) return

  loading.value = true
  try {
    await assetsStore.deleteServiceRecord(selectedAsset.value.id, serviceId)
    success('Service record deleted successfully')
  } catch (err) {
    console.error('Error deleting service record:', err)
    error(err.response?.data?.message || err.message || 'Error deleting service record. Please try again.')
  } finally {
    loading.value = false
  }
}

const openAttachmentsModal = async (asset) => {
  selectedAsset.value = asset
  await assetsStore.fetchAttachments(asset.id)
  showAttachmentsModal.value = true
}

const closeAttachmentsModal = () => {
  showAttachmentsModal.value = false
  selectedAsset.value = null
}

const deleteAttachment = async (attachmentId) => {
  if (!selectedAsset.value) return

  const confirmed = await confirm({
    title: 'Delete Attachment',
    message: 'Are you sure you want to delete this attachment?',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (!confirmed) return

  loading.value = true
  try {
    await assetsStore.deleteAttachment(selectedAsset.value.id, attachmentId)
    success('Attachment deleted successfully')
  } catch (err) {
    console.error('Error deleting attachment:', err)
    error(err.response?.data?.message || err.message || 'Error deleting attachment. Please try again.')
  } finally {
    loading.value = false
  }
}

const handleFilesUploaded = async (files) => {
  if (!selectedAsset.value || !files || files.length === 0) return

  uploading.value = true
  try {
    // Files are already uploaded to the server, now create attachment records for each file
    for (const file of files) {
      const attachmentData = {
        fileName: file.fileName,
        originalName: file.originalName,
        filePath: file.filePath,
        fileUrl: file.fileUrl,
        fileSize: file.fileSize,
        mimeType: file.mimeType,
        attachmentType: attachmentType.value,
        description: attachmentDescription.value
      }

      await assetsStore.addAttachment(selectedAsset.value.id, attachmentData)
    }

    // Refresh attachments list
    await assetsStore.fetchAttachments(selectedAsset.value.id)

    // Clear the file upload component
    fileUploadRef.value.clearFiles()

    alert('Attachments saved successfully!')
  } catch (error) {
    console.error('Error saving attachments:', error)
    alert('Error saving attachments. Please try again.')
  } finally {
    uploading.value = false
  }
}

const handleEditFormFilesUploaded = async (files) => {
  if (!editingAssetId.value || !files || files.length === 0) return

  uploading.value = true
  try {
    // Files are already uploaded to the server, now create attachment records for each file
    for (const file of files) {
      const attachmentData = {
        fileName: file.fileName,
        originalName: file.originalName,
        filePath: file.filePath,
        fileUrl: file.fileUrl,
        fileSize: file.fileSize,
        mimeType: file.mimeType,
        attachmentType: attachmentType.value,
        description: attachmentDescription.value
      }

      await assetsStore.addAttachment(editingAssetId.value, attachmentData)
    }

    // Refresh attachments list in the form
    await assetsStore.fetchAttachments(editingAssetId.value)
    editFormAttachments.value = assetsStore.getAttachmentsForAsset(editingAssetId.value)

    // Clear the file upload component
    if (editFileUploadRef.value) {
      editFileUploadRef.value.clearFiles()
    }

    // Reset attachment metadata
    attachmentType.value = 'photo'
    attachmentDescription.value = ''
  } catch (error) {
    console.error('Error saving attachments:', error)
    alert('Error saving attachments. Please try again.')
  } finally {
    uploading.value = false
  }
}

const deleteEditFormAttachment = async (attachmentId) => {
  if (!editingAssetId.value) return

  const confirmed = await confirm({
    title: 'Delete Attachment',
    message: 'Are you sure you want to delete this attachment?',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (!confirmed) return

  loading.value = true
  try {
    await assetsStore.deleteAttachment(editingAssetId.value, attachmentId)
    // Refresh the attachments list
    editFormAttachments.value = assetsStore.getAttachmentsForAsset(editingAssetId.value)
    success('Attachment deleted successfully')
  } catch (err) {
    console.error('Error deleting attachment:', err)
    error(err.response?.data?.message || err.message || 'Error deleting attachment. Please try again.')
  } finally {
    loading.value = false
  }
}

const viewAttachment = (attachment) => {
  viewingAttachment.value = attachment
  showViewAttachmentModal.value = true
}

const closeViewAttachmentModal = () => {
  showViewAttachmentModal.value = false
  viewingAttachment.value = null
}

const getAttachmentUrl = (attachment) => {
  if (!attachment) return ''
  return attachment.fileUrl || attachment.fileURL || ''
}

const isImageAttachment = (attachment) => {
  if (!attachment) return false
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  const fileName = attachment.originalName || attachment.fileName || ''
  const ext = '.' + fileName.split('.').pop().toLowerCase()
  return imageExtensions.includes(ext)
}

const formatAttachmentType = (type) => {
  if (!type) return 'File'
  const types = {
    photo: 'Photo',
    receipt: 'Receipt',
    warranty_document: 'Warranty',
    manual: 'Manual',
    other: 'Other'
  }
  return types[type] || type
}

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.bg-light-success {
  background-color: #e8f5e9;
  border: 1px solid #c8e6c9;
  color: #2e7d32;
}

.bg-light-danger {
  background-color: #ffebee;
  border: 1px solid #ffcdd2;
  color: #c62828;
}

.warranty-info-card {
  border: 1px solid;
}
</style>
