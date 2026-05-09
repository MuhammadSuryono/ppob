package errors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	TraceID    string                 `json:"trace_id"`
	Timestamp  string                 `json:"timestamp"`
	StatusCode int                    `json:"-"`
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
		StatusCode: httpStatusForCode(code),
	}
}

func httpStatusForCode(code string) int {
	switch {
	case strings.HasPrefix(code, "AUTH_"):
		if code == "AUTH_ACCOUNT_LOCKED" || code == "AUTH_DEVICE_NOT_TRUSTED" || code == "AUTH_INSUFFICIENT_PERMISSION" {
			return http.StatusForbidden
		}
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "VALIDATION_"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "SYSTEM_"):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func RespondWithError(c *gin.Context, err *AppError) {
	if err.TraceID == "" {
		err.TraceID = uuid.New().String()
	}
	if err.Timestamp == "" {
		err.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	c.JSON(err.StatusCode, gin.H{
		"error": gin.H{
			"code":      err.Code,
			"message":   err.Message,
			"details":   err.Details,
			"trace_id":  err.TraceID,
			"timestamp": err.Timestamp,
		},
	})
}

func BadRequest(c *gin.Context, code, message string) {
	RespondWithError(c, NewError(code, message, nil))
}

func Unauthorized(c *gin.Context, code, message string) {
	RespondWithError(c, NewError(code, message, nil))
}

func Forbidden(c *gin.Context, code, message string) {
	RespondWithError(c, NewError(code, message, nil))
}

func NotFound(c *gin.Context, code, message string) {
	RespondWithError(c, NewError(code, message, nil))
}

func InternalError(c *gin.Context, code, message string) {
	RespondWithError(c, NewError(code, message, nil))
}

var ErrorMessages = map[string]string{
	"AUTH_INVALID_CREDENTIALS":   "Nomor HP atau PIN/Password salah",
	"AUTH_TOKEN_EXPIRED":          "Sesi berakhir, silakan login ulang",
	"AUTH_TOKEN_REVOKED":          "Token tidak valid",
	"AUTH_OTP_EXPIRED":            "Kode OTP sudah kadaluarsa",
	"AUTH_OTP_INVALID":           "Kode OTP tidak valid",
	"AUTH_OTP_RATE_LIMIT":        "Terlalu banyak percobaan OTP, coba lagi dalam 1 menit",
	"AUTH_ACCOUNT_LOCKED":        "Akun diblokir karena terlalu banyak percobaan PIN",
	"AUTH_DEVICE_NOT_TRUSTED":    "Perangkat tidak terpercaya, silakan Otentikasi ulang",
	"AUTH_INSUFFICIENT_PERMISSION": "Anda tidak memiliki izin untuk mengakses sumber daya ini",
	"VALIDATION_PHONE_FORMAT":    "Format nomor HP tidak valid (contoh: +6281234567890)",
	"VALIDATION_PIN_FORMAT":      "PIN harus 6 digit angka",
	"VALIDATION_PIN_SEQUENTIAL":  "PIN tidak boleh berurutan (123456, 654321)",
	"VALIDATION_PASSWORD_WEAK":    "Password minimal 8 karakter, mengandung huruf besar & angka",
	"VALIDATION_MISSING_FIELD":    "Field '%s' wajib diisi",
	"VALIDATION_JSON_INVALID":     "Format JSON tidak valid",
	"SYSTEM_INTERNAL":           "Terjadi kesalahan internal",
	"SYSTEM_TIMEOUT":             "Request timeout",
	"SYSTEM_DB_UNAVAILABLE":       "Database tidak tersedia, coba beberapa saat lagi",
	"SYSTEM_REDIS_UNAVAILABLE":    "Cache tidak tersedia",
}

func GetErrorMessage(code string) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "Terjadi kesalahan"
}