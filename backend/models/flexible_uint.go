package models

import (
	"encoding/json"
	"strconv"
)

// FlexibleUint is a custom type that can unmarshal from both string and number in JSON
type FlexibleUint uint

// UnmarshalJSON handles unmarshaling from both string and number
func (f *FlexibleUint) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as uint first
	var num uint
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexibleUint(num)
		return nil
	}

	// Try to unmarshal as string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Empty string means 0
	if str == "" {
		*f = 0
		return nil
	}

	// Parse string as uint
	num64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		*f = 0 // Default to 0 if parsing fails
		return nil
	}

	*f = FlexibleUint(num64)
	return nil
}

// MarshalJSON marshals as a number
func (f FlexibleUint) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint(f))
}

// ToUint converts FlexibleUint to uint
func (f FlexibleUint) ToUint() uint {
	return uint(f)
}

// NewFlexibleUint creates a FlexibleUint from a uint
func NewFlexibleUint(v uint) FlexibleUint {
	return FlexibleUint(v)
}
