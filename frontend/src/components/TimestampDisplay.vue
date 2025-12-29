<template>
  <span :class="containerClass" :title="fullTimestamp">
    <i v-if="showIcon" :class="iconClass" class="me-1"></i>
    {{ formattedDate }}
  </span>
</template>

<script>
import { computed } from 'vue'
import { formatTimestamp, formatDateTime, formatDateTimeLong, getRelativeTime } from '@/utils/dateUtils'

export default {
  name: 'TimestampDisplay',
  props: {
    // The timestamp to display (ISO string or Date object)
    timestamp: {
      type: [String, Date],
      default: null
    },
    // Display format: 'auto', 'relative', 'datetime', 'long'
    format: {
      type: String,
      default: 'auto',
      validator: (value) => ['auto', 'relative', 'datetime', 'long'].includes(value)
    },
    // Whether to show time with date
    showTime: {
      type: Boolean,
      default: true
    },
    // Whether to show an icon
    showIcon: {
      type: Boolean,
      default: false
    },
    // Icon class (e.g., 'bi bi-clock' for Bootstrap Icons)
    icon: {
      type: String,
      default: 'bi bi-clock'
    },
    // Custom CSS class for the container
    customClass: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    const formattedDate = computed(() => {
      if (!props.timestamp) return 'N/A'

      switch (props.format) {
        case 'relative':
          return getRelativeTime(props.timestamp)
        case 'datetime':
          return formatDateTime(props.timestamp, props.showTime)
        case 'long':
          return formatDateTimeLong(props.timestamp)
        case 'auto':
        default:
          return formatTimestamp(props.timestamp, props.showTime)
      }
    })

    const fullTimestamp = computed(() => {
      if (!props.timestamp) return ''
      return formatDateTimeLong(props.timestamp)
    })

    const containerClass = computed(() => {
      const classes = ['timestamp-display']
      if (props.customClass) {
        classes.push(props.customClass)
      }
      return classes.join(' ')
    })

    const iconClass = computed(() => {
      return props.icon
    })

    return {
      formattedDate,
      fullTimestamp,
      containerClass,
      iconClass
    }
  }
}
</script>

<style scoped>
.timestamp-display {
  font-size: 0.875rem;
  color: #6c757d;
}

.timestamp-display i {
  font-size: 0.875rem;
}
</style>
