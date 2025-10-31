/**
 * Date utility functions for consistent date handling across the application
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
