<template>
  <Teleport to="body">
    <div class="toast-container position-fixed top-0 end-0 p-3" style="z-index: 9999;">
      <TransitionGroup name="toast">
        <div
          v-for="alert in alerts"
          :key="alert.id"
          class="toast show"
          role="alert"
        >
          <div class="toast-header" :class="getHeaderClass(alert.type)">
            <span class="me-2">{{ getIcon(alert.type) }}</span>
            <strong class="me-auto">{{ getTitle(alert.type) }}</strong>
            <button
              type="button"
              class="btn-close"
              :class="{ 'btn-close-white': ['danger', 'success'].includes(alert.type) }"
              @click="removeAlert(alert.id)"
            ></button>
          </div>
          <div class="toast-body">
            {{ alert.message }}
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { ref } from 'vue'

const alerts = ref([])
let alertIdCounter = 0

const getIcon = (type) => {
  const icons = {
    success: '✅',
    error: '❌',
    warning: '⚠️',
    info: 'ℹ️'
  }
  return icons[type] || icons.info
}

const getTitle = (type) => {
  const titles = {
    success: 'Success',
    error: 'Error',
    warning: 'Warning',
    info: 'Info'
  }
  return titles[type] || titles.info
}

const getHeaderClass = (type) => {
  const classes = {
    success: 'toast-success',
    error: 'toast-error',
    warning: 'toast-warning',
    info: 'toast-info'
  }
  return classes[type] || classes.info
}

const addAlert = (message, type = 'info', duration = 4000) => {
  const id = alertIdCounter++
  alerts.value.push({ id, message, type })

  if (duration > 0) {
    setTimeout(() => {
      removeAlert(id)
    }, duration)
  }
}

const removeAlert = (id) => {
  const index = alerts.value.findIndex(alert => alert.id === id)
  if (index > -1) {
    alerts.value.splice(index, 1)
  }
}

defineExpose({
  addAlert,
  success: (message, duration) => addAlert(message, 'success', duration),
  error: (message, duration) => addAlert(message, 'error', duration),
  warning: (message, duration) => addAlert(message, 'warning', duration),
  info: (message, duration) => addAlert(message, 'info', duration)
})
</script>

<style scoped>
.toast {
  min-width: 300px;
  max-width: 350px;
  margin-bottom: 0.5rem;
}

/* Professional toast colors */
.toast-success {
  background-color: #dbeafe;
  color: #1e40af;
  border-left: 4px solid #3b82f6;
}

.toast-error {
  background-color: #fee2e2;
  color: #991b1b;
  border-left: 4px solid #ef4444;
}

.toast-warning {
  background-color: #fef3c7;
  color: #92400e;
  border-left: 4px solid #f59e0b;
}

.toast-info {
  background-color: #e0e7ff;
  color: #3730a3;
  border-left: 4px solid #6366f1;
}

/* Dark mode support */
.dark-mode .toast-success {
  background-color: #1e3a5f;
  color: #93c5fd;
  border-left-color: #3b82f6;
}

.dark-mode .toast-error {
  background-color: #5f1e1e;
  color: #fca5a5;
  border-left-color: #ef4444;
}

.dark-mode .toast-warning {
  background-color: #5f4e1e;
  color: #fcd34d;
  border-left-color: #f59e0b;
}

.dark-mode .toast-info {
  background-color: #312e5f;
  color: #a5b4fc;
  border-left-color: #6366f1;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
