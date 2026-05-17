package errors

import (
	"time"

	"github.com/google/uuid"
)

type ErrorCode struct {
	Code        string
	Message     string
	HTTPStatus  int
	Retryable   bool
	Description string
}

var ErrorCatalog = map[string]ErrorCode{
	// AUTH Errors
	"AUTH_INVALID_CREDENTIALS":    {Code: "AUTH_INVALID_CREDENTIALS", Message: "Nomor HP atau PIN/Password salah", HTTPStatus: 401, Retryable: false},
	"AUTH_TOKEN_EXPIRED":           {Code: "AUTH_TOKEN_EXPIRED", Message: "Sesi berakhir, silakan login ulang", HTTPStatus: 401, Retryable: false},
	"AUTH_TOKEN_REVOKED":           {Code: "AUTH_TOKEN_REVOKED", Message: "Token tidak valid", HTTPStatus: 401, Retryable: false},
	"AUTH_OTP_EXPIRED":            {Code: "AUTH_OTP_EXPIRED", Message: "Kode OTP sudah kadaluarsa", HTTPStatus: 400, Retryable: false},
	"AUTH_OTP_INVALID":           {Code: "AUTH_OTP_INVALID", Message: "Kode OTP tidak valid", HTTPStatus: 400, Retryable: false},
	"AUTH_OTP_RATE_LIMIT":        {Code: "AUTH_OTP_RATE_LIMIT", Message: "Terlalu banyak percobaan OTP, coba lagi dalam 1 menit", HTTPStatus: 429, Retryable: true},
	"AUTH_ACCOUNT_LOCKED":        {Code: "AUTH_ACCOUNT_LOCKED", Message: "Akun diblokir karena terlalu banyak percobaan PIN", HTTPStatus: 403, Retryable: false},
	"AUTH_DEVICE_NOT_TRUSTED":    {Code: "AUTH_DEVICE_NOT_TRUSTED", Message: "Perangkat tidak terpercaya, silakan Otentikasi ulang", HTTPStatus: 403, Retryable: false},
	"AUTH_INSUFFICIENT_PERMISSION": {Code: "AUTH_INSUFFICIENT_PERMISSION", Message: "Anda tidak memiliki izin untuk mengakses sumber daya ini", HTTPStatus: 403, Retryable: false},
	"AUTH_AUTHORIZE_INVALID":      {Code: "AUTH_AUTHORIZE_INVALID", Message: "Otorisasi transaksi tidak valid atau kedaluwarsa", HTTPStatus: 401, Retryable: false},

	// SYSTEM Errors

	"TRANSACTION_INSUFFICIENT_BALANCE":   {Code: "TRANSACTION_INSUFFICIENT_BALANCE", Message: "Saldo tidak mencukupi", HTTPStatus: 400, Retryable: false},
	"TRANSACTION_DAILY_LIMIT_EXCEEDED":     {Code: "TRANSACTION_DAILY_LIMIT_EXCEEDED", Message: "Limit harian transaksi tercapai", HTTPStatus: 403, Retryable: false},
	"TRANSACTION_PRODUCT_INACTIVE":       {Code: "TRANSACTION_PRODUCT_INACTIVE", Message: "Produk tidak aktif", HTTPStatus: 400, Retryable: false},
	"TRANSACTION_PRICE_BELOW_COST":        {Code: "TRANSACTION_PRICE_BELOW_COST", Message: "Harga jual minimal Rp %s", HTTPStatus: 400, Retryable: false},
	"TRANSACTION_HOLD_FAILED":             {Code: "TRANSACTION_HOLD_FAILED", Message: "Gagal mengunci saldo", HTTPStatus: 500, Retryable: true},
	"TRANSACTION_ALREADY_PROCESSING":     {Code: "TRANSACTION_ALREADY_PROCESSING", Message: "Transaksi sedang diproses", HTTPStatus: 409, Retryable: false},
	"TRANSACTION_CANCELLED":              {Code: "TRANSACTION_CANCELLED", Message: "Transaksi dibatalkan", HTTPStatus: 400, Retryable: false},
	"TRANSACTION_NOT_FOUND":              {Code: "TRANSACTION_NOT_FOUND", Message: "Transaksi tidak ditemukan", HTTPStatus: 404, Retryable: false},
	"TRANSACTION_EXPIRED":                {Code: "TRANSACTION_EXPIRED", Message: "Transaksi sudah kadaluarsa", HTTPStatus: 410, Retryable: false},
	"TRANSACTION_INQUIRY_EXPIRED":        {Code: "TRANSACTION_INQUIRY_EXPIRED", Message: "Pengecekan tagihan kadaluarsa, silakan cek ulang", HTTPStatus: 400, Retryable: false},

	// DIGIFLAZZ Errors
	"DIGIFLAZZ_SUCCESS":             {Code: "DIGIFLAZZ_SUCCESS", Message: "Transaksi berhasil", HTTPStatus: 200, Retryable: false},
	"DIGIFLAZZ_TIMEOUT":            {Code: "DIGIFLAZZ_TIMEOUT", Message: "Timeout dari provider, mencoba ulang...", HTTPStatus: 502, Retryable: true},
	"DIGIFLAZZ_PENDING":            {Code: "DIGIFLAZZ_PENDING", Message: "Transaksi sedang diproses", HTTPStatus: 202, Retryable: false},
	"DIGIFLAZZ_GENERAL_FAILURE":    {Code: "DIGIFLAZZ_GENERAL_FAILURE", Message: "Transaksi gagal", HTTPStatus: 502, Retryable: false},
	"DIGIFLAZZ_PAYLOAD_ERROR":      {Code: "DIGIFLAZZ_PAYLOAD_ERROR", Message: "Format request tidak valid", HTTPStatus: 400, Retryable: false},
	"DIGIFLAZZ_INVALID_SIGNATURE": {Code: "DIGIFLAZZ_INVALID_SIGNATURE", Message: "Signature tidak valid", HTTPStatus: 401, Retryable: false},
	"DIGIFLAZZ_SELLER_NOT_FOUND":   {Code: "DIGIFLAZZ_SELLER_NOT_FOUND", Message: "Gagal memproses API seller", HTTPStatus: 502, Retryable: false},
	"DIGIFLAZZ_SKU_NOT_FOUND":     {Code: "DIGIFLAZZ_SKU_NOT_FOUND", Message: "Produk tidak ditemukan", HTTPStatus: 404, Retryable: false},
	"DIGIFLAZZ_INSUFFICIENT_DEPOSIT": {Code: "DIGIFLAZZ_INSUFFICIENT_DEPOSIT", Message: "Saldo platform tidak cukup — hubungi admin", HTTPStatus: 502, Retryable: false},
	"DIGIFLAZZ_IP_BLOCKED":         {Code: "DIGIFLAZZ_IP_BLOCKED", Message: "IP tidak dikenali", HTTPStatus: 403, Retryable: false},
	"DIGIFLAZZ_REF_ID_DUPLICATE":  {Code: "DIGIFLAZZ_REF_ID_DUPLICATE", Message: "ID transaksi sudah digunakan", HTTPStatus: 409, Retryable: false},
	"DIGIFLAZZ_NUMBER_BLOCKED":   {Code: "DIGIFLAZZ_NUMBER_BLOCKED", Message: "Nomor tujuan diblokir", HTTPStatus: 400, Retryable: false},
	"DIGIFLAZZ_PREFIX_MISMATCH":    {Code: "DIGIFLAZZ_PREFIX_MISMATCH", Message: "Prefix nomor tidak sesuai operator", HTTPStatus: 400, Retryable: false},
	"DIGIFLAZZ_PRODUCT_UNAVAILABLE": {Code: "DIGIFLAZZ_PRODUCT_UNAVAILABLE", Message: "Produk tidak tersedia saat ini", HTTPStatus: 503, Retryable: false},
	"DIGIFLAZZ_INVALID_NUMBER":    {Code: "DIGIFLAZZ_INVALID_NUMBER", Message: "Nomor tujuan salah", HTTPStatus: 400, Retryable: false},
	"DIGIFLAZZ_PRODUCT_DISRUPTION": {Code: "DIGIFLAZZ_PRODUCT_DISRUPTION", Message: "Produk sedang gangguan", HTTPStatus: 503, Retryable: true},
	"DIGIFLAZZ_CUT_OFF":            {Code: "DIGIFLAZZ_CUT_OFF", Message: "Transaksi cut-off, coba lagi dalam 15 menit", HTTPStatus: 503, Retryable: true},
	"DIGIFLAZZ_BILL_NOT_AVAILABLE": {Code: "DIGIFLAZZ_BILL_NOT_AVAILABLE", Message: "Tagihan belum tersedia", HTTPStatus: 404, Retryable: false},
	"DIGIFLAZZ_SELLER_NOT_VERIFIED": {Code: "DIGIFLAZZ_SELLER_NOT_VERIFIED", Message: "Seller belum ter-verifikasi", HTTPStatus: 403, Retryable: false},
	"DIGIFLAZZ_OUT_OF_STOCK":       {Code: "DIGIFLAZZ_OUT_OF_STOCK", Message: "Stok habis", HTTPStatus: 503, Retryable: false},
	"DIGIFLAZZ_PRICE_MISMATCH":     {Code: "DIGIFLAZZ_PRICE_MISMATCH", Message: "Harga telah berubah, silakan sinkronisasi ulang", HTTPStatus: 400, Retryable: true},
	"DIGIFLAZZ_BILLER_TIMEOUT":     {Code: "DIGIFLAZZ_BILLER_TIMEOUT", Message: "Timeout dari biller, mencoba ulang...", HTTPStatus: 502, Retryable: true},
	"DIGIFLAZZ_PRODUCT_UNSTABLE":  {Code: "DIGIFLAZZ_PRODUCT_UNSTABLE", Message: "Produk tidak stabil", HTTPStatus: 502, Retryable: true},
	"DIGIFLAZZ_PRICELIST_LIMIT":    {Code: "DIGIFLAZZ_PRICELIST_LIMIT", Message: "Limit API pencarian harga, coba dalam 4 menit", HTTPStatus: 429, Retryable: true},
	"DIGIFLAZZ_RATE_LIMIT":         {Code: "DIGIFLAZZ_RATE_LIMIT", Message: "Terlalu banyak permintaan, coba 1 menit lagi", HTTPStatus: 429, Retryable: true},
	"DIGIFLAZZ_PLN_INQUIRY_LIMIT":  {Code: "DIGIFLAZZ_PLN_INQUIRY_LIMIT", Message: "Limit cek PLN, coba lagi nanti", HTTPStatus: 429, Retryable: true},
	"DIGIFLAZZ_EMONEY_MULTIPLE":   {Code: "DIGIFLAZZ_EMONEY_MULTIPLE", Message: "Nominal e-money harus kelipatan Rp 1.000", HTTPStatus: 400, Retryable: false},
	"DIGIFLAZZ_ACCOUNT_BLOCKED":  {Code: "DIGIFLAZZ_ACCOUNT_BLOCKED", Message: "Akun diblokir oleh operator", HTTPStatus: 403, Retryable: false},
	"DIGIFLAZZ_UNKNOWN_ERROR":     {Code: "DIGIFLAZZ_UNKNOWN_ERROR", Message: "Terjadi kesalahan dari provider", HTTPStatus: 502, Retryable: false},

	// VALIDATION Errors
	"VALIDATION_PHONE_FORMAT":       {Code: "VALIDATION_PHONE_FORMAT", Message: "Format nomor HP tidak valid (contoh: +6281234567890)", HTTPStatus: 400, Retryable: false},
	"VALIDATION_PIN_FORMAT":         {Code: "VALIDATION_PIN_FORMAT", Message: "PIN harus 6 digit angka", HTTPStatus: 400, Retryable: false},
	"VALIDATION_PIN_SEQUENTIAL":     {Code: "VALIDATION_PIN_SEQUENTIAL", Message: "PIN tidak boleh berurutan (123456, 654321)", HTTPStatus: 400, Retryable: false},
	"VALIDATION_PASSWORD_WEAK":      {Code: "VALIDATION_PASSWORD_WEAK", Message: "Password minimal 8 karakter, mengandung huruf besar & angka", HTTPStatus: 400, Retryable: false},
	"VALIDATION_CUSTOMER_NO_FORMAT": {Code: "VALIDATION_CUSTOMER_NO_FORMAT", Message: "Format nomor tujuan tidak sesuai", HTTPStatus: 400, Retryable: false},
	"VALIDATION_AMOUNT_NEGATIVE":   {Code: "VALIDATION_AMOUNT_NEGATIVE", Message: "Jumlah tidak boleh negatif", HTTPStatus: 400, Retryable: false},
	"VALIDATION_MISSING_FIELD":      {Code: "VALIDATION_MISSING_FIELD", Message: "Field '%s' wajib diisi", HTTPStatus: 400, Retryable: false},
	"VALIDATION_JSON_INVALID":      {Code: "VALIDATION_JSON_INVALID", Message: "Format JSON tidak valid", HTTPStatus: 400, Retryable: false},

	// SYSTEM Errors
	"SYSTEM_DB_UNAVAILABLE":     {Code: "SYSTEM_DB_UNAVAILABLE", Message: "Database tidak tersedia, coba beberapa saat lagi", HTTPStatus: 503, Retryable: true},
	"SYSTEM_REDIS_UNAVAILABLE":   {Code: "SYSTEM_REDIS_UNAVAILABLE", Message: "Cache tidak tersedia", HTTPStatus: 503, Retryable: true},
	"SYSTEM_VAULT_UNAVAILABLE":   {Code: "SYSTEM_VAULT_UNAVAILABLE", Message: "Konfigurasi tidak dapat diakses", HTTPStatus: 503, Retryable: true},
	"SYSTEM_INTERNAL":           {Code: "SYSTEM_INTERNAL", Message: "Terjadi kesalahan internal", HTTPStatus: 500, Retryable: false},
	"SYSTEM_TIMEOUT":             {Code: "SYSTEM_TIMEOUT", Message: "Request timeout", HTTPStatus: 504, Retryable: true},
	"SYSTEM_RATE_LIMIT":           {Code: "SYSTEM_RATE_LIMIT", Message: "Server sedang sibuk, coba lagi nanti", HTTPStatus: 429, Retryable: true},
	"SYSTEM_CIRCUIT_OPEN":        {Code: "SYSTEM_CIRCUIT_OPEN", Message: "Layanan sedang gangguan", HTTPStatus: 503, Retryable: true},
	"SYSTEM_MAINTENANCE":         {Code: "SYSTEM_MAINTENANCE", Message: "Sistem dalam pemeliharaan", HTTPStatus: 503, Retryable: false},
}

func GetError(code string) ErrorCode {
	if err, ok := ErrorCatalog[code]; ok {
		return err
	}
	return ErrorCatalog["SYSTEM_INTERNAL"]
}

func NewAppError(code string, details map[string]interface{}) *AppError {
	errInfo := GetError(code)
	return &AppError{
		Code:       code,
		Message:    errInfo.Message,
		Details:    details,
		TraceID:    GenerateTraceID(),
		StatusCode: errInfo.HTTPStatus,
		Retryable:  errInfo.Retryable,
		RetryAfter: RetryAfterForCode(code),
		Timestamp:  CurrentTimestamp(),
	}
}

func GenerateTraceID() string {
	return uuid.New().String()
}

func CurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}