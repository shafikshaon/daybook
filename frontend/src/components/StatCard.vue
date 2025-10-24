<template>
  <div class="stat-card">
    <div class="stat-icon" :class="iconColor">
      {{ icon }}
    </div>
    <div class="stat-value">{{ value }}</div>
    <div class="stat-label">{{ label }}</div>
    <div v-if="change" class="stat-change" :class="changeClass">
      {{ change }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  icon: {
    type: String,
    required: true
  },
  iconColor: {
    type: String,
    default: 'purple',
    validator: (value) => ['purple', 'green', 'blue', 'orange', 'red'].includes(value)
  },
  value: {
    type: [String, Number],
    required: true
  },
  label: {
    type: String,
    required: true
  },
  change: {
    type: String,
    default: ''
  },
  changeType: {
    type: String,
    default: 'positive',
    validator: (value) => ['positive', 'negative', 'neutral'].includes(value)
  }
})

const changeClass = computed(() => {
  if (!props.change) return ''
  return props.changeType
})
</script>
