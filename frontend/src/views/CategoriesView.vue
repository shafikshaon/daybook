<template>
  <div class="categories-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Categories</h1>
      <div class="d-flex gap-2">
        <button
          class="btn"
          :class="isReorderMode ? 'btn-secondary' : 'btn-outline-secondary'"
          @click="toggleReorderMode"
        >
          {{ isReorderMode ? 'Done' : 'Reorder' }}
        </button>
        <button class="btn btn-primary" @click="showAddModal = true">
          + Add Category
        </button>
      </div>
    </div>

    <!-- Info Alert -->
    <div class="alert mb-4" :class="isReorderMode ? 'alert-warning' : 'alert-info'" role="alert">
      <strong v-if="!isReorderMode">Category Management:</strong>
      <strong v-else>Reorder Mode:</strong>
      <span v-if="!isReorderMode">
        Manage your income and expense categories. Categories help you organize and track your financial transactions.
      </span>
      <span v-else>
        Drag and drop categories to reorder them within each type. Changes are saved automatically.
        <span v-if="filter === 'all'" class="text-danger">Please select a specific type to enable reordering.</span>
      </span>
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
                <tr
                  v-for="(category, index) in filteredCategories"
                  :key="category.id"
                  :draggable="isReorderMode && canDragCategory(category)"
                  @dragstart="handleDragStart($event, category, index)"
                  @dragover="handleDragOver($event)"
                  @drop="handleDrop($event, category, index)"
                  @dragend="handleDragEnd"
                  :class="{
                    'dragging': draggedCategory?.id === category.id,
                    'drag-over': dragOverIndex === index,
                    'draggable-row': isReorderMode && canDragCategory(category)
                  }"
                >
                  <td>
                    <span v-if="isReorderMode && canDragCategory(category)" class="drag-handle me-2">⋮⋮</span>
                    <span class="fs-4">{{ category.icon }}</span>
                  </td>
                  <td>
                    <strong>{{ category.name }}</strong>
                    <span v-if="category.isDefault" class="badge bg-secondary ms-2">Default</span>
                  </td>
                  <td>
                    <span class="badge" :class="getCategoryTypeBadgeClass(category.type)">
                      {{ getCategoryTypeLabel(category.type) }}
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
                      v-if="!isReorderMode"
                      class="btn btn-sm btn-primary me-1"
                      @click="editCategory(category)"
                    >
                      Edit
                    </button>
                    <button
                      v-if="!isReorderMode"
                      class="btn btn-sm btn-danger"
                      @click="confirmDelete(category)"
                    >
                      Delete
                    </button>
                    <span v-else class="text-muted">
                      <small>Order: {{ category.order }}</small>
                    </span>
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

// Drag and drop state
const isReorderMode = ref(false)
const draggedCategory = ref(null)
const draggedIndex = ref(null)
const dragOverIndex = ref(null)

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
  let cats = categories.value
  if (filter.value !== 'all') {
    cats = cats.filter(cat => cat.type === filter.value)
  }
  // Ensure categories are sorted by order
  return [...cats].sort((a, b) => (a.order || 0) - (b.order || 0))
})

// Helper functions for badges
const getCategoryTypeBadgeClass = (type) => {
  switch(type) {
    case 'income': return 'bg-success'
    case 'expense': return 'bg-danger'
    case 'transfer': return 'bg-info'
    default: return 'bg-secondary'
  }
}

const getCategoryTypeLabel = (type) => {
  switch(type) {
    case 'income': return 'Income'
    case 'expense': return 'Expense'
    case 'transfer': return 'Transfer'
    default: return type
  }
}

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

// Drag and drop functions
const toggleReorderMode = () => {
  isReorderMode.value = !isReorderMode.value
  if (!isReorderMode.value) {
    // Clean up drag state when exiting reorder mode
    draggedCategory.value = null
    draggedIndex.value = null
    dragOverIndex.value = null
  }
}

const canDragCategory = (category) => {
  // Only allow dragging within the filtered type
  // If showing all, user must filter by type to reorder
  if (filter.value === 'all') return false
  return true
}

const handleDragStart = (event, category, index) => {
  if (!canDragCategory(category)) return

  draggedCategory.value = category
  draggedIndex.value = index
  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/html', event.target)
}

const handleDragOver = (event) => {
  if (!isReorderMode.value) return
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
  return false
}

const handleDrop = async (event, targetCategory, targetIndex) => {
  if (!isReorderMode.value || !draggedCategory.value) return

  event.stopPropagation()
  event.preventDefault()

  const draggedCat = draggedCategory.value
  const fromIndex = draggedIndex.value
  const toIndex = targetIndex

  if (fromIndex === toIndex) {
    handleDragEnd()
    return
  }

  // Can only reorder within same type
  if (draggedCat.type !== targetCategory.type) {
    error('Cannot reorder categories across different types')
    handleDragEnd()
    return
  }

  try {
    // Get all categories of this type (currently filtered)
    const typedCategories = [...filteredCategories.value]

    // Reorder the array locally
    const [movedItem] = typedCategories.splice(fromIndex, 1)
    typedCategories.splice(toIndex, 0, movedItem)

    // Recalculate order values (1-indexed)
    const categoryOrders = typedCategories.map((cat, idx) => ({
      id: cat.id,
      order: idx + 1
    }))

    // Save to backend
    await categoriesStore.reorderCategories(categoryOrders)

    success('Categories reordered successfully')
  } catch (err) {
    error(err.message || 'Failed to reorder categories')
    // Refresh categories to restore correct order
    await categoriesStore.fetchCategories()
  } finally {
    handleDragEnd()
  }
}

const handleDragEnd = () => {
  draggedCategory.value = null
  draggedIndex.value = null
  dragOverIndex.value = null
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

/* Drag and drop styles */
.draggable-row {
  cursor: grab;
  user-select: none;
}

.draggable-row:active {
  cursor: grabbing;
}

.draggable-row.dragging {
  opacity: 0.5;
  background-color: #f8f9fa;
}

.draggable-row.drag-over {
  border-top: 2px solid #0d6efd;
}

.drag-handle {
  color: #6c757d;
  cursor: grab;
  font-size: 1.2rem;
  line-height: 1;
  vertical-align: middle;
}

.drag-handle:active {
  cursor: grabbing;
}

/* Smooth transitions for reordering */
tbody tr {
  transition: background-color 0.2s ease;
}
</style>
