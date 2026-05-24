package errors

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	TraceID   string                 `json:"trace_id"`
	Timestamp string                 `json:"timestamp"`
}

func RespondWithError(c *gin.Context, appErr *AppError) {
	if appErr.TraceID == "" {
		appErr.TraceID = GenerateTraceID()
	}
	if appErr.Timestamp == "" {
		appErr.Timestamp = CurrentTimestamp()
	}

	log.Printf("ERROR: code=%s trace_id=%s message=%s", appErr.Code, appErr.TraceID, appErr.Message)

	if appErr.RetryAfter > 0 {
		c.Header("Retry-After", string(rune(appErr.RetryAfter)))
	}

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Details:   appErr.Details,
			TraceID:   appErr.TraceID,
			Timestamp: appErr.Timestamp,
		},
	}

	c.JSON(appErr.StatusCode, response)
}

func RespondWithErrorJSON(w http.ResponseWriter, appErr *AppError) {
	if appErr.TraceID == "" {
		appErr.TraceID = GenerateTraceID()
	}
	if appErr.Timestamp == "" {
		appErr.Timestamp = CurrentTimestamp()
	}

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Details:   appErr.Details,
			TraceID:   appErr.TraceID,
			Timestamp: appErr.Timestamp,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if appErr.RetryAfter > 0 {
		w.Header().Set("Retry-After", string(rune(appErr.RetryAfter)))
	}
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func BadRequest(c *gin.Context, code, message string) {
	appErr := NewAppError(code, map[string]interface{}{"message": message})
	if message != "" {
		appErr.Message = message
	}
	RespondWithError(c, appErr)
}

func Unauthorized(c *gin.Context, code, message string) {
	appErr := NewAppError(code, nil).WithStatus(http.StatusUnauthorized)
	if message != "" {
		appErr.Message = message
	}
	RespondWithError(c, appErr)
}

func Forbidden(c *gin.Context, code, message string) {
	appErr := NewAppError(code, nil).WithStatus(http.StatusForbidden)
	if message != "" {
		appErr.Message = message
	}
	RespondWithError(c, appErr)
}

func NotFound(c *gin.Context, code, message string) {
	appErr := NewAppError(code, nil).WithStatus(http.StatusNotFound)
	if message != "" {
		appErr.Message = message
	}
	RespondWithError(c, appErr)
}

func InternalError(c *gin.Context, code, message string) {
	appErr := NewAppError(code, nil).WithStatus(http.StatusInternalServerError)
	if message != "" {
		appErr.Message = message
	}
	RespondWithError(c, appErr)
}

func ServiceUnavailable(c *gin.Context, code, message string) {
	RespondWithError(c, NewAppError(code, nil).WithStatus(http.StatusServiceUnavailable))
}