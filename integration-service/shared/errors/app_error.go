package errors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AppError struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	TraceID     string                 `json:"trace_id"`
	StatusCode  int                    `json:"-"`
	Retryable   bool                   `json:"-"`
	RetryAfter  int                    `json:"-"`
	Timestamp   string                 `json:"timestamp"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewError(code, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		TraceID:    uuid.New().String(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		StatusCode: HTTPStatusForCode(code),
		Retryable:  IsRetryable(code),
		RetryAfter: RetryAfterForCode(code),
	}
}

func (e *AppError) WithStatus(status int) *AppError {
	e.StatusCode = status
	return e
}

func (e *AppError) WithRetryable(retryable bool) *AppError {
	e.Retryable = retryable
	return e
}

func (e *AppError) WithRetryAfter(seconds int) *AppError {
	e.RetryAfter = seconds
	return e
}

func HTTPStatusForCode(code string) int {
	switch {
	case strings.HasPrefix(code, "AUTH_"):
		if code == "AUTH_ACCOUNT_LOCKED" || code == "AUTH_DEVICE_NOT_TRUSTED" || code == "AUTH_INSUFFICIENT_PERMISSION" {
			return http.StatusForbidden
		}
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "VALIDATION_"):
		return http.StatusBadRequest
	case code == "TRANSACTION_NOT_FOUND":
		return http.StatusNotFound
	case code == "TRANSACTION_EXPIRED":
		return http.StatusGone
	case code == "TRANSACTION_DAILY_LIMIT_EXCEEDED":
		return http.StatusForbidden
	case strings.HasPrefix(code, "TRANSACTION_"), strings.HasPrefix(code, "DIGIFLAZZ_"):
		if code == "DIGIFLAZZ_PENDING" {
			return http.StatusAccepted
		}
		if code == "DIGIFLAZZ_RATE_LIMIT" || code == "DIGIFLAZZ_PRICELIST_LIMIT" || code == "DIGIFLAZZ_PLN_INQUIRY_LIMIT" {
			return http.StatusTooManyRequests
		}
		if code == "DIGIFLAZZ_TIMEOUT" || code == "DIGIFLAZZ_BILLER_TIMEOUT" || code == "DIGIFLAZZ_PRODUCT_UNSTABLE" || strings.Contains(code, "UNAVAILABLE") || code == "DIGIFLAZZ_INSUFFICIENT_DEPOSIT" {
			return http.StatusBadGateway
		}
		return http.StatusBadRequest
	case strings.HasPrefix(code, "SYSTEM_"):
		if code == "SYSTEM_RATE_LIMIT" {
			return http.StatusTooManyRequests
		}
		if code == "SYSTEM_TIMEOUT" {
			return http.StatusGatewayTimeout
		}
		if code == "SYSTEM_DB_UNAVAILABLE" || code == "SYSTEM_REDIS_UNAVAILABLE" || code == "SYSTEM_CIRCUIT_OPEN" || code == "SYSTEM_MAINTENANCE" {
			return http.StatusServiceUnavailable
		}
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func IsRetryable(code string) bool {
	retryableCodes := map[string]bool{
		"DIGIFLAZZ_TIMEOUT":             true,
		"DIGIFLAZZ_BILLER_TIMEOUT":     true,
		"DIGIFLAZZ_PRODUCT_UNSTABLE":    true,
		"DIGIFLAZZ_RATE_LIMIT":          true,
		"DIGIFLAZZ_PRICELIST_LIMIT":    true,
		"DIGIFLAZZ_PLN_INQUIRY_LIMIT":   true,
		"DIGIFLAZZ_CUT_OFF":             true,
		"DIGIFLAZZ_PRICE_MISMATCH":     true,
		"SYSTEM_DB_UNAVAILABLE":         true,
		"SYSTEM_REDIS_UNAVAILABLE":      true,
		"SYSTEM_VAULT_UNAVAILABLE":     true,
		"SYSTEM_TIMEOUT":                true,
		"SYSTEM_RATE_LIMIT":            true,
		"SYSTEM_CIRCUIT_OPEN":           true,
		"TRANSACTION_HOLD_FAILED":       true,
	}
	return retryableCodes[code]
}

func RetryAfterForCode(code string) int {
	retryAfterMap := map[string]int{
		"DIGIFLAZZ_TIMEOUT":           2,
		"DIGIFLAZZ_BILLER_TIMEOUT":   2,
		"DIGIFLAZZ_PRODUCT_UNSTABLE": 5,
		"DIGIFLAZZ_RATE_LIMIT":       60,
		"DIGIFLAZZ_PRICELIST_LIMIT":  240,
		"DIGIFLAZZ_PLN_INQUIRY_LIMIT": 60,
		"DIGIFLAZZ_CUT_OFF":          900,
		"DIGIFLAZZ_PRICE_MISMATCH":   10,
		"SYSTEM_RATE_LIMIT":          30,
		"SYSTEM_CIRCUIT_OPEN":        10,
	}
	if val, ok := retryAfterMap[code]; ok {
		return val
	}
	return 0
}