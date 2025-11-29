package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name        string         `gorm:"not null" json:"name" binding:"required"`
	Type        string         `gorm:"not null;index" json:"type" binding:"required"` // income, expense
	Icon        string         `gorm:"not null" json:"icon" binding:"required"`       // Icon identifier
	Color       string         `gorm:"default:#3B82F6" json:"color"`                  // Hex color code
	Description string         `json:"description"`
	IsDefault   bool           `gorm:"default:false" json:"isDefault"` // System default categories
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// DefaultCategories returns a list of default categories for new users
func GetDefaultCategories(userID uuid.UUID) []Category {
	return []Category{
		// Income categories - using varied green/blue tones
		{UserID: userID, Name: "Salary", Type: "income", Icon: "💰", Color: "#10B981", IsDefault: true},                // Green
		{UserID: userID, Name: "Freelance", Type: "income", Icon: "💼", Color: "#14B8A6", IsDefault: true},             // Teal
		{UserID: userID, Name: "Business", Type: "income", Icon: "🏢", Color: "#06B6D4", IsDefault: true},              // Cyan
		{UserID: userID, Name: "Gift", Type: "income", Icon: "🎁", Color: "#0EA5E9", IsDefault: true},                  // Light Blue
		{UserID: userID, Name: "Bonus", Type: "income", Icon: "🎉", Color: "#22C55E", IsDefault: true},                 // Bright Green
		{UserID: userID, Name: "Refund", Type: "income", Icon: "🔄", Color: "#3B82F6", IsDefault: true},                // Blue
		{UserID: userID, Name: "Cashback", Type: "income", Icon: "💳", Color: "#8B5CF6", IsDefault: true},              // Purple
		{UserID: userID, Name: "Insurance Settlement", Type: "income", Icon: "🛡️", Color: "#059669", IsDefault: true}, // Dark Green
		{UserID: userID, Name: "Other Income", Type: "income", Icon: "💵", Color: "#10B981", IsDefault: true},          // Green

		// Expense categories - using varied warm tones (red, orange, amber, pink)
		{UserID: userID, Name: "Food & Dining", Type: "expense", Icon: "🍔", Color: "#F59E0B", IsDefault: true},     // Amber
		{UserID: userID, Name: "Transportation", Type: "expense", Icon: "🚗", Color: "#3B82F6", IsDefault: true},    // Blue
		{UserID: userID, Name: "Groceries", Type: "expense", Icon: "🛒", Color: "#84CC16", IsDefault: true},         // Lime
		{UserID: userID, Name: "Shopping", Type: "expense", Icon: "🛍️", Color: "#EC4899", IsDefault: true},         // Pink
		{UserID: userID, Name: "Entertainment", Type: "expense", Icon: "🎬", Color: "#8B5CF6", IsDefault: true},     // Purple
		{UserID: userID, Name: "Healthcare", Type: "expense", Icon: "🏥", Color: "#EF4444", IsDefault: true},        // Red
		{UserID: userID, Name: "Bills & Utilities", Type: "expense", Icon: "💡", Color: "#F97316", IsDefault: true}, // Orange
		{UserID: userID, Name: "Rent", Type: "expense", Icon: "🏠", Color: "#DC2626", IsDefault: true},              // Dark Red
		{UserID: userID, Name: "Insurance", Type: "expense", Icon: "🛡️", Color: "#7C3AED", IsDefault: true},        // Violet
		{UserID: userID, Name: "Education", Type: "expense", Icon: "📚", Color: "#0EA5E9", IsDefault: true},         // Sky Blue
		{UserID: userID, Name: "Travel", Type: "expense", Icon: "✈️", Color: "#14B8A6", IsDefault: true},           // Teal
		{UserID: userID, Name: "Fitness", Type: "expense", Icon: "💪", Color: "#10B981", IsDefault: true},           // Green
		{UserID: userID, Name: "Personal Care", Type: "expense", Icon: "💅", Color: "#F472B6", IsDefault: true},     // Light Pink
		{UserID: userID, Name: "Subscriptions", Type: "expense", Icon: "📱", Color: "#6366F1", IsDefault: true},     // Indigo
		{UserID: userID, Name: "Gifts & Donations", Type: "expense", Icon: "🎀", Color: "#A855F7", IsDefault: true}, // Purple
		{UserID: userID, Name: "Other Expense", Type: "expense", Icon: "💸", Color: "#6B7280", IsDefault: true},     // Gray
	}
}

// AvailableIcons returns a list of available icons for categories
func GetAvailableIcons() map[string][]string {
	return map[string][]string{
		"income": {
			"💰", "💵", "💴", "💶", "💷", "💸", "💳", "💎",
			"📈", "📊", "💼", "🏢", "🏦", "🎁", "🎉", "⭐",
			"✨", "💫", "🌟", "🔥", "↩️", "✅", "💯", "🎯",
		},
		"expense": {
			"🍔", "🍕", "🍜", "🍱", "🛒", "🛍️", "🎬", "🎮",
			"🚗", "🚕", "🚌", "✈️", "🏠", "🏡", "💡", "💧",
			"📱", "💻", "🏥", "💊", "💪", "📚", "🎓", "🎵",
			"🎸", "🎨", "👕", "👗", "👠", "💅", "💇", "🛡️",
			"📺", "🎧", "⚽", "🏀", "🎾", "🎭", "🍿", "☕",
			"🍺", "🍷", "🎂", "🌮", "🍣", "🍦", "💸", "💳",
		},
	}
}
