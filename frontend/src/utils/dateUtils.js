/**
 * Date and DateTime utility functions for consistent date/time handling across the application
 */

/**
 * Converts a date string or Date object to ISO 8601 format for backend API
 * Handles date-only strings (YYYY-MM-DD) and converts them to full datetime
 *
 * @param {string|Date} date - Date string or Date object
 * @returns {string} ISO 8601 formatted datetime string
 */
export function toISOString(date) {
  if (!date) return null

  // If it's already a Date object, just convert
  if (date instanceof Date) {
    return date.toISOString()
  }

  // If it's a string
  if (typeof date === 'string') {
    // If it's already in ISO format (contains 'T'), return as-is
    if (date.includes('T')) {
      return new Date(date).toISOString()
    }

    // If it's a date-only string (YYYY-MM-DD), convert to midnight UTC
    // This prevents timezone issues when sending to backend
    const dateObj = new Date(date + 'T00:00:00.000Z')
    return dateObj.toISOString()
  }

  // Fallback: try to create a Date object
  return new Date(date).toISOString()
}

/**
 * Converts a Date object to YYYY-MM-DD format for date inputs
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} Date string in YYYY-MM-DD format
 */
export function toDateInputValue(date) {
  if (!date) return ''

  const dateObj = date instanceof Date ? date : new Date(date)
  return dateObj.toISOString().split('T')[0]
}

/**
 * Formats object dates for API submission
 * Converts all date fields to ISO format
 *
 * @param {Object} data - Object with date fields
 * @param {Array<string>} dateFields - Array of field names that contain dates
 * @returns {Object} Object with formatted dates
 */
export function formatDatesForAPI(data, dateFields = []) {
  const formatted = { ...data }

  dateFields.forEach(field => {
    if (formatted[field]) {
      formatted[field] = toISOString(formatted[field])
    }
  })

  return formatted
}

/**
 * Formats a date/time to a readable format
 * Examples: "Dec 28, 2025", "Dec 28, 2025 at 3:45 PM"
 *
 * @param {Date|string} date - Date object or ISO string
 * @param {boolean} includeTime - Whether to include time (default: false)
 * @returns {string} Formatted date string
 */
export function formatDateTime(date, includeTime = false) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    // Check if date is valid
    if (isNaN(dateObj.getTime())) return 'Invalid Date'

    const dateOptions = {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }

    const dateStr = dateObj.toLocaleDateString('en-US', dateOptions)

    if (includeTime) {
      const timeStr = dateObj.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })
      return `${dateStr} at ${timeStr}`
    }

    return dateStr
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Invalid Date'
  }
}

/**
 * Formats a date to short format
 * Example: "12/28/2025"
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} Formatted date string
 */
export function formatDateShort(date) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    if (isNaN(dateObj.getTime())) return 'Invalid Date'

    return dateObj.toLocaleDateString('en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Invalid Date'
  }
}

/**
 * Formats a date to long format with full date and time
 * Example: "December 28, 2025 at 3:45:30 PM"
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} Formatted datetime string
 */
export function formatDateTimeLong(date) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    if (isNaN(dateObj.getTime())) return 'Invalid Date'

    const dateStr = dateObj.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })

    const timeStr = dateObj.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      second: '2-digit',
      hour12: true
    })

    return `${dateStr} at ${timeStr}`
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Invalid Date'
  }
}

/**
 * Formats time only
 * Example: "3:45 PM"
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} Formatted time string
 */
export function formatTimeOnly(date) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    if (isNaN(dateObj.getTime())) return 'Invalid Time'

    return dateObj.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    })
  } catch (error) {
    console.error('Error formatting time:', error)
    return 'Invalid Time'
  }
}

/**
 * Gets relative time (e.g., "2 hours ago", "3 days ago")
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {string} Relative time string
 */
export function getRelativeTime(date) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    if (isNaN(dateObj.getTime())) return 'Invalid Date'

    const now = new Date()
    const diffMs = now - dateObj
    const diffSeconds = Math.floor(diffMs / 1000)
    const diffMinutes = Math.floor(diffSeconds / 60)
    const diffHours = Math.floor(diffMinutes / 60)
    const diffDays = Math.floor(diffHours / 24)
    const diffWeeks = Math.floor(diffDays / 7)
    const diffMonths = Math.floor(diffDays / 30)
    const diffYears = Math.floor(diffDays / 365)

    if (diffSeconds < 60) return 'Just now'
    if (diffMinutes < 60) return `${diffMinutes} minute${diffMinutes !== 1 ? 's' : ''} ago`
    if (diffHours < 24) return `${diffHours} hour${diffHours !== 1 ? 's' : ''} ago`
    if (diffDays < 7) return `${diffDays} day${diffDays !== 1 ? 's' : ''} ago`
    if (diffWeeks < 4) return `${diffWeeks} week${diffWeeks !== 1 ? 's' : ''} ago`
    if (diffMonths < 12) return `${diffMonths} month${diffMonths !== 1 ? 's' : ''} ago`
    return `${diffYears} year${diffYears !== 1 ? 's' : ''} ago`
  } catch (error) {
    console.error('Error calculating relative time:', error)
    return 'Invalid Date'
  }
}

/**
 * Formats created/updated timestamps for display
 * Shows relative time if recent, otherwise shows full date
 *
 * @param {Date|string} date - Date object or ISO string
 * @param {boolean} showTime - Whether to show time (default: true)
 * @returns {string} Formatted timestamp
 */
export function formatTimestamp(date, showTime = true) {
  if (!date) return 'N/A'

  try {
    const dateObj = date instanceof Date ? date : new Date(date)

    if (isNaN(dateObj.getTime())) return 'Invalid Date'

    const now = new Date()
    const diffDays = Math.floor((now - dateObj) / (1000 * 60 * 60 * 24))

    // If less than 7 days ago, show relative time
    if (diffDays < 7) {
      return getRelativeTime(dateObj)
    }

    // Otherwise show formatted date/time
    return formatDateTime(dateObj, showTime)
  } catch (error) {
    console.error('Error formatting timestamp:', error)
    return 'Invalid Date'
  }
}

/**
 * Check if a date is today
 *
 * @param {Date|string} date - Date object or ISO string
 * @returns {boolean} True if date is today
 */
export function isToday(date) {
  if (!date) return false

  try {
    const dateObj = date instanceof Date ? date : new Date(date)
    const today = new Date()

    return dateObj.getDate() === today.getDate() &&
           dateObj.getMonth() === today.getMonth() &&
           dateObj.getFullYear() === today.getFullYear()
  } catch (error) {
    return false
  }
}
