package com.yonotech.ppob.data.remote.model

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

// ============== Auth DTOs ==============

@JsonClass(generateAdapter = true)
data class RegisterRequest(
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "name") val name: String
)

@JsonClass(generateAdapter = true)
data class VerifyOtpRequest(
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "otp_code") val otpCode: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String
)

@JsonClass(generateAdapter = true)
data class LoginRequest(
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String,
    @Json(name = "device_fingerprint") val deviceFingerprint: DeviceFingerprint? = null,
    @Json(name = "otp_code") val otpCode: String? = null
)

@JsonClass(generateAdapter = true)
data class DeviceFingerprint(
    @Json(name = "device_id") val deviceId: String,
    @Json(name = "user_agent") val userAgent: String,
    @Json(name = "app_version") val appVersion: String,
    @Json(name = "install_ts") val installTs: Long,
    @Json(name = "last_login_ts") val lastLoginTs: Long? = null
)

@JsonClass(generateAdapter = true)
data class RefreshTokenRequest(
    @Json(name = "refresh_token") val refreshToken: String
)

@JsonClass(generateAdapter = true)
data class ChangePinRequest(
    @Json(name = "current_pin") val currentPin: String? = null,
    @Json(name = "password") val password: String? = null,
    @Json(name = "otp_code") val otpCode: String? = null,
    @Json(name = "new_pin") val newPin: String,
    @Json(name = "confirm_pin") val confirmPin: String
)

@JsonClass(generateAdapter = true)
data class AuthResponse(
    @Json(name = "access_token") val accessToken: String,
    @Json(name = "refresh_token") val refreshToken: String,
    @Json(name = "token_type") val tokenType: String,
    @Json(name = "expires_in") val expiresIn: Int,
    @Json(name = "user") val user: UserResponse
)

@JsonClass(generateAdapter = true)
data class UserResponse(
    @Json(name = "user_id") val userId: String,
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "name") val name: String,
    @Json(name = "roles") val roles: List<RoleResponse>,
    @Json(name = "active_role") val activeRole: RoleResponse,
    @Json(name = "wallet_id") val walletId: String
)

@JsonClass(generateAdapter = true)
data class RoleResponse(
    @Json(name = "role_id") val roleId: String,
    @Json(name = "role_name") val roleName: String
)

// ============== ApiResponse Wrapper ==============

@JsonClass(generateAdapter = true)
data class ApiResponse<T>(
    @Json(name = "success") val success: Boolean? = null,
    @Json(name = "data") val data: T? = null,
    @Json(name = "message") val message: String? = null,
    @Json(name = "error") val error: ApiError? = null
)

@JsonClass(generateAdapter = true)
data class ApiError(
    @Json(name = "code") val code: String,
    @Json(name = "message") val message: String,
    @Json(name = "details") val details: Map<String, Any>? = null,
    @Json(name = "trace_id") val traceId: String? = null,
    @Json(name = "timestamp") val timestamp: String? = null
)

@JsonClass(generateAdapter = true)
data class ApiErrorResponse(
    @Json(name = "error") val error: ApiError
)

// ============== User DTOs ==============

@JsonClass(generateAdapter = true)
data class UserProfileResponse(
    @Json(name = "user_id") val userId: String,
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "name") val name: String,
    @Json(name = "roles") val roles: List<RoleResponse>,
    @Json(name = "active_role") val activeRole: RoleResponse,
    @Json(name = "wallet") val wallet: WalletBalanceResponse
)

@JsonClass(generateAdapter = true)
data class SwitchRoleRequest(
    @Json(name = "role_id") val roleId: String
)

@JsonClass(generateAdapter = true)
data class SwitchRoleResponse(
    @Json(name = "access_token") val accessToken: String,
    @Json(name = "refresh_token") val refreshToken: String,
    @Json(name = "expires_in") val expiresIn: Int,
    @Json(name = "active_role") val activeRole: RoleResponse,
    @Json(name = "wallet_id") val walletId: String
)

// ============== Staff DTOs ==============

@JsonClass(generateAdapter = true)
data class AddStaffRequest(
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "name") val name: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String,
    @Json(name = "margin_scheme") val marginScheme: String = "FixedAllowance",
    @Json(name = "margin_value") val marginValue: Double = 10000.0,
    @Json(name = "daily_txn_limit") val dailyTxnLimit: Int = 50,
    @Json(name = "daily_amount_limit") val dailyAmountLimit: Double = 5000000.0
)

@JsonClass(generateAdapter = true)
data class StaffItem(
    @Json(name = "user_id") val userId: String,
    @Json(name = "name") val name: String,
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "wallet_balance") val walletBalance: Double,
    @Json(name = "daily_txn_count") val dailyTxnCount: Int,
    @Json(name = "daily_txn_amount") val dailyTxnAmount: Double,
    @Json(name = "margin_scheme") val marginScheme: String,
    @Json(name = "margin_value") val marginValue: Double,
    @Json(name = "is_active") val isActive: Boolean
)

@JsonClass(generateAdapter = true)
data class StaffListResponse(
    @Json(name = "staff") val staff: List<StaffItem>,
    @Json(name = "pagination") val pagination: Pagination
)

@JsonClass(generateAdapter = true)
data class UpdateStaffRequest(
    @Json(name = "name") val name: String? = null,
    @Json(name = "margin_scheme") val marginScheme: String? = null,
    @Json(name = "margin_value") val marginValue: Double? = null,
    @Json(name = "daily_txn_limit") val dailyTxnLimit: Int? = null,
    @Json(name = "daily_amount_limit") val dailyAmountLimit: Double? = null,
    @Json(name = "is_active") val isActive: Boolean? = null
)

@JsonClass(generateAdapter = true)
data class StaffDetailResponse(
    @Json(name = "user_id") val userId: String,
    @Json(name = "phone_number") val phoneNumber: String,
    @Json(name = "name") val name: String,
    @Json(name = "margin_scheme") val marginScheme: String,
    @Json(name = "margin_value") val marginValue: Double,
    @Json(name = "daily_txn_limit") val dailyTxnLimit: Int,
    @Json(name = "daily_amount_limit") val dailyAmountLimit: Double,
    @Json(name = "is_active") val isActive: Boolean
)

@JsonClass(generateAdapter = true)
data class Pagination(
    @Json(name = "total") val total: Int,
    @Json(name = "limit") val limit: Int,
    @Json(name = "offset") val offset: Int,
    @Json(name = "has_more") val hasMore: Boolean
)

// ============== Product DTOs ==============

@JsonClass(generateAdapter = true)
data class Category(
    @Json(name = "category_id") val categoryId: String,
    @Json(name = "category_name") val categoryName: String,
    @Json(name = "icon_url") val iconUrl: String,
    @Json(name = "display_order") val displayOrder: Int
)

@JsonClass(generateAdapter = true)
data class CategoriesResponse(
    @Json(name = "categories") val categories: List<Category>
)

@JsonClass(generateAdapter = true)
data class Product(
    @Json(name = "product_id") val productId: String,
    @Json(name = "buyer_sku_code") val buyerSkuCode: String,
    @Json(name = "product_name") val productName: String,
    @Json(name = "base_price") val basePrice: Double,
    @Json(name = "platform_price") val platformPrice: Double,
    @Json(name = "is_prepaid") val isPrepaid: Boolean,
    @Json(name = "mitra_selling_price") val mitraSellingPrice: Double? = null
)

@JsonClass(generateAdapter = true)
data class ProductsResponse(
    @Json(name = "products") val products: List<Product>,
    @Json(name = "pagination") val pagination: Pagination? = null
)

@JsonClass(generateAdapter = true)
data class SyncStatusResponse(
    @Json(name = "last_sync_prepaid") val lastSyncPrepaid: String? = null,
    @Json(name = "last_sync_postpaid") val lastSyncPostpaid: String? = null
)

// ============== Wallet DTOs ==============

@JsonClass(generateAdapter = true)
data class WalletBalanceResponse(
    @Json(name = "wallet_id") val walletId: String,
    @Json(name = "balance_available") val balanceAvailable: Double,
    @Json(name = "balance_held") val balanceHeld: Double,
    @Json(name = "balance_total") val balanceTotal: Double? = null,
    @Json(name = "currency") val currency: String = "IDR",
    @Json(name = "updated_at") val updatedAt: String
)

@JsonClass(generateAdapter = true)
data class TopUpRequest(
    @Json(name = "amount") val amount: Double
)

@JsonClass(generateAdapter = true)
data class TopUpResponse(
    @Json(name = "message") val message: String,
    @Json(name = "staff_wallet") val staffWallet: WalletBalanceResponse,
    @Json(name = "mitra_wallet") val mitraWallet: WalletBalanceResponse? = null
)

// ============== Transaction DTOs ==============

@JsonClass(generateAdapter = true)
data class InitiateTransactionRequest(
    @Json(name = "product_id") val productId: String,
    @Json(name = "customer_no") val customerNo: String,
    @Json(name = "selling_price") val sellingPrice: Double,
    @Json(name = "pin") val pin: String,
    @Json(name = "device_fingerprint") val deviceFingerprint: String? = null
)

@JsonClass(generateAdapter = true)
data class TransactionResponse(
    @Json(name = "transaction_id") val transactionId: String,
    @Json(name = "ref_id") val refId: String? = null,
    @Json(name = "status") val status: String,
    @Json(name = "message") val message: String? = null,
    @Json(name = "selling_price") val sellingPrice: Double,
    @Json(name = "platform_price") val platformPrice: Double,
    @Json(name = "commission_amount") val commissionAmount: Double? = null
)

@JsonClass(generateAdapter = true)
data class TransactionSummary(
    @Json(name = "id") val id: String,
    @Json(name = "ref_id") val refId: String?,
    @Json(name = "product_name") val productName: String,
    @Json(name = "customer_no") val customerNumber: String,
    @Json(name = "status") val status: String,
    @Json(name = "selling_price") val sellingPrice: Double,
    @Json(name = "platform_price") val platformPrice: Double,
    @Json(name = "margin_amount") val marginAmount: Double?,
    @Json(name = "commission_amount") val commissionAmount: Double?,
    @Json(name = "created_at") val createdAt: String,
    @Json(name = "completed_at") val completedAt: String? = null
)

@JsonClass(generateAdapter = true)
data class TransactionHistoryResponse(
    @Json(name = "transactions") val transactions: List<TransactionSummary>,
    @Json(name = "pagination") val pagination: Pagination? = null
)

@JsonClass(generateAdapter = true)
data class TransactionDetailResponse(
    @Json(name = "transaction_id") val transactionId: String,
    @Json(name = "ref_id") val refId: String?,
    @Json(name = "status") val status: String,
    @Json(name = "message") val message: String?,
    @Json(name = "serial_number") val serialNumber: String? = null,
    @Json(name = "details") val details: TransactionDetailData? = null
)

@JsonClass(generateAdapter = true)
data class TransactionDetailData(
    @Json(name = "product_name") val productName: String,
    @Json(name = "customer_no") val customerNumber: String,
    @Json(name = "amount") val amount: Double,
    @Json(name = "selling_price") val sellingPrice: Double,
    @Json(name = "commission_amount") val commissionAmount: Double?,
    @Json(name = "margin_amount") val marginAmount: Double?,
    @Json(name = "created_at") val createdAt: String,
    @Json(name = "completed_at") val completedAt: String? = null
)

// ============== Report DTOs ==============

@JsonClass(generateAdapter = true)
data class ReportKPIs(
    @Json(name = "total_sales") val totalSales: Double,
    @Json(name = "platform_profit") val platformProfit: Double,
    @Json(name = "staff_count") val staffCount: Int,
    @Json(name = "success_rate") val successRate: Double
)

@JsonClass(generateAdapter = true)
data class SalesTrendData(
    @Json(name = "date") val date: String,
    @Json(name = "total_amount") val totalAmount: Double,
    @Json(name = "transaction_count") val transactionCount: Int
)

@JsonClass(generateAdapter = true)
data class StaffPerformanceData(
    @Json(name = "staff_name") val staffName: String,
    @Json(name = "transaction_count") val transactionCount: Int,
    @Json(name = "total_amount") val totalAmount: Double
)

// ============== Webhook DTOs ==============

@JsonClass(generateAdapter = true)
data class DigiflazzWebhook(
    @Json(name = "ref_id") val refId: String,
    @Json(name = "rc") val rc: String,
    @Json(name = "rc_message") val rcMessage: String,
    @Json(name = "customer_no") val customerNo: String,
    @Json(name = "sn") val serialNumber: String?,
    @Json(name = "price") val price: Double? = null,
    @Json(name = "selling_price") val sellingPrice: Double? = null
)

// ============== Notification ==============

@JsonClass(generateAdapter = true)
data class NotificationData(
    @Json(name = "type") val type: String,
    @Json(name = "transaction_id") val transactionId: String? = null,
    @Json(name = "amount") val amount: Double? = null,
    @Json(name = "selling_price") val sellingPrice: Double? = null,
    @Json(name = "product_name") val productName: String? = null,
    @Json(name = "customer_no") val customerNo: String? = null,
    @Json(name = "deep_link") val deepLink: String? = null
)