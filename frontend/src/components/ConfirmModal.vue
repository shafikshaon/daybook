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
  return props.variant === 'danger' ? 'bg-danger text-white' : ''
})

const confirmButtonClass = computed(() => {
  return `btn btn-${props.variant}`
})

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
}
</script>
