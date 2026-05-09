package services

type ErrorCode string

const (
	ErrSuccess       ErrorCode = "00"
	ErrPending       ErrorCode = "03"
	ErrInvalidCode   ErrorCode = "01"
	ErrInvalidPhone  ErrorCode = "02"
	ErrInsufficient  ErrorCode = "39"
	ErrProductInactive ErrorCode = "06"
	ErrSystemError   ErrorCode = "69"
	ErrTimeout       ErrorCode = "99"
	ErrUnknown       ErrorCode = "XX"
)

type DigiflazzError struct {
	Code           ErrorCode `json:"code"`
	RC             string    `json:"rc"`
	Message        string    `json:"message"`
	UserMessageID  string    `json:"user_message_id"`
	SuggestedAction string   `json:"suggested_action"`
}

var errorMapping = map[ErrorCode]DigiflazzError{
	ErrSuccess: {
		Code:           ErrSuccess,
		RC:             "00",
		Message:        "Transaction successful",
		UserMessageID:  "success",
		SuggestedAction: "Transaction completed successfully",
	},
	ErrPending: {
		Code:           ErrPending,
		RC:             "03",
		Message:        "Transaction is being processed",
		UserMessageID:  "pending",
		SuggestedAction: "Please wait while we process your transaction",
	},
	ErrInvalidCode: {
		Code:           ErrInvalidCode,
		RC:             "01",
		Message:        "Invalid product code",
		UserMessageID:  "invalid_product",
		SuggestedAction: "Please check the product code and try again",
	},
	ErrInvalidPhone: {
		Code:           ErrInvalidPhone,
		RC:             "02",
		Message:        "Invalid phone number or customer ID",
		UserMessageID:  "invalid_phone",
		SuggestedAction: "Please verify the phone number or customer ID",
	},
	ErrInsufficient: {
		Code:           ErrInsufficient,
		RC:             "39",
		Message:        "Insufficient balance in provider",
		UserMessageID:  "insufficient_balance",
		SuggestedAction: "Please try again later or use a different product",
	},
	ErrProductInactive: {
		Code:           ErrProductInactive,
		RC:             "06",
		Message:        "Product is currently inactive",
		UserMessageID:  "product_inactive",
		SuggestedAction: "Please choose a different product",
	},
	ErrSystemError: {
		Code:           ErrSystemError,
		RC:             "69",
		Message:        "System error occurred",
		UserMessageID:  "system_error",
		SuggestedAction: "Please try again in a few minutes",
	},
	ErrTimeout: {
		Code:           ErrTimeout,
		RC:             "99",
		Message:        "Transaction timeout",
		UserMessageID:  "timeout",
		SuggestedAction: "Please try again or contact support",
	},
}

func GetErrorInfo(rc string) DigiflazzError {
	code := ErrorCode(rc)
	if err, exists := errorMapping[code]; exists {
		return err
	}

	return DigiflazzError{
		Code:           ErrUnknown,
		RC:             rc,
		Message:        "Unknown error occurred",
		UserMessageID:  "unknown_error",
		SuggestedAction: "Please contact customer support",
	}
}

func MapRCToUserMessage(rc string) (userMessage string, suggestedAction string, status string) {
	errInfo := GetErrorInfo(rc)
	return errInfo.UserMessageID, errInfo.SuggestedAction, mapStatusFromRC(ErrorCode(rc))
}

func mapStatusFromRC(code ErrorCode) string {
	switch code {
	case ErrSuccess:
		return "success"
	case ErrPending:
		return "pending"
	default:
		return "failed"
	}
}

type ErrorCatalog struct {
	Errors []DigiflazzError `json:"errors"`
}

func NewErrorCatalog() *ErrorCatalog {
	errors := make([]DigiflazzError, 0, len(errorMapping))
	for _, err := range errorMapping {
		errors = append(errors, err)
	}
	return &ErrorCatalog{Errors: errors}
}

func (e *ErrorCatalog) GetError(rc string) DigiflazzError {
	return GetErrorInfo(rc)
}

func (e *ErrorCatalog) GetAllErrors() []DigiflazzError {
	return e.Errors
}