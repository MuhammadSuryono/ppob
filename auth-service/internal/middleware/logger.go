package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoggerConfig struct {
	PIIFields []string
}

var defaultLoggerConfig = LoggerConfig{
	PIIFields: []string{"password", "pin", "token", "refresh_token", "secret"},
}

func StructuredLogger() gin.HandlerFunc {
	return StructuredLoggerWithConfig(defaultLoggerConfig)
}

func StructuredLoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		latency := time.Since(startTime)

		logEntry := map[string]interface{}{
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"level":       getLogLevel(c.Writer.Status()),
			"service":     "auth-service",
			"trace_id":    traceID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		if userID, exists := c.Get("user_id"); exists {
			logEntry["user_id"] = userID
		}

		if len(requestBody) > 0 {
			redactedBody := redactSensitiveData(requestBody, config.PIIFields)
			logEntry["request_body"] = string(redactedBody)
		}

		if len(c.Errors) > 0 {
			logEntry["errors"] = c.Errors.String()
		}

		logJSON, _ := json.Marshal(logEntry)
		println(string(logJSON))
	}
}

func getLogLevel(statusCode int) string {
	if statusCode >= 500 {
		return "ERROR"
	}
	if statusCode >= 400 {
		return "WARN"
	}
	return "INFO"
}

func redactSensitiveData(data []byte, piiFields []string) []byte {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return data
	}

	for _, field := range piiFields {
		if _, ok := result[field]; ok {
			result[field] = "***REDACTED***"
		}
	}

	redacted, _ := json.Marshal(result)
	return redacted
}