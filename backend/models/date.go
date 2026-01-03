package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type that handles date-only strings in JSON
type Date struct {
	time.Time
}

// UnmarshalJSON handles unmarshaling of date strings in multiple formats
func (d *Date) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}

	// Remove quotes from JSON string
	str := string(b)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Trim whitespace
	str = strings.TrimSpace(str)

	// Handle empty strings after quote removal
	if str == "" {
		return nil
	}

	// Try parsing as date-only format first (2006-01-02)
	t, err := time.Parse("2006-01-02", str)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as RFC3339 format (2006-01-02T15:04:05Z07:00)
	t, err = time.Parse(time.RFC3339, str)
	if err == nil {
		d.Time = t
		return nil
	}

	// Try parsing as RFC3339Nano format
	t, err = time.Parse(time.RFC3339Nano, str)
	if err == nil {
		d.Time = t
		return nil
	}

	return fmt.Errorf("unable to parse date: %s", str)
}

// MarshalJSON marshals the date as a date-only string
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

// Value implements the driver.Valuer interface for database storage
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
}

// GormDataType tells GORM to use DATE type in the database
func (Date) GormDataType() string {
	return "timestamp"
}
