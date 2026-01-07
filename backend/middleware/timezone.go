package middleware

import (
	"bytes"
	"encoding/json"
	"time"

	"daybook-backend/timezone"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return len(b), nil
}

// TimezoneMiddleware converts all timestamps in JSON responses to UTC+6
func TimezoneMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create custom response writer to capture response body
		blw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Process request
		c.Next()

		// Only convert timestamps for successful JSON responses
		contentType := blw.Header().Get("Content-Type")
		statusOK := blw.Status() >= 200 && blw.Status() < 300
		isJSON := contentType == "application/json" || contentType == "application/json; charset=utf-8"

		if statusOK && isJSON && blw.body.Len() > 0 {
			// Parse response body
			var data interface{}
			if err := json.Unmarshal(blw.body.Bytes(), &data); err == nil {
				// Convert timestamps recursively
				convertedData := convertTimestamps(data)

				// Marshal back to JSON
				if convertedBytes, err := json.Marshal(convertedData); err == nil {
					// Write converted response
					blw.ResponseWriter.Header().Set("Content-Length", string(rune(len(convertedBytes))))
					blw.ResponseWriter.Write(convertedBytes)
					return
				}
			}
		}

		// If conversion failed or not applicable, write original response
		blw.ResponseWriter.Write(blw.body.Bytes())
	}
}

// convertTimestamps recursively converts all RFC3339 timestamp strings to UTC+6
func convertTimestamps(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertTimestamps(value)
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = convertTimestamps(value)
		}
		return result

	case string:
		// Try to parse as RFC3339 timestamp
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return timezone.ToAppTimezone(t).Format(time.RFC3339)
		}
		// Try to parse as RFC3339Nano timestamp
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return timezone.ToAppTimezone(t).Format(time.RFC3339Nano)
		}
		return v

	default:
		return v
	}
}
