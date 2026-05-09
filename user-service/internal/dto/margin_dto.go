package dto

import (
	"strconv"
)

type SetMarginRequest struct {
	StaffID   uint    `json:"staff_id" binding:"required"`
	SchemeType string  `json:"scheme_type" binding:"required,oneof=MarginShare FixedAllowance"`
	Value     float64 `json:"value" binding:"required,gt=0"`
}

type SetProductMarginRequest struct {
	StaffID     uint    `json:"staff_id" binding:"required"`
	ProductID   uint    `json:"product_id" binding:"required"`
	ProductCode string  `json:"product_code" binding:"required"`
	SchemeType  string  `json:"scheme_type" binding:"required,oneof=MarginShare FixedAllowance"`
	Value       float64 `json:"value" binding:"required,gt=0"`
}

type SwitchRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type SwitchRoleResponse struct {
	Message   string `json:"message"`
	RoleID    uint   `json:"role_id"`
	RoleName  string `json:"role_name"`
	WalletID  uint   `json:"wallet_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	ErrorCode string `json:"error_code"`
}

type ValidationErrorResponse struct {
	Error    string            `json:"error"`
	Message  string            `json:"message"`
	Errors   []ValidationError `json:"errors"`
}

func ParseUint(s string, target *uint) (bool, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return false, err
	}
	*target = uint(v)
	return true, nil
}

var ValidationErrorMessages = map[string]string{
	"INVALID_PHONE":         "Nomor HP tidak valid",
	"WEAK_PASSWORD":         "Password minimal 8 karakter, huruf besar & angka",
	"INVALID_PIN":           "PIN tidak valid",
	"INVALID_CUSTOMER_NO":   "Nomor tujuan salah",
	"PRICE_BELOW_COST":       "Harga jual minimal Rp %s",
	"INVALID_OTP":           "OTP tidak valid",
	"DAILY_TXN_LIMIT":       "Limit harian transaksi tercapai",
	"DAILY_AMOUNT_LIMIT":     "Limit harian nilai transaksi tercapai",
	"INSUFFICIENT_BALANCE":  "Saldo tidak mencukupi",
	"PRODUCT_INACTIVE":      "Produk tidak aktif",
	"INQUIRY_EXPIRED":       "Inquiry sudah expired, silakan cek ulang",
	"DUPLICATE_REF_ID":       "ID transaksi sudah digunakan",
}

func GetValidationMessage(errorCode string, params ...interface{}) string {
	template, exists := ValidationErrorMessages[errorCode]
	if !exists {
		return "Terjadi kesalahan validasi"
	}
	if len(params) > 0 {
		return sprintf(template, params...)
	}
	return template
}

func sprintf(format string, a ...interface{}) string {
	result := format
	for _, v := range a {
		switch val := v.(type) {
		case string:
			result = replaceFirst(result, val)
		case float64:
			result = replaceFirst(result, formatNumber(val))
		case int:
			result = replaceFirst(result, formatNumber(float64(val)))
		}
	}
	return result
}

func replaceFirst(s string, replacement string) string {
	if len(s) > 0 && s[0] == '%' {
		return replacement + s[1:]
	}
	return replacement
}

func formatNumber(n float64) string {
	if n == float64(int64(n)) {
		return formatWhole(n)
	}
	return formatDecimal(n)
}

func formatWhole(n float64) string {
	intPart := int64(n)
	result := ""
	for intPart > 0 {
		result = string(rune('0'+intPart%10)) + result
		intPart /= 10
	}
	if result == "" {
		return "0"
	}
	return result
}

func formatDecimal(n float64) string {
	return sprintf("%.0f", n)
}