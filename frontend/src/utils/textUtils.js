/**
 * Formats slug-like text to human-readable format
 * Example: "fully_received" -> "Fully Received"
 *          "partially_paid" -> "Partially Paid"
 * @param {string} text - The text to format
 * @returns {string} - The formatted text
 */
export const formatStatus = (text) => {
  if (!text) return ''

  return text
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}
