# Timestamp Usage Guide

This guide explains how to properly display `createdAt` and `updatedAt` timestamps throughout the application.

## Backend (Go)

### Model Structure

All GORM models automatically include `createdAt` and `updatedAt` fields:

```go
type Account struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
    UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
    Name        string         `gorm:"not null" json:"name"`
    Balance     float64        `gorm:"default:0" json:"balance"`
    CreatedAt   time.Time      `json:"createdAt"`   // Automatically set by GORM on create
    UpdatedAt   time.Time      `json:"updatedAt"`   // Automatically updated by GORM on save
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**GORM Handles These Automatically:**
- `CreatedAt` is set when `db.Create()` is called
- `UpdatedAt` is updated when `db.Save()` or `db.Update()` is called
- No manual intervention needed!

### API Response Format

Timestamps are automatically serialized to RFC3339 format (ISO 8601):

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Checking Account",
  "balance": 1500.50,
  "createdAt": "2025-12-28T15:30:00Z",
  "updatedAt": "2025-12-28T16:45:30Z"
}
```

## Frontend (Vue 3)

### Utility Functions

Located in `src/utils/dateUtils.js`:

#### 1. **formatDateTime(date, includeTime)**
Basic date formatting with optional time.

```javascript
import { formatDateTime } from '@/utils/dateUtils'

// Without time
formatDateTime('2025-12-28T15:30:00Z')
// Output: "Dec 28, 2025"

// With time
formatDateTime('2025-12-28T15:30:00Z', true)
// Output: "Dec 28, 2025 at 3:30 PM"
```

#### 2. **formatDateTimeLong(date)**
Full date and time format.

```javascript
import { formatDateTimeLong } from '@/utils/dateUtils'

formatDateTimeLong('2025-12-28T15:30:45Z')
// Output: "December 28, 2025 at 3:30:45 PM"
```

#### 3. **getRelativeTime(date)**
Shows how long ago something happened.

```javascript
import { getRelativeTime } from '@/utils/dateUtils'

getRelativeTime('2025-12-28T14:30:00Z') // If now is 14:35
// Output: "5 minutes ago"

getRelativeTime('2025-12-27T14:30:00Z')
// Output: "1 day ago"
```

#### 4. **formatTimestamp(date, showTime)** ⭐ Recommended
Smart formatting - shows relative time if recent, otherwise full date.

```javascript
import { formatTimestamp } from '@/utils/dateUtils'

// Recent (< 7 days): Shows relative time
formatTimestamp('2025-12-28T14:30:00Z')
// Output: "2 hours ago"

// Older (> 7 days): Shows full date
formatTimestamp('2025-11-15T10:00:00Z', true)
// Output: "Nov 15, 2025 at 10:00 AM"
```

#### 5. **formatTimeOnly(date)**
Shows only the time portion.

```javascript
import { formatTimeOnly } from '@/utils/dateUtils'

formatTimeOnly('2025-12-28T15:30:00Z')
// Output: "3:30 PM"
```

#### 6. **formatDateShort(date)**
Shows date in short numeric format.

```javascript
import { formatDateShort } from '@/utils/dateUtils'

formatDateShort('2025-12-28T15:30:00Z')
// Output: "12/28/2025"
```

### TimestampDisplay Component

Reusable component for displaying timestamps.

**Location:** `src/components/TimestampDisplay.vue`

#### Basic Usage

```vue
<template>
  <TimestampDisplay :timestamp="account.createdAt" />
</template>

<script>
import TimestampDisplay from '@/components/TimestampDisplay.vue'

export default {
  components: { TimestampDisplay }
}
</script>
```

#### Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `timestamp` | String/Date | null | The timestamp to display |
| `format` | String | 'auto' | Format type: 'auto', 'relative', 'datetime', 'long' |
| `showTime` | Boolean | true | Include time with date |
| `showIcon` | Boolean | false | Show clock icon |
| `icon` | String | 'bi bi-clock' | Icon class |
| `customClass` | String | '' | Additional CSS classes |

#### Advanced Usage

```vue
<!-- Auto format (smart relative/absolute) -->
<TimestampDisplay :timestamp="item.createdAt" />

<!-- Always show relative time -->
<TimestampDisplay :timestamp="item.updatedAt" format="relative" />

<!-- Full date and time -->
<TimestampDisplay :timestamp="item.createdAt" format="long" />

<!-- Date only, no time -->
<TimestampDisplay :timestamp="item.createdAt" :showTime="false" />

<!-- With icon -->
<TimestampDisplay
  :timestamp="item.updatedAt"
  :showIcon="true"
  icon="bi bi-clock-history"
/>

<!-- Custom styling -->
<TimestampDisplay
  :timestamp="item.createdAt"
  customClass="text-success fw-bold"
/>
```

## Implementation Examples

### 1. Table Row with Timestamps

```vue
<template>
  <table class="table">
    <thead>
      <tr>
        <th>Name</th>
        <th>Balance</th>
        <th>Created</th>
        <th>Last Updated</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="account in accounts" :key="account.id">
        <td>{{ account.name }}</td>
        <td>{{ formatCurrency(account.balance) }}</td>
        <td>
          <TimestampDisplay :timestamp="account.createdAt" />
        </td>
        <td>
          <TimestampDisplay
            :timestamp="account.updatedAt"
            format="relative"
            :showIcon="true"
          />
        </td>
        <td>
          <button @click="edit(account)">Edit</button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import TimestampDisplay from '@/components/TimestampDisplay.vue'

export default {
  components: { TimestampDisplay }
}
</script>
```

### 2. Detail View with Timestamps

```vue
<template>
  <div class="detail-view">
    <h2>{{ account.name }}</h2>

    <div class="metadata">
      <div class="metadata-row">
        <span class="label">Created:</span>
        <TimestampDisplay
          :timestamp="account.createdAt"
          format="long"
        />
      </div>
      <div class="metadata-row">
        <span class="label">Last Updated:</span>
        <TimestampDisplay
          :timestamp="account.updatedAt"
          format="long"
        />
      </div>
    </div>
  </div>
</template>
```

### 3. Card Footer with Timestamps

```vue
<template>
  <div class="card">
    <div class="card-header">
      <h5>{{ transaction.description }}</h5>
    </div>
    <div class="card-body">
      <p class="amount">{{ formatCurrency(transaction.amount) }}</p>
    </div>
    <div class="card-footer text-muted small">
      <div class="d-flex justify-content-between">
        <span>
          Created <TimestampDisplay :timestamp="transaction.createdAt" format="relative" />
        </span>
        <span v-if="transaction.updatedAt !== transaction.createdAt">
          Updated <TimestampDisplay :timestamp="transaction.updatedAt" format="relative" />
        </span>
      </div>
    </div>
  </div>
</template>
```

### 4. List View with Relative Times

```vue
<template>
  <div class="activity-list">
    <div v-for="activity in activities" :key="activity.id" class="activity-item">
      <div class="activity-header">
        <strong>{{ activity.action }}</strong>
        <TimestampDisplay
          :timestamp="activity.createdAt"
          format="relative"
          customClass="text-muted"
        />
      </div>
      <div class="activity-body">
        {{ activity.description }}
      </div>
    </div>
  </div>
</template>
```

### 5. Using Utility Functions Directly

When you need more control or custom formatting:

```vue
<template>
  <div>
    <!-- Simple usage -->
    <small class="text-muted">
      Last updated: {{ formatTimestamp(item.updatedAt) }}
    </small>

    <!-- Conditional formatting -->
    <div>
      <span v-if="isRecent(item.createdAt)">
        🆕 {{ getRelativeTime(item.createdAt) }}
      </span>
      <span v-else>
        {{ formatDateTime(item.createdAt, true) }}
      </span>
    </div>

    <!-- Tooltip with full timestamp -->
    <span :title="formatDateTimeLong(item.updatedAt)">
      {{ getRelativeTime(item.updatedAt) }}
    </span>
  </div>
</template>

<script>
import {
  formatTimestamp,
  formatDateTime,
  formatDateTimeLong,
  getRelativeTime
} from '@/utils/dateUtils'

export default {
  methods: {
    formatTimestamp,
    formatDateTime,
    formatDateTimeLong,
    getRelativeTime,

    isRecent(date) {
      const now = new Date()
      const dateObj = new Date(date)
      const diffDays = (now - dateObj) / (1000 * 60 * 60 * 24)
      return diffDays < 7
    }
  }
}
</script>
```

## Best Practices

### ✅ DO:

1. **Use the TimestampDisplay component** for consistency
2. **Show relative times** for recent items (< 7 days)
3. **Include full timestamp in tooltips** (title attribute)
4. **Use formatTimestamp()** for auto-smart formatting
5. **Show both Created and Updated** in detail views
6. **Hide Updated if same as Created** (newly created items)

### ❌ DON'T:

1. **Don't format dates manually** - use utilities
2. **Don't show seconds** unless absolutely necessary
3. **Don't use cryptic formats** - keep it human-readable
4. **Don't forget timezone handling** - utilities handle this
5. **Don't show timestamps without labels** - always add context

## Common Patterns

### Pattern 1: Table Column
```vue
<th>Last Updated</th>
<td><TimestampDisplay :timestamp="item.updatedAt" /></td>
```

### Pattern 2: Info Text
```vue
<small class="text-muted">
  Created {{ getRelativeTime(item.createdAt) }}
</small>
```

### Pattern 3: Tooltip
```vue
<span
  :title="formatDateTimeLong(item.updatedAt)"
  class="cursor-help"
>
  {{ getRelativeTime(item.updatedAt) }}
</span>
```

### Pattern 4: Conditional Display
```vue
<div v-if="item.updatedAt !== item.createdAt">
  <small>Edited <TimestampDisplay :timestamp="item.updatedAt" format="relative" /></small>
</div>
```

### Pattern 5: Activity Feed
```vue
<div class="timeline-item">
  <div class="timeline-time">
    <TimestampDisplay :timestamp="activity.createdAt" format="relative" />
  </div>
  <div class="timeline-content">
    {{ activity.description }}
  </div>
</div>
```

## Migration Checklist

When updating existing views to include timestamps:

- [ ] Import `TimestampDisplay` component or utility functions
- [ ] Add `Created` column/field to tables
- [ ] Add `Last Updated` column/field to tables
- [ ] Show full timestamps in detail/edit modals
- [ ] Add timestamps to card footers
- [ ] Include timestamps in export/print views
- [ ] Test with different date ranges (today, week old, year old)
- [ ] Verify tooltip shows full timestamp
- [ ] Check mobile responsive display
- [ ] Update sorting to work with timestamps

## Troubleshooting

### Timestamp shows "Invalid Date"
- Check that the backend is sending valid ISO 8601 format
- Verify the field name matches (createdAt vs created_at)
- Ensure the value is not null/undefined

### Timezone issues
- Backend sends UTC timestamps
- Utilities automatically convert to user's local timezone
- No manual timezone handling needed

### Performance with large lists
- Utilities are optimized for performance
- Use component for consistency and caching
- Virtual scrolling for very large lists (1000+ items)

## Summary

1. **Backend**: GORM handles timestamps automatically ✅
2. **Frontend**: Use `TimestampDisplay` component or utility functions
3. **Smart formatting**: Recent = relative time, Old = full date
4. **Always provide context**: Label timestamps appropriately
5. **Tooltips**: Show full timestamp on hover for clarity

Need help? Check the utility functions in `src/utils/dateUtils.js` and the component in `src/components/TimestampDisplay.vue`.
