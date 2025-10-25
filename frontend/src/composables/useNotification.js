import { ref } from 'vue'

// Global state for alerts
const alertInstance = ref(null)

// Global state for confirm modals
const confirmState = ref({
  show: false,
  title: '',
  message: '',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'danger',
  onConfirm: null,
  onCancel: null
})

export function useNotification() {
  // Register the alert component instance
  const registerAlertInstance = (instance) => {
    alertInstance.value = instance
  }

  // Alert/Toast methods
  const showAlert = (message, type = 'info', duration = 4000) => {
    if (alertInstance.value) {
      alertInstance.value.addAlert(message, type, duration)
    }
  }

  const success = (message, duration = 4000) => {
    showAlert(message, 'success', duration)
  }

  const error = (message, duration = 4000) => {
    showAlert(message, 'error', duration)
  }

  const warning = (message, duration = 4000) => {
    showAlert(message, 'warning', duration)
  }

  const info = (message, duration = 4000) => {
    showAlert(message, 'info', duration)
  }

  // Confirm modal methods
  const confirm = (options) => {
    return new Promise((resolve) => {
      confirmState.value = {
        show: true,
        title: options.title || 'Confirm Action',
        message: options.message,
        confirmText: options.confirmText || 'Confirm',
        cancelText: options.cancelText || 'Cancel',
        variant: options.variant || 'danger',
        onConfirm: () => {
          confirmState.value.show = false
          resolve(true)
        },
        onCancel: () => {
          confirmState.value.show = false
          resolve(false)
        }
      }
    })
  }

  return {
    // Alert instance registration
    registerAlertInstance,

    // Alert methods
    showAlert,
    success,
    error,
    warning,
    info,

    // Confirm methods
    confirm,
    confirmState
  }
}
