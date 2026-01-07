package timezone

import (
	"time"
)

// AppTimezone is UTC+6 (Bangladesh Standard Time)
var AppTimezone = time.FixedZone("UTC+6", 6*60*60)

// ToAppTimezone converts a UTC time to app timezone (UTC+6)
func ToAppTimezone(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(AppTimezone)
}

// ToUTC converts a time in app timezone to UTC
func ToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

// NowInAppTimezone returns current time in app timezone
func NowInAppTimezone() time.Time {
	return time.Now().In(AppTimezone)
}

// ParseInAppTimezone parses a date string and interprets it as being in app timezone
func ParseInAppTimezone(layout, value string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, AppTimezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil // Store as UTC in database
}
