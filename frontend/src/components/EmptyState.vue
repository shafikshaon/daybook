<template>
  <div class="empty-state text-center py-5">
    <div class="empty-state-icon mb-3" v-if="icon">
      {{ icon }}
    </div>
    <h5 class="empty-state-title mb-2">{{ title }}</h5>
    <p class="text-muted mb-4" v-if="description">{{ description }}</p>
    <button
      v-if="actionText"
      class="btn btn-primary"
      @click="handleAction"
    >
      {{ actionText }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  icon: {
    type: String,
    default: '📭'
  },
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: ''
  },
  actionText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['action'])

const handleAction = () => {
  emit('action')
}
</script>

<style scoped>
.empty-state-icon {
  font-size: 4rem;
  opacity: 0.5;
}

.empty-state-title {
  color: #6b7280;
}
</style>
