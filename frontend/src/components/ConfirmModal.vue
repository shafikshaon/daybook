<template>
  <div
    class="modal fade"
    :class="{ 'show d-block': show }"
    style="background-color: rgba(0,0,0,0.5);"
    @click.self="handleCancel"
    v-if="show"
  >
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header" :class="headerClass">
          <h5 class="modal-title">
            <span class="me-2">{{ icon }}</span>
            {{ title }}
          </h5>
          <button
            type="button"
            class="btn-close"
            :class="{ 'btn-close-white': variant === 'danger' }"
            @click="handleCancel"
          ></button>
        </div>
        <div class="modal-body">
          <p class="mb-0">{{ message }}</p>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-secondary"
            @click="handleCancel"
          >
            {{ cancelText }}
          </button>
          <button
            type="button"
            :class="confirmButtonClass"
            @click="handleConfirm"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  cancelText: {
    type: String,
    default: 'Cancel'
  },
  variant: {
    type: String,
    default: 'danger',
    validator: (value) => ['danger', 'warning', 'primary', 'success'].includes(value)
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const icon = computed(() => {
  const icons = {
    danger: '⚠️',
    warning: '⚠️',
    primary: 'ℹ️',
    success: '✅'
  }
  return icons[props.variant]
})

const headerClass = computed(() => {
  const classes = {
    danger: 'modal-header-danger',
    warning: 'modal-header-warning',
    primary: 'modal-header-primary',
    success: 'modal-header-success'
  }
  return classes[props.variant] || ''
})

const confirmButtonClass = computed(() => {
  const classes = {
    danger: 'btn btn-modal-danger',
    warning: 'btn btn-modal-warning',
    primary: 'btn btn-primary',
    success: 'btn btn-modal-success'
  }
  return classes[props.variant] || 'btn btn-primary'
})

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
/* Professional modal header colors */
.modal-header-danger {
  background-color: #fee2e2;
  color: #991b1b;
  border-bottom: 2px solid #ef4444;
}

.modal-header-warning {
  background-color: #fef3c7;
  color: #92400e;
  border-bottom: 2px solid #f59e0b;
}

.modal-header-primary {
  background-color: #e0e7ff;
  color: #3730a3;
  border-bottom: 2px solid #6366f1;
}

.modal-header-success {
  background-color: #dbeafe;
  color: #1e40af;
  border-bottom: 2px solid #3b82f6;
}

/* Professional button colors */
.btn-modal-danger {
  color: #ffffff;
  background-color: #ef4444;
  border-color: #ef4444;
}

.btn-modal-danger:hover {
  background-color: #dc2626;
  border-color: #dc2626;
}

.btn-modal-warning {
  color: #ffffff;
  background-color: #f59e0b;
  border-color: #f59e0b;
}

.btn-modal-warning:hover {
  background-color: #d97706;
  border-color: #d97706;
}

.btn-modal-success {
  color: #ffffff;
  background-color: #3b82f6;
  border-color: #3b82f6;
}

.btn-modal-success:hover {
  background-color: #2563eb;
  border-color: #2563eb;
}

/* Dark mode support */
.dark-mode .modal-header-danger {
  background-color: #5f1e1e;
  color: #fca5a5;
  border-bottom-color: #ef4444;
}

.dark-mode .modal-header-warning {
  background-color: #5f4e1e;
  color: #fcd34d;
  border-bottom-color: #f59e0b;
}

.dark-mode .modal-header-primary {
  background-color: #312e5f;
  color: #a5b4fc;
  border-bottom-color: #6366f1;
}

.dark-mode .modal-header-success {
  background-color: #1e3a5f;
  color: #93c5fd;
  border-bottom-color: #3b82f6;
}
</style>
