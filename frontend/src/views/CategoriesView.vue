<template>
  <div class="categories-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Categories</h1>
      <button class="btn btn-primary" @click="showAddModal = true">
        + Add Category
      </button>
    </div>

    <!-- Info Alert -->
    <div class="alert alert-info mb-4" role="alert">
      <strong>Category Management:</strong> Manage your income and expense categories. Categories help you organize and track your financial transactions.
    </div>

    <!-- Filter Tabs -->
    <ul class="nav nav-tabs mb-4">
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'all' }" @click="filter = 'all'" style="cursor: pointer;">
          All Categories
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'income' }" @click="filter = 'income'" style="cursor: pointer;">
          Income
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'expense' }" @click="filter = 'expense'" style="cursor: pointer;">
          Expense
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" :class="{ active: filter === 'transfer' }" @click="filter = 'transfer'" style="cursor: pointer;">
          Transfer
        </a>
      </li>
    </ul>

    <!-- Loading State -->
    <div v-if="loading" class="text-center p-4">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Categories List -->
    <div v-else>
      <div class="card">
        <div class="card-body p-0">
          <div v-if="filteredCategories.length === 0" class="p-4 text-center text-muted">
            <p>No categories yet.</p>
          </div>
          <div v-else class="table-responsive">
            <table class="table table-hover mb-0">
              <thead>
                <tr>
                  <th>Icon</th>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Color</th>
                  <th>Description</th>
                  <th class="text-center">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="category in filteredCategories" :key="category.id">
                  <td>
                    <span class="fs-4">{{ category.icon }}</span>
                  </td>
                  <td>
                    <strong>{{ category.name }}</strong>
                    <span v-if="category.isDefault" class="badge bg-secondary ms-2">Default</span>
                  </td>
                  <td>
                    <span class="badge" :class="category.type === 'income' ? 'bg-success' : 'bg-danger'">
                      {{ category.type === 'income' ? 'Income' : 'Expense' }}
                    </span>
                  </td>
                  <td>
                    <span class="d-inline-flex align-items-center">
                      <span class="color-preview" :style="{ backgroundColor: category.color }"></span>
                      <span class="ms-2">{{ category.color }}</span>
                    </span>
                  </td>
                  <td>
                    <span class="text-muted">{{ category.description || '-' }}</span>
                  </td>
                  <td class="text-center">
                    <button
                      class="btn btn-sm btn-primary me-1"
                      @click="editCategory(category)"
                    >
                      Edit
                    </button>
                    <button
                      class="btn btn-sm btn-danger"
                      @click="confirmDelete(category)"
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
              {{ showEditModal ? 'Edit Category' : 'Add Category' }}
            </h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveCategory">
              <div class="mb-3">
                <label class="form-label">Name *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.name"
                  required
                  placeholder="e.g., Food & Dining"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Type *</label>
                <select class="form-select" v-model="form.type" required @change="updateIconOptions">
                  <option value="">Select type...</option>
                  <option value="income">Income</option>
                  <option value="expense">Expense</option>
                  <option value="transfer">Transfer</option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Icon *</label>
                <select class="form-select" v-model="form.icon" required :disabled="!form.type">
                  <option value="">Select icon...</option>
                  <option v-for="icon in currentIconOptions" :key="icon" :value="icon">
                    {{ icon }}
                  </option>
                </select>
                <div v-if="form.icon" class="mt-2 p-3 border rounded text-center bg-light">
                  <span class="fs-1">{{ form.icon }}</span>
                  <div class="text-muted small mt-1">Preview</div>
                </div>
              </div>

              <div class="mb-3">
                <label class="form-label">Color *</label>
                <div class="d-flex align-items-center gap-2">
                  <input
                    type="color"
                    class="form-control form-control-color"
                    v-model="form.color"
                    required
                  />
                  <input
                    type="text"
                    class="form-control"
                    v-model="form.color"
                    required
                    placeholder="#3B82F6"
                  />
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

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  {{ showEditModal ? 'Update' : 'Create' }} Category
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
import { useCategoriesStore } from '@/stores/categories'
import { useNotification } from '@/composables/useNotification'

const categoriesStore = useCategoriesStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingCategory = ref(null)
const filter = ref('all')

const form = ref({
  name: '',
  type: '',
  icon: '',
  color: '#3B82F6',
  description: ''
})

// Icon options based on category type
const iconOptions = {
  income: [
    "💰", "💵", "💴", "💶", "💷", "💸", "💳", "💎",
    "📈", "📊", "💼", "🏢", "🏦", "🎁", "🎉", "⭐",
    "✨", "💫", "🌟", "🔥", "↩️", "✅", "💯", "🎯"
  ],
  expense: [
    "🍔", "🍕", "🍜", "🍱", "🛒", "🛍️", "🎬", "🎮",
    "🚗", "🚕", "🚌", "✈️", "🏠", "🏡", "💡", "💧",
    "📱", "💻", "🏥", "💊", "💪", "📚", "🎓", "🎵",
    "🎸", "🎨", "👕", "👗", "👠", "💅", "💇", "🛡️",
    "📺", "🎧", "⚽", "🏀", "🎾", "🎭", "🍿", "☕",
    "🍺", "🍷", "🎂", "🌮", "🍣", "🍦", "💸", "💳"
  ],
  transfer: [
      '→'
  ]
}

const currentIconOptions = computed(() => {
  return form.value.type ? iconOptions[form.value.type] || [] : []
})

const categories = computed(() => categoriesStore.allCategories)
const loading = computed(() => categoriesStore.loading)

const filteredCategories = computed(() => {
  if (filter.value === 'all') {
    return categories.value
  }
  return categories.value.filter(cat => cat.type === filter.value)
})

const updateIconOptions = () => {
  // Reset icon when type changes
  if (form.value.icon && currentIconOptions.value.length > 0) {
    if (!currentIconOptions.value.includes(form.value.icon)) {
      form.value.icon = ''
    }
  }

  // Set default color based on type
  if (form.value.type === 'income') {
    form.value.color = '#10B981'
  } else if (form.value.type === 'expense') {
    form.value.color = '#EF4444'
  }
}

const editCategory = (category) => {
  editingCategory.value = category
  form.value = {
    name: category.name,
    type: category.type,
    icon: category.icon || '',
    color: category.color || '#3B82F6',
    description: category.description || ''
  }
  showEditModal.value = true
}

const confirmDelete = async (category) => {
  const confirmed = await confirm({
    title: 'Delete Category',
    message: `Are you sure you want to delete "${category.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await categoriesStore.deleteCategory(category.id)
      success('Category deleted successfully')
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting category')
    }
  }
}

const saveCategory = async () => {
  try {
    if (showEditModal.value) {
      await categoriesStore.updateCategory(editingCategory.value.id, form.value)
      success('Category updated successfully')
    } else {
      await categoriesStore.createCategory(form.value)
      success('Category created successfully')
    }
    closeModal()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving category')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingCategory.value = null
  form.value = {
    name: '',
    type: '',
    icon: '',
    color: '#3B82F6',
    description: ''
  }
}

onMounted(async () => {
  await categoriesStore.fetchCategories()
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

.color-preview {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #dee2e6;
}

.nav-tabs .nav-link {
  cursor: pointer;
}
</style>
