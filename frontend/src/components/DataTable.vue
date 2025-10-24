<template>
  <div class="card">
    <div class="card-header" v-if="title || $slots.header">
      <div class="d-flex justify-content-between align-items-center">
        <h5 class="mb-0">{{ title }}</h5>
        <slot name="header"></slot>
      </div>
    </div>
    <div class="card-body p-0">
      <div v-if="!data || data.length === 0" class="p-4 text-center text-muted">
        <slot name="empty">
          <EmptyState
            :icon="emptyIcon"
            :title="emptyTitle"
            :description="emptyDescription"
          />
        </slot>
      </div>
      <div v-else class="table-responsive">
        <table class="table table-hover mb-0">
          <thead v-if="columns && columns.length > 0">
            <tr>
              <th
                v-for="column in columns"
                :key="column.key"
                :class="column.headerClass"
                @click="column.sortable && handleSort(column.key)"
                :style="{ cursor: column.sortable ? 'pointer' : 'default' }"
              >
                {{ column.label }}
                <span v-if="column.sortable && sortBy === column.key">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <slot name="body" :data="sortedData">
              <tr v-for="(item, index) in sortedData" :key="item.id || index">
                <td v-for="column in columns" :key="column.key" :class="column.cellClass">
                  <slot :name="`cell-${column.key}`" :item="item" :value="item[column.key]">
                    {{ formatCell(item[column.key], column) }}
                  </slot>
                </td>
              </tr>
            </slot>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card-footer" v-if="$slots.footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import EmptyState from './EmptyState.vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  data: {
    type: Array,
    required: true
  },
  columns: {
    type: Array,
    default: () => []
  },
  emptyIcon: {
    type: String,
    default: '📭'
  },
  emptyTitle: {
    type: String,
    default: 'No data available'
  },
  emptyDescription: {
    type: String,
    default: ''
  }
})

const sortBy = ref('')
const sortOrder = ref('asc')

const sortedData = computed(() => {
  if (!sortBy.value) return props.data

  return [...props.data].sort((a, b) => {
    const aVal = a[sortBy.value]
    const bVal = b[sortBy.value]

    if (aVal === bVal) return 0

    const comparison = aVal > bVal ? 1 : -1
    return sortOrder.value === 'asc' ? comparison : -comparison
  })
})

const handleSort = (key) => {
  if (sortBy.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = key
    sortOrder.value = 'asc'
  }
}

const formatCell = (value, column) => {
  if (column.formatter) {
    return column.formatter(value)
  }
  return value
}
</script>
