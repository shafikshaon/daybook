<template>
  <div>
    <div class="d-flex justify-content-between mb-1" v-if="showLabels">
      <span>{{ leftLabel }}</span>
      <span :class="labelClass">{{ rightLabel }}</span>
    </div>
    <div class="progress" :style="{ height: height }">
      <div
        class="progress-bar"
        :class="barClass"
        :style="barStyle"
        role="progressbar"
        :aria-valuenow="percentage"
        aria-valuemin="0"
        aria-valuemax="100"
      >
        <span v-if="showPercentage">{{ Math.round(percentage) }}%</span>
      </div>
    </div>
    <small v-if="bottomText" class="text-muted d-block mt-1">{{ bottomText }}</small>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  percentage: {
    type: Number,
    required: true,
    validator: (value) => value >= 0 && value <= 100
  },
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'success', 'warning', 'danger', 'info', 'purple'].includes(value)
  },
  height: {
    type: String,
    default: '10px'
  },
  showPercentage: {
    type: Boolean,
    default: false
  },
  showLabels: {
    type: Boolean,
    default: false
  },
  leftLabel: {
    type: String,
    default: ''
  },
  rightLabel: {
    type: String,
    default: ''
  },
  bottomText: {
    type: String,
    default: ''
  },
  striped: {
    type: Boolean,
    default: false
  },
  animated: {
    type: Boolean,
    default: false
  }
})

const barClass = computed(() => {
  const classes = []

  if (props.variant === 'purple') {
    classes.push('bg-purple')
  } else {
    classes.push(`bg-${props.variant}`)
  }

  if (props.striped) classes.push('progress-bar-striped')
  if (props.animated) classes.push('progress-bar-animated')

  return classes
})

const barStyle = computed(() => {
  return {
    width: `${Math.min(props.percentage, 100)}%`
  }
})

const labelClass = computed(() => {
  if (props.percentage >= 100) return 'text-danger fw-bold'
  if (props.percentage >= 80) return 'text-warning fw-bold'
  return 'text-muted'
})
</script>
