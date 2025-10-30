<template>
  <div class="account-types-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Account Types</h1>
      <button class="btn btn-primary" @click="showAddModal = true">
        + Add Account Type
      </button>
    </div>

    <!-- Info Alert -->
    <div class="alert alert-info mb-4" role="alert">
      <strong>Account Type Management:</strong> Manage your account types. You can edit, delete, or add new types. Note: You cannot delete types that are currently in use by existing accounts.
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center p-4">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Account Types List -->
    <div v-else>
      <div class="card">
        <div class="card-body p-0">
          <div v-if="accountTypes.length === 0" class="p-4 text-center text-muted">
            <p>No account types yet.</p>
          </div>
          <div v-else class="table-responsive">
            <table class="table table-hover mb-0">
              <thead>
                <tr>
                  <th>Icon</th>
                  <th>Name</th>
                  <th>Order</th>
                  <th>Status</th>
                  <th>Description</th>
                  <th class="text-center">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="type in accountTypes" :key="type.id">
                  <td>
                    <span class="fs-4">{{ type.icon }}</span>
                  </td>
                  <td>
                    <strong>{{ type.name }}</strong>
                  </td>
                  <td>
                    <strong>{{ type.sortOrder }}</strong>
                  </td>
                  <td>
                    <strong>{{ type.active ? 'Active' : 'Inactive' }}</strong>
                  </td>
                  <td>
                    <span class="text-muted">{{ type.description || '-' }}</span>
                  </td>
                  <td class="text-center">
                    <button
                      class="btn btn-sm btn-primary me-1"
                      @click="editType(type)"
                    >
                      Edit
                    </button>
                    <button
                      class="btn btn-sm btn-danger"
                      @click="confirmDelete(type)"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      tabindex="-1"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">
              {{ showEditModal ? 'Edit Account Type' : 'Add Account Type' }}
            </h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveType">
              <div class="mb-3">
                <label class="form-label">Name *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.name"
                  required
                  placeholder="e.g., Savings Account"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Icon *</label>
                <select class="form-select" v-model="form.icon" required>
                  <option value="">Select icon...</option>
                  <option v-for="iconOption in iconOptions" :key="iconOption.emoji" :value="iconOption.emoji">
                    {{ iconOption.emoji }} {{ iconOption.label }}
                  </option>
                </select>
                <div v-if="form.icon" class="mt-2 p-3 border rounded text-center bg-light">
                  <span class="fs-1">{{ form.icon }}</span>
                  <div class="text-muted small mt-1">Preview</div>
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <textarea
                  class="form-control"
                  v-model="form.description"
                  rows="2"
                  placeholder="Optional description"
                ></textarea>
              </div>

              <div class="mb-3">
                <label class="form-label">Sort Order</label>
                <input
                  type="number"
                  class="form-control"
                  v-model.number="form.sortOrder"
                  placeholder="0"
                />
                <small class="text-muted">Lower numbers appear first</small>
              </div>

              <div class="form-check mb-3">
                <input
                  type="checkbox"
                  class="form-check-input"
                  id="activeCheck"
                  v-model="form.active"
                />
                <label class="form-check-label" for="activeCheck">
                  Active
                </label>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  {{ showEditModal ? 'Update' : 'Create' }} Type
                </button>
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
import { useAccountTypesStore } from '@/stores/accountTypes'
import { useNotification } from '@/composables/useNotification'

const accountTypesStore = useAccountTypesStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingType = ref(null)

const form = ref({
  name: '',
  icon: '',
  description: '',
  active: true,
  sortOrder: 0
})

// Finance-related icon options
const iconOptions = [
  { emoji: '💵', label: 'Cash (Dollar Bill)' },
  { emoji: '💴', label: 'Cash (Yen)' },
  { emoji: '💶', label: 'Cash (Euro)' },
  { emoji: '💷', label: 'Cash (Pound)' },
  { emoji: '💰', label: 'Money Bag' },
  { emoji: '💸', label: 'Money with Wings' },
  { emoji: '🏦', label: 'Bank' },
  { emoji: '🏧', label: 'ATM' },
  { emoji: '💳', label: 'Credit Card' },
  { emoji: '💎', label: 'Gem Stone' },
  { emoji: '📱', label: 'Mobile Phone' },
  { emoji: '📊', label: 'Bar Chart' },
  { emoji: '📈', label: 'Chart Increasing' },
  { emoji: '📉', label: 'Chart Decreasing' },
  { emoji: '💹', label: 'Chart with Yen' },
  { emoji: '🪙', label: 'Coin' },
  { emoji: '💲', label: 'Dollar Sign' },
  { emoji: '🤑', label: 'Money Face' },
  { emoji: '💼', label: 'Briefcase' },
  { emoji: '🏪', label: 'Convenience Store' },
  { emoji: '🏬', label: 'Department Store' },
  { emoji: '🏠', label: 'House' },
  { emoji: '🏡', label: 'House with Garden' },
  { emoji: '🏘️', label: 'Houses' },
  { emoji: '🏗️', label: 'Building Construction' },
  { emoji: '🏢', label: 'Office Building' },
  { emoji: '🏛️', label: 'Classical Building' },
  { emoji: '⚡', label: 'Lightning (Energy)' },
  { emoji: '🔑', label: 'Key' },
  { emoji: '🎓', label: 'Graduation Cap' },
  { emoji: '🚗', label: 'Car' },
  { emoji: '✈️', label: 'Airplane' },
  { emoji: '🎁', label: 'Gift' },
  { emoji: '🛒', label: 'Shopping Cart' },
  { emoji: '🛍️', label: 'Shopping Bags' },
  { emoji: '📋', label: 'Clipboard' },
  { emoji: '📝', label: 'Memo' },
  { emoji: '🔒', label: 'Lock' },
  { emoji: '🔓', label: 'Unlock' },
  { emoji: '⭐', label: 'Star' },
  { emoji: '💡', label: 'Light Bulb' },
  { emoji: '📦', label: 'Package' },
  { emoji: '🎯', label: 'Target' },
  { emoji: '🔔', label: 'Bell' },
  { emoji: '⏰', label: 'Alarm Clock' },
  { emoji: '📅', label: 'Calendar' },
  { emoji: '🌐', label: 'Globe' },
  { emoji: '₿', label: 'Bitcoin' },
  { emoji: '🍎', label: 'Apple' },
  { emoji: '💙', label: 'Blue Heart' },
  { emoji: '🧾', label: 'Receipt' },
  { emoji: '📄', label: 'Document' }
]

const accountTypes = computed(() => accountTypesStore.allAccountTypes)
const loading = computed(() => accountTypesStore.loading)

const editType = (type) => {
  editingType.value = type
  form.value = {
    name: type.name,
    icon: type.icon || '',
    description: type.description || '',
    active: type.active,
    sortOrder: type.sortOrder || 0
  }
  showEditModal.value = true
}

const confirmDelete = async (type) => {
  const confirmed = await confirm({
    title: 'Delete Account Type',
    message: `Are you sure you want to delete "${type.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await accountTypesStore.deleteAccountType(type.id)
      success('Account type deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting account type')
    }
  }
}

const saveType = async () => {
  try {
    if (showEditModal.value) {
      await accountTypesStore.updateAccountType(editingType.value.id, form.value)
      success('Account type updated successfully')
    } else {
      await accountTypesStore.createAccountType(form.value)
      success('Account type created successfully')
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving account type')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingType.value = null
  form.value = {
    name: '',
    icon: '',
    description: '',
    active: true,
    sortOrder: 0
  }
}

onMounted(async () => {
  await accountTypesStore.fetchAccountTypes()
})
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
